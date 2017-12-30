package main

import (
	"fmt"
	"os"
	"strings"
)

func newparser(items chan item) *parser {
	p := &parser{
		nexttok: item{typ: itemStart},
		tok:     item{typ: itemStart},
		itemc:   items,
	}
	return p
}

type parser struct {
	itemc   chan item
	nexttok item
	tok     item
	pop     itemType
	nests   int
	err     error
	level   int
}

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

func (p *parser) next() item {
	if p.tok.typ == itemError {
		Printf("itemError tok=%#v next=%#v\n", p.tok, p.nexttok)
		panic(p.tok)
	}
	p.tok = p.nexttok
	p.nexttok = <-p.itemc
	Printf("tok=%#v next=%#v\n", p.tok, p.nexttok)
	return p.tok
}

func (p *parser) peek() item {
	return p.nexttok
}

func (p *parser) error(fm string, i ...interface{}) {
	Println("!!!!!!!!!!!!!!!!!!!!!!!")
	Printf(fm, i...)
	Println("!!!!!!!!!!!!!!!!!!!!!!!")
}

func (p *parser) parseInit() (c Cmd) {
	defer un(tracef("parseInit"))
	defer func() {
		e := recover()
		switch e := e.(type) {
		case item:
			if e.typ == itemError {
				return
			}
		case interface{}:
			panic(e)
		}
	}()
	// init here
	p.next()
	p.next()
	return p.parseCmd()
}
func (p *parser) parseCmd() (cmd Cmd) {
	defer un(tracef("parseCmd"))
	Printf(" %#v\n", p.tok)
	switch p.tok.typ {
	case itemIf:
		cmd = p.parseIfStmt()
	case itemText:
		Printf("itemText %#v\n", p.tok)
		cmd = p.parseSimpleCmd()
	default:
		fmt.Printf("default %#v\n", p.tok)
	}
	return
}
func (p *parser) parseSimpleCmd() (sc *SimpleCmd) {
	tracef("parseSimpleCmd")
	Printf("parseSimpleCmd %#v\n", p.tok)
	defer func() {
		Printf("SimpleCmd: %+#v\n", sc)
		un("")
	}()
	if p.tok.typ != itemText {
		p.error("parseSimpleCmd: expected text, got %#v\n", p.tok)
		panic("!")
		return nil
	}
	name := TextArg{Text: p.tok.val}

	Printf("1 parseSimpleCmd %#v\n", p.tok)
	p.next()
	Printf("2 parseSimpleCmd %#v\n", p.tok)

	if p.terminus(p.tok) {
		Printf("3 parseSimpleCmd %#v\n", p.tok)
		p.next()
		Printf("4 parseSimpleCmd %#v\n", p.tok)
		return &SimpleCmd{
			Name: name,
		}
	}

	var redir []RedirDecl
	if redirector(p.tok) {
		r := p.parseRedirect()
		if r != nil {
			redir = append(redir, *r)
		}
		return &SimpleCmd{
			Name:   name,
			Redirs: redir,
		}
	}

	arglist := p.parseArgs()

	// check for trailing redirects
	if redirector(p.tok) {
		r := p.parseRedirect()
		if r != nil {
			redir = append(redir, *r)
		}
	}

	return &SimpleCmd{
		Name:   name,
		Args:   arglist,
		Redirs: redir,
	}
}

func (p *parser) parseRedirect() *RedirDecl {
	defer un(tracef("parseRedirect"))
	//TODO(as): <>[1=0]
	return p.parseSimpleRedirect()
}
func (p *parser) parseSimpleRedirect() *RedirDecl {
	defer un(tracef("parseSimpleRedirect"))
	if p.terminus(p.peek()) {
		// EOF is invalid in any of these
		p.error("parseSimpleRedirect: EOF reading redirect rhs\n")
		panic("!")
		return nil
	}
	flags := 0

	switch p.tok.typ {
	case itemGreatGreat:
		flags |= os.O_APPEND
	case itemGreat:
		flags |= os.O_WRONLY
		p.next()
		rhs := p.parsePath()
		//TODO(as): all strings can be expanded
		return &RedirDecl{Dst: FD{
			fd:    1,
			name:  rhs,
			flags: flags,
		}}
	case itemLess:
	case itemLessLess:
		// This is an invalid operator. Only defined for
		// clear error messages
	}
	p.error("parseSimpleRedirect: unknown/unsupported redirect rhs: %q", p.tok.val)
	panic("!")
	return nil

}

func (p *parser) parseArgs() ArgList {
	defer un(tracef("parseArgs"))
	list := make([]Arg, 0)
	// TODO(as): handle semicolons, &&, ||, etc
	if p.tok.typ != itemText {
		p.error("expected itemText, got %#v\n", p.tok)
		panic("!")
	}
	for p.tok.typ == itemText {
		// TODO(as): handle expansions and nested commands
		list = append(list, TextArg{Text: p.tok.val})
		p.next()
	}
	// An empty argument list is a valid argument list
	Printf("ret list %#v\n", list)
	return ArgList{Args: list}
}

func (p *parser) parseIfStmt() *IfStmt {
	defer un(tracef("parseIfStmt"))
	if p.tok.typ != itemIf {
		p.error("bad token: %#v\n", p.tok)
		return nil
	}
	p.next()
	if p.tok.typ != itemLeftParen {
		p.error("if: expected '(' got %q", p.tok)
		return nil
	}
	cond := p.parseCmdList()

	tok := p.peek()
	if tok.typ == itemLeftBrace {
		body := p.parseBraceStmt()
		return &IfStmt{
			Cond: cond,
			body: body,
		}
	}
	body := p.parseSimpleCmd()
	return &IfStmt{
		Cond: cond,
		body: body,
	}
}
func (p *parser) parseBraceStmt() BraceStmt {
	defer un(tracef("parseBraceStmt"))
	return BraceStmt{CmdList: CmdList{}}
}
func (p *parser) parseCmdList() (c *CmdList) {
	defer un(tracef("parseCmdList"))
	if p.tok.typ != itemLeftParen {
		p.error("parseCmdList: expected '(' got %#v", p.tok)
		return nil
	}
	p.pop = itemRightParen
	cmd := p.parseCmds()
	if p.peek().typ != itemRightParen {
		p.error("parseCmdList: expected ')' got %#v", p.tok)
		return nil
	}
	return &CmdList{cmd}
}

func redirector(tok item) bool {
	switch tok.typ {
	case itemLeftBrace, itemLess, itemLessLess, itemGreat, itemGreatGreat, itemDiamond:
		return true
	}
	return false
}

func (p *parser) parsePath() (path string) {
	defer un(tracef("parsePath"))
	if redirector(p.tok) {
		p.error("redirect in path: %#v", p.tok)
		panic("!")
		return ""
	}
	//TODO(as): stop piggybacking on parseArgs
	arglist := p.parseArgs()

	//TODO(as): pointer the arglist
	//	if arglist == nil{
	//		p.error("parsePath: bad arglist from parseArgs")
	//		return
	//	}
	//TODO(as): don't call this here, we can't map identifiers
	// to their values at runtime since we aren't executing as
	// we go
	return strings.Join(arglist.Resolve(), "")
}

func (p *parser) parseCmds() (c Cmd) {
	defer un(tracef("parseCmds"))
	Println(p.peek())
	sc := p.parseSimpleCmd()
	if sc == nil {
		return
	}
	ptr := sc
	p.nests = 1
	for !p.terminus(p.peek()) {
		if p.nests > 25 {
			panic("too many commands")
		}
		p.next()
		sc0 := p.parseSimpleCmd()
		if sc0 == nil {
			break
		}
		if p.chain(p.peek()) {
			p.next()
			sc0.Op = p.tok
		}
		ptr.Next = sc0
		ptr = sc0
	}
	return c
}
func (p *parser) parseSimple() {}
func (p *parser) parseBrace()  {}
func (p *parser) parseName()   {}
func (p *parser) parseArg()    {}
func (p *parser) parseAssign() {}
func (p *parser) parseVar()    {}
func (p *parser) parseOp()     {}
func (p *parser) parseCat()    {}
func (p *parser) parseSub()    {}
func (p *parser) parseFor()    {}
func (p *parser) parseWhile()  {}
func (p *parser) parseSwitch() {}
func (p *parser) parseIf()     {}
func (p *parser) parseFn()     {}
func (p *parser) parseAt()     {}
