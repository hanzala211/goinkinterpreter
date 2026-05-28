package evaluator

import (
	"fmt"

	"github.com/hanzala211/goinkinterpreter/expr"
	"github.com/hanzala211/goinkinterpreter/stmt"
	"github.com/hanzala211/goinkinterpreter/token"
)

type Evaluator struct {
	env *Environment
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		env: NewEnvironment(),
	}
}

func (e *Evaluator) Interpret(stmts []stmt.Stmt) error {
	for _, stmt := range stmts {
		err := e.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ev *Evaluator) execute(statement stmt.Stmt) error {
	switch s := statement.(type) {
	case *stmt.ExprStmt:
		_, err := ev.Eval(s.Expr)
		return err
	case *stmt.VarStmt:
		var value any = nil
		var err error
		if s.Initalizer != nil {
			value, err = ev.Eval(s.Initalizer)
			if err != nil {
				return err
			}
		}
		ev.env.Set(s.Name.Lexeme, value)
		return nil
	case *stmt.PrintStmt:
		value, err := ev.Eval(s.Expr)
		if err != nil {
			return err
		}
		fmt.Println(value)
		return nil
	case *stmt.BlockStmt:
		oldEnv := ev.env
		ev.env = NewEnvironmentWithParent(oldEnv)
		for _, stmt := range s.Statements {
			err := ev.execute(stmt)
			if err != nil {
				ev.env = oldEnv
				return err
			}
		}
		ev.env = oldEnv
		return nil
	}
	return nil
}

func (ev *Evaluator) Eval(e expr.Expr) (any, error) {
	switch e := e.(type) {
	case *expr.LiteralExpr:
		return e.Value, nil
	case *expr.VarExpr:
		return ev.env.Get(e.Name.Lexeme), nil
	case *expr.AssignExpr:
		value, err := ev.Eval(e.Value)
		if err != nil {
			return nil, err
		}
		ev.env.Assign(e.Name.Lexeme, value)
		return value, nil
	case *expr.Grouping:
		return ev.Eval(e.Expression)
	case *expr.UnaryExpr:
		right, err := ev.Eval(e.Right)
		if err != nil {
			return nil, err
		}
		switch e.Operator.Type {
		case token.TokenType_Bang:
			return !isTruthy(right), nil
		case token.TokenType_Minus:
			if val, ok := right.(float64); ok {
				return -val, nil
			}
			return nil, fmt.Errorf("operand must be a number")
		}
	case *expr.BinaryExpr:
		left, err := ev.Eval(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := ev.Eval(e.Right)
		if err != nil {
			return nil, err
		}
		switch e.Operator.Type {
		case token.TokenType_Plus:
			if left, ok := left.(float64); ok {
				if right, ok := right.(float64); ok {
					return left + right, nil
				}
			}
			if leftStr, ok1 := left.(string); ok1 {
				if rightStr, ok2 := right.(string); ok2 {
					return leftStr + rightStr, nil
				}
			}
			return nil, fmt.Errorf("operands must be numbers")
		case token.TokenType_Minus:
			if left, ok := left.(float64); ok {
				if right, ok := right.(float64); ok {
					return left - right, nil
				}
			}
			return nil, fmt.Errorf("operands must be numbers")
		case token.TokenType_Star:
			if left, ok := left.(float64); ok {
				if right, ok := right.(float64); ok {
					return left * right, nil
				}
			}
			return nil, fmt.Errorf("operands must be numbers")
		case token.TokenType_Slash:
			if left, ok := left.(float64); ok {
				if right, ok := right.(float64); ok {
					return left / right, nil
				}
			}
			return nil, fmt.Errorf("operands must be numbers")
		case token.TokenType_Greater:
			if left, ok := left.(float64); ok {
				if right, ok := right.(float64); ok {
					return left > right, nil
				}
			}
			return nil, fmt.Errorf("operands must be numbers")
		case token.TokenType_Less:
			if left, ok := left.(float64); ok {
				if right, ok := right.(float64); ok {
					return left < right, nil
				}
			}
			return nil, fmt.Errorf("operands must be numbers")
		case token.TokenType_GreaterEqual:
			if left, ok := left.(float64); ok {
				if right, ok := right.(float64); ok {
					return left >= right, nil
				}
			}
			return nil, fmt.Errorf("operands must be numbers")
		case token.TokenType_LessEqual:
			if left, ok := left.(float64); ok {
				if right, ok := right.(float64); ok {
					return left <= right, nil
				}
			}
			return nil, fmt.Errorf("operands must be numbers")
		case token.TokenType_EqualEqual:
			return left == right, nil
		case token.TokenType_BangEqual:
			return left != right, nil
		default:
			return nil, fmt.Errorf("invalid operator")
		}
	case *expr.LogicalExpr: // it uses short circuit evaluation for example if left side is true and operator is or then right side is not evaluated if left side is false and operator is and then right side is not evaluated it will be evaluated in the opposite cases of above 2
		left, err := ev.Eval(e.Left)
		if err != nil {
			return nil, err
		}
		if e.Operator.Type == token.TokenType_Or {
			if isTruthy(left) {
				return left, nil
			}
		}
		if e.Operator.Type == token.TokenType_And {
			if !isTruthy(left) {
				return left, nil
			}
		}
		return ev.Eval(e.Right)
	}
	return nil, nil
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if b, ok := value.(bool); ok {
		return b
	}
	return true
}
