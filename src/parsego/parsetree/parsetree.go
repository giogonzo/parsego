package pt

type ParseTree struct {
	Value    []byte
	Type     int
	Parent   *ParseTree
	Children []*ParseTree
}

type Walker func(level int, node *ParseTree)

func (self *ParseTree) Walk(level int, walker Walker) {
	if self == nil {
		return
	}
	walker(level, self)
	for _, child := range self.Children {
		child.Walk(level+1, walker)
	}
}
