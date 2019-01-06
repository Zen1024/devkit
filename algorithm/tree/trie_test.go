package tree

import (
	"fmt"
	"testing"
)

func TestTrie(t *testing.T) {
	tree := &TrieTree{
		root: &trieNode{},
	}
	tree.Insert("hello")
	fmt.Printf("%s\n\n", tree.String())
	tree.Insert("world")
	fmt.Printf("%s\n\n", tree.String())
}
