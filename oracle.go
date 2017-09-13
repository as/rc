package main

func (p *parser) terminus(tok item) bool {
	if p.pop == tok.typ {
		p.pop = itemEOF
		return true
	}
	switch tok.typ {
	case itemSemi, itemNL, itemAmp, itemEOF:
		return true
	}
	return false
}

func (p *parser) chain(tok item) bool {
	switch tok.typ {
	case itemAnd, itemOr, itemPipe:
		return true
	}
	return false
}
