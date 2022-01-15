package parser

import (
	"strings"
	"unicode/utf8"
)

const (
	eof              rune = -1
	tokenChBufSize        = 10
	runeStackBufSize      = 64
)

type stateFunc func(*lexer) stateFunc

type tokenType int

type token struct {
	Type       tokenType
	Value      string
	Start, End int
}

type lexer struct {
	source               string
	startPos, currentPos int
	runeStack            []rune

	startState stateFunc
	tokenCh    chan token

	err          error
	errorHandler func(err error)
}

// newLexer creates a returns a lexer ready to parse the given source code.
func newLexer(src string, start stateFunc, errHand func(err error)) *lexer {
	return &lexer{
		source:       src,
		startState:   start,
		startPos:     0,
		currentPos:   0,
		errorHandler: errHand,
		runeStack:    make([]rune, runeStackBufSize),
	}
}

func (l *lexer) Start() {
	l.tokenCh = make(chan token, tokenChBufSize)

	go l.start()
}

func (l *lexer) start() {
	state := l.startState
	for state != nil {
		state = state(l)
	}
	close(l.tokenCh)
}

// NextToken returns the next token from the lexer and a value to denote whether
// or not the token is finished.
func (l *lexer) NextToken() (*token, bool) {
	tok, ok := <-l.tokenCh
	if ok {
		return &tok, false
	}
	return nil, true
}

// Error if an errorHandler is given, sets lex.Err with given error and calls errorHandler
// if no errorHandler is given, then it panics with given error.
func (l *lexer) Error(e error) {
	if l.errorHandler == nil {
		panic(e)
	}

	l.err = e
	l.errorHandler(e)
}

// Current returns the value being being analyzed at this moment.
func (l *lexer) Current() string {
	return l.source[l.startPos:l.currentPos]
}

// Current returns the value being being analyzed at this moment.
func (l *lexer) Get(startPos, endPos int) string {
	return l.source[startPos:endPos]
}

func (l *lexer) Err() error {
	return l.err
}

// Emit will receive a token type and push a new token with the current analyzed
// value into the tokens channel.
func (l *lexer) Emit(t tokenType) {
	tok := token{
		Type:  t,
		Value: l.Current(),
		Start: l.startPos,
		End:   l.currentPos,
	}
	l.tokenCh <- tok
	l.startPos = l.currentPos
	l.clearRune()
}

// Ignore clears the rewind stack and then sets the current beginning position
// to the current position in the source which effectively ignores the section
// of the source being analyzed.
func (l *lexer) Ignore() {
	l.clearRune()
	l.startPos = l.currentPos
}

// Peek performs a Next operation immediately followed by a Rewind returning the
// peeked rune.
func (l *lexer) Peek() rune {
	r := l.Next()
	l.Rewind()

	return r
}

// Rewind will take the last rune read (if any) and rewind back. Rewinds can
// occur more than once per call to Next but you can never rewind past the
// last point a token was emitted.
func (l *lexer) Rewind() {
	r := l.popRune()
	if r > eof {
		size := utf8.RuneLen(r)
		l.currentPos -= size
		if l.currentPos < l.startPos {
			l.currentPos = l.startPos
		}
	}
}

// Next pulls the next rune from the Lexer and returns it, moving the position
// forward in the source.
func (l *lexer) Next() rune {
	str := l.source[l.currentPos:]
	if str == "" {
		l.pushRune(eof)
		return eof
	}

	r, size := utf8.DecodeRuneInString(str)
	l.currentPos += size
	l.pushRune(r)

	return r
}

// Take receives a string containing all acceptable strings and will contine
// over each consecutive character in the source until a token not in the given
// string is encountered. This should be used to quickly pull token parts.
func (l *lexer) Take(chars string) {
	r := l.Next()
	for strings.ContainsRune(chars, r) {
		r = l.Next()
	}
	l.Rewind() // last next wasn't a match
}

// TakeNext is similar to Take but takes if next rune matches
func (l *lexer) TakeNext(ch rune) {
	r := l.Next()

	if ch != r {
		l.Rewind() // last next wasn't a match
	}
}

func (l *lexer) pushRune(r rune) {
	l.runeStack = append(l.runeStack, r)
}

func (l *lexer) popRune() rune {
	r := l.runeStack[len(l.runeStack)-1]
	l.runeStack = l.runeStack[:len(l.runeStack)-1]
	return r
}

func (l *lexer) clearRune() {
	l.runeStack = l.runeStack[:0]
}
