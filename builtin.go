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
	"cd":    builtinCd,
	"pwd":   builtinPwd,
}

func builtinShift(cmd *exec.Cmd) (err error) {
	if len(cmd.Args) == 0 {
		return nil
	}
	copy(cmd.Args, cmd.Args[1:])
	cmd.Args = cmd.Args[:len(cmd.Args)-1]
	return nil
}

func builtinCd(cmd *exec.Cmd) (err error) {
	builtinShift(cmd)
	return os.Chdir(cmd.Args[0])
}
func builtinPwd(cmd *exec.Cmd) (err error) {
	builtinShift(cmd)
	wd, err := os.Getwd()
	if err == nil {
		fmt.Fprintf(os.Stdout, "%s\n", wd)
	}
	return err
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
