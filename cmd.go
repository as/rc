package main

type Cmd interface {
	Exec() error
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

type Namespace struct {
	// The Windows PEB holds a handle to the current directory so that
	// it can't be removed. We don't need this type of hand holding here
	// so just store the name
	CurrentDir string

	// Should we inherit the environment from the operating system?
	Env map[string]string

	// Open file descriptors
	Fd map[int]interface{}
}
