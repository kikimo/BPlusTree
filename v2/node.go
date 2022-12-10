package v2

import (
	"fmt"

	"github.com/golang/glog"
)

var ErrDupKey error = fmt.Errorf("duplicate key")
var ErrKeyNotFound error = fmt.Errorf("key not found")

type tNode struct {
	isLeaf bool
	// parent points to parent pointer when should parent pointer
	// being updated:
	// 	1. update children's parent pointers when an entry is being
	//     merged or splited
	// 	2. update parent pointer when an entry is being inserted
	parent  *tNode
	entries []Entry
}

type Entry struct {
	key     int
	pointer interface{}
}

func newTNode(isLeaf bool, maxSize int) *tNode {
	n := &tNode{
		isLeaf:  isLeaf,
		entries: make([]Entry, 0, maxSize+1),
	}

	if isLeaf {
		n.entries = n.entries[:1]
	}

	return n
}

// findInsertPos find smallest index such that tn.entries[index].key >= key
func (tn *tNode) findInsertPos(key int, s, e int) int {
	for s < e {
		m := (s + e) / 2
		if tn.entries[m].key >= key {
			e = m
		} else { // tn.entries[m].key < key
			s = m + 1
		}
	}

	return s
}

func (tn *tNode) findLeafInsertPos(key int) int {
	return tn.findInsertPos(key, 0, len(tn.entries)-1)
}

func (tn *tNode) findInternalInsertPos(key int) int {
	return tn.findInsertPos(key, 1, len(tn.entries))
}

func (tn *tNode) insertAt(pos int, e *Entry) {
	// expand tn.entries by one
	sz := len(tn.entries)
	tn.entries = tn.entries[:sz+1]
	copy(tn.entries[pos+1:], tn.entries[pos:sz])
	tn.entries[pos] = *e
}

func (tn *tNode) insertLeaf(e *Entry) error {
	sz := len(tn.entries)

	// check invariant
	if sz+1 > cap(tn.entries) {
		glog.Fatalf("leaf entry overflow(maxsize: %d) inserting new entry: %+v", cap(tn.entries), e)
	}

	pos := tn.findLeafInsertPos(e.key)
	if pos < sz-1 && tn.entries[pos].key == e.key {
		return ErrDupKey
	}

	tn.insertAt(pos, e)
	return nil
}

// split nodes
func (tn *tNode) splitInternalNode() *Entry {
	// 4 -> 2
	// 5 -> 2
	sz := len(tn.entries)
	pos := (sz + 1) / 2
	newN := newTNode(false, sz-1)
	newN.parent = tn.parent
	// glog.Infof("pos: %d, e: %+v, entries: %+v", pos, tn.entries[pos], tn.entries)

	// split pointers
	newN.entries = newN.entries[:len(tn.entries[pos:])]
	copy(newN.entries, tn.entries[pos:])
	tn.entries = tn.entries[:pos]

	// adjust parent for splited children
	for _, p := range newN.entries {
		child := p.pointer.(*tNode)
		child.parent = newN
	}

	// insert newEntry into parent
	ne := &Entry{key: newN.entries[0].key, pointer: newN}
	newN.entries[0].key = 0
	return ne
}

func (tn *tNode) splitLeafNode() *Entry {
	glog.V(2).Infof("spliting leaf node: %s", tn.ChildrenStr())
	sz := len(tn.entries)
	// 4 -> 2
	// 5 -> 2
	pos := sz / 2

	newN := newTNode(true, sz-1)
	newN.parent = tn.parent
	newN.entries = newN.entries[:len(tn.entries[pos:])]
	copy(newN.entries, tn.entries[pos:])

	// leave one extra space to connect to sibling
	tn.entries = tn.entries[:pos+1]
	// connect to sibling
	tn.entries[pos] = Entry{pointer: newN}

	glog.V(2).Infof("new entry after split leaf: %s", newN.ChildrenStr())
	return &Entry{key: newN.entries[0].key, pointer: newN}
}

// merge nodes
func (tn *tNode) mergeNodes(key int, right *tNode) bool {
	if tn.isLeaf && right.isLeaf {
		return tn.mergeLeaves(right)
	} else if !tn.isLeaf && !right.isLeaf {
		return tn.mergeInternalNodes(key, right)
	}

	glog.Fatalf("merge leaf to internal node: %s, %s", tn.ChildrenStr(), right.ChildrenStr())
	panic("unreachable")
}

func (tn *tNode) mergeLeaves(right *tNode) bool {
	sz := len(tn.entries) + len(right.entries) - 1

	// unable to merge
	if sz > cap(tn.entries)-1 {
		return false
	}

	glog.V(2).Infof("merge internal node, left: %s, right: %s", tn.ChildrenStr(), right.ChildrenStr())
	start := len(tn.entries) - 1
	tn.entries = tn.entries[:sz]
	for i := range right.entries {
		tn.entries[start+i] = right.entries[i]
	}

	return true
}

// mergeInternalNodes mrege children of right into tn
func (tn *tNode) mergeInternalNodes(key int, right *tNode) bool {
	sz := len(tn.entries) + len(right.entries)

	// unable to merge
	if sz > cap(tn.entries)-1 {
		return false
	}

	glog.V(2).Infof("merge %s into %s", right.ChildrenStr(), tn.ChildrenStr())
	// update parent of right children
	for _, e := range right.entries {
		c := e.pointer.(*tNode)
		c.parent = tn
	}

	right.entries[0].key = key
	start := len(tn.entries)
	tn.entries = tn.entries[:sz]
	for i := range right.entries {
		tn.entries[start+i] = right.entries[i]
	}

	return true
}

// delete entry with key
func (tn *tNode) deleteEntry(key int) error {
	var pos int
	if !tn.isLeaf {
		pos = tn.findInternalInsertPos(key)
	} else {
		pos = tn.findLeafInsertPos(key)
	}

	if pos >= len(tn.entries) || tn.entries[pos].key != key {
		return ErrKeyNotFound
	}

	tn.deleteEntryAt(pos)
	return nil
}

// delete entry at pos
func (tn *tNode) deleteEntryAt(pos int) {
	// delete entry at from leaf
	glog.V(2).Infof("deleting entry at %d: %+v", pos, tn.entries)
	copy(tn.entries[pos:], tn.entries[pos+1:])
	tn.entries = tn.entries[:len(tn.entries)-1]
}

func (tn *tNode) tooFewPointers() bool {
	if tn.isLeaf {
		return len(tn.entries) < (cap(tn.entries)+1)/2
	}

	return len(tn.entries) < cap(tn.entries)/2
}

func borrowFromLeft(left *tNode, key *int, right *tNode) {
	if left.isLeaf && right.isLeaf {
		leafBorrowFromLeft(left, key, right)
		return
	}

	if !left.isLeaf && !right.isLeaf {
		internalBorrowFromLeft(left, key, right)
		return
	}

	glog.Fatalf("leaf cannot borrow from internal node(vice vesa), left: %s, right: %s", left.ChildrenStr(), right.ChildrenStr())
	panic("unreachable")
}

func borrowFromRight(left *tNode, key *int, right *tNode) {
	if left.isLeaf && right.isLeaf {
		leafBorrowFromRight(left, key, right)
		return
	}

	if !left.isLeaf && !right.isLeaf {
		internalBorrowFromRight(left, key, right)
		return
	}

	glog.Fatalf("leaf cannot borrow from internal node(vice vesa), left: %s, right: %s", left.ChildrenStr(), right.ChildrenStr())
	panic("unreachable")
}

func leafBorrowFromLeft(left *tNode, key *int, right *tNode) {
	sz := len(left.entries)
	e := left.entries[sz-2]

	// shrink left by one
	left.entries[sz-2] = left.entries[sz-1]
	left.entries = left.entries[:sz-1]

	*key = e.key

	// prepend entry e to right
	// expand right first
	sz = len(right.entries)
	right.entries = right.entries[:sz+1]
	copy(right.entries[1:], right.entries[:sz-1])
	right.entries[0] = e
}

func leafBorrowFromRight(left *tNode, key *int, right *tNode) {
	sz := len(right.entries)
	e := right.entries[0]

	// shrink right by one
	copy(right.entries[:sz-1], right.entries[1:])
	right.entries = right.entries[:sz-1]

	*key = e.key

	// append entry (k, p) to left
	// expand left first
	sz = len(left.entries)
	left.entries = left.entries[:sz+1]
	left.entries[sz] = left.entries[sz-1]
	left.entries[sz-1] = e
}

func internalBorrowFromLeft(left *tNode, key *int, right *tNode) {
	glog.V(2).Infof("borrow one key from left %s to right %s", left.ChildrenStr(), right.ChildrenStr())
	sz := len(left.entries)
	e := left.entries[sz-1]

	// swap key and e.key
	*key, right.entries[0].key = e.key, *key
	e.key = 0

	// shrink left by one
	left.entries = left.entries[:sz-1]

	// // update children
	e.pointer.(*tNode).parent = right

	// prepend entry (k, p) to right
	// expand right first
	sz = len(right.entries)
	right.entries = right.entries[:sz+1]
	copy(right.entries[1:], right.entries[:sz])
	right.entries[0] = e
}

func internalBorrowFromRight(left *tNode, key *int, right *tNode) {
	glog.V(2).Infof("borrow one key from right %s to left %s", right.ChildrenStr(), left.ChildrenStr())
	sz := len(right.entries)
	e := right.entries[0]
	e.key = right.entries[1].key

	// swap key and e.key
	*key, e.key = e.key, *key

	// // update children
	e.pointer.(*tNode).parent = left

	// shrink right by one
	copy(right.entries[:sz-1], right.entries[1:])
	right.entries = right.entries[:sz-1]
	right.entries[0].key = 0

	// append entry (k, p) to left
	sz = len(left.entries)
	left.entries = left.entries[:sz+1]
	left.entries[sz] = e
}
