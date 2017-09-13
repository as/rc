package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	leftMeta  = "{"
	rightMeta = "}"
	backTick  = "`"
	eof       = '‡'
)

const (
	itemError itemType = iota
	itemStart
	itemDot
	itemEOF
	itemLHS
	itemRHS
	itemEquals
	itemFor
	itemWhile
	itemSwtich
	itemFn
	itemAmp
	itemAnd
	itemOr
	itemPipe
	itemSemi
	itemLeftMeta
	itemRightMeta
	itemLeftParen
	itemRightParen
	itemNumber
	itemText
	itemBackTick
	itemHereString
	itemFnStart
	itemFnInside
	itemBreak
	itemContinue
	itemLeftBrace
	itemRightBrace
	itemNL
	itemEnv
	itemIf
)

type itemType int

type statefn func(*lexer) statefn

type item struct {
	typ itemType
	val string
}

type lexer struct {
	name  string
	input string
	start int
	pos   int
	width int
	items chan item
}

func (i item) String() string {
	return fmt.Sprintf("%d %s", i.typ, i.val)
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run() // run state machine
	return l, l.items
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) errorf(format string, args ...interface{}) statefn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func (l *lexer) current() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) emit(t itemType) {
	log.Printf("emit: %#v\n", item{t, l.current()})
	l.items <- item{t, l.current()}
	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func lexIdentifier(l *lexer) statefn {
	l.acceptRun("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad identifier syntax: %q",
			l.input[l.start:l.pos])
	}
	l.emit(itemText)
	return lexInsideAction
}

func (l *lexer) acceptWord() {
	l.acceptRun("/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`{}")
}

func (l *lexer) acceptBasicText() bool {
	if !l.accept("/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") {
		return false
	}
	l.acceptRun("/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	return true
}

func lexParen(l *lexer) statefn {
	ignoreSpaces(l)
	if !l.accept("(") {
		return l.errorf("lexParen: want '(' have %q", l.current())
	}
	l.emit(itemLeftParen)
	ignoreSpaces(l)
	for l.peek() != ')' {
		l.next()
		ignoreSpaces(l)
		ok := l.acceptBasicText()
		if !ok {
			if l.peek() == '$' {
				l.emit(itemEnv)
				continue
			}
			break
		}
		l.emit(itemText)
	}
	l.next()
	ignoreSpaces(l)
	if l.accept(")") {
		l.emit(itemError)
		return nil
	}
	l.emit(itemRightParen)
	l.next()
	return lexText
}

func space(r rune) bool {
	return unicode.IsSpace(r)
}

func ignoreSpaces(l *lexer) {
	if l.accept(" 	") {
		l.acceptRun(" 	")
		l.ignore()
	}
}

func lexText(l *lexer) statefn {
	ignoreSpaces(l)
	l.acceptRun("/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if l.pos == l.start && l.next() == eof {
		l.emit(itemEOF)
		return nil
	}
	switch l.current() {
	case "if":
		l.emit(itemIf)
		return lexParen
	case "switch":
		l.emit(itemIf)
		return lexParen
	case "for":
		l.emit(itemFor)
		return lexParen
	case "while":
		l.emit(itemWhile)
		return lexParen
	case "fn":
		l.emit(itemFn)
		return lexParen
	case "break":
		l.emit(itemBreak)
	case "continue":
		l.emit(itemContinue)
	case ";":
		l.emit(itemSemi)
	case "\n":
		l.emit(itemNL)
	case "&":
		l.emit(itemAmp)
	case "$":
		return lexEnv
	case "{":
		l.emit(itemLeftBrace)
	case "}":
		l.emit(itemRightBrace)
	default:
		l.emit(itemText)
	}
	switch r := l.peek(); {
	case r == eof:
		if l.pos == l.start {
			println("itemEOF")
			l.emit(itemEOF)
			return nil
		}
	}
	l.next()
	return lexText
}

func lexEquals(l *lexer) statefn {
	l.acceptRun("/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	l.emit(itemText)
	return lexText
}

func lexEnv(l *lexer) statefn {
	if !l.accept("$") {
		return l.errorf("Invalid variable", l.input[:])
	}
	l.ignore()
	l.acceptRun("/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	l.emit(itemEnv)
	return lexText
}

func lexLeftMeta(l *lexer) statefn {
	l.pos += len(leftMeta)
	l.emit(itemLeftMeta)
	return lexInsideAction
}

func lexRightMeta(l *lexer) statefn {
	l.pos += len(rightMeta)
	l.emit(itemRightMeta)
	return lexText
}

func lexInsideAction(l *lexer) statefn {
	// Either num, string, or id
	for {
		if strings.HasPrefix(l.input[l.pos:], rightMeta) {
			return lexRightMeta
		}
		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed action")
		case unicode.IsSpace(r):
			l.ignore()
		case r == '|':
			l.emit(itemPipe)
		case isAlphaNumeric(r):
			l.backup()
			return lexIdentifier
		}
	}
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func prompt() {
	fmt.Print("; ")
}

func (cl cmdline) push(i item, fn func()) {
	cl = append(cl, fn)
}

type cmdline []cmd

type cmd func()

func assign(r, l *item) {
	fmt.Printf("assign %v to %v\n", r, l)
	err := os.Setenv(r.val, l.val)
	if err != nil {
		fmt.Println(err)
	}
}

func extract(i item) item {
	return item{itemText, os.Getenv(i.val)}
}

//
// Execution

func bcd(i item) {
	//	pl("change dir: %q\n", i.val)
	err := os.Chdir(i.val)
	if err != nil {
		fmt.Println(err)
	}
	panic("err")
}

func becho(i string) {
	fmt.Println(i)
}

func bexit() {
	//os.Exit(0)
}

func main() {
	l, c := lex("main", "if ( ls ) { echo listed }")
	p := newparser(c)
	log.Printf("final value: %#v\n", p.parseInit())
	l = l
}
