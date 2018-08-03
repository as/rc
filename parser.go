package main

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

func (p *parser) next() item {
	if p.current(itemError) {
		Printf("itemError tok=%v next=%v\n", p.tok, p.nexttok)
		panic(p.tok)
	}
	p.tok = p.nexttok
	p.nexttok = <-p.itemc
	Printf("next: tok=%v next=%v\n", p.tok, p.nexttok)
	return p.tok
}

func (p *parser) verify(current ...itemType) bool {
	for i, current := range current {
		if i != 0 {
			p.next()
		}
		if p.tok.typ != current {
			p.error("bad token: have %v want %v\n", p.tok, current)
			return false
		}
	}
	return true
}

func (p *parser) current(it itemType) bool {
	return p.tok.typ == it
}

func (p *parser) peek() item {
	return p.nexttok
}

func (p *parser) error(fm string, i ...interface{}) {
	Println("!!!!!!!!!!!!!!!!!!!!!!!")
	Printf(fm, i...)
	Println("!!!!!!!!!!!!!!!!!!!!!!!")
}

func redirector(tok item) bool {
	switch tok.typ {
	case itemLeftBrace, itemLess, itemLessLess, itemGreat, itemGreatGreat, itemDiamond:
		return true
	}
	return false
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
