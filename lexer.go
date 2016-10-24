package lexer

import "strings"

const EOF = -1

// Lexer implements a shellray lexer for the PHP language.
type Lexer struct {
	Tokens []int
	input  string
	length int
	pos    int // current position in the input
	start  int // start position of this item

	state StateFn // the next lexing function to enter
	//width  int     // width of last rune read from input
}

type StateFn func(*Lexer) StateFn

// NewLexer creates a new lexer instance for the given data.
func NewLexer(data string) *Lexer {
	return &Lexer{
		Tokens: make([]int, 0),
		input:  data,
		length: len(data),
		pos:    0,
		start:  0,
	}
}

// Reset resets the lexer so it can be run again.
func (l *Lexer) Reset() {
	l.pos = 0
	l.start = 0
	l.Tokens = make([]int, 0)
}

// Run runs the state machine for the lexer.
func (l *Lexer) Run(fn StateFn) {
	for l.state = fn; l.state != nil; {
		l.state = l.state(l)
	}
}

// Next returns the next rune in the input.
func (l *Lexer) Next() rune {
	r := l.Peek()
	l.pos++
	return r
}

// Peek returns but does not consume the next rune in the input.
func (l *Lexer) Peek() rune {
	if l.pos >= l.length {
		return EOF
	}
	r := rune(l.input[l.pos])
	return r
}

// Advance advances the current pointer by n.
func (l *Lexer) Advance(n int) {
	l.pos = l.pos + n
}

// Backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) Backup() {
	l.pos--
}

// Reverts resets the current pointer to the start of the current token.
func (l *Lexer) Revert() {
	l.pos = l.start
}

// Emit passes an item back to the client.
func (l *Lexer) Emit(t int) {
	l.Tokens = append(l.Tokens, t)
	l.start = l.pos
}

// Ignore sets the token start point to the current position.
func (l *Lexer) Ignore() {
	l.start = l.pos
}

// Accept consumes the next rune if it's from the valid set.
func (l *Lexer) Accept(valid string) bool {
	if strings.IndexRune(valid, l.Next()) >= 0 {
		return true
	}
	l.Backup()
	return false
}

// AcceptRun consumes a run of runes from the valid set.
func (l *Lexer) AcceptRun(valid string) {
	for strings.IndexRune(valid, l.Next()) >= 0 {
	}
	l.Backup()
}

// LexemLength returns the length of the current token.
func (l *Lexer) TokenLength() int {
	return l.pos - l.start
}

// TokenValue returns the value of the current token.
func (l *Lexer) TokenValue() string {
	// it is possible to advance beyound the input, so be careful here.
	end := l.pos
	if end > len(l.input) {
		end := len(l.input)
	}

	return l.input[l.start:end]
}
