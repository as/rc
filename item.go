package main

type itemType int

const (
	leftMeta  = "{"
	rightMeta = "}"
	backTick  = "`"
	eof       = '‡'
	runText = "/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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