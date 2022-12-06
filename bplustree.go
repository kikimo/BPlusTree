package bplustree

import (
	"fmt"

	"github.com/golang/glog"
)

var ErrKeyNotFound error = fmt.Errorf("key not found")

// TODO: check parent pointer invariant
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

func (t *BPlusTree) findLeaf(key int) *tnode {
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
	newEntry, err := t.doInsert(t.root, key, p)
	if err != nil {
		return err
	}

	if newEntry != nil {
		newRoot := newTNode(false, t.n)
		newRoot.pointers = newRoot.pointers[:2]
		newRoot.pointers[0] = t.root
		newRoot.pointers[1] = newEntry.node
		newRoot.keys = newRoot.keys[:1]
		newRoot.keys[0] = newEntry.key
		t.root = newRoot
	}

	return nil
}

func (t *BPlusTree) splitInternalNode(node *tnode) *entry {
	glog.V(2).Infof("spliting internal node: %+v", node)
	// split at ceil((n + 1) / 2)
	// (1) n == 3 -> 2
	// (2) n == 4 -> 3
	pos := t.n/2 + 1
	newN := newTNode(false, t.n)
	newN.parent = node.parent

	// split pointers
	newN.pointers = newN.pointers[:len(node.pointers[pos:])]
	copy(newN.pointers, node.pointers[pos:])
	// adjust parent for splited children
	for _, p := range newN.pointers {
		child := p.(*tnode)
		child.parent = newN
	}
	node.pointers = node.pointers[:pos]

	// split keys
	newKey := node.keys[pos-1]
	newN.keys = newN.keys[:len(node.keys[pos:])]
	copy(newN.keys, node.keys[pos:])
	node.keys = node.keys[:pos-1]

	glog.V(2).Infof("new entry after split internal node: key(%d), %+v", newKey, newN)
	return &entry{key: newKey, node: newN}
}

func (t *BPlusTree) splitLeafNode(node *tnode) *entry {
	glog.V(2).Infof("spliting leaf node: %+v", node)
	// split leaf at ceil(n/2) -> [0, ceil(n/2)] /CUP [ceil(n/2)+1:]
	pos := (t.n + 1) / 2 // (tn + 1) / 2 == ceil(n / 2)

	newN := newTNode(true, t.n)
	newN.parent = node.parent
	newN.pointers = newN.pointers[:len(node.pointers[pos:])]
	copy(newN.pointers, node.pointers[pos:])
	node.pointers = node.pointers[:pos+1]

	newN.keys = newN.keys[:len(node.keys[pos:])]
	copy(newN.keys, node.keys[pos:])
	node.keys = node.keys[:pos]

	// connect to sibling
	node.pointers[pos] = newN

	glog.V(2).Infof("new entry after split leaf: %+v", newN)
	return &entry{key: newN.keys[0], node: newN}
}

type entry struct {
	key  int
	node *tnode
}

func (t *BPlusTree) doInsert(root *tnode, key int, p interface{}) (*entry, error) {
	if root.isLeaf {
		if err := root.insertLeaf(key, p); err != nil {
			return nil, err
		}

		var newEntry *entry = nil
		if len(root.pointers) == t.n+1 {
			newEntry = t.splitLeafNode(root)
		}

		// invariant check
		if len(root.pointers) > t.n+1 {
			glog.Fatalf("illegal node pointer size: %+v", root)
		}

		return newEntry, nil
	}

	pos := root.findInsertPos(key)
	if pos < len(root.keys) && root.keys[pos] == key {
		// 1. pos >= len(root.keys) -> key is greater than all values in root.keys
		// 2. root.keys[pos] == key -> key resides in root.pointer[pos+1]
		pos += 1
	}

	newChild, err := t.doInsert(root.pointers[pos].(*tnode), key, p)
	if err != nil {
		return nil, err
	}

	if newChild == nil {
		return nil, nil
	}

	// insert newNode after pos
	root.insertNonLeafAt(pos, newChild.key, newChild.node)
	if len(root.pointers) <= t.n {
		return nil, nil
	}

	// invariant check
	if len(root.pointers) != t.n+1 {
		glog.Fatalf("illegal node pointer size: %+v", root)
	}

	newSibling := t.splitInternalNode(root)
	return newSibling, nil
}

// TODO: unit test me
func (t *BPlusTree) Find(key int) (interface{}, error) {
	leaf := t.findLeaf(key)
	pos := leaf.findInsertPos(key)
	if pos >= len(leaf.keys) || leaf.keys[pos] != key {
		return nil, ErrKeyNotFound
	}

	return leaf.pointers[pos], nil
}

func (t *BPlusTree) mergeNodes(left *tnode, key int, right *tnode) bool {
	if left.isLeaf && right.isLeaf {
		return t.mergeLeaves(left, right)
	} else if !left.isLeaf && !right.isLeaf {
		return t.mergeInternalNodes(left, key, right)
	}

	glog.Fatalf("merge leaf to internal node: %+v, %+v", left, right)
	panic("unreachable")
}

func (t *BPlusTree) mergeLeaves(left, right *tnode) bool {
	if len(left.pointers)+len(right.pointers)-1 > t.n {
		return false
	}

	left.pointers = left.pointers[:len(left.pointers)-1]
	// FIXME: ok to append() here?
	left.pointers = append(left.pointers, right.pointers...)
	left.keys = append(left.keys, right.keys...)

	if len(left.keys)+1 != len(left.pointers) {
		glog.Fatalf("illegal node status: %+v", left)
	}

	return true
}

func (t *BPlusTree) mergeInternalNodes(left *tnode, key int, right *tnode) bool {
	glog.V(2).Infof("merge internal node, left: %+v, right: %+v", left, right)
	if len(left.pointers)+len(right.pointers) > t.n {
		return false
	}

	// adjust parent pointer
	for _, p := range right.pointers {
		c := p.(*tnode)
		c.parent = left
	}

	// FIXME: ok to append() here?
	left.pointers = append(left.pointers, right.pointers...)
	left.keys = append(left.keys, key)
	left.keys = append(left.keys, right.keys...)

	if len(left.keys)+1 != len(left.pointers) {
		glog.Fatalf("illegal node status: %+v", left)
	}

	return true
}

func (t *BPlusTree) Delete(key int) error {
	deleted, err := t.deleteEntry(t.root, key)
	if err != nil {
		return fmt.Errorf("error deleting key %d: %+v", key, err)
	}

	if !deleted {
		return nil
	}

	if len(t.root.pointers) == 1 {
		if t.root.isLeaf {
			t.root.pointers = t.root.pointers[:0]
		} else {
			t.root = t.root.pointers[0].(*tnode)
		}
	}

	return nil
}

func (t *BPlusTree) borrowFromLeft(left *tnode, key *int, right *tnode) {
	if left.isLeaf && right.isLeaf {
		t.leafBorrowFromLeft(left, key, right)
		return
	}

	if !left.isLeaf && !right.isLeaf {
		t.internalNodeBorrowFromLeft(left, key, right)
		return
	}

	glog.Fatalf("leaf cannot borrow from internal node(vice vesa), left: %+v, right: %+v", left, right)
	panic("unreachable")
}

func (t *BPlusTree) borrowFromRight(left *tnode, key *int, right *tnode) {
	if left.isLeaf && right.isLeaf {
		t.leafBorrowFromRight(left, key, right)
		return
	}

	if !left.isLeaf && !right.isLeaf {
		t.internalNodeBorrowFromRight(left, key, right)
		return
	}

	glog.Fatalf("leaf cannot borrow from internal node(vice vesa), left: %+v, right: %+v", left, right)
	panic("unreachable")
}

func (t *BPlusTree) leafBorrowFromLeft(left *tnode, key *int, right *tnode) {
	sz := len(left.keys)
	k := left.keys[sz-1]
	p := left.pointers[sz-1]

	// shrink left by one
	left.keys = left.keys[:sz-1]
	left.pointers[sz-1] = left.pointers[sz]
	left.pointers = left.pointers[:sz]

	*key = k

	// prepend entry (k, p) to right
	// expand right first
	sz = len(right.keys)
	right.keys = right.keys[:sz+1]
	if sz > 0 {
		copy(right.keys[1:], right.keys[:sz-1])
	}
	right.pointers = right.pointers[:sz+2]
	copy(right.pointers[1:], right.pointers[:sz])
	right.keys[0] = k
	right.pointers[0] = p
}

func (t *BPlusTree) leafBorrowFromRight(left *tnode, key *int, right *tnode) {
	sz := len(right.keys)
	k := right.keys[0]
	p := right.pointers[0]

	// shrink right by one
	copy(right.keys, right.keys[1:])
	right.keys = right.keys[:sz-1]
	copy(right.pointers, right.pointers[1:])
	right.pointers = right.pointers[:sz]

	*key = k

	// append entry (k, p) to left
	// expand left first
	sz = len(left.keys)
	left.keys = left.keys[:sz+1]
	left.keys[sz] = k
	left.pointers = left.pointers[:sz+2]
	left.pointers[sz+1] = left.pointers[sz]
	left.pointers[sz] = p
}

func (t *BPlusTree) internalNodeBorrowFromLeft(left *tnode, key *int, right *tnode) {
	sz := len(left.keys)
	k := left.keys[sz-1]
	p := left.pointers[sz]

	// swap key and k
	*key, k = k, *key

	// shrink left by one
	left.keys = left.keys[:sz-1]
	left.pointers = left.pointers[:sz]

	// prepend entry (k, p) to right
	// expand right first
	sz = len(right.keys)
	right.keys = right.keys[:sz+1]
	if sz > 0 {
		copy(right.keys[1:], right.keys[:sz-1])
	}
	right.pointers = right.pointers[:sz+2]
	copy(right.pointers[1:], right.pointers[:sz])
	right.keys[0] = k
	right.pointers[0] = p
}

func (t *BPlusTree) internalNodeBorrowFromRight(left *tnode, key *int, right *tnode) {
	sz := len(right.keys)
	k := right.keys[0]
	p := right.pointers[0]

	// swap key and k
	*key, k = k, *key

	// shrink right by one
	copy(right.keys, right.keys[1:])
	right.keys = right.keys[:sz-1]
	copy(right.pointers, right.pointers[1:])
	right.pointers = right.pointers[:sz]

	// append entry (k, p) to left
	glog.Infof("left before borrow: %+v", left)
	sz = len(left.keys)
	left.keys = left.keys[:sz+1]
	left.keys[sz] = k
	left.pointers = left.pointers[:sz+2]
	left.pointers[sz+1] = p
}

func (t *BPlusTree) deleteEntry(root *tnode, key int) (bool, error) {
	if root.isLeaf {
		return true, root.deleteEntry(key)
	}

	pos := root.findInsertPos(key)
	if pos < len(root.keys) && root.keys[pos] == key {
		pos += 1
	}

	child := root.pointers[pos].(*tnode)
	deleted, err := t.deleteEntry(child, key)
	if err != nil || !deleted {
		return false, err
	}

	if !child.tooFewPointers() {
		return false, nil
	}

	// too few pointers, try merge entries
	if pos-1 >= 0 {
		if t.mergeNodes(root.pointers[pos-1].(*tnode), root.keys[pos-1], child) {
			root.deleteEntryAt(pos - 1)
			return true, nil
		}
	}

	// TODO: when pos + 1 >= len(root.pointers) holds
	if pos+1 < len(root.pointers) {
		if t.mergeNodes(child, root.keys[pos], root.pointers[pos+1].(*tnode)) {
			root.deleteEntryAt(pos)
			return true, nil
		}
	}

	// now try redistribute entries
	if pos-1 >= 0 {
		t.borrowFromLeft(root.pointers[pos-1].(*tnode), &root.keys[pos-1], child)
		return false, nil
	}

	if pos+1 < len(root.pointers) {
		t.borrowFromRight(child, &root.keys[pos], root.pointers[pos+1].(*tnode))
		return false, nil
	}

	glog.Fatalf("unable to delete key %d from %+v", key, root)
	panic("unreachable")
}

func (tn *tnode) tooFewPointers() bool {
	if tn.isLeaf {
		return len(tn.pointers) < cap(tn.keys)/2+1
	}

	return len(tn.pointers) < (cap(tn.keys)+1)/2
}

// findInsertPos find smallest k in keys greater or equal to key
func (tn *tnode) findInsertPos(key int) int {
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

// delete entry at pos
func (tn *tnode) deleteEntryAt(pos int) {
	// delete entry at from leaf
	glog.V(2).Infof("deleting entry at %d: %+v", pos, tn)
	copy(tn.keys[pos:], tn.keys[pos+1:])
	tn.keys = tn.keys[:len(tn.keys)-1]

	ppos := pos
	if !tn.isLeaf {
		ppos += 1
	}
	copy(tn.pointers[ppos:], tn.pointers[ppos+1:])
	tn.pointers = tn.pointers[:len(tn.pointers)-1]
}

// delete entry with key
func (tn *tnode) deleteEntry(key int) error {
	pos := tn.findInsertPos(key)
	if pos >= len(tn.keys) || tn.keys[pos] != key {
		return ErrKeyNotFound
	}

	tn.deleteEntryAt(pos)
	return nil
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

// insertLeafAt insert entry (key, p) at index
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
