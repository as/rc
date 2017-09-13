package main

type Decl interface {
	Scope() *Scope
}
type FnDecl struct {
	Name  string
	Body  BraceStmt
	Scope *Scope
}
type VarDecl struct {
	Name  string
	Value Arg
	Scope *Scope
}
type RedirDecl struct {
	Src   FD
	Dst   FD
	Scope *Scope
}
