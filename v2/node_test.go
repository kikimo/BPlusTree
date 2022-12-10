package v2

import (
	"reflect"
	"testing"
)

func TestLeafInsert(t *testing.T) {
	leaf := newTNode(true, 4)
	keys := []int{5, 1, 4}
	for _, key := range keys {
		e := &Entry{
			key:     key,
			pointer: key,
		}
		if err := leaf.insertLeaf(e); err != nil {
			t.Fatalf("error inserting key %d to leaf: %+v", key, err)
		}
	}

	if err := leaf.insertLeaf(&Entry{key: 4, pointer: 4}); err != ErrDupKey {
		t.Fatalf("expect error %+v but got %+v", ErrDupKey, err)
	}

	if err := leaf.insertLeaf(&Entry{key: 2, pointer: 2}); err != nil {
		t.Fatalf("error insert key %d: %+v", 2, err)
	}

	wentries := []Entry{
		{1, 1},
		{2, 2},
		{4, 4},
		{5, 5},
		{0, nil},
	}

	if !reflect.DeepEqual(wentries, leaf.entries) {
		t.Fatalf("want %+v, but got %+v", wentries, leaf.entries)
	}

	// leaf.insertLeaf(&Entry{key: 6, pointer: 6})
	t.Logf("leaft after insert: %s", leaf.ToString())
}

func TestSplitLeafNodeEven(t *testing.T) {
	leaf := newTNode(true, 4)
	for i := 4; i >= 1; i-- {
		leaf.insertLeaf(&Entry{key: i, pointer: i})
	}
	t.Logf("leaf before split: %+v", leaf.ToString())

	ne := leaf.splitLeafNode()
	right := ne.pointer.(*tNode)
	t.Logf("after split, left: %s, right: %s", leaf.ToString(), right.ToString())

	lsz, rsz := len(leaf.entries), len(right.entries)
	if lsz != 3 || rsz != 3 {
		t.Fatalf("expect both entry sizes to be 3 after split but got left: %d and right: %d", lsz, rsz)
	}

	if leaf.entries[lsz-1].pointer.(*tNode) != right {
		t.Fatalf("expect sibling connected to %+v but got %+v", right, leaf.entries[lsz-1].pointer.(*tNode))
	}
}

func TestSplitLeafNodeOdd(t *testing.T) {
	leaf := newTNode(true, 5)
	for i := 5; i >= 1; i-- {
		leaf.insertLeaf(&Entry{key: i, pointer: i})
	}
	t.Logf("leaf before split:\n%+v", leaf.ToString())

	ne := leaf.splitLeafNode()
	right := ne.pointer.(*tNode)
	t.Logf("after split:\nleft:\n%s\nright:\n%s", leaf.ToString(), right.ToString())

	lsz, rsz := len(leaf.entries), len(right.entries)
	if lsz != 4 || rsz != 3 {
		t.Fatalf("expect both entry sizes to be 3 after split but got left: %d and right: %d", lsz, rsz)
	}

	if leaf.entries[lsz-1].pointer.(*tNode) != right {
		t.Fatalf("expect sibling connected to %+v but got %+v", right, leaf.entries[lsz-1].pointer.(*tNode))
	}
}

func TestSplitInternalNodeEven(t *testing.T) {
	inode := newTNode(false, 4)
	inode.entries = inode.entries[:1]
	inode.entries[0] = Entry{pointer: &tNode{parent: inode}}
	for i := 4; i >= 1; i-- {
		pos := inode.findInternalInsertPos(i)
		t.Logf("insert key %d at %d", i, pos)
		inode.insertAt(pos, &Entry{key: i, pointer: &tNode{parent: inode}})
	}

	ne := inode.splitInternalNode()
	t.Logf("inode entries after split: %+v", inode.entries)
	if ne.key != 3 {
		t.Fatalf("expect new entry key 3 but got %d", ne.key)
	}

	if len(inode.entries) != 3 {
		t.Fatalf("expect inode entry of size 3 but got %d", len(inode.entries))
	}

	newChild := ne.pointer.(*tNode)
	t.Logf("new child entries: %+v", newChild.entries)
	if len(newChild.entries) != 2 {
		t.Fatalf("expect new child entry of size 2 but got %d", len(inode.entries))
	}
}

func TestSplitInternalNodeOdd(t *testing.T) {
	inode := newTNode(false, 5)
	inode.entries = inode.entries[:1]
	inode.entries[0] = Entry{pointer: &tNode{parent: inode}}
	for i := 5; i >= 1; i-- {
		pos := inode.findInternalInsertPos(i)
		t.Logf("insert key %d at %d", i, pos)
		inode.insertAt(pos, &Entry{key: i, pointer: &tNode{parent: inode}})
	}

	ne := inode.splitInternalNode()
	t.Logf("inode entries after split: %+v", inode.entries)
	if ne.key != 3 {
		t.Fatalf("expect new entry key 3 but got %d", ne.key)
	}

	if len(inode.entries) != 3 {
		t.Fatalf("expect inode entry of size 3 but got %d", len(inode.entries))
	}

	newChild := ne.pointer.(*tNode)
	t.Logf("new child entries: %+v", newChild.entries)
	if len(newChild.entries) != 3 {
		t.Fatalf("expect new child entry of size 3 but got %d", len(inode.entries))
	}
}

func TestMergeLeafNodes(t *testing.T) {
	left := newTNode(true, 4)
	left.insertLeaf(&Entry{key: 1, pointer: 1})
	left.insertLeaf(&Entry{key: 2, pointer: 2})

	right := newTNode(true, 4)
	right.insertLeaf(&Entry{key: 3, pointer: 3})
	if !left.mergeLeaves(right) {
		t.Fatalf("expect leaves merged but not")
	}

	wentries := []Entry{
		{1, 1},
		{2, 2},
		{3, 3},
		{0, nil},
	}

	if !reflect.DeepEqual(wentries, left.entries) {
		t.Fatalf("expect entries %+v but got %+v", wentries, left.entries)
	}

	right = newTNode(true, 4)
	right.insertLeaf(&Entry{key: 4, pointer: 4})
	if left.mergeLeaves(right) {
		t.Fatalf("should no be able to merge leaves(too many pointers)")
	}
}

func TestMergeInternalNode(t *testing.T) {
	left := newTNode(false, 4)
	left.insertAt(0, &Entry{key: 0, pointer: nil})
	left.insertAt(1, &Entry{key: 1, pointer: 1})
	left.insertAt(2, &Entry{key: 2, pointer: 2})

	right := newTNode(false, 4)
	rightChild := newTNode(true, 4)
	right.insertAt(0, &Entry{key: 0, pointer: rightChild})
	if !left.mergeInternalNodes(3, right) {
		t.Fatalf("expect internal node merged but not")
	}

	wentries := []Entry{
		{0, nil},
		{1, 1},
		{2, 2},
		{3, rightChild},
	}

	if !reflect.DeepEqual(wentries, left.entries) {
		t.Fatalf("expect entries %+v but got %+v", wentries, left.entries)
	}

	right = newTNode(false, 4)
	right.insertAt(0, &Entry{key: 0, pointer: 4})
	if left.mergeInternalNodes(4, right) {
		t.Fatalf("should no be able to merge internal nodes(too many pointers)")
	}
}

func TestDeleteInternalNode(t *testing.T) {
	root := newTNode(false, 5)
	wentries := []Entry{
		{0, nil},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
	}
	root.entries = wentries

	if err := root.deleteEntry(2); err != nil {
		t.Fatal(err)
	}
	wentries = []Entry{
		{0, nil},
		{1, 1},
		{3, 3},
		{4, 4},
	}
	if !reflect.DeepEqual(wentries, root.entries) {
		t.Fatalf("expect entries %+v but got %+v", wentries, root.entries)
	}

	if err := root.deleteEntry(1); err != nil {
		t.Fatal(err)
	}
	wentries = []Entry{
		{0, nil},
		{3, 3},
		{4, 4},
	}
	if !reflect.DeepEqual(wentries, root.entries) {
		t.Fatalf("expect entries %+v but got %+v", wentries, root.entries)
	}

	if err := root.deleteEntry(4); err != nil {
		t.Fatal(err)
	}
	wentries = []Entry{
		{0, nil},
		{3, 3},
	}
	if !reflect.DeepEqual(wentries, root.entries) {
		t.Fatalf("expect entries %+v but got %+v", wentries, root.entries)
	}
}

func TestDeleteLeafNode(t *testing.T) {
	tr := newTree(t, 5, 4, 1)
	root := tr.root
	wentries := []Entry{
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
		{0, nil},
	}
	if !reflect.DeepEqual(wentries, root.entries) {
		t.Fatalf("expect entries %+v but got %+v", wentries, root.entries)
	}

	if err := root.deleteEntry(2); err != nil {
		t.Fatal(err)
	}
	wentries = []Entry{
		{1, 1},
		{3, 3},
		{4, 4},
		{0, nil},
	}
	if !reflect.DeepEqual(wentries, root.entries) {
		t.Fatalf("expect entries %+v but got %+v", wentries, root.entries)
	}

	if err := root.deleteEntry(1); err != nil {
		t.Fatal(err)
	}
	wentries = []Entry{
		{3, 3},
		{4, 4},
		{0, nil},
	}
	if !reflect.DeepEqual(wentries, root.entries) {
		t.Fatalf("expect entries %+v but got %+v", wentries, root.entries)
	}

	if err := root.deleteEntry(4); err != nil {
		t.Fatal(err)
	}
	wentries = []Entry{
		{3, 3},
		{0, nil},
	}
	if !reflect.DeepEqual(wentries, root.entries) {
		t.Fatalf("expect entries %+v but got %+v", wentries, root.entries)
	}
}

func TestTooFewPointers(t *testing.T) {
	cases := []struct {
		isLeaf         bool
		maxSize        int
		keys           []int
		tooFewPointers bool
	}{
		{
			isLeaf:         true,
			maxSize:        4,
			keys:           []int{1},
			tooFewPointers: true,
		},
		{
			isLeaf:         true,
			maxSize:        4,
			keys:           []int{2, 1},
			tooFewPointers: false,
		},
		{
			isLeaf:         true,
			maxSize:        5,
			keys:           []int{2},
			tooFewPointers: true,
		},
		{
			isLeaf:         true,
			maxSize:        5,
			keys:           []int{2, 1},
			tooFewPointers: false,
		},
		{
			isLeaf:         false,
			maxSize:        4,
			keys:           []int{1},
			tooFewPointers: false,
		},
		{
			isLeaf:         false,
			maxSize:        4,
			keys:           []int{},
			tooFewPointers: true,
		},
		{
			isLeaf:         false,
			maxSize:        5,
			keys:           []int{2},
			tooFewPointers: true,
		},
		{
			isLeaf:         false,
			maxSize:        5,
			keys:           []int{2, 1},
			tooFewPointers: false,
		},
	}

	for i, tc := range cases {
		node := newTNode(tc.isLeaf, tc.maxSize)
		if !tc.isLeaf {
			node.entries = node.entries[:1]
			node.entries[0] = Entry{key: 0}
		}

		for _, k := range tc.keys {
			if node.isLeaf {
				node.insertLeaf(&Entry{key: k})
			} else {
				pos := node.findInternalInsertPos(k)
				node.insertAt(pos, &Entry{key: k})
			}
		}

		if node.tooFewPointers() != tc.tooFewPointers {
			t.Fatalf("test case %d, expect too few pointers: %t but got: %t, entry: %+v", i, tc.tooFewPointers, node.tooFewPointers(), node.entries)
		}
	}
}
