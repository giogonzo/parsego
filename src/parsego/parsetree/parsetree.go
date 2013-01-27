package pt

type ParseTree struct {
	Value       []byte
	Type        int
	Children    []*ParseTree
	Position    InputPosition
	ActualValue interface{}
	ActualType  interface{}
	ActualId    string
}

type InputPosition struct {
	StartPosition int
	EndPosition   int
	StartLine     int
	EndLine       int
}

type Walker func(level int, node *ParseTree, env interface{}) bool

func (self *ParseTree) Walk(level int, walkerDown, walkerUp Walker, env interface{}) {
	if self == nil {
		return
	}
	if walkerDown(level, self, env) {
		return
	}
	for _, child := range self.Children {
		child.Walk(level+1, walkerDown, walkerUp, env)
	}
	walkerUp(level, self, env)
}
