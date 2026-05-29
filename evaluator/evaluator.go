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

type InkCallable interface {
	Call(ev *Evaluator, args []any) (any, error)
	Arity() int
}

type ReturnError struct {
	Value any
}

func (e *ReturnError) Error() string {
	return "return"
}

type InkFunc struct {
	Declaration *stmt.FuncStmt
	Closure     *Environment
}

func (f *InkFunc) Arity() int {
	return len(f.Declaration.Params)
}

func (f *InkFunc) Call(ev *Evaluator, arguments []any) (any, error) {
	// 1. Create a brand new memory scope for the function
	environment := NewEnvironment(f.Closure)

	// 2. Bind all the arguments to the parameter names
	for i, param := range f.Declaration.Params {
		environment.Set(param.Lexeme, arguments[i])
	}

	// 3. Execute the body block in this new environment
	// (Assuming your evaluator has a method to execute a slice of statements in a specific env)
	err := ev.executeBlock(f.Declaration.Body, environment)
	// 4. Catch the return value if the user triggered a ReturnStmt
	if retErr, ok := err.(*ReturnError); ok {
		return retErr.Value, nil
	}

	// If it wasn't a return error, pass the real error up (or nil if successful)
	return nil, err
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		env: NewEnvironment(nil),
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
		err := ev.executeBlock(s.Statements, NewEnvironment(ev.env))
		return err
	case *stmt.IfStmt:
		condition, err := ev.Eval(s.Condition)
		if err != nil {
			return err
		}
		if isTruthy(condition) {
			return ev.execute(s.Then)
		}
		if s.Else != nil {
			return ev.execute(s.Else)
		}
		return nil
	case *stmt.WhileStmt:
		for {
			condition, err := ev.Eval(s.Condition)
			if err != nil {
				return err
			}
			if !isTruthy(condition) {
				break
			}
			err = ev.execute(s.Body)
			if err != nil {
				return err
			}
		}
		return nil
	case *stmt.FuncStmt:
		inkFunc := &InkFunc{
			Declaration: s,
			Closure:     ev.env,
		}
		ev.env.Set(s.Name.Lexeme, inkFunc)
		return nil
	case *stmt.ReturnStmt:
		var value any = nil
		if s.Value != nil {
			var err error
			value, err = ev.Eval(s.Value)
			if err != nil {
				return err
			}
		}
		return &ReturnError{
			Value: value,
		}
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
	case *expr.CallExpr:
		// 1. Evaluate the identifier (looks up the function in memory)
		callee, err := ev.Eval(e.Callee)
		if err != nil {
			return nil, err
		}

		// 2. Evaluate all the arguments
		var arguments []any
		for _, argExpr := range e.Args {
			argValue, err := ev.Eval(argExpr)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, argValue)
		}

		// 3. Ensure the thing we pulled from memory is actually a function
		function, ok := callee.(InkCallable)
		if !ok {
			return nil, fmt.Errorf("can only call functions and classes")
		}

		// 4. Ensure the Arity matches
		if len(arguments) != function.Arity() {
			return nil, fmt.Errorf("expected %d arguments but got %d", function.Arity(), len(arguments))
		}

		// 5. Fire the Call!
		return function.Call(ev, arguments)
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

func (ev *Evaluator) executeBlock(stmts []stmt.Stmt, env *Environment) error {
	oldEnv := ev.env
	ev.env = env
	for _, stmt := range stmts {
		err := ev.execute(stmt)
		if err != nil {
			ev.env = oldEnv
			return err
		}
	}
	ev.env = oldEnv
	return nil
}
