package v2

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func newTree(t *testing.T, n int, numKeys int, step int) *BPlusTree {
	tr, _ := NewTree(n)
	for i := 1; i <= numKeys; i++ {
		key := (i-1)*step + 1
		if err := tr.Insert(&Entry{key: key, pointer: key}); err != nil {
			t.Fatalf("error inserting key: %d: %+v", i, err)
		}

		t.Logf("tree after insert key: %d\n%s", i, tr.ToString())
	}

	return tr
}

func TestBTreeNewTree(t *testing.T) {
	if _, err := NewTree(2); err == nil {
		t.Fatalf("expect error but got none")
	}

	if _, err := NewTree(3); err != nil {
		t.Fatalf("expect no error but got: %+v", err)
	}
}

func TestBTreeInsertLeaf(t *testing.T) {
	tr := newTree(t, 4, 0, 0)
	tr.Insert(&Entry{2, 2})
	tr.Insert(&Entry{1, 1})
	tr.Insert(&Entry{3, 3})

	t.Logf("%+v", tr.root.entries)
	wentries := []Entry{
		{1, 1},
		{2, 2},
		{3, 3},
		{0, nil},
	}

	if !reflect.DeepEqual(tr.root.entries, wentries) {
		t.Fatalf("expect keys: %+v but got: %+v", wentries, tr.root.entries)
	}
}

// TODO: SplitInternalRoot
func TestBTreeSplitLeafRoot(t *testing.T) {
	tr, _ := NewTree(6)
	tr.Insert(&Entry{4, 4})
	tr.Insert(&Entry{1, 1})
	tr.Insert(&Entry{2, 2})
	tr.Insert(&Entry{5, 5})
	tr.Insert(&Entry{3, 3})

	t.Logf("tree before split: %+v", tr.root.ToString())

	tr.Insert(&Entry{6, 6})
	if len(tr.root.entries) != 2 {
		t.Fatalf("expect root with 2 pointer but got %d", len(tr.root.entries))
	}

	if tr.root.entries[1].key != 4 {
		t.Fatalf("expect first key in root node to be 4 but got %d", tr.root.entries[1].key)
	}

	c1 := tr.root.entries[0].pointer.(*tNode)
	c2 := tr.root.entries[1].pointer.(*tNode)
	wentries := []Entry{
		{1, 1},
		{2, 2},
		{3, 3},
		{0, c2},
	}
	if !reflect.DeepEqual(wentries, c1.entries) {
		t.Fatalf("expect keys %+v but got %+v", wentries, c1.entries)
	}

	wentries = []Entry{
		{4, 4},
		{5, 5},
		{6, 6},
		{0, nil},
	}
	if !reflect.DeepEqual(wentries, c2.entries) {
		t.Fatalf("expect keys %+v but got %+v", wentries, c2.entries)
	}

	t.Logf("tree after split:\n%s", tr.root.ToString())
}

func TestBTreeSplitInternalNode(t *testing.T) {
	tr, _ := NewTree(3)
	keys := []int{3, 1, 2, 4, 5, 6, 7}

	for _, key := range keys {
		if err := tr.Insert(&Entry{key: key, pointer: key}); err != nil {
			t.Fatalf("error inserting key %d: %+v", key, err)
		}
		t.Logf("tree after insert key: %d\n%s\n", key, tr.ToString())
	}

	c2 := tr.root.entries[1].pointer.(*tNode)
	if c2.isLeaf {
		t.Fatalf("expect internal node but got leaf: %+v", c2)
	}

	if len(c2.entries) != 2 {
		t.Fatalf("expect 2 children but got %d", len(c2.entries))
	}

	for _, e := range c2.entries {
		c := e.pointer.(*tNode)
		if c.parent != c2 {
			t.Fatalf("expect parent %+v but got %+v", c2, c.parent)
		}
	}
}

func TestBTreeInsertNode(t *testing.T) {
	tr, _ := NewTree(3)
	keys := []int{3, 1, 2, 4, 5, 6, 7, 20, 18, 19, 13, 10, 12, 11, 17, 16, 14, 15, 9, 8}

	for _, key := range keys {
		if err := tr.Insert(&Entry{key: key, pointer: key}); err != nil {
			t.Fatalf("error inserting key %d: %+v", key, err)
		}
		t.Logf("tree after insert key: %d\n%s\n", key, tr.ToString())
	}
}

func TestBTreeInsertDuplicate(t *testing.T) {
	tr := newTree(t, 3, 7, 1)
	if err := tr.Insert(&Entry{1, 1}); err != ErrDupKey {
		t.Fatalf("expect err %+v but got %+v", ErrDupKey, err)
	}
}

func TestBTreeDeleteKeyNoExist(t *testing.T) {
	numKeys := 3
	tr := newTree(t, 4, numKeys, 1)
	err := tr.Delete(4)
	if err == nil || !strings.Contains(err.Error(), "key not found") {
		t.Fatalf("expect err %+v but got %+v", ErrKeyNotFound, err)
	}
}

func TestBTreeDeleteRootLeaf(t *testing.T) {
	numKeys := 3
	tr := newTree(t, 4, numKeys, 1)
	for i := 1; i <= numKeys; i += 1 {
		if err := tr.Delete(i); err != nil {
			t.Fatalf("error deleting key %d: %+v", i, err)
		}

		t.Logf("tree after deleting key %d:\n%s", i, tr.ToString())
	}
}

func TestBTreeDeleteMergeLeftLeaf(t *testing.T) {
	numKeys := 5
	tr := newTree(t, 4, numKeys, 1)

	if err := tr.Delete(4); err != nil {
		t.Fatalf("error deleting key %d: %+v", 4, err)
	}
	t.Logf("b tree after deleting key %d:\n%s\n", 4, tr.ToString())

	// will merge (1, 2) with (5) after deleting 3
	if err := tr.Delete(3); err != nil {
		t.Fatalf("error deleting key %d: %+v", 3, err)
	}
	t.Logf("b tree after deleting key %d:\n%s\n", 3, tr.ToString())
}

func TestBTreeDeleteMergeRightLeaf(t *testing.T) {
	numKeys := 5
	tr := newTree(t, 4, numKeys, 1)

	if err := tr.Delete(4); err != nil {
		t.Fatalf("error deleting key %d: %+v", 4, err)
	}
	t.Logf("b tree after deleting key %d:\n%s\n", 4, tr.ToString())

	// will merge (1, 2) with (5) after deleting 3
	if err := tr.Delete(2); err != nil {
		t.Fatalf("error deleting key %d: %+v", 3, err)
	}
	t.Logf("b tree after deleting key %d:\n%s\n", 2, tr.ToString())
}

func TestBTreeDeleteMergeLeftInternal(t *testing.T) {
	numKeys := 12
	tr := newTree(t, 4, numKeys, 1)

	for i := 12; i > 9; i-- {
		if err := tr.Delete(i); err != nil {
			t.Fatalf("error deleting key %d: %+v", i, err)
		}
		t.Logf("b tree after deleting key %d:\n%s\n", i, tr.ToString())
	}
}

func TestBTreeDeleteMergeRightInternal(t *testing.T) {
	numKeys := 12
	tr := newTree(t, 4, numKeys, 1)

	for i := 6; i > 3; i-- {
		if err := tr.Delete(i); err != nil {
			t.Fatalf("error deleting key %d: %+v", i, err)
		}
		t.Logf("b tree after deleting key %d:\n%s\n", i, tr.ToString())
	}
}

func TestBTreeDeleteBorrowLeftLeaf(t *testing.T) {
	numKeys := 1
	tr := newTree(t, 4, numKeys, 1)
	tr.Insert(&Entry{key: 8, pointer: 8})
	tr.Insert(&Entry{key: 4, pointer: 4})
	tr.Insert(&Entry{key: 9, pointer: 9})
	tr.Insert(&Entry{key: 5, pointer: 5})
	t.Logf("b tree:\n%s", tr.ToString())

	if err := tr.Delete(8); err != nil {
		t.Fatalf("error deleting key %d: %+v", 8, err)
	}
	t.Logf("b tree after deleting 8:\n%s\n", tr.ToString())
}

func TestBTreeDeleteBorrwoRightLeaf(t *testing.T) {
	numKeys := 1
	tr := newTree(t, 4, numKeys, 1)
	tr.Insert(&Entry{7, 7})
	tr.Insert(&Entry{4, 4})
	tr.Insert(&Entry{9, 9})
	tr.Insert(&Entry{8, 8})
	t.Logf("b tree:\n%s\n", tr.ToString())

	if err := tr.Delete(4); err != nil {
		t.Fatalf("error deleting key %d: %+v", 4, err)
	}
	t.Logf("b tree after deleting 4:\n%s\n", tr.ToString())
}

func TestBTreeDeleteBorrowRightInternal(t *testing.T) {
	numKeys := 14
	tr := newTree(t, 4, numKeys, 1)

	for i := 1; i < 4; i++ {
		if err := tr.Delete(i); err != nil {
			t.Fatalf("error deleting key %d: %+v", i, err)
		}
		t.Logf("b tree after deleting %d:\n%s\n", i, tr.ToString())
	}

	if tr.root.entries[1].key != 9 {
		t.Fatalf("exptect root first key 9 but got %d", tr.root.entries[1].key)
	}
}

func TestBTreeDeleteBorrowLeftInternal(t *testing.T) {
	numKeys := 10
	tr := newTree(t, 4, numKeys, 3)
	tr.Insert(&Entry{11, 11})
	tr.Insert(&Entry{12, 12})
	t.Logf("b tree:\n%s\n", tr.ToString())

	if err := tr.Delete(19); err != nil {
		t.Fatalf("error deleting key %d: %+v", 19, err)
	}
	t.Logf("b tree after deleting key 19:\n%s\n", tr.ToString())
	c2 := tr.root.entries[1].pointer.(*tNode)
	if c2.entries[1].key != 19 {
		t.Fatalf("expect first key 19 after borrowing from left but got %d", c2.entries[1].key)
	}
}

func doCheckBPlusTreeInvariant(parent *tNode, tn *tNode) error {
	// 1. check parent pointer
	if parent != tn.parent {
		return fmt.Errorf("expect parent of %v to be %+v but got: %+v", tn.ChildrenStr(), parent.ChildrenStr(), tn.parent.ChildrenStr())
	}

	// 2. check entry size
	if len(tn.entries) >= cap(tn.entries) {
		return fmt.Errorf("max entry size %d but got %d entires: %+v", cap(tn.entries)-1, len(tn.entries), tn.entries)
	}

	if tn.tooFewPointers() {
		return fmt.Errorf("max entry size %d, too few entrys: %d", cap(tn.entries)-1, tn.entries)
	}

	// 3. check value invariant
	// TODO

	if !tn.isLeaf {
		for _, e := range tn.entries {
			child := e.pointer.(*tNode)
			if err := doCheckBPlusTreeInvariant(tn, child); err != nil {
				return err
			}
		}
	}

	return nil
}

func checkBPlusTreeInvariant(tr *BPlusTree) error {
	return doCheckBPlusTreeInvariant(nil, tr.root)
}

func TestBTreeUpdateParent(t *testing.T) {
	cases := []struct {
		maxEntries int
		step       int
		numKeys    int
		extraKeys  []int
		deleteKeys []int
	}{
		// internal merge left
		{
			maxEntries: 4,
			step:       1,
			numKeys:    10,
			deleteKeys: []int{9},
		},
		// internal merge right
		{
			maxEntries: 4,
			step:       1,
			numKeys:    10,
			deleteKeys: []int{1, 2, 3},
		},
		// internal borrow from left
		{
			maxEntries: 4,
			step:       2,
			numKeys:    10,
			extraKeys:  []int{2, 4},
			deleteKeys: []int{15},
		},
		// internal borrow from right
		{
			maxEntries: 4,
			step:       2,
			numKeys:    10,
			extraKeys:  []int{14, 16, 18, 20},
			deleteKeys: []int{1, 3, 5},
		},
		// update root parent after delete
		{
			maxEntries: 4,
			step:       1,
			numKeys:    10,
			deleteKeys: []int{7},
		},
		// update root parent after insert
		{
			maxEntries: 4,
			step:       1,
			numKeys:    4,
		},
	}
	for i, tc := range cases {
		t.Logf("test case %d: %+v", i, tc)
		tr := newTree(t, tc.maxEntries, tc.numKeys, tc.step)

		for _, key := range tc.extraKeys {
			if err := tr.Insert(&Entry{key, key}); err != nil {
				t.Fatalf("error inserting key %d: %+v", key, err)
			}
			t.Logf("b tree after inserting key %d:\n%s", key, tr.ToString())
		}

		// tr := newTree(t, tc.maxEntries, tc.numKeys, 1)
		for _, key := range tc.deleteKeys {
			if err := tr.Delete(key); err != nil {
				t.Fatalf("case %d error deleting key %d: %+v", i, key, err)
			}
			t.Logf("b tree after deleting key %d:\n%s", key, tr.ToString())
		}

		if err := checkBPlusTreeInvariant(tr); err != nil {
			t.Fatalf("b tree invariant check failed: %+v", err)
		}
	}
}
