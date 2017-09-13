package main

type Stmt interface {
	Body() Cmd
}

type IfStmt struct {
	Cond Cmd
	body Cmd
}
type SwitchStmt struct {
	Cond Cmd
	body BraceStmt
}
type WhileStmt struct {
	Cond Cmd
	body BraceStmt
}
type ForStmt struct {
	Range ArgList
	body  BraceStmt
}
type BraceStmt struct {
	CmdList
}

func (s IfStmt) Body() Cmd     { return s.body }
func (s SwitchStmt) Body() Cmd { return s.body.CmdList }
func (s WhileStmt) Body() Cmd  { return s.body.CmdList }
func (s ForStmt) Body() Cmd    { return s.body.CmdList }
func (s BraceStmt) Body() Cmd  { return s.CmdList }

func (s IfStmt) Exec() int { return s.Body().Exec() }
