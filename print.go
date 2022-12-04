package bplustree

import (
	"bytes"
	"strconv"
	"strings"
)

type pnode struct {
	size     int
	val      string
	children []*pnode
}

func traverseTree(root *tnode) *pnode {
	keys := make([]string, len(root.keys))
	for i, k := range root.keys {
		keys[i] = strconv.Itoa(k)
	}
	val := "(" + strings.Join(keys, ",") + ")"
	proot := &pnode{val: val}

	if !root.isLeaf {
		for _, c := range root.pointers {
			child := c.(*tnode)
			cp := traverseTree(child)
			proot.size += cp.size
			proot.children = append(proot.children, cp)
		}
	} else {
		// proot.size = len(proot.val) + 2
		proot.size = len(proot.val) + 4
	}

	return proot
}

func printTree(pn *pnode, buf *bytes.Buffer) {
	que := []*pnode{pn}

	for len(que) > 0 {
		newQue := []*pnode{}

		for _, n := range que {
			padding := (n.size - len(n.val)) / 2

			// center align
			// pad left
			for i := 0; i < padding; i++ {
				buf.WriteByte(' ')
			}

			buf.WriteString(n.val)

			// pad right
			for i := 0; i < padding; i++ {
				buf.WriteByte(' ')
			}

			if len(n.children) > 0 {
				newQue = append(newQue, n.children...)
			}
		}

		buf.WriteByte('\n')
		que = newQue
	}
}

func (tr *BPlusTree) String() string {
	buf := bytes.NewBuffer(nil)
	root := traverseTree(tr.root)
	printTree(root, buf)

	return buf.String()
}
