package tree

type NodeVal interface {
	Compare(other NodeVal) int
	EqualMerge(other NodeVal)
	String() string
}

type Tree interface {
	Insert(node NodeVal)
	Find(key interface{}) NodeVal
	Delete(key interface{})
}
