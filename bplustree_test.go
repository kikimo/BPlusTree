package bplustree

import (
	"reflect"
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
