package v2

import (
	"fmt"

	"github.com/golang/glog"
)

type BPlusTree struct {
	maxSize int // max pointer in a node
	root    *tNode
}

func (tr *BPlusTree) Find(key int) (interface{}, error) {
	tn := tr.root

	for !tn.isLeaf {
		pos := tn.findInternalInsertPos(key)
		if pos >= len(tn.entries) || tn.entries[pos].key > key {
			pos -= 1
		}

		tn = tn.entries[pos].pointer.(*tNode)
	}

	pos := tn.findLeafInsertPos(key)
	if pos >= len(tn.entries) || tn.entries[pos].key != key {
		return nil, ErrKeyNotFound
	}

	return tn.entries[pos].pointer, nil
}

func (tr *BPlusTree) Insert(e *Entry) error {
	ne, err := tr.doInsert(tr.root, e)
	if err != nil {
		return err
	}

	if ne == nil {
		return nil
	}

	newRoot := newTNode(false, tr.maxSize)
	newRoot.entries = newRoot.entries[:2]
	newRoot.entries[0] = Entry{pointer: tr.root}
	newRoot.entries[1] = *ne
	tr.root.parent = newRoot
	ne.pointer.(*tNode).parent = newRoot
	tr.root = newRoot

	return nil
}

// doInsert insert Entry e into root, a new entry is returned and
// insert to parent node if root is splited
func (tr *BPlusTree) doInsert(root *tNode, e *Entry) (*Entry, error) {
	// insert leaf node
	if root.isLeaf {
		glog.V(2).Infof("entries size: %d, cap: %d, entries: %+v", len(root.entries), cap(root.entries), root.ChildrenStr())
		if err := root.insertLeaf(e); err != nil {
			return nil, err
		}

		glog.V(2).Infof("entries size: %d, cap: %d, entries: %+v", len(root.entries), cap(root.entries), root.ChildrenStr())
		if len(root.entries) < cap(root.entries) {
			return nil, nil
		}

		ne := root.splitLeafNode()
		glog.V(2).Infof("entries size: %d, cap: %d, entries: %+v", len(root.entries), cap(root.entries), root.ChildrenStr())
		glog.V(2).Infof("ne: %+v", *ne)
		return ne, nil
	}

	// insert internal node
	pos := root.findInternalInsertPos(e.key)
	glog.V(2).Infof("internal insert pos: %d", pos)
	if pos >= len(root.entries) || root.entries[pos].key > e.key {
		pos -= 1
	}

	// nce: new child entry
	nce, err := tr.doInsert(root.entries[pos].pointer.(*tNode), e)
	if err != nil {
		return nil, err
	}

	if nce == nil {
		return nil, nil
	}

	// invariant check
	if len(root.entries) >= cap(root.entries) {
		glog.Errorf("cap entries: %d, size entries: %d, entries: %+v", cap(root.entries), len(root.entries), root.ChildrenStr())
		glog.Fatalf("illegal node entry size:\n %+v", root)
	}

	// insert newNode after pos
	root.insertAt(pos+1, nce)
	if len(root.entries) < cap(root.entries) {
		return nil, nil
	}

	ne := root.splitInternalNode()
	return ne, nil
}

func (t *BPlusTree) Delete(key int) error {
	deleted, err := t.deleteEntry(t.root, key)
	if err != nil {
		return fmt.Errorf("error deleting key %d: %+v", key, err)
	}

	if !deleted {
		return nil
	}

	if len(t.root.entries) == 1 {
		if !t.root.isLeaf {
			t.root = t.root.entries[0].pointer.(*tNode)
			t.root.parent = nil
		}
	}

	return nil
}

func (t *BPlusTree) deleteEntry(root *tNode, key int) (bool, error) {
	if root.isLeaf {
		return true, root.deleteEntry(key)
	}

	// pos points to index into which key will be inserted
	pos := root.findInternalInsertPos(key)
	if pos >= len(root.entries) || root.entries[pos].key > key {
		pos -= 1
	}

	de := root.entries[pos]
	child := de.pointer.(*tNode)
	deleted, err := t.deleteEntry(child, key)
	if err != nil || !deleted {
		return false, err
	}

	glog.Infof("children after deletion: %+v", child.ChildrenStr())
	if !child.tooFewPointers() {
		return false, nil
	}

	// too few pointers, try merge entries
	if pos-1 >= 0 {
		left := root.entries[pos-1].pointer.(*tNode)
		if left.mergeNodes(de.key, child) {
			glog.Infof("deleting entry at %d from %+v", pos, root.ChildrenStr())
			root.deleteEntryAt(pos)
			return true, nil
		}
	}

	if pos+1 < len(root.entries) {
		right := root.entries[pos+1].pointer.(*tNode)
		if child.mergeNodes(root.entries[pos+1].key, right) {
			glog.Infof("deleting entry at %d from %+v", pos+1, root.ChildrenStr())
			root.deleteEntryAt(pos + 1)
			return true, nil
		}
	}

	// now try redistribute entries
	if pos-1 >= 0 {
		borrowFromLeft(root.entries[pos-1].pointer.(*tNode), &root.entries[pos].key, child)
		return false, nil
	}

	if pos+1 < len(root.entries) {
		borrowFromRight(child, &root.entries[pos+1].key, root.entries[pos+1].pointer.(*tNode))
		return false, nil
	}

	glog.Fatalf("unable to delete key %d from %+v", key, root.entries)
	panic("unreachable")
}

func NewTree(maxSize int) (*BPlusTree, error) {
	if maxSize < 3 {
		return nil, fmt.Errorf("BPlusTree maxSize should be greater than 3: %d", maxSize)
	}

	tr := &BPlusTree{
		maxSize: maxSize,
		root:    newTNode(true, maxSize),
	}

	return tr, nil
}
