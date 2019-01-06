package tree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type intNode int

func (n intNode) Compare(other NodeVal) int {
	val, _ := other.(intNode)
	return int(n) - int(val)
}

func (n intNode) String() string {
	return fmt.Sprintf("%4d", n)
}

func (n intNode) EqualMerge(other NodeVal) {}

func randArr(length, max int) []intNode {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	re := []intNode{}
	for i := 0; i < length; i++ {
		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
		re = append(re, intNode(rnd.Intn(max)))
		time.Sleep(time.Millisecond)
	}
	return re
}

func TestInsertDelete(t *testing.T) {
	for {
		array := randArr(10, 1000)
		tree := newBrTree()
		for _, elem := range array {
			tree.Insert(elem)
			if !tree.root.valid() {
				fmt.Println(array)
				fmt.Println(tree.String())
				return
			}
		}
		for _, elem := range array {
			tree.Delete(elem)
			if !tree.root.valid() {
				fmt.Println(array)
				fmt.Println(tree.String())
				return
			}
		}

	}
}
