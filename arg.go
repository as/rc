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

func (a TextArg) Resolve() string { return a.Text }
func (a ArgList) Resolve() (args []string) {
	for _, v := range a.Args {
		args = append(args, v.Resolve())
	}
	return args
}
