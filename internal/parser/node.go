// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Parse nodes.

package parse

import (
	"fmt"
	"strings"
)

var textFormat = "%s" // Changed to "%q" in tests for better error messages.

// A Node is an element in the parse tree. The interface is trivial.
// The interface contains an unexported method so that only
// types local to this package can satisfy it.
type Node interface {
	Type() NodeType
	String() string
	// Copy does a deep copy of the Node and all its components.
	// To avoid type assertions, some XxxNodes also have specialized
	// CopyXxx methods that return *XxxNode.
	Copy() Node
	Position() Pos // byte position of start of node in full original input string
	// tree returns the containing *Tree.
	// It is unexported so all implementations of Node are in this package.
	tree() *Tree
	// writeTo writes the String output to the builder.
	writeTo(*strings.Builder)
}

// NodeType identifies the type of a parse tree node.
type NodeType int

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

func (p Pos) Position() Pos {
	return p
}

// Type returns itself and provides an easy default implementation
// for embedding in a Node. Embedded in all non-trivial Nodes.
func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeList NodeType = iota // a list of Nodes
	NodeModel
	NodeTable
	NodeView
	NodePartitionKey
	NodeClusteringKey
)

// Nodes.

// ListNode holds a sequence of nodes.
type ListNode struct {
	NodeType
	Pos
	tr    *Tree
	Nodes []Node // The element nodes in lexical order.
}

func (t *Tree) newList(pos Pos) *ListNode {
	return &ListNode{tr: t, NodeType: NodeList, Pos: pos}
}

func (l *ListNode) append(n Node) {
	l.Nodes = append(l.Nodes, n)
}

func (l *ListNode) tree() *Tree {
	return l.tr
}

func (l *ListNode) String() string {
	var sb strings.Builder
	l.writeTo(&sb)
	return sb.String()
}

func (l *ListNode) writeTo(sb *strings.Builder) {
	for _, n := range l.Nodes {
		n.writeTo(sb)
	}
}

func (l *ListNode) CopyList() *ListNode {
	if l == nil {
		return l
	}
	n := l.tr.newList(l.Pos)
	for _, elem := range l.Nodes {
		n.append(elem.Copy())
	}
	return n
}

func (l *ListNode) Copy() Node {
	return l.CopyList()
}

// ModelNode holds plain text.
type ModelNode struct {
	NodeType
	Pos
	tr   *Tree
	Text []byte // The text; may span newlines.
}

func (t *Tree) newModel(pos Pos, text string) *ModelNode {
	return &ModelNode{tr: t, NodeType: NodeModel, Pos: pos, Text: []byte(text)}
}

func (t *ModelNode) String() string {
	return fmt.Sprintf(textFormat, t.Text)
}

func (t *ModelNode) writeTo(sb *strings.Builder) {
	sb.WriteString(t.String())
}

func (t *ModelNode) tree() *Tree {
	return t.tr
}

func (t *ModelNode) Copy() Node {
	return &ModelNode{tr: t.tr, NodeType: NodeModel, Pos: t.Pos, Text: append([]byte{}, t.Text...)}
}

type TableNode struct {
	NodeType
	Pos
	tr             *Tree
	PrimaryKeys    []*PartitionKeyNode
	ClusteringKeys []*ClusteringKeyNode
}

type PartitionKeyNode struct {
	NodeType
	Pos
	tr   *Tree
	Name []byte
}

type ClusteringKeyNode struct {
	NodeType
	Pos
	tr   *Tree
	Name []byte
}

// // AssignNode holds a list of variable names, possibly with chained field
// // accesses. The dollar sign is part of the (first) name.
// type VariableNode struct {
// 	NodeType
// 	Pos
// 	tr    *Tree
// 	Ident []string // Variable name and fields in lexical order.
// }
//
// func (t *Tree) newVariable(pos Pos, ident string) *VariableNode {
// 	return &VariableNode{tr: t, NodeType: NodeVariable, Pos: pos, Ident: strings.Split(ident, ".")}
// }
//
// func (v *VariableNode) String() string {
// 	var sb strings.Builder
// 	v.writeTo(&sb)
// 	return sb.String()
// }
//
// func (v *VariableNode) writeTo(sb *strings.Builder) {
// 	for i, id := range v.Ident {
// 		if i > 0 {
// 			sb.WriteByte('.')
// 		}
// 		sb.WriteString(id)
// 	}
// }
//
// func (v *VariableNode) tree() *Tree {
// 	return v.tr
// }
//
// func (v *VariableNode) Copy() Node {
// 	return &VariableNode{tr: v.tr, NodeType: NodeVariable, Pos: v.Pos, Ident: append([]string{}, v.Ident...)}
// }
//
