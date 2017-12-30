package main

import "fmt"

var tracer int

func tracef(fm string, i ...interface{}) string {
	Printf(fm, i...)
	fmt.Println("(")
	tracer++
	return ""
}
func un(s string) {
	tracer--
	bar()
	fmt.Println(")")
}

func Printf(fm string, i ...interface{}) {
	bar()
	fmt.Printf(fm, i...)
}
func Println(i ...interface{}) {
	bar()
	fmt.Println(i...)
}

func bar() {
	for i := 0; i < tracer; i++ {
		fmt.Print(". . . ")
	}
}
