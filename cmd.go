package main

type Cmd interface {
	Exec(n *Ns) error
}
type CmdList struct {
	Cmd Cmd
}
type SimpleCmd struct {
	Vars   []VarDecl
	Redirs []RedirDecl
	Name   Arg
	Args   ArgList
	Op     item
	Next   Cmd
}
