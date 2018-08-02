package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func (c CmdList) Exec(n *Ns) error {
	return c.Cmd.Exec(n)
}

func (c SimpleCmd) Exec(n *Ns) error {
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
				n.Fd[v.Dst.fd] = fd
				continue
			}

			fd, err := os.OpenFile(v.Dst.name, v.Dst.flags, 0666)
			if err != nil {
				return fmt.Errorf("open: %s\n", err)
			}

			n.Fd[v.Dst.fd] = fd

		}
	}

	name := c.Name.Resolve()
	cmd := exec.Command(name, c.Args.Resolve()...)
	cmd.Stdout = n.Fd[1].(io.WriteCloser)
	cmd.Stdin = n.Fd[0].(io.ReadCloser)
	cmd.Stderr = n.Fd[2].(io.WriteCloser)
	if c.Op.typ == itemPipe {
		pr, pw := io.Pipe()
		cmd.Stdout = pw
		ns := ns.Clone()
		ns.Fd[0] = pr
		c.Next.Exec(ns)
		go func() {
			cmd.Run()
			pw.Close()
		}()
		return nil
	}

	// CreationFlags: 0x00000008,
	if bt, ok := builtinTab[name]; ok {
		return bt(cmd)
	}
	return cmd.Start()
}

func (s IfStmt) Exec(n *Ns) error {
	err := s.Cond.Exec(n)
	if err != nil {
		return err
	}
	return s.Body().Exec(n)
}
