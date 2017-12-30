package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	runcmd = flag.String("c", "", "run command and exit")
)

func init() {
	flag.Parse()
}

func repl(line string) {
	_, c := lex("main", line)
	p := newparser(c)
	cmd := p.parseInit()
	if cmd != nil {
		//fmt.Printf("@@@%+v\n", cmd)
		cmd.Exec()
	} else {
		log.Println("no")
	}
}

// read from stdin
func main() {
	if *runcmd != "" {
		repl(*runcmd)
	} else {
		sc := bufio.NewScanner(os.Stdin)
		fmt.Print("; ")
		for sc.Scan() {
			repl(sc.Text())
			fmt.Print("; ")
		}
	}
}

func mainz() {
	l, c := lex("main", "if ( ls ) { echo listed }")
	p := newparser(c)
	log.Printf("final value: %#v\n", p.parseInit())
	l = l
}
