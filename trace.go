package main

import (
	"fmt"
	"os"
)

func init() {
	if os.Getenv("rctrace") == "" {
		tracer = -1
	}
}

var tracer int

func tracef(fm string, i ...interface{}) string {
	if tracer < 0 {
		return ""
	}
	Printf(fm, i...)
	fmt.Println("(")
	tracer++
	return ""
}
func un(s string) {
	if tracer < 0 {
		return
	}
	tracer--
	bar()
	fmt.Println(")")
}

func Printf(fm string, i ...interface{}) {
	if tracer < 0 {
		return
	}
	bar()
	fmt.Printf(fm, i...)
}
func Println(i ...interface{}) {
	if tracer < 0 {
		return
	}
	bar()
	fmt.Println(i...)
}

func bar() {
	if tracer < 0 {
		return
	}
	for i := 0; i < tracer; i++ {
		fmt.Print(". . ")
	}
}
