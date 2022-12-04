package bplustree

import (
	"fmt"

	"github.com/golang/glog"
)

type BPlusTree struct {
	n    int // n paramater of BPlusTree
	root *tnode
}

type tnode struct {
	parent   *tnode
	isLeaf   bool
	keys     []int
	pointers []interface{}
}

func newTNode(isLeaf bool, capcity int) *tnode {
	n := &tnode{
		isLeaf:   isLeaf,
		keys:     make([]int, 0, capcity),
		pointers: make([]interface{}, 1, capcity+1),
	}

	return n
}

func (t *BPlusTree) insertInParent(n *tnode, key int, newN *tnode) error {
	glog.V(2).Infof("insert into parent, key: %d, node: %+v", key, newN)
	if n == t.root {
		newRoot := newTNode(false, t.n)
		n.parent = newRoot
		newN.parent = newRoot
		newRoot.pointers = make([]interface{}, 0, t.n+1)
		newRoot.pointers = append(newRoot.pointers, n, newN)
		newRoot.keys = append(newRoot.keys, key)
		t.root = newRoot

		return nil
	}

	parent := n.parent
	parent.insertNonLeaf(key, newN)
	if len(parent.pointers) <= t.n {
		return nil
	}

	// split at ceil((n + 1) / 2)
	// (1) n == 3 -> 2
	// (2) n == 4 -> 3
	pos := t.n/2 + 1
	newInternalN := newTNode(false, t.n)
	newInternalN.parent = parent.parent

	// split pointers
	newInternalN.pointers = newInternalN.pointers[:len(parent.pointers[pos:])]
	copy(newInternalN.pointers, parent.pointers[pos:])
	// adjust parent for splited children
	for _, p := range newInternalN.pointers {
		child := p.(*tnode)
		child.parent = newInternalN
	}
	parent.pointers = parent.pointers[:pos]

	// split keys
	newKey := parent.keys[pos-1]
	newInternalN.keys = newInternalN.keys[:len(parent.keys[pos:])]
	copy(newInternalN.keys, parent.keys[pos:])
	parent.keys = parent.keys[:pos-1]

	return t.insertInParent(parent, newKey, newInternalN)
}

func (t *BPlusTree) findLeafForInsert(key int) *tnode {
	r := t.root

	for !r.isLeaf {
		pos := r.findInsertPos(key)
		if pos == len(r.keys) {
			r = r.pointers[pos].(*tnode)
		} else if r.keys[pos] > key {
			r = r.pointers[pos].(*tnode)
		} else { // r.keys[pos] == key
			r = r.pointers[pos+1].(*tnode)
		}
	}

	return r
}

func (t *BPlusTree) Insert(key int, p interface{}) error {
	// just insert root
	n := t.findLeafForInsert(key)
	glog.V(2).Infof("insert key %d to leaf: %+v, parent: %+v", key, n, n.parent)
	if err := n.insertLeaf(key, p); err != nil {
		return fmt.Errorf("error insert key: %d, err: %+v", key, err)
	}

	if len(n.pointers) <= t.n {
		return nil
	}

	// len(n.pointers) == t.n + 1
	newN := newTNode(true, t.n)
	newN.parent = n.parent

	// split leaf at ceil(n/2) -> [0, ceil(n/2)] /CUP [ceil(n/2)+1:]
	pos := (t.n + 1) / 2 // (tn + 1) / 2 == ceil(n / 2)

	newN.pointers = newN.pointers[:len(n.pointers[pos:])]
	copy(newN.pointers, n.pointers[pos:])
	n.pointers = n.pointers[:pos+1]

	newN.keys = newN.keys[:len(n.keys[pos:])]
	copy(newN.keys, n.keys[pos:])
	n.keys = n.keys[:pos]

	// connect to sibling
	n.pointers[pos] = newN

	return t.insertInParent(n, newN.keys[0], newN)
}

func (t *BPlusTree) Find(key int) (interface{}, error) {
	// TODO: find node
	return t.root, nil
}

func (t *BPlusTree) Delete(key int) error {
	panic("no impl")
}

func (tn *tnode) find(key int) (int, error) {
	s, e := 0, len(tn.keys)
	for s < e {
		m := (s + e) / 2
		if tn.keys[m] > key {
			e = m
		} else if tn.keys[m] < key {
			s = m + 1
		} else { // tn.keys[m] == key
			return m, nil
		}
	}

	return 0, fmt.Errorf("key %d not found", key)
}

func (tn *tnode) insert(key int, p interface{}) error {
	if tn.isLeaf {
		return tn.insertLeaf(key, p)
	}

	return tn.insertNonLeaf(key, p)
}

func (tn *tnode) insertNonLeaf(key int, p interface{}) error {
	pos := tn.findInsertPos(key)
	if pos < len(tn.keys) && tn.keys[pos] == key {
		return fmt.Errorf("duplicate key %d in internal node, keys: %+v", key, tn.keys)
	}

	tn.insertNonLeafAt(pos, key, p)
	return nil
}

func (tn *tnode) findInsertPos(key int) int {
	// find smallest k in tn.keys greater or equal to key
	s, e := 0, len(tn.keys)
	for s < e {
		m := (s + e) / 2
		if tn.keys[m] < key {
			s = m + 1
		} else { // tn.keys[m] >= key
			e = m
		}
	}

	return s
}

func (tn *tnode) insertLeaf(key int, p interface{}) error {
	pos := tn.findInsertPos(key)
	if pos < len(tn.keys) && tn.keys[pos] == key {
		return fmt.Errorf("duplicate key in leaf node: %d, keys: %+v", key, tn.keys)
	}

	tn.insertLeafAt(pos, key, p)
	return nil
}

func (tn *tnode) insertNonLeafAt(index int, key int, p interface{}) {
	nsz := len(tn.keys) + 1
	if nsz > cap(tn.keys) {
		glog.Fatalf("node key size overflow: %d vs %d", nsz, cap(tn.keys))
	}

	tn.keys = tn.keys[:nsz]
	copy(tn.keys[index+1:], tn.keys[index:nsz-1])
	tn.keys[index] = key

	// one more element thant keys
	tn.pointers = tn.pointers[:nsz+1]
	copy(tn.pointers[index+2:], tn.pointers[index+1:nsz])
	tn.pointers[index+1] = p
}

func (tn *tnode) insertLeafAt(index int, key int, p interface{}) {
	nsz := len(tn.keys) + 1
	if nsz > cap(tn.keys) {
		glog.Fatalf("node key size overflow: %d vs %d", nsz, cap(tn.keys))
	}

	tn.keys = tn.keys[:nsz]
	copy(tn.keys[index+1:], tn.keys[index:nsz-1])
	tn.keys[index] = key

	// one more element thant keys
	tn.pointers = tn.pointers[:nsz+1]
	copy(tn.pointers[index+1:], tn.pointers[index:nsz])
	tn.pointers[index] = p
}

func NewTree(n int) (*BPlusTree, error) {
	if n < 3 {
		return nil, fmt.Errorf("illegal n of BPlusTree, should be greater than 3: %d", n)
	}

	t := &BPlusTree{
		n:    n,
		root: newTNode(true, n),
	}

	return t, nil
}
