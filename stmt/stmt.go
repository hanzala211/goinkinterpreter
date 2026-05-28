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

func (e *ExprStmt) stmtNode()  {}
func (v *VarStmt) stmtNode()   {}
func (p *PrintStmt) stmtNode() {}
func (b *BlockStmt) stmtNode() {}
