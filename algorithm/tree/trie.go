package tree

import (
	"bytes"
	"fmt"
)

type trieNode struct {
	p      *trieNode
	childs []*trieNode
	end    bool
	ch     rune
}

type TrieTree struct {
	root *trieNode
}

func (t *TrieTree) Insert(word string) {
	if t == nil || t.root == nil || len(word) == 0 {
		return
	}
	runes := []rune(word)
	iterator := t.root
	wordlen := len(runes)

	for i, r := range runes {
		childscnt := len(iterator.childs)
		end := (i == wordlen-1)
		if childscnt == 0 {
			node := &trieNode{
				p:      iterator,
				childs: []*trieNode{},
				end:    end,
				ch:     r,
			}
			iterator.childs = append(iterator.childs, node)
			iterator = node
			continue
		} else {
			exist := false
			for _, c := range iterator.childs {
				if c.ch == r {
					iterator = c
					if end {
						c.end = true
						return
					}
					exist = true
					break
				}
			}
			if exist {
				continue
			}
			node := &trieNode{
				p:      iterator,
				childs: []*trieNode{},
				end:    end,
			}
			iterator.childs = append(iterator.childs, node)
			iterator = node
		}
	}

}

func (t *TrieTree) Get(word string) bool {
	if t == nil || t.root == nil || word == "" {
		return true
	}
	iterator := t.root
	wordlen := len(word)

	for i, r := range word {
		if iterator == nil {
			return false
		}
		for _, c := range iterator.childs {
			if c.ch == r {
				if i == wordlen-1 {
					return true
				}
				iterator = c
				break
			}
		}
	}
	return false
}

func (t *TrieTree) Delete(word string) {
	if !t.Get(word) {
		return
	}
}

func (t *TrieTree) String() string {
	if t == nil || t.root == nil || len(t.root.childs) == 0 {
		return ""
	}
	buf := bytes.NewBuffer([]byte{})
	nodes := t.root.childs
	for {
		node_append := []*trieNode{}
		for _, node := range nodes {
			fmt.Fprint(buf, string(node.ch))
			if len(node.childs) != 0 {
				node_append = append(node_append, node.childs...)
			}
		}
		buf.WriteString("\n")
		nodes = node_append
		if len(nodes) == 0 {
			return buf.String()
		}
	}
}
