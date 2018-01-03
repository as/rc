package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func (c CmdList) Exec() error {
	return c.Cmd.Exec()
}

func (c SimpleCmd) Exec() error {
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
					return fmt.Errorf("create: %s\n", err)
				}
				fdtab[v.Dst.fd] = fd
				continue
			}

			fd, err := os.OpenFile(v.Dst.name, v.Dst.flags, 0666)
			if err != nil {
				return fmt.Errorf("open: %s\n", err)
			}

			fdtab[v.Dst.fd] = fd

		}
	}

	name := c.Name.Resolve()
	cmd := exec.Command(name, c.Args.Resolve()...)
	cmd.Stdin = fdtab[0].(io.ReadCloser)
	cmd.Stdout = fdtab[1].(io.WriteCloser)
	cmd.Stderr = fdtab[2].(io.WriteCloser)
	//SysProcAttr: &syscall.SysProcAttr{
	// By running the process as detached, we can avoid
	// the nasty conhost.
	// CreationFlags: 0x00000008,
	if bt, ok := builtinTab[name]; ok {
		return bt(cmd)
	}
	return cmd.Run()
}

func (s IfStmt) Exec() error {
	err := s.Cond.Exec()
	if err != nil {
		return err
	}
	return s.Body().Exec()
}
