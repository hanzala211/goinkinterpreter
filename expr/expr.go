package expr

import (
	"fmt"
	"strings"

	"github.com/hanzala211/goinkinterpreter/token"
)

type Expr interface {
	String() string
	exprNode()
}

type BinaryExpr struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

type UnaryExpr struct {
	Operator *token.Token
	Right    Expr
}

type LiteralExpr struct {
	Value any
}

type Grouping struct {
	Expression Expr
}

type VarExpr struct {
	Name *token.Token
}

type AssignExpr struct {
	Name  *token.Token
	Value Expr
}

type LogicalExpr struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

type CallExpr struct {
	Callee Expr
	Paren  *token.Token
	Args   []Expr
}

func (c CallExpr) String() string {
	return paranthesized("call", c.Callee)
}

func (l LogicalExpr) String() string {
	return paranthesized(l.Operator.Lexeme, l.Left, l.Right)
}

func (a AssignExpr) String() string {
	return fmt.Sprintf("%s = %v", a.Name.Lexeme, a.Value)
}

func (v VarExpr) String() string {
	return v.Name.Lexeme
}

func paranthesized(name string, expr ...Expr) string {
	var builder strings.Builder
	builder.WriteRune('(')
	builder.WriteString(name)

	for _, e := range expr {
		builder.WriteRune(' ')
		builder.WriteString(e.String())
	}
	builder.WriteRune(')')
	return builder.String()
}

func (g Grouping) String() string {
	return paranthesized("group", g.Expression)
}

func (b BinaryExpr) String() string {
	return paranthesized(b.Operator.Lexeme, b.Left, b.Right)
}

func (u UnaryExpr) String() string {
	return paranthesized(u.Operator.Lexeme, u.Right)
}

func (l LiteralExpr) String() string {
	if l.Value == nil {
		return "nil"
	}
	switch l.Value.(type) {
	case string:
		return fmt.Sprintf("%q", l.Value)
	default:
		return fmt.Sprintf("%v", l.Value)

	}
}

func (e BinaryExpr) exprNode()  {}
func (u UnaryExpr) exprNode()   {}
func (v VarExpr) exprNode()     {}
func (a AssignExpr) exprNode()  {}
func (l LiteralExpr) exprNode() {}
func (l LogicalExpr) exprNode() {}
func (g Grouping) exprNode()    {}
func (c CallExpr) exprNode()    {}
