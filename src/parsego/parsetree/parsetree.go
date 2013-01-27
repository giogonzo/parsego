package pt

type ParseTree struct {
	Value       []byte
	Type        int
	Children    []*ParseTree
	Position    InputPosition
	ActualValue interface{}
	ActualType  interface{}
}

type InputPosition struct {
	StartPosition int
	EndPosition   int
	StartLine     int
	EndLine       int
}

type Walker func(level int, node *ParseTree, env interface{})

func (self *ParseTree) Walk(level int, walkerDown, walkerUp Walker, env interface{}) {
	if self == nil {
		return
	}
	walkerDown(level, self, env)
	for _, child := range self.Children {
		child.Walk(level+1, walkerDown, walkerUp, env)
	}
	walkerUp(level, self, env)
}
