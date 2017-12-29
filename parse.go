package main

import (
	"fmt"
	"log"
)

func newparser(items chan item) *parser {
	p := &parser{
		itemc: items,
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
}

func (p *parser) next() item {
	p.tok = p.nexttok
	p.nexttok = <-p.itemc
	fmt.Printf("tok=%#v next=%#v\n", p.tok, p.nexttok)
	if p.tok.typ == itemEOF {
		panic(p.tok)
	}
	return p.tok
}

func (p *parser) peek() item {
	return p.nexttok
}

func (p *parser) parseInit() (c Cmd) {
	defer func() {
		e := recover()
		switch e := e.(type) {
		case item:
			if e.typ == itemEOF {
				return
			}
		case interface{}:
			panic(e)
		}
	}()
	// init here
	p.next()
	return p.parseCmd()
}
func (p *parser) parseCmd() (cmd Cmd) {
	p.next()
	switch p.tok.typ {
	case itemIf:
		println("got itemIf")
		cmd = p.parseIfStmt()
	case itemText:
		cmd = p.parseSimpleCmd()
		println("got itemText")
	default:
		println("got default")
		fmt.Printf("%#v\n", p.tok)
	}
	return
}
func (p *parser) error(fm string, i ...interface{}) {
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Printf(fm, i...)
}

func (p *parser) parseIfStmt() *IfStmt {
	println("parseIf")
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
	return BraceStmt{CmdList: CmdList{}}
}
func (p *parser) parseCmdList() (c *CmdList) {
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
func (p *parser) parseSimpleCmd() *SimpleCmd {
	// TODO(as): we're not filling out any type information here, maybe that's for
	// the best. But the tokens need to be refactored to show their irrelevance for
	// the parser's functionality.
	str := func() TextArg {
		return TextArg{Text: p.tok.val}
	}
	name := str()

	list := make([]Arg, 0)
	// TODO(as): handle semicolons, &&, ||, etc
	for !p.terminus(p.peek()) {
		p.next()
		list = append(list, str())
	}

	return &SimpleCmd{
		Name: name,
		Args: ArgList{Args: list},
	}
}
func (p *parser) parseCmds() (c Cmd) {
	log.Println(p.peek())
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
func (p *parser) parseSimple()   {}
func (p *parser) parseBrace()    {}
func (p *parser) parseName()     {}
func (p *parser) parseArg()      {}
func (p *parser) parseAssign()   {}
func (p *parser) parseVar()      {}
func (p *parser) parseRedirect() {}
func (p *parser) parseOp()       {}
func (p *parser) parseCat()      {}
func (p *parser) parseSub()      {}
func (p *parser) parseFor()      {}
func (p *parser) parseWhile()    {}
func (p *parser) parseSwitch()   {}
func (p *parser) parseIf()       {}
func (p *parser) parseFn()       {}
func (p *parser) parseAt()       {}
