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
	root      *brNode
	emptyNode *brNode
	height    int
}

type brNode struct {
	empty  bool
	color  int
	parent *brNode
	left   *brNode
	right  *brNode
	elem   NodeVal
}

type fakeNode int

func (n fakeNode) Compare(other NodeVal) int {
	val, _ := other.(fakeNode)
	return int(n) - int(val)
}

func (n fakeNode) String() string {
	return fmt.Sprintf("%4d", n)
}

func (n fakeNode) EqualMerge(other NodeVal) {}

func newBrTree() *BrTree {
	emptyNode := &brNode{
		empty: true,
		elem:  fakeNode(-100),
		color: colorBlack,
	}
	emptyNode.parent = emptyNode
	emptyNode.left = emptyNode
	emptyNode.right = emptyNode
	return &BrTree{
		root:      emptyNode,
		emptyNode: emptyNode,
	}
}

func colorStr(s string, color int) string {
	if color == colorRed {
		return fmt.Sprintf("\x1b[31m%s\x1b[0m", s)
	}
	return fmt.Sprintf("\x1b[34m%s\x1b[0m", s)
}

func (t *BrTree) String() string {
	return t.root.String(0)
}

func (n *brNode) leftRotate() {
	if n.empty {
		return
	}
	nright := n.right
	if !n.parent.empty {
		if n.parent.left == n {
			n.parent.left = nright
		} else {
			n.parent.right = nright
		}
	}
	nright.parent = n.parent
	n.parent = nright
	if !nright.left.empty {
		nright.left.parent = n
	}
	n.right = nright.left
	nright.left = n
}

func (n *brNode) rightRotate() {
	if n.empty {
		return
	}
	nleft := n.left
	if !n.parent.empty {
		if n.parent.left == n {
			n.parent.left = nleft
		} else {
			n.parent.right = nleft
		}
	}
	nleft.parent = n.parent
	n.parent = nleft
	if !nleft.right.empty {
		nleft.right.parent = n
	}
	n.left = nleft.right
	nleft.right = n
}

func (n *brNode) String(depth int) string {
	if n.empty {
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
	if n.empty {
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
		if iterator.empty {
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
	if n.left.empty || n.right.empty {
		*leaves = append(*leaves, n)
		return
	}
	if !n.left.empty {
		n.left.iteratorLeaves(leaves)
	}
	if !n.right.empty {
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
	if n.empty {
		return true
	}

	if n.parent.empty && n.color != colorBlack {
		fmt.Println(1)
		return false
	}
	if n.color == colorRed {
		if n.left.color == colorRed || n.right.color == colorRed {
			fmt.Println(2)
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
	if t == nil || t.root.empty {
		return true
	}
	return t.root.valid()
}

func (t *BrTree) rawInsert(val NodeVal) *brNode {
	node := &brNode{
		elem:  val,
		color: colorRed,
		left:  t.emptyNode,
		right: t.emptyNode,
	}

	if t.root.empty {
		node.color = colorBlack
		t.root = node
		node.parent = t.emptyNode
		return node
	}

	iterator := t.root
	for {
		cmp := iterator.elem.Compare(val)
		if cmp == 0 {
			iterator.elem.EqualMerge(val)
			return nil
		} else if cmp > 0 {
			if iterator.left.empty {
				node.parent = iterator
				iterator.left = node
				return node
			} else {
				iterator = iterator.left
			}
		} else {
			if iterator.right.empty {
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
		if iterator.empty {
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
		if p.color == colorBlack {
			if p.empty {
				iterator.color = colorBlack
				t.root = iterator
			}
			return
		}
		pp = p.parent
		if pp.empty {
			p.color = colorBlack
			t.root = p
			return
		}
		uncle := pp.left
		if uncle == p {
			uncle = pp.right
		}
		if uncle.color == colorRed {
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
				if p.parent.empty {
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
				if p.parent.empty {
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
	if n.empty {
		return nil
	}
	if n.left.empty {
		return n
	}
	iterator := n
	for {
		if iterator.left.empty {
			return iterator
		}
		iterator = iterator.left
	}
}

func (n *brNode) replaceBySub(sub *brNode) *brNode {
	if n.empty {
		return sub
	}
	if sub != nil {
		sub.parent = n.parent
	}
	p := n.parent
	if !p.empty {
		if p.left == n {
			p.left = sub
		} else {
			p.right = sub
		}
	}
	n.parent = nil
	return sub
}

func (t *BrTree) Delete(val NodeVal) {
	node := t.get(val)
	if node == nil {
		return
	}
	if node.parent.empty {
		if node.left == node.right {
			node.parent = nil
			t.root = t.emptyNode
		}
		return
	}

	orig := node
	orig_color := orig.color
	var del_node *brNode

	p := node.parent
	if node.left.empty {
		del_node = node.right
		node.replaceBySub(node.right)
		node.left.parent = p

	} else if node.right.empty {
		del_node = node.left
		node.replaceBySub(node.left)
		node.right.parent = p
	} else {
		orig = node.right.minimum()
		orig_color = orig.color
		del_node = orig.right
		if orig.parent == node {
			del_node.parent = orig
		} else {
			orig.replaceBySub(orig.right)
			orig.right = node.right
			orig.right.parent = orig
		}
		node.replaceBySub(orig)
		orig.left = node.left
		orig.left.parent = orig
		orig.color = node.color

	}
	if orig_color == colorBlack {
		t.fixDelete(del_node, val)
		t.emptyNode.parent = t.emptyNode
	}
}

func (t *BrTree) fixDelete(del_node *brNode, val NodeVal) {
	//新填充的为红色并且:
	//为原来的子节点,所以只需变色
	iterator := del_node
	for {
		if iterator.color == colorRed {
			iterator.color = colorBlack
			return
		}
		p := iterator.parent
		bro := p.left
		if bro == iterator {
			bro = p.right
		}
		if bro == p.right {
			if bro.color == colorRed {
				if p.parent.empty {
					t.root = iterator
				}
				p.leftRotate()
				p.color = colorRed
				bro.color = colorBlack
				continue
			}

			if bro.right.color == colorRed {
				if p.parent.empty {
					t.root = iterator
				}
				p.leftRotate()
				bro.color = p.color
				p.color = colorBlack
				bro.right.color = colorBlack
				if bro.parent.empty {
					t.root = bro
					bro.color = colorBlack
				}
				return
			}

			if bro.left.color == colorRed {
				bro.rightRotate()
				bro.parent.color = colorBlack
				bro.color = colorRed
				continue
			}

			if bro.left.color == colorBlack && bro.right.color == colorBlack {
				bro.color = colorRed
				iterator = p
				continue
			}
		} else {
			if bro.color == colorRed {
				if p.parent.empty {
					t.root = iterator
				}
				p.rightRotate()
				p.color = colorRed
				bro.color = colorBlack
				continue
			}

			if bro.left.color == colorRed {
				if p.parent.empty {
					t.root = iterator
				}
				p.rightRotate()
				bro.color = p.color
				p.color = colorBlack

				bro.left.color = colorBlack
				if bro.parent.empty {
					t.root = bro
					bro.color = colorBlack
				}
				return
			}

			if bro.right.color == colorRed {
				bro.leftRotate()
				bro.parent.color = colorBlack
				bro.color = colorRed
				continue
			}

			if bro.left.color == colorBlack && bro.right.color == colorBlack {
				bro.color = colorRed
				iterator = p
				continue
			}
		}
	}
}
