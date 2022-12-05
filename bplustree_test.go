package bplustree

import (
	"reflect"
	"strings"
	"testing"
)

func TestInsertLeaf(t *testing.T) {
	tn := newTNode(true, 5)
	tn.insert(4, 4)
	tn.insert(1, 1)
	tn.insert(2, 2)
	tn.insert(5, 5)
	tn.insert(3, 3)
	t.Logf("%+v", tn.keys)
	t.Logf("%+v", tn.pointers)

	wkeys := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(tn.keys, wkeys) {
		t.Fatalf("expect keys: %+v but got: %+v", wkeys, tn.keys)
	}

	wpointers := []interface{}{1, 2, 3, 4, 5, nil}
	if !reflect.DeepEqual(tn.pointers, wpointers) {
		t.Fatalf("expect pointers: %+v but got: %+v", wpointers, tn.pointers)
	}
}

func TestInsertNonLeaf(t *testing.T) {
	tn := newTNode(false, 5)
	tn.insert(4, 4)
	tn.insert(1, 1)
	tn.insert(2, 2)
	tn.insert(5, 5)
	tn.insert(3, 3)
	t.Logf("%+v", tn.keys)
	t.Logf("%+v", tn.pointers)

	wkeys := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(tn.keys, wkeys) {
		t.Fatalf("expect keys: %+v but got: %+v", wkeys, tn.keys)
	}

	wpointers := []interface{}{nil, 1, 2, 3, 4, 5}
	if !reflect.DeepEqual(tn.pointers, wpointers) {
		t.Fatalf("expect pointers: %+v but got: %+v", wpointers, tn.pointers)
	}
}

func TestSplitLeafRoot(t *testing.T) {
	tr, _ := NewTree(6)
	tr.Insert(4, 4)
	tr.Insert(1, 1)
	tr.Insert(2, 2)
	tr.Insert(5, 5)
	tr.Insert(3, 3)

	t.Logf("keys: %+v", tr.root.keys)
	t.Logf("pointers: %+v", tr.root.pointers)

	tr.Insert(6, 6)
	if len(tr.root.pointers) != 2 {
		t.Fatalf("expect root with 2 pointer but got %d", len(tr.root.pointers))
	}

	if len(tr.root.keys) != 1 {
		t.Fatalf("expect root with 1 key but got %d", len(tr.root.keys))
	}

	if tr.root.keys[0] != 4 {
		t.Fatalf("expect first key in root node to be 4 but got %d", tr.root.keys[0])
	}

	c1 := tr.root.pointers[0].(*tnode)
	c2 := tr.root.pointers[1].(*tnode)
	wkeys := []int{1, 2, 3}
	if !reflect.DeepEqual(wkeys, c1.keys) {
		t.Fatalf("expect keys %+v but got %+v", wkeys, c1.keys)
	}

	wpointers := []interface{}{1, 2, 3, c2}
	if !reflect.DeepEqual(wpointers, c1.pointers) {
		t.Fatalf("expect pointers %+v but got %+v", wpointers, c1.pointers)
	}

	wkeys = []int{4, 5, 6}
	if !reflect.DeepEqual(wkeys, c2.keys) {
		t.Fatalf("expect keys %+v but got %+v", wkeys, c2.keys)
	}

	wpointers = []interface{}{4, 5, 6, nil}
	if !reflect.DeepEqual(wpointers, c2.pointers) {
		t.Fatalf("expect pointers %+v but got %+v", wpointers, c2.pointers)
	}
}

func TestSplitInternalNode(t *testing.T) {
	tr, _ := NewTree(3)
	keys := []int{3, 1, 2, 4, 5, 6, 7}

	for _, key := range keys {
		if err := tr.Insert(key, key); err != nil {
			t.Fatalf("error inserting key %d: %+v", key, err)
		}
		t.Logf("tree after insert key: %d\n%s\n", key, tr.String())
	}

	c2 := tr.root.pointers[1].(*tnode)
	if c2.isLeaf {
		t.Fatalf("expect internal node but got leaf: %+v", c2)
	}

	if len(c2.pointers) != 2 {
		t.Fatalf("expect 2 children but got %d", len(c2.pointers))
	}

	for _, p := range c2.pointers {
		c := p.(*tnode)
		if c.parent != c2 {
			t.Fatalf("expect parent %+v but got %+v", c2, c.parent)
		}
	}
}

func TestInsertNode(t *testing.T) {
	tr, _ := NewTree(3)
	keys := []int{3, 1, 2, 4, 5, 6, 7, 20, 18, 19, 13, 10, 12, 11, 17, 16, 14, 15, 9, 8}

	for _, key := range keys {
		if err := tr.Insert(key, key); err != nil {
			t.Fatalf("error inserting key %d: %+v", key, err)
		}
		t.Logf("tree after insert key: %d\n%s\n", key, tr.String())
	}
}

func newTree(t *testing.T, n int, numKeys int) *BPlusTree {
	tr, _ := NewTree(n)
	for i := 1; i <= numKeys; i++ {
		if err := tr.Insert(i, i); err != nil {
			t.Fatalf("error inserting key: %d: %+v", i, err)
		}

		// t.Logf("tree after insert key: %d\n%s", i, tr.String())
	}

	return tr
}

func TestDeleteKeyNoExist(t *testing.T) {
	numKeys := 3
	tr := newTree(t, 4, numKeys)
	err := tr.Delete(4)
	if err == nil || !strings.Contains(err.Error(), "key not found") {
		t.Fatalf("expect err %+v but got %+v", ErrKeyNotFound, err)
	}
}

func TestDeleteRootLeaf(t *testing.T) {
	numKeys := 3
	tr := newTree(t, 4, numKeys)
	t.Logf("b tree:\n%s\n", tr.String())
	for i := 1; i <= numKeys; i += 1 {
		if err := tr.Delete(i); err != nil {
			t.Fatalf("error deleting key %d: %+v", i, err)
		}

		t.Logf("tree after deleting key %d:\n%s", i, tr.String())
	}
}

func TestDeleteMergeLeftLeaf(t *testing.T) {
	numKeys := 5
	tr := newTree(t, 4, numKeys)
	t.Logf("b tree:\n%s\n", tr.String())

	if err := tr.Delete(4); err != nil {
		t.Fatalf("error deleting key %d: %+v", 4, err)
	}
	t.Logf("b tree after deleting key %d:\n%s\n", 4, tr.String())

	// will merge (1, 2) with (5) after deleting 3
	if err := tr.Delete(3); err != nil {
		t.Fatalf("error deleting key %d: %+v", 3, err)
	}
	t.Logf("b tree after deleting key %d:\n%s\n", 3, tr.String())
}

func TestDeleteMergeRightLeaf(t *testing.T) {
	numKeys := 5
	tr := newTree(t, 4, numKeys)
	t.Logf("b tree:\n%s\n", tr.String())

	if err := tr.Delete(4); err != nil {
		t.Fatalf("error deleting key %d: %+v", 4, err)
	}
	t.Logf("b tree after deleting key %d:\n%s\n", 4, tr.String())

	// will merge (1, 2) with (5) after deleting 3
	if err := tr.Delete(2); err != nil {
		t.Fatalf("error deleting key %d: %+v", 3, err)
	}
	t.Logf("b tree after deleting key %d:\n%s\n", 2, tr.String())
}

func TestDeleteMergeLeftInternalNode(t *testing.T) {
	numKeys := 12
	tr := newTree(t, 4, numKeys)
	t.Logf("b tree:\n%s\n", tr.String())

	for i := 12; i > 9; i-- {
		if err := tr.Delete(i); err != nil {
			t.Fatalf("error deleting key %d: %+v", i, err)
		}
		t.Logf("b tree after deleting key %d:\n%s\n", i, tr.String())
	}
}

func TestDeleteMergeRightInternalNode(t *testing.T) {
	numKeys := 12
	tr := newTree(t, 4, numKeys)
	t.Logf("b tree:\n%s\n", tr.String())

	for i := 6; i > 3; i-- {
		if err := tr.Delete(i); err != nil {
			t.Fatalf("error deleting key %d: %+v", i, err)
		}
		t.Logf("b tree after deleting key %d:\n%s\n", i, tr.String())
	}
}

func TestDeleteLeafBorrwoLeft(t *testing.T) {
	numKeys := 1
	tr := newTree(t, 4, numKeys)
	tr.Insert(8, 8)
	tr.Insert(4, 4)
	tr.Insert(9, 9)
	tr.Insert(5, 5)
	t.Logf("b tree:\n%s\n", tr.String())

	if err := tr.Delete(8); err != nil {
		t.Fatalf("error deleting key %d: %+v", 8, err)
	}
	t.Logf("b tree after deleting 8:\n%s\n", tr.String())
}

func TestDeleteLeafBorrwoRight(t *testing.T) {
	numKeys := 1
	tr := newTree(t, 4, numKeys)
	tr.Insert(7, 7)
	tr.Insert(4, 4)
	tr.Insert(9, 9)
	tr.Insert(8, 8)
	t.Logf("b tree:\n%s\n", tr.String())

	if err := tr.Delete(4); err != nil {
		t.Fatalf("error deleting key %d: %+v", 4, err)
	}
	t.Logf("b tree after deleting 4:\n%s\n", tr.String())
}

func TestDeleteInternalNodeBorrowLeft(t *testing.T) {
	numKeys := 14
	tr := newTree(t, 4, numKeys)
	t.Logf("b tree:\n%s\n", tr.String())

	for i := 1; i < 4; i++ {
		if err := tr.Delete(i); err != nil {
			t.Fatalf("error deleting key %d: %+v", i, err)
		}
		t.Logf("b tree after deleting %d:\n%s\n", i, tr.String())
	}
}

func TestDeleteInternalNodeBorrowRight(t *testing.T) {
	end := 8
	tr, _ := NewTree(4)
	for i := 1; i < end; i++ {
		if err := tr.Insert(i, i); err != nil {
			t.Fatalf("error insertiing key %d: %+v", i, err)
		}

		t.Logf("b tree after inserting key %d:\n%s", i, tr.String())
	}

	for i := 16; i >= 10; i-- {
		if err := tr.Insert(i, i); err != nil {
			t.Fatalf("error insertiing key %d: %+v", i, err)
		}

		t.Logf("b tree after inserting key %d:\n%s", i, tr.String())
	}

	for i := 1; i < 4; i++ {
		if err := tr.Delete(i); err != nil {
			t.Fatalf("error deleting key %d: %+v", i, err)
		}

		t.Logf("b tree after deleting key %d:\n%s", i, tr.String())
	}
}
