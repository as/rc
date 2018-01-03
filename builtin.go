package main

import (
	"fmt"
	"os"
	"os/exec"
)

type builtin func(cmd *exec.Cmd) error

var builtinTab = map[string]builtin{
	"touch": builtinTouch,
	"echo":  builtinEcho,
}

func builtinShift(cmd *exec.Cmd) (err error) {
	if len(cmd.Args) == 0 {
		return nil
	}
	copy(cmd.Args, cmd.Args[1:])
	cmd.Args = cmd.Args[:len(cmd.Args)-1]
	return nil
}

func builtinEcho(cmd *exec.Cmd) (err error) {
	builtinShift(cmd)
	sep := ""
	for _, v := range cmd.Args {
		fmt.Fprintf(cmd.Stdout, "%s%s", sep, v)
		sep = " "
	}
	fmt.Fprintln(cmd.Stdout)
	return nil
}

func builtinTouch(cmd *exec.Cmd) (err error) {
	builtinShift(cmd)
	var fd *os.File
	for _, v := range cmd.Args {
		fd, err = os.Create(v)
		if err != nil {
			break
		}
		fd.Close()
	}
	return err
}
