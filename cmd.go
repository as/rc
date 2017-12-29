package main

import (
	"os"
	"os/exec"
)

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

func (c CmdList) Exec() int { return 0 }
func (c *SimpleCmd) Exec() int {
	cmd := exec.Cmd{
		Path:   c.Name.Resolve(),
		Args:   c.Args.Resolve(),
		Stderr: os.Stderr,
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
	}
	cmd.Run()
	return 0
}
