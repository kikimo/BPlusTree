package v2

import (
	"bytes"
	"strconv"
	"strings"
)

func (tr *BPlusTree) String() string {
	if tr.root != nil && len(tr.root.entries) > 1 {
		return tr.root.String()
	}

	return "()"
}

func (tn *tNode) String() string {
	buf := bytes.NewBuffer(nil)
	root := traverseTree(tn)
	printTree(root, buf)

	return buf.String()
}

type pnode struct {
	size     int
	val      string
	children []*pnode
}

func traverseTree(root *tNode) *pnode {
	keys := make([]string, 0, len(root.entries)-1)
	for i, k := range root.entries {
		if (root.isLeaf && i == len(root.entries)-1) || (!root.isLeaf && i == 0) {
			continue
		}

		keys = append(keys, strconv.Itoa(k.key))
	}
	val := "(" + strings.Join(keys, ",") + ")"
	proot := &pnode{val: val}

	if !root.isLeaf {
		for _, c := range root.entries {
			child := c.pointer.(*tNode)
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
