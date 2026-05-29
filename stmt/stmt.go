package stmt

import (
	"github.com/hanzala211/goinkinterpreter/expr"
	"github.com/hanzala211/goinkinterpreter/token"
)

type Stmt interface {
	stmtNode()
}

type ExprStmt struct {
	Expr expr.Expr
}

type VarStmt struct {
	Initalizer expr.Expr
	Name       *token.Token
}

type PrintStmt struct {
	Expr expr.Expr
}

type BlockStmt struct {
	Statements []Stmt
}

type IfStmt struct {
	Condition expr.Expr
	Then      Stmt
	Else      Stmt
}

type WhileStmt struct {
	Condition expr.Expr
	Body      Stmt
}

type FuncStmt struct {
	Name   *token.Token
	Params []*token.Token
	Body   []Stmt
}

type ReturnStmt struct {
	Keyword *token.Token
	Value   expr.Expr
}

func (e *ExprStmt) stmtNode()   {}
func (w *WhileStmt) stmtNode()  {}
func (i *IfStmt) stmtNode()     {}
func (v *VarStmt) stmtNode()    {}
func (p *PrintStmt) stmtNode()  {}
func (b *BlockStmt) stmtNode()  {}
func (f *FuncStmt) stmtNode()   {}
func (r *ReturnStmt) stmtNode() {}
