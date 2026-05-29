package parser

import (
	"github.com/hanzala211/goinkinterpreter/expr"
	"github.com/hanzala211/goinkinterpreter/stmt"
	"github.com/hanzala211/goinkinterpreter/token"
)

type vm interface {
	ReportParserError(err ParserError)
}

type Parser struct {
	vm      vm
	Tokens  []*token.Token
	current int
}

type ParserError struct {
	Token   *token.Token
	Message string
}

func (e ParserError) Error() string {
	return e.Message
}

func NewParser(tokens []*token.Token, vm vm) *Parser {
	return &Parser{
		vm:      vm,
		Tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ([]stmt.Stmt, error) {
	var stmts []stmt.Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *Parser) declaration() (stmt.Stmt, error) {
	if p.match(token.TokenType_Var) {
		return p.varDeclaration()
	}
	if p.match(token.TokenType_Fun) {
		return p.functionDeclaration("function")
	}
	return p.statement()
}

func (p *Parser) functionDeclaration(keyword string) (stmt.Stmt, error) {
	name, err := p.consume(token.TokenType_Identifier, "Expect "+keyword+" name.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.TokenType_LeftParen, "Expect '(' after "+keyword+" name.")
	if err != nil {
		return nil, err
	}
	var params []*token.Token
	for !p.check(token.TokenType_RightParen) {
		for {
			param, err := p.consume(token.TokenType_Identifier, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			params = append(params, param)
			if len(params) >= 255 {
				return nil, ParserError{
					Token:   p.peek(),
					Message: "Maximum number of parameters exceeded.",
				}
			}
			if !p.match(token.TokenType_Comma) {
				break
			}
		}
	}
	_, err = p.consume(token.TokenType_RightParen, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.TokenType_LeftBrace, "Expect '{' before function body.")
	if err != nil {
		return nil, err
	}
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}
	return &stmt.FuncStmt{
		Name:   name,
		Params: params,
		Body:   body.(*stmt.BlockStmt).Statements,
	}, nil
}

func (p *Parser) varDeclaration() (stmt.Stmt, error) {
	name, err := p.consume(token.TokenType_Identifier, "Expected variable name")
	if err != nil {
		return nil, err
	}
	var initializer expr.Expr
	if p.match(token.TokenType_Equal) {
		initializer, err = p.expression() // this is the right side of the equal
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.TokenType_Semicolon, "Expected ';' after variable declaration")
	if err != nil {
		return nil, err
	}
	return &stmt.VarStmt{
		Name:       name,
		Initalizer: initializer,
	}, nil
}

func (p *Parser) statement() (stmt.Stmt, error) {
	if p.match(token.TokenType_Print) {
		return p.printStatement()
	} else if p.match(token.TokenType_LeftBrace) {
		return p.blockStatement()
	} else if p.match(token.TokenType_If) {
		return p.ifStatement()
	} else if p.match(token.TokenType_While) {
		return p.whileStatement()
	} else if p.match(token.TokenType_For) {
		return p.forStatement()
	} else if p.match(token.TokenType_Return) {
		return p.returnStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) returnStatement() (stmt.Stmt, error) {
	keyword := p.previous()
	var value expr.Expr = nil
	if !p.check(token.TokenType_Semicolon) {
		var err error
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err := p.consume(token.TokenType_Semicolon, "Expected ';' after return statement")
	if err != nil {
		return nil, err
	}
	return &stmt.ReturnStmt{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *Parser) forStatement() (stmt.Stmt, error) {
	_, err := p.consume(token.TokenType_LeftParen, "Expected '(' after 'for'")
	if err != nil {
		return nil, err
	}
	var initializer stmt.Stmt
	if p.match(token.TokenType_Semicolon) { // No initializer e.g for (; ...)
		initializer = nil
	} else if p.match(token.TokenType_Var) { // e.g for (var ...)
		initializer, err = p.varDeclaration()
	} else { // e.g for (i = 0...)
		initializer, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}
	var condition expr.Expr = nil
	if !p.check(token.TokenType_Semicolon) { // for checking if there is a condition e.g for (var i = 0; i < 10; ...) if not then skip fetching the condition
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.TokenType_Semicolon, "Expected ';' after for condition")
	if err != nil {
		return nil, err
	}
	var increment expr.Expr = nil
	if !p.check(token.TokenType_RightParen) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.TokenType_RightParen, "Expected ')' after for increment")
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	if increment != nil {
		body = &stmt.BlockStmt{
			Statements: []stmt.Stmt{
				body,
				&stmt.ExprStmt{
					Expr: increment,
				},
			},
		}
	}
	if condition == nil {
		condition = &expr.LiteralExpr{
			Value: true,
		}
	}
	body = &stmt.WhileStmt{
		Condition: condition,
		Body:      body,
	}
	if initializer != nil {
		body = &stmt.BlockStmt{
			Statements: []stmt.Stmt{
				initializer,
				body,
			},
		}
	}
	return body, nil
}

func (p *Parser) whileStatement() (stmt.Stmt, error) {
	_, err := p.consume(token.TokenType_LeftParen, "Expected '(' after 'while'")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.TokenType_RightParen, "Expected ')' after 'while' condition")
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &stmt.WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) ifStatement() (stmt.Stmt, error) {
	_, err := p.consume(token.TokenType_LeftParen, "Expected '(' after 'if'")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.TokenType_RightParen, "Expected ')' after 'if' condition")
	if err != nil {
		return nil, err
	}
	thenStmt, err := p.statement() // this is called for block stmt and what it does that it allows one liner if we wanted only braces if control flow we would call p.blockStatement()
	if err != nil {
		return nil, err
	}
	var elseStmt stmt.Stmt
	if p.match(token.TokenType_Else) {
		elseStmt, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &stmt.IfStmt{
		Condition: condition,
		Then:      thenStmt,
		Else:      elseStmt,
	}, nil
}

func (p *Parser) blockStatement() (stmt.Stmt, error) {
	var statements []stmt.Stmt
	for !p.check(token.TokenType_RightBrace) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	_, err := p.consume(token.TokenType_RightBrace, "Expected '}' after block")
	if err != nil {
		return nil, err
	}
	return &stmt.BlockStmt{
		Statements: statements,
	}, nil
}

func (p *Parser) printStatement() (stmt.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.TokenType_Semicolon, "Expected ';' after print statement")
	if err != nil {
		return nil, err
	}
	return &stmt.PrintStmt{
		Expr: value,
	}, nil
}

func (p *Parser) expressionStatement() (stmt.Stmt, error) {
	ex, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.TokenType_Semicolon, "Expected ';' after expression statement")
	if err != nil {
		return nil, err
	}
	return &stmt.ExprStmt{
		Expr: ex,
	}, nil
}

func (p *Parser) expression() (expr.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (expr.Expr, error) {
	ex, err := p.or()
	if err != nil {
		return ex, err
	}
	if p.match(token.TokenType_Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return value, err
		}
		if _, ok := ex.(*expr.VarExpr); ok {
			return &expr.AssignExpr{
				Name:  ex.(*expr.VarExpr).Name,
				Value: value,
			}, nil
		}
		return nil, ParserError{
			Token:   equals,
			Message: "invalid assignment target",
		}
	}
	return ex, nil
}

func (p *Parser) or() (expr.Expr, error) {
	ex, err := p.and()
	if err != nil {
		return ex, err
	}
	for p.match(token.TokenType_Or) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return right, err
		}
		ex = &expr.LogicalExpr{
			Left:     ex,
			Operator: operator,
			Right:    right,
		}
	}
	return ex, nil
}

func (p *Parser) and() (expr.Expr, error) {
	ex, err := p.equality()
	if err != nil {
		return ex, err
	}
	for p.match(token.TokenType_And) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return right, err
		}
		ex = &expr.LogicalExpr{
			Left:     ex,
			Operator: operator,
			Right:    right,
		}
	}
	return ex, nil
}

func (p *Parser) equality() (expr.Expr, error) {
	ex, err := p.comparison()
	if err != nil {
		return ex, err
	}
	for p.match(token.TokenType_EqualEqual, token.TokenType_BangEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return right, err
		}
		ex = &expr.BinaryExpr{
			Left:     ex,
			Operator: operator,
			Right:    right,
		}
	}
	return ex, nil
}

func (p *Parser) comparison() (expr.Expr, error) {
	ex, err := p.term()
	if err != nil {
		return ex, err
	}
	for p.match(token.TokenType_Greater, token.TokenType_Less, token.TokenType_GreaterEqual, token.TokenType_LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return right, err
		}
		ex = &expr.BinaryExpr{
			Left:     ex,
			Operator: operator,
			Right:    right,
		}
	}
	return ex, nil
}

func (p *Parser) term() (expr.Expr, error) {
	ex, err := p.factor()
	if err != nil {
		return ex, err
	}
	for p.match(token.TokenType_Plus, token.TokenType_Minus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return right, err
		}
		ex = &expr.BinaryExpr{
			Left:     ex,
			Operator: operator,
			Right:    right,
		}
	}
	return ex, nil
}

func (p *Parser) factor() (expr.Expr, error) {
	ex, err := p.unary()
	if err != nil {
		return ex, err
	}
	for p.match(token.TokenType_Star, token.TokenType_Slash) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return right, err
		}
		ex = &expr.BinaryExpr{
			Left:     ex,
			Operator: operator,
			Right:    right,
		}
	}
	return ex, nil
}

func (p *Parser) unary() (expr.Expr, error) {
	if p.match(token.TokenType_Bang, token.TokenType_Minus) {
		operator := p.previous()
		ex, err := p.unary()
		if err != nil {
			return ex, err
		}
		return &expr.UnaryExpr{
			Operator: operator,
			Right:    ex,
		}, nil
	}
	return p.call()
}

func (p *Parser) call() (expr.Expr, error) {
	ex, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if p.match(token.TokenType_LeftParen) {
			ex, err = p.finishCall(ex)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return ex, nil
}

func (p *Parser) finishCall(callee expr.Expr) (expr.Expr, error) {
	var arguments []expr.Expr
	if !p.check(token.TokenType_RightParen) {
		for {
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)

			if len(arguments) >= 255 {
				return nil, ParserError{Token: p.peek(), Message: "Can't have more than 255 arguments."}
			}
			if !p.match(token.TokenType_Comma) {
				break
			}
		}
	}

	paren, err := p.consume(token.TokenType_RightParen, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &expr.CallExpr{
		Callee: callee,
		Paren:  paren,
		Args:   arguments,
	}, nil
}

func (p *Parser) primary() (expr.Expr, error) {
	if p.match(token.TokenType_False) {
		return &expr.LiteralExpr{
			Value: false,
		}, nil
	}
	if p.match(token.TokenType_True) {
		return &expr.LiteralExpr{
			Value: true,
		}, nil
	}
	if p.match(token.TokenType_Number) {
		return &expr.LiteralExpr{
			Value: p.previous().Literal,
		}, nil
	}
	if p.match(token.TokenType_String) {
		return &expr.LiteralExpr{
			Value: p.previous().Literal,
		}, nil
	}
	if p.match(token.TokenType_Nil) {
		return &expr.LiteralExpr{
			Value: nil,
		}, nil
	}
	if p.match(token.TokenType_Identifier) {
		return &expr.VarExpr{
			Name: p.previous(),
		}, nil
	}
	if p.match(token.TokenType_LeftParen) {
		ex, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.TokenType_RightParen, "Expected ')' after expression")
		if err != nil {
			return nil, err
		}
		return &expr.Grouping{
			Expression: ex,
		}, nil
	}
	return nil, ParserError{
		Token:   p.peek(),
		Message: "unexpected token",
	}
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.TokenType_EOF
}

func (p *Parser) peek() *token.Token {
	return p.Tokens[p.current]
}

func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) previous() *token.Token {
	if p.current == 0 {
		return nil
	}
	return p.Tokens[p.current-1]
}

func (p *Parser) consume(t token.TokenType, message string) (*token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return nil, ParserError{
		Token:   p.peek(),
		Message: message,
	}
}
