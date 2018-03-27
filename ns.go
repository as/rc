package main

import "os"

var ns = Ns{
	CurrentDir: ".",
	Env:        map[string]string{},
	Fd: map[int]interface{}{
		0: os.Stdin,
		1: os.Stdout,
		2: os.Stderr,
	},
}

type Ns struct {
	// The Windows PEB holds a handle to the current directory so that
	// it can't be removed. We don't need this type of hand holding here
	// so just store the name
	CurrentDir string

	// Should we inherit the environment from the operating system?
	Env map[string]string

	// Open file descriptors
	Fd map[int]interface{}
}

func (n *Ns) Clone() *Ns {
	m := *n
	for k, v := range n.Fd {
		m.Fd[k] = v
	}
	return &m
}
