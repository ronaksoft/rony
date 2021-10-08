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
	NodeText
	NodeModel
	NodeTable
	NodeView
	NodeCounter
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

// TextNode holds plain text.
type TextNode struct {
	NodeType
	Pos
	tr   *Tree
	Text string // The text; may span newlines.
}

func (t *Tree) newText(pos Pos, text string) *TextNode {
	return &TextNode{tr: t, NodeType: NodeText, Pos: pos, Text: text}
}

func (t *TextNode) String() string {
	return fmt.Sprintf(textFormat, t.Text)
}

func (t *TextNode) writeTo(sb *strings.Builder) {
	sb.WriteString(t.String())
}

func (t *TextNode) tree() *Tree {
	return t.tr
}

func (t *TextNode) Copy() Node {
	return &TextNode{tr: t.tr, NodeType: NodeText, Pos: t.Pos, Text: t.Text}
}

// ModelNode holds model
type ModelNode struct {
	NodeType
	Pos
	tr   *Tree
	Text string // The text; may span newlines.
}

func (t *Tree) newModel(pos Pos, text string) *ModelNode {
	return &ModelNode{tr: t, NodeType: NodeModel, Pos: pos, Text: text}
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
	return &ModelNode{tr: t.tr, NodeType: NodeModel, Pos: t.Pos, Text: t.Text}
}

// TableNode holds table
type TableNode struct {
	NodeType
	Pos
	tr             *Tree
	PartitionKeys  []string
	ClusteringKeys []string
}

func (t *Tree) newTable(pos Pos, pks, cks []string) *TableNode {
	return &TableNode{tr: t, NodeType: NodeTable, Pos: pos, PartitionKeys: pks, ClusteringKeys: cks}
}

func (t *TableNode) String() string {
	return fmt.Sprintf("PKs: %v, CKs: %v", t.PartitionKeys, t.ClusteringKeys)
}

func (t *TableNode) writeTo(sb *strings.Builder) {
	sb.WriteString(t.String())
}

func (t *TableNode) tree() *Tree {
	return t.tr
}

func (t *TableNode) Copy() Node {
	return &TableNode{
		tr:             t.tr,
		NodeType:       NodeTable,
		Pos:            t.Pos,
		PartitionKeys:  t.PartitionKeys,
		ClusteringKeys: t.ClusteringKeys,
	}
}

// ViewNode holds view
type ViewNode struct {
	NodeType
	Pos
	tr             *Tree
	PartitionKeys  []string
	ClusteringKeys []string
}

func (t *Tree) newView(pos Pos, pks, cks []string) *ViewNode {
	return &ViewNode{tr: t, NodeType: NodeView, Pos: pos, PartitionKeys: pks, ClusteringKeys: cks}
}

func (t *ViewNode) String() string {
	return fmt.Sprintf("PKs: %v, CKs: %v", t.PartitionKeys, t.ClusteringKeys)
}

func (t *ViewNode) writeTo(sb *strings.Builder) {
	sb.WriteString(t.String())
}

func (t *ViewNode) tree() *Tree {
	return t.tr
}

func (t *ViewNode) Copy() Node {
	return &ViewNode{
		tr:             t.tr,
		NodeType:       NodeView,
		Pos:            t.Pos,
		PartitionKeys:  t.PartitionKeys,
		ClusteringKeys: t.ClusteringKeys,
	}
}

// CounterNode holds counter
type CounterNode struct {
	NodeType
	Pos
	tr   *Tree
	Text string // The text; may span newlines.
}

func (t *Tree) newCounter(pos Pos, text string) *CounterNode {
	return &CounterNode{tr: t, NodeType: NodeCounter, Pos: pos, Text: text}
}

func (t *CounterNode) String() string {
	return fmt.Sprintf(textFormat, t.Text)
}

func (t *CounterNode) writeTo(sb *strings.Builder) {
	sb.WriteString(t.String())
}

func (t *CounterNode) tree() *Tree {
	return t.tr
}

func (t *CounterNode) Copy() Node {
	return &CounterNode{tr: t.tr, NodeType: NodeCounter, Pos: t.Pos, Text: t.Text}
}
