package tools

import (
	"fmt"
	"strings"
	"sync"
)

/*
   Creation Time: 2019 - Jun - 20
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	nodePool sync.Pool
)

func acquireNode() *Node {
	n, ok := nodePool.Get().(*Node)
	if !ok {
		return &Node{}
	}
	return n
}

func releaseNode(n *Node) {
	*n = Node{}
	nodePool.Put(n)
}

type Node struct {
	next *Node
	prev *Node
	data interface{}
}

func (n Node) GetData() interface{} {
	return n.data
}

type LinkedList struct {
	lock sync.RWMutex
	head *Node
	tail *Node
	size int32
}

func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

func (ll *LinkedList) Append(data interface{}) {
	ll.lock.Lock()
	ll.size += 1
	n := acquireNode()
	n.data = data
	if ll.tail == nil {
		ll.tail, ll.head = n, n
	} else {
		o := ll.tail
		ll.tail = n
		ll.tail.prev = o
		o.next = ll.tail
	}
	ll.lock.Unlock()
}

func (ll *LinkedList) Prepend(data interface{}) {
	ll.lock.Lock()
	ll.size += 1
	n := acquireNode()
	n.data = data
	if ll.head == nil {
		ll.tail, ll.head = n, n
	} else {
		o := ll.head
		ll.head = n
		ll.head.next = o
		o.prev = ll.head
	}
	ll.lock.Unlock()
}

func (ll *LinkedList) Size() int32 {
	ll.lock.RLock()
	n := ll.size
	ll.lock.RUnlock()
	return n
}

func (ll *LinkedList) Head() *Node {
	return ll.head
}

func (ll *LinkedList) PickHeadData() interface{} {
	ll.lock.Lock()
	if ll.head == nil {
		ll.lock.Unlock()
		return nil
	}

	n := ll.head
	if ll.head.next != nil {
		ll.head = ll.head.next
		ll.head.prev = nil
	} else {
		ll.head, ll.tail = nil, nil
	}
	ll.size--
	ll.lock.Unlock()
	data := n.data
	releaseNode(n)
	return data
}

func (ll *LinkedList) Tail() *Node {
	return ll.tail
}

func (ll *LinkedList) PickTailData() interface{} {
	ll.lock.Lock()
	if ll.tail == nil {
		ll.lock.Unlock()
		return nil
	}
	n := ll.tail
	if ll.tail.prev != nil {
		ll.tail = ll.tail.prev
		ll.tail.next = nil
	} else {
		ll.head, ll.tail = nil, nil
	}
	ll.size -= 1
	ll.lock.Unlock()
	data := n.data
	releaseNode(n)
	return data
}

func (ll *LinkedList) Get(index int32) (n *Node) {
	ll.lock.RLock()
	if index < ll.size<<1 {
		n = ll.head
		for index > 0 {
			n = n.next
			index--
		}

	} else {
		n = ll.tail
		for index > 0 {
			n = n.prev
			index--
		}
	}
	ll.lock.RUnlock()
	return n
}

func (ll *LinkedList) RemoveAt(index int32) {
	n := ll.Get(index)
	ll.lock.Lock()
	if n.next != nil {
		n.next.prev = n.prev
	}
	if n.prev != nil {
		n.prev.next = n.next
	}
	ll.size--
	ll.lock.Unlock()
}

func (ll *LinkedList) Reset() {
	for ll.PickHeadData() != nil {
	}
}

func (ll *LinkedList) String() string {
	sb := strings.Builder{}
	n := ll.head
	idx := 0
	for n != nil {
		sb.WriteString(fmt.Sprintf("%d. %v\n", idx, n.data))
		n = n.next
		idx++
	}
	return sb.String()
}
