package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
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
	fdtab := make(map[int]interface{})
	fdtab[0] = os.Stdin
	fdtab[1] = os.Stdout
	fdtab[2] = os.Stderr
	for _, v := range c.Redirs {
		//TODO(as): handle all redirection cases
		// half duplex file descriptors are a special case
		if v.Dst.name != "" {
			//TODO(as): namespace the operation

			_, err := os.Stat(v.Dst.name)
			if err != nil {
				// TODO(as): this handles the case where the file doesn't exist but ignores other errors like perm issues. fix it.
				fd, err := os.Create(v.Dst.name)
				if err != nil {
					log.Printf("SimpleCmd: create: %s\n", err)
					return 1
				}
				fdtab[v.Dst.fd] = fd
				continue
			}

			fd, err := os.OpenFile(v.Dst.name, v.Dst.flags, 0666)
			if err != nil {
				log.Printf("SimpleCmd: %s\n", err)
				return 1
			}

			fdtab[v.Dst.fd] = fd

		}
	}

	cmd := exec.Cmd{
		SysProcAttr: &syscall.SysProcAttr{
		// By running the process as detached, we can avoid
		// the nasty conhost.
		// CreationFlags: 0x00000008,
		},
		Path:   c.Name.Resolve(),
		Args:   c.Args.Resolve(),
		Stdin:  fdtab[0].(io.ReadCloser),
		Stdout: fdtab[1].(io.WriteCloser),
		Stderr: fdtab[2].(io.WriteCloser),
	}
	cmd.Run()
	return 0
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
