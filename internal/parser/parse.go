// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package parse builds parse trees for templates as defined by text/template
// and html/template. Clients should use those packages to construct templates
// rather than this one, which provides shared internal data structures not
// intended for general use.
package parse

import (
	"fmt"
	"runtime"
)

// Tree is the representation of a single parsed template.
type Tree struct {
	Name      string    // name of the template represented by the tree.
	ParseName string    // name of the top-level template during parsing, for error messages.
	Root      *ListNode // top-level root of the tree.
	text      string    // text parsed to create the template (or its parent)
	// Parsing only; cleared after parse.
	lex       *lexer
	token     [3]tokenItem // three-token lookahead for parser.
	peekCount int
}

// Copy returns a copy of the Tree. Any parsing state is discarded.
func (t *Tree) Copy() *Tree {
	if t == nil {
		return nil
	}
	return &Tree{
		Name:      t.Name,
		ParseName: t.ParseName,
		Root:      t.Root.CopyList(),
		text:      t.text,
	}
}

// Parse returns a map from template name to parse.Tree, created by parsing the
// templates described in the argument string. The top-level template will be
// given the specified name. If an error is encountered, parsing stops and an
// empty map is returned with the error.
func Parse(name, text string) (*Tree, error) {
	t := New(name)
	t.text = text
	_, err := t.Parse(text)
	return t, err
}

// next returns the next token.
func (t *Tree) next() tokenItem {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextItem()
	}
	return t.token[t.peekCount]
}

// backup backs the input stream up one token.
func (t *Tree) backup() {
	t.peekCount++
}

// backup2 backs the input stream up two tokens.
// The zeroth token is already there.
func (t *Tree) backup2(t1 tokenItem) {
	t.token[1] = t1
	t.peekCount = 2
}

// backup3 backs the input stream up three tokens
// The zeroth token is already there.
func (t *Tree) backup3(t2, t1 tokenItem) { // Reverse order: we're pushing back.
	t.token[1] = t1
	t.token[2] = t2
	t.peekCount = 3
}

// peek returns but does not consume the next token.
func (t *Tree) peek() tokenItem {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.lex.nextItem()
	return t.token[0]
}

// nextNonSpace returns the next non-space token.
func (t *Tree) nextNonSpace() (token tokenItem) {
	for {
		token = t.next()
		if token.tok != SPACE {
			break
		}
	}
	return token
}

// peekNonSpace returns but does not consume the next non-space token.
func (t *Tree) peekNonSpace() tokenItem {
	token := t.nextNonSpace()
	t.backup()
	return token
}

// New allocates a new parse tree with the given name.
func New(name string) *Tree {
	return &Tree{
		Name: name,
	}
}

// errorf formats the error and terminates processing.
func (t *Tree) errorf(format string, args ...interface{}) {
	t.Root = nil
	format = fmt.Sprintf("template: %s:%d: %s", t.ParseName, t.token[0].line, format)
	panic(fmt.Errorf(format, args...))
}

// error terminates processing.
func (t *Tree) error(err error) {
	t.errorf("%s", err)
}

// expect consumes the next token and guarantees it has the required type.
func (t *Tree) expect(expected token, context string) tokenItem {
	token := t.nextNonSpace()
	if token.tok != expected {
		t.unexpected(token, context)
	}
	return token
}

// expectOneOf consumes the next token and guarantees it has one of the required types.
func (t *Tree) expectOneOf(expected1, expected2 token, context string) tokenItem {
	token := t.nextNonSpace()
	if token.tok != expected1 && token.tok != expected2 {
		t.unexpected(token, context)
	}
	return token
}

// unexpected complains about the token and terminates processing.
func (t *Tree) unexpected(token tokenItem, context string) {
	t.errorf("unexpected %s in %s", token, context)
}

// recover is the handler that turns panics into returns from the top level of Parse.
func (t *Tree) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		if t != nil {
			t.lex.drain()
			t.stopParse()
		}
		*errp = e.(error)
	}
}

// startParse initializes the parser, using the lexer.
func (t *Tree) startParse(lex *lexer) {
	t.Root = nil
	t.lex = lex
}

// stopParse terminates parsing.
func (t *Tree) stopParse() {
	t.lex = nil
}

// Parse parses the template definition string to construct a representation of
// the template for execution. If either action delimiter string is empty, the
// default ("{{" or "}}") is used. Embedded template definitions are added to
// the treeSet map.
func (t *Tree) Parse(text string) (tree *Tree, err error) {
	defer t.recover(&err)
	t.ParseName = t.Name
	t.startParse(lex(t.Name, text))
	t.text = text
	t.parse()
	t.stopParse()
	return t, nil
}

// parse is the top-level parser for a template, essentially the same
// as itemList except it also parses {{define}} actions.
// It runs to EOF.
func (t *Tree) parse() {
	t.Root = t.newList(t.peek().pos)
	for t.peek().tok != EOF {
		switch item := t.nextNonSpace(); item.tok {
		case TEXT:
			t.Root.append(t.newText(item.pos, item.val))
		case L_DELIM:
			item = t.nextNonSpace()
			t.Root.append(t.action())
		default:
			t.unexpected(item, "input")
		}
	}
}

// Action:
//	model, table, view, counter , ...
// Left delim is past. Now get actions.
func (t *Tree) action() (n Node) {
	switch item := t.nextNonSpace(); item.tok {
	case MODEL:
		n = t.model(item.pos)
	case TABLE:
		n = t.table(item.pos)
	case VIEW:
		n = t.view(item.pos)
	case COUNTER:
		n = t.counter(item.pos)
	default:
		t.unexpected(item, "action")
	}
	if item := t.nextNonSpace(); item.tok != R_DELIM {
		t.unexpected(item, "action right delmic")
	}
	return
}

func (t *Tree) model(pos Pos) (n Node) {
	switch item := t.nextNonSpace(); item.tok {
	case IDENT:
		return t.newModel(pos, item.val)
	default:
		t.unexpected(item, "model")
	}

	panic("unreachable code")
}

func (t *Tree) idents() []string {
	var args []string
Loop:
	for {
		item := t.nextNonSpace()
		switch item.tok {
		case IDENT:
			args = append(args, item.val)
		case COMMA:
		case R_PAREN:
			break Loop
		default:
			t.unexpected(item, "idents")
		}
	}
	return args
}

func (t *Tree) table(pos Pos) (n Node) {
	var (
		pks, cks []string
	)
	// we must see at least one '('
	if item := t.nextNonSpace(); item.tok != L_PAREN {
		t.unexpected(item, "table")
	}

	// if we see second '(' then we have composed partition keys
	switch item := t.nextNonSpace(); item.tok {
	case L_PAREN:
		pks = t.idents()
		cks = t.idents()
	default:
		// put back to buffer, since we did not pick '('
		t.backup()

		ks := t.idents()
		switch len(ks) {
		case 0:
			t.errorf("no primary key for table: (%d:%d)", item.line, item.pos)
		case 1:
			pks = append(pks, ks[0])
		default:
			pks = append(pks, ks[0])
			cks = append(cks, ks[1:]...)
		}
	}
	return t.newTable(pos, pks, cks)
}

func (t *Tree) view(pos Pos) (n Node) {
	var (
		pks, cks []string
	)
	// we must see at least one '('
	if item := t.nextNonSpace(); item.tok != L_PAREN {
		t.unexpected(item, "table")
	}

	// if we see second '(' then we have composed partition keys
	switch item := t.nextNonSpace(); item.tok {
	case L_PAREN:
		pks = t.idents()
		cks = t.idents()
	default:
		// put back to buffer, since we did not pick '('
		t.backup()

		ks := t.idents()
		switch len(ks) {
		case 0:
			t.errorf("no primary key for view: (%d:%d)", item.line, item.pos)
		case 1:
			pks = append(pks, ks[0])
		default:
			pks = append(pks, ks[0])
			cks = append(cks, ks[1:]...)
		}
	}
	return t.newView(pos, pks, cks)
}

func (t *Tree) counter(pos Pos) (n Node) {
	switch item := t.nextNonSpace(); item.tok {
	case IDENT:
		return t.newCounter(pos, item.val)
	default:
		t.unexpected(item, "counter")
	}

	panic("unreachable code")
}
