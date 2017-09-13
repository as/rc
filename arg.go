package main

type Arg interface {
	Resolve() string
}

type ArgList struct {
	Args []Arg
}
type TextArg struct {
	Text string
}
type VarArg struct {
	Arg Arg
}
type SubArg struct {
	Cmd BraceStmt
}
