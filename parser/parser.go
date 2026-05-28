package parser

import (
	"github.com/hanzala211/goinkinterpreter/expr"
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

func (p *Parser) Parse() (expr.Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (expr.Expr, error) {
	return p.equality()
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
		return expr.UnaryExpr{
			Operator: operator,
			Right:    ex,
		}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (expr.Expr, error) {
	if p.match(token.TokenType_False) {
		return expr.LiteralExpr{
			Value: false,
		}, nil
	}
	if p.match(token.TokenType_True) {
		return expr.LiteralExpr{
			Value: true,
		}, nil
	}
	if p.match(token.TokenType_Number) {
		return expr.LiteralExpr{
			Value: p.previous().Literal,
		}, nil
	}
	if p.match(token.TokenType_String) {
		return expr.LiteralExpr{
			Value: p.previous().Literal,
		}, nil
	}
	if p.match(token.TokenType_Nil) {
		return expr.LiteralExpr{
			Value: nil,
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
		return expr.Grouping{
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
