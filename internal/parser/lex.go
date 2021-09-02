// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// tokenItem represents a token or text string returned from the scanner.
type tokenItem struct {
	tok  token  // The type of this tokenItem.
	pos  Pos    // The starting position, in bytes, of this tokenItem in the input string.
	val  string // The value of this tokenItem.
	line int    // The line number at the start of this tokenItem.
}

func (i tokenItem) String() string {
	switch {
	case i.tok == EOF:
		return "EOF"
	case i.tok == ERROR:
		return i.val
	case i.tok > keyword_beg && i.tok < keyword_end:
		return fmt.Sprintf("%d: <%s>", i.tok, i.val)
	}

	return fmt.Sprintf("%d: %q", i.tok, i.val)
}

var key = map[string]token{
	"model": MODEL,
	"tab":   TABLE,
	"view":  VIEW,
	"cnt":   COUNTER,
}

// state functions
const (
	leftDelim  = "{{"
	rightDelim = "}}"
	eof        = -1
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name       string         // the name of the input; used only for error reports
	input      string         // the string being scanned
	pos        Pos            // current position in the input
	start      Pos            // start position of this tokenItem
	width      Pos            // width of last rune read from input
	items      chan tokenItem // channel of scanned items
	parenDepth int            // nesting depth of ( ) exprs
	line       int            // 1+number of newlines seen
	startLine  int            // start line of this tokenItem
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0

		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	if r == '\n' {
		l.line++
	}

	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()

	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// emit passes an tokenItem back to the client.
func (l *lexer) emit(t token) {
	l.items <- tokenItem{t, l.start, l.input[l.start:l.pos], l.startLine}
	l.start = l.pos
	l.startLine = l.line
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.line += strings.Count(l.input[l.start:l.pos], "\n")
	l.start = l.pos
	l.startLine = l.line
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- tokenItem{ERROR, l.start, fmt.Sprintf(format, args...), l.startLine}

	return nil
}

// nextItem returns the next tokenItem from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() tokenItem {
	return <-l.items
}

// drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) drain() {
	for range l.items {
	}
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// atRightDelim reports whether the lexer is at a right delimiter, possibly preceded by a trim marker.
func (l *lexer) atRightDelim() (delim bool) {
	return strings.HasPrefix(l.input[l.pos:], rightDelim)
}

// atTerminator reports whether the input is at valid termination character to
// appear after an identifier. Breaks .X.Y into two pieces. Also catches cases
// like "$x+2" not being acceptable without a space, in case we decide one
// day to implement arithmetic.
func (l *lexer) atTerminator() bool {
	r := l.peek()
	if isSpace(r) || isEndOfLine(r) {
		return true
	}
	switch r {
	case eof, '.', ',', '|', ':', ')', '(':
		return true
	}
	// Does r start the delimiter? This can be ambiguous (with delim=="//", $x/2 will
	// succeed but should fail) but only in extremely rare cases caused by willfully
	// bad choice of delimiter.
	if rd, _ := utf8.DecodeRuneInString(rightDelim); rd == r {
		return true
	}

	return false
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:      name,
		input:     input,
		items:     make(chan tokenItem),
		line:      1,
		startLine: 1,
	}
	go l.run()

	return l
}

// lexText scans until an opening action delimiter, "{{".
func lexText(l *lexer) stateFn {
	l.width = 0
	if x := strings.Index(l.input[l.pos:], leftDelim); x >= 0 {
		l.pos += Pos(x)
		if l.pos > l.start {
			l.line += strings.Count(l.input[l.start:l.pos], "\n")
			l.emit(TEXT)
		}
		l.ignore()

		return lexLeftDelim
	}
	l.pos = Pos(len(l.input))
	// Correctly reached EOF.
	if l.pos > l.start {
		l.line += strings.Count(l.input[l.start:l.pos], "\n")
		l.emit(TEXT)
		l.ignore()
	}
	l.emit(EOF)

	return nil
}

// lexLeftDelim scans the left delimiter, which is known to be present
func lexLeftDelim(l *lexer) stateFn {
	l.pos += Pos(len(leftDelim))
	l.emit(L_DELIM)
	l.ignore()
	l.parenDepth = 0

	return lexInsideAction
}

// lexRightDelim scans the right delimiter, which is known to be present
func lexRightDelim(l *lexer) stateFn {
	l.pos += Pos(len(rightDelim))
	l.emit(R_DELIM)
	l.ignore()

	return lexText
}

// lexInsideAction scans the elements inside action delimiters.
func lexInsideAction(l *lexer) stateFn {
	// Either number, quoted string, or identifier.
	// Spaces separate arguments; runs of spaces turn into itemSpace.
	// Pipe symbols separate and are emitted.
	if delim := l.atRightDelim(); delim {
		if l.parenDepth == 0 {
			return lexRightDelim
		}

		return l.errorf("unclosed left paren")
	}
	switch r := l.next(); {
	case r == eof || isEndOfLine(r):
		return l.errorf("unclosed action")
	case isSpace(r):
		l.backup() // Put space back in case we have " -}}".

		return lexSpace
	case r == '@':
		l.emit(AT_SIGN)

		return lexIdentifier
	case r == ',':
		l.emit(COMMA)
	case isAlphaNumeric(r):
		l.backup()

		return lexIdentifier
	case r == '(':
		l.emit(L_PAREN)
		l.parenDepth++
	case r == ')':
		l.emit(R_PAREN)
		l.parenDepth--
		if l.parenDepth < 0 {
			return l.errorf("unexpected right paren %#U", r)
		}
	default:

		return l.errorf("unrecognized character in action: %#U", r)
	}

	return lexInsideAction
}

// lexSpace scans a run of space characters.
// We have not consumed the first space, which is known to be present.
func lexSpace(l *lexer) stateFn {
	var r rune
	for {
		r = l.peek()
		if !isSpace(r) {
			break
		}
		l.next()
	}
	l.emit(SPACE)

	return lexInsideAction
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			if !l.atTerminator() {
				return l.errorf("bad character %#U", r)
			}
			switch {
			case key[word] > keyword_beg && key[word] < keyword_end:
				l.emit(key[word])
			default:
				l.emit(IDENT)
			}

			break Loop
		}
	}

	return lexInsideAction
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || r == '-' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
