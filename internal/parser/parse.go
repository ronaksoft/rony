// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package parse builds parse trees for templates as defined by text/template
// and html/template. Clients should use those packages to construct templates
// rather than this one, which provides shared internal data structures not
// intended for general use.
package parse

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// Tree is the representation of a single parsed template.
type Tree struct {
	Name      string    // name of the template represented by the tree.
	ParseName string    // name of the top-level template during parsing, for error messages.
	Root      *ListNode // top-level root of the tree.
	text      string    // text parsed to create the template (or its parent)
	// Parsing only; cleared after parse.
	funcs     []map[string]interface{}
	lex       *lexer
	token     [3]tokenItem // three-token lookahead for parser.
	peekCount int
	vars      []string // variables defined at the moment.
	treeSet   map[string]*Tree
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
// func Parse(name, text, leftDelim, rightDelim string, funcs ...map[string]interface{}) (map[string]*Tree, error) {
// 	treeSet := make(map[string]*Tree)
// 	t := New(name)
// 	t.text = text
// 	_, err := t.Parse(text, leftDelim, rightDelim, treeSet, funcs...)
// 	return treeSet, err
// }

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

// Parsing.

// New allocates a new parse tree with the given name.
func New(name string, funcs ...map[string]interface{}) *Tree {
	return &Tree{
		Name:  name,
		funcs: funcs,
	}
}

// ErrorContext returns a textual representation of the location of the node in the input text.
// The receiver is only used when the node does not have a pointer to the tree inside,
// which can occur in old code.
func (t *Tree) ErrorContext(n Node) (location, context string) {
	pos := int(n.Position())
	tree := n.tree()
	if tree == nil {
		tree = t
	}
	text := tree.text[:pos]
	byteNum := strings.LastIndex(text, "\n")
	if byteNum == -1 {
		byteNum = pos // On first line.
	} else {
		byteNum++ // After the newline.
		byteNum = pos - byteNum
	}
	lineNum := 1 + strings.Count(text, "\n")
	context = n.String()
	return fmt.Sprintf("%s:%d:%d", tree.ParseName, lineNum, byteNum), context
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
	token := t.next()
	if token.tok != expected {
		t.unexpected(token, context)
	}
	return token
}

// expectOneOf consumes the next token and guarantees it has one of the required types.
func (t *Tree) expectOneOf(expected1, expected2 token, context string) tokenItem {
	token := t.next()
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
func (t *Tree) startParse(funcs []map[string]interface{}, lex *lexer, treeSet map[string]*Tree) {
	t.Root = nil
	t.lex = lex
	t.vars = []string{"$"}
	t.funcs = funcs
	t.treeSet = treeSet
}

// stopParse terminates parsing.
func (t *Tree) stopParse() {
	t.lex = nil
	t.vars = nil
	t.funcs = nil
	t.treeSet = nil
}

// Parse parses the template definition string to construct a representation of
// the template for execution. If either action delimiter string is empty, the
// default ("{{" or "}}") is used. Embedded template definitions are added to
// the treeSet map.
// func (t *Tree) Parse(text string, treeSet map[string]*Tree, funcs ...map[string]interface{}) (tree *Tree, err error) {
// 	defer t.recover(&err)
// 	t.ParseName = t.Name
// 	t.startParse(funcs, lex(t.Name, text), treeSet)
// 	t.text = text
// 	t.parse()
// 	t.add()
// 	t.stopParse()
// 	return t, nil
// }

// add adds tree to t.treeSet.
func (t *Tree) add() {
	tree := t.treeSet[t.Name]
	if tree == nil || IsEmptyTree(tree.Root) {
		t.treeSet[t.Name] = t
		return
	}
	if !IsEmptyTree(t.Root) {
		t.errorf("template: multiple definition of template %q", t.Name)
	}
}

// IsEmptyTree reports whether this tree (node) is empty of everything but space.
func IsEmptyTree(n Node) bool {
	switch n := n.(type) {
	case nil:
		return true
	case *ListNode:
		for _, node := range n.Nodes {
			if !IsEmptyTree(node) {
				return false
			}
		}
		return true
	case *ModelNode:
		return len(bytes.TrimSpace(n.Text)) == 0
	default:
		panic("unknown node: " + n.String())
	}
	// return false
}

// // parse is the top-level parser for a template, essentially the same
// // as itemList except it also parses {{define}} actions.
// // It runs to EOF.
// func (t *Tree) parse() {
// 	t.Root = t.newList(t.peek().pos)
// 	for t.peek().tok != EOF {
// 		if t.peek().tok == L_DELIM {
// 			delim := t.next()
// 			if t.next().tok == itemDefine {
// 				newT := New("definition") // name will be updated once we know it.
// 				newT.text = t.text
// 				newT.ParseName = t.ParseName
// 				newT.startParse(t.funcs, t.lex, t.treeSet)
// 				newT.parseDefinition()
// 				continue
// 			}
// 			t.backup2(delim)
// 		}
// 		switch n := t.textOrAction(); n.Type() {
// 		case nodeEnd, nodeElse:
// 			t.errorf("unexpected %s", n)
// 		default:
// 			t.Root.append(n)
// 		}
// 	}
// }

// // parseDefinition parses a {{define}} ...  {{end}} template definition and
// // installs the definition in t.treeSet. The "define" keyword has already
// // been scanned.
// func (t *Tree) parseDefinition() {
// 	const context = "define clause"
// 	name := t.expectOneOf(itemString, itemRawString, context)
// 	var err error
// 	t.Name, err = strconv.Unquote(name.val)
// 	if err != nil {
// 		t.error(err)
// 	}
// 	t.expect(itemRightDelim, context)
// 	var end Node
// 	t.Root, end = t.itemList()
// 	if end.Type() != nodeEnd {
// 		t.errorf("unexpected %s in %s", end, context)
// 	}
// 	t.add()
// 	t.stopParse()
// }

// // itemList:
// //	textOrAction*
// // Terminates at {{end}} or {{else}}, returned separately.
// func (t *Tree) itemList() (list *ListNode, next Node) {
// 	list = t.newList(t.peekNonSpace().pos)
// 	for t.peekNonSpace().tok != itemEOF {
// 		n := t.textOrAction()
// 		switch n.Type() {
// 		case nodeEnd, nodeElse:
// 			return list, n
// 		}
// 		list.append(n)
// 	}
// 	t.errorf("unexpected EOF")
// 	return
// }

// // textOrAction:
// //	text | action
// func (t *Tree) textOrAction() Node {
// 	switch token := t.nextNonSpace(); token.tok {
// 	case itemText:
// 		return t.newModel(token.pos, token.val)
// 	case itemLeftDelim:
// 		return t.action()
// 	default:
// 		t.unexpected(token, "input")
// 	}
// 	return nil
// }

// // Action:
// //	control
// //	command ("|" command)*
// // Left delim is past. Now get actions.
// // First word could be a keyword such as range.
// func (t *Tree) action() (n Node) {
// 	switch token := t.nextNonSpace(); token.tok {
// 	case itemBlock:
// 		return t.blockControl()
// 	case itemElse:
// 		return t.elseControl()
// 	case itemEnd:
// 		return t.endControl()
// 	case itemIf:
// 		return t.ifControl()
// 	case itemRange:
// 		return t.rangeControl()
// 	case itemTemplate:
// 		return t.templateControl()
// 	case itemWith:
// 		return t.withControl()
// 	}
// 	t.backup()
// 	token := t.peek()
// 	// Do not pop variables; they persist until "end".
// 	return t.newAction(token.pos, token.line, t.pipeline("command"))
// }

// If:
//	{{if pipeline}} itemList {{end}}
//	{{if pipeline}} itemList {{else}} itemList {{end}}
// If keyword is past.
// func (t *Tree) ifControl() Node {
// 	return t.newIf(t.parseControl(true, "if"))
// }

// Range:
//	{{range pipeline}} itemList {{end}}
//	{{range pipeline}} itemList {{else}} itemList {{end}}
// Range keyword is past.
// func (t *Tree) rangeControl() Node {
// 	return t.newRange(t.parseControl(false, "range"))
// }

// // Else:
// //	{{else}}
// // Else keyword is past.
// func (t *Tree) elseControl() Node {
// 	// Special case for "else if".
// 	peek := t.peekNonSpace()
// 	if peek.tok == itemIf {
// 		// We see "{{else if ... " but in effect rewrite it to {{else}}{{if ... ".
// 		return t.newElse(peek.pos, peek.line)
// 	}
// 	token := t.expect(itemRightDelim, "else")
// 	return t.newElse(token.pos, token.line)
// }
//
// // Block:
// //	{{block stringValue pipeline}}
// // Block keyword is past.
// // The name must be something that can evaluate to a string.
// // The pipeline is mandatory.
// func (t *Tree) blockControl() Node {
// 	const context = "block clause"
//
// 	token := t.nextNonSpace()
// 	name := t.parseTemplateName(token, context)
// 	pipe := t.pipeline(context)
//
// 	block := New(name) // name will be updated once we know it.
// 	block.text = t.text
// 	block.ParseName = t.ParseName
// 	block.startParse(t.funcs, t.lex, t.treeSet)
// 	var end Node
// 	block.Root, end = block.itemList()
// 	if end.Type() != nodeEnd {
// 		t.errorf("unexpected %s in %s", end, context)
// 	}
// 	block.add()
// 	block.stopParse()
//
// 	return t.newTemplate(token.pos, token.line, name, pipe)
// }

// // term:
// //	literal (number, string, nil, boolean)
// //	function (identifier)
// //	.
// //	.Field
// //	$
// //	'(' pipeline ')'
// // A term is a simple "expression".
// // A nil return means the next tokenItem is not a term.
// func (t *Tree) term() Node {
// 	switch token := t.nextNonSpace(); token.tok {
// 	case itemError:
// 		t.errorf("%s", token.val)
// 	case itemIdentifier:
// 		if !t.hasFunction(token.val) {
// 			t.errorf("function %q not defined", token.val)
// 		}
// 		return NewIdentifier(token.val).SetTree(t).SetPos(token.pos)
// 	case itemDot:
// 		return t.newDot(token.pos)
// 	case itemNil:
// 		return t.newNil(token.pos)
// 	case itemVariable:
// 		return t.useVar(token.pos, token.val)
// 	case itemField:
// 		return t.newField(token.pos, token.val)
// 	case itemBool:
// 		return t.newBool(token.pos, token.val == "true")
// 	case itemCharConstant, itemComplex, itemNumber:
// 		number, err := t.newNumber(token.pos, token.val, token.tok)
// 		if err != nil {
// 			t.error(err)
// 		}
// 		return number
// 	case itemLeftParen:
// 		pipe := t.pipeline("parenthesized pipeline")
// 		if token := t.next(); token.tok != itemRightParen {
// 			t.errorf("unclosed right paren: unexpected %s", token)
// 		}
// 		return pipe
// 	case itemString, itemRawString:
// 		s, err := strconv.Unquote(token.val)
// 		if err != nil {
// 			t.error(err)
// 		}
// 		return t.newString(token.pos, token.val, s)
// 	}
// 	t.backup()
// 	return nil
// }
