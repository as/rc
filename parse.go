package main

import (
	"fmt"
	"os"
	"strings"
)

func (p *parser) parseCmd() (cmd Cmd) {
	defer un(tracef("parseCmd"))
	Printf("%v\n", p.tok)
	switch p.tok.typ {
	case itemLeftParen:
		cmd = p.parseCmdList()
	case itemLeftBrace:
		cmd = p.parseBraceStmt()
	case itemWhile:
		cmd = p.parseWhileStmt()
	case itemIf:
		cmd = p.parseIfStmt()
	case itemText:
		cmd = p.parseSimpleCmd()
	default:
		fmt.Printf("default %#v\n", p.tok)
	}
	return
}

func (p *parser) parseSimpleCmd() (sc *SimpleCmd) {
	tracef("parseSimpleCmd")
	Printf("parseSimpleCmd %v\n", p.tok)
	defer func() {
		Printf("SimpleCmd: %+#v\n", sc)
		un("")
	}()

	if !p.current(itemText) {
		return nil
	}
	name := TextArg{Text: p.tok.val}
	p.next()

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

	var arglist ArgList
	if p.current(itemText) {
		arglist = p.parseArgs()
	}

	// check for trailing redirects
	if redirector(p.tok) {
		r := p.parseRedirect()
		if r != nil {
			redir = append(redir, *r)
		}
	}

	var next Cmd
	var op item
	if p.current(itemPipe) {
		op = p.tok
		p.next()
		next = p.parseSimpleCmd()
	}

	return &SimpleCmd{
		Op:     op,
		Next:   next,
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
		fallthrough
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

func (p *parser) parseParity(push, pop itemType) (c *CmdList) {
	defer un(tracef("parseParity"))
	if !p.verify(push) {
		return nil
	}
	p.next()

	p.pop = pop
	cmd := p.parseSimpleCmd()
	if !p.verify(pop) {
		return nil
	}
	p.next()
	return &CmdList{cmd}
}

func (p *parser) parseCmdList() (c *CmdList) {
	defer un(tracef("parseCmdList"))
	return p.parseParity(itemLeftParen, itemRightParen)
}

func (p *parser) parseBraceStmt() *BraceStmt {
	defer un(tracef("parseBraceStmt"))
	list := p.parseParity(itemLeftBrace, itemRightBrace)
	if list == nil {
		return nil
	}
	return &BraceStmt{CmdList: *list}
}

func (p *parser) parseArgs() ArgList {
	defer un(tracef("parseArgs"))
	list := []Arg{}
	for p.current(itemText) {
		list = append(list, TextArg{Text: p.tok.val})
		p.next()
	}
	// an empty argument list is valid
	Printf("ret list %#v\n", list)
	return ArgList{Args: list}
}

func (p *parser) parseIfStmt() *IfStmt {
	defer un(tracef("parseIfStmt"))
	if !p.verify(itemIf, itemLeftParen) {
		return nil
	}

	cond := p.parseCmdList()

	if p.current(itemLeftBrace) {
		br := p.parseBraceStmt()
		return &IfStmt{cond, br}
	}
	return &IfStmt{cond, p.parseSimpleCmd()}
}

func (p *parser) parseWhileStmt() *WhileStmt {
	defer un(tracef("parseWhileStmt"))
	if !p.verify(itemWhile, itemLeftParen) {
		return nil
	}

	cond := p.parseCmdList()
	var body Cmd
	if p.current(itemLeftBrace) {
		body := p.parseBraceStmt()
		if body == nil {
			return nil
		}
		return &WhileStmt{cond, *body}
	}
	body = p.parseSimpleCmd()
	if body == nil {
		return nil
	}
	return &WhileStmt{cond, body}
}

func (p *parser) parsePath() (path string) {
	defer un(tracef("parsePath"))
	if redirector(p.tok) {
		p.error("redirect in path: %v", p.tok)
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
	p.next()
	ptr := sc
	p.nests = 1
	for !p.terminus(p.tok) {
		if p.nests > 25 {
			panic("too many commands")
		}
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
