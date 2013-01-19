package pt

type ParseTree struct {
	Value    []byte
	Type     int
	Parent   *ParseTree
	Children []*ParseTree
}
