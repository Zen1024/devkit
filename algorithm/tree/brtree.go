package tree

import (
	"bytes"
	"fmt"
)

const (
	colorBlack = iota
	colorRed
)

type BrTree struct {
	root   *brNode
	height int
}

type brNode struct {
	color  int
	parent *brNode
	left   *brNode
	right  *brNode
	elem   NodeVal
}

func colorStr(s string, color int) string {
	if color == colorRed {
		return fmt.Sprintf("\x1b[31m%s\x1b[0m", s)
	}
	return fmt.Sprintf("\x1b[34m%s\x1b[0m", s)
}

func (n *brNode) leftRotate() {
	if n == nil || n.right == nil {
		return
	}
	nright := n.right
	if n.parent != nil {
		if n.parent.left == n {
			n.parent.left = nright
		} else {
			n.parent.right = nright
		}
	}
	nright.parent = n.parent
	n.parent = nright
	if nright.left != nil {
		nright.left.parent = n
	}
	n.right = nright.left
	nright.left = n
}

func (n *brNode) rightRotate() {
	if n == nil || n.left == nil {
		return
	}
	nleft := n.left
	if n.parent != nil {
		if n.parent.left == n {
			n.parent.left = nleft
		} else {
			n.parent.right = nleft
		}
	}
	nleft.parent = n.parent
	n.parent = nleft
	if nleft.right != nil {
		nleft.right.parent = n
	}
	n.left = nleft.right
	nleft.right = n
}

func (n *brNode) String(depth int) string {
	if n == nil {
		return ""
	}
	buf := bytes.NewBufferString("")
	right := n.right.String(depth + 1)
	buf.WriteString(right)
	buf.WriteString("\n")

	for i := 0; i < depth; i++ {
		buf.WriteString("    ")
	}
	buf.WriteString(colorStr(n.elem.String(), n.color))
	buf.WriteString("\n")

	left := n.left.String(depth + 1)
	buf.WriteString(left)
	return buf.String()
}

//black height from n to leaf
func (n *brNode) blackHeight(leaf NodeVal) int {
	if n == nil {
		return 0
	}
	var iterator *brNode
	cmp := n.elem.Compare(leaf)
	if cmp < 0 {
		iterator = n.right
	} else if cmp == 0 {
		return 0
	} else {
		iterator = n.left
	}

	h := 0
	for {
		//leaf not found
		if iterator == nil {
			return -1
		}

		if iterator.color == colorBlack {
			h++
		}

		cmp := iterator.elem.Compare(leaf)
		if cmp == 0 {
			return h
		}
		if cmp > 0 {
			iterator = iterator.left
		} else {
			iterator = iterator.right
		}
	}
}

func (n *brNode) getleaves() []*brNode {
	leaves := []*brNode{}
	n.iteratorLeaves(&leaves)
	return leaves
}

func (n *brNode) iteratorLeaves(leaves *[]*brNode) {
	if n.left == nil || n.right == nil {
		*leaves = append(*leaves, n)
		return
	}
	if n.left != nil {
		n.left.iteratorLeaves(leaves)
	}
	if n.right != nil {
		n.right.iteratorLeaves(leaves)
	}
	return
}

func (n *brNode) validBlackHeight() bool {
	leaves := n.getleaves()
	height := n.blackHeight(leaves[0].elem)
	for _, leaf := range leaves {
		h := n.blackHeight(leaf.elem)
		if h != height {
			fmt.Println(height, n.elem, leaf.elem, h)
			return false
		}
	}
	return true
}

func (n *brNode) valid() bool {
	if n == nil {
		return true
	}
	if n.parent == nil && n.color != colorBlack {
		fmt.Println(1)
		return false
	}
	if n.color == colorRed {
		if n.left != nil && n.left.color == colorRed {
			fmt.Println(2)
			return false
		}
		if n.right != nil && n.right.color == colorRed {
			fmt.Println(3)
			return false
		}
	}
	if !n.validBlackHeight() {
		fmt.Println(3)
		return false
	}
	return n.left.valid() && n.right.valid()
}

func (t *BrTree) valid() bool {
	if t == nil || t.root == nil {
		return true
	}
	return t.root.valid()
}

func (t *BrTree) rawInsert(val NodeVal) *brNode {
	node := &brNode{
		elem:  val,
		color: colorRed,
	}
	if t.root == nil {
		node.color = colorBlack
		t.root = node
		return node
	}

	iterator := t.root
	for {
		cmp := iterator.elem.Compare(val)
		if cmp == 0 {
			iterator.elem.EqualMerge(val)
			return nil
		} else if cmp > 0 {
			if iterator.left == nil {
				node.parent = iterator
				iterator.left = node
				return node
			} else {
				iterator = iterator.left
			}
		} else {
			if iterator.right == nil {
				node.parent = iterator
				iterator.right = node
				return node
			} else {
				iterator = iterator.right
			}
		}
	}
}

func (t *BrTree) get(key NodeVal) *brNode {
	if t == nil {
		return nil
	}
	iterator := t.root
	for {
		if iterator == nil {
			return nil
		}
		cmp := iterator.elem.Compare(key)
		if cmp == 0 {
			return iterator
		}
		if cmp > 0 {
			iterator = iterator.left
		} else {
			iterator = iterator.right
		}
	}
}

func (t *BrTree) insertFix(new_node *brNode) {
	var p, pp *brNode
	iterator := new_node
	for {
		p = iterator.parent
		if p == nil || p.color == colorBlack {
			if p == nil {
				iterator.color = colorBlack
				t.root = iterator
			}
			return
		}
		pp = p.parent
		if pp == nil {
			p.color = colorBlack
			t.root = p
			return
		}
		uncle := pp.left
		if uncle == p {
			uncle = pp.right
		}
		if uncle != nil && uncle.color == colorRed {
			uncle.color = colorBlack
			p.color = colorBlack
			pp.color = colorRed
			iterator = pp
		} else {
			if pp.right == uncle {
				if iterator == p.right {
					p.leftRotate()
					p = iterator
				}
				pp.rightRotate()
				p.color = colorBlack
				pp.color = colorRed
				if p.parent == nil {
					t.root = p
				}
				return
			} else {
				if iterator == p.left {
					p.rightRotate()
					p = iterator
				}
				pp.leftRotate()
				p.color = colorBlack
				pp.color = colorRed
				if p.parent == nil {
					t.root = p
				}
				return
			}
		}
	}

}

func (t *BrTree) Insert(val NodeVal) {
	if t == nil {
		return
	}

	new_node := t.rawInsert(val)
	if new_node != nil {
		if new_node.color == colorRed {
			if new_node.parent.color == colorRed {
				t.insertFix(new_node)
			}
		}
	}
}

func (n *brNode) minimum() *brNode {
	if n == nil {
		return nil
	}
	if n.left == nil {
		return n
	}
	iterator := n
	for {
		if iterator.left == nil {
			return iterator
		}
		iterator = iterator.left
	}
}

func (n *brNode) replaceBySub(sub *brNode) *brNode {
	if n == nil {
		return sub
	}
	if sub != nil {
		sub.parent = n.parent
	}
	p := n.parent
	if p != nil {
		if p.left == n {
			p.left = sub
		} else {
			p.right = sub
		}
	}
	n.parent = nil
	return sub
}

func (t *BrTree) rawDelete(val NodeVal) {
	node := t.get(val)
	if node == nil {
		return
	}
	orig := node
	orig_color := node.color
	var del_node, del_node_p *brNode
	isroot := (node.parent == nil)
	if node.left == nil {
		del_node = node.right
		del_node_p = node.parent
		new_node := node.replaceBySub(del_node)
		if isroot {
			t.root = new_node
		}
	} else if node.right == nil {
		del_node = node.left
		del_node_p = node.parent
		new_node := node.replaceBySub(del_node)
		if isroot {
			t.root = new_node
		}
	} else {
		orig = node.right.minimum()
		orig_color = orig.color
		del_node = orig.right
		del_node_p = orig.parent
		orig.replaceBySub(orig.right)
		if node.right != nil {
			orig.right = node.right
			node.right.parent = orig
		}
		new_node := node.replaceBySub(orig)
		orig.left = node.left
		if node.left != nil {
			node.left.parent = orig
		}
		orig.color = node.color
		if isroot {
			t.root = new_node
		}
	}
	if orig_color == colorBlack {
		t.fixDelete(del_node, del_node_p)
	}
}

func (t *BrTree) fixDelete(node, pnode *brNode) {

}
