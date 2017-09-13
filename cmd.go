package main

type Cmd interface {
	Exec() int
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

func (c CmdList) Exec() int    { return 0 }
func (c *SimpleCmd) Exec() int { return 0 }
