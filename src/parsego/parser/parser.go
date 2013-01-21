package pg

import (
	"parsego/parsetree"
	"regexp"
)

type State interface {
	Next() (int, bool)
	SetInput(in string)
	GetInput() string
	GetPosition() int
	SetPosition(position int)
}

type Parser func(in State) (*pt.ParseTree, bool)

const (
	TYPE_UNDEFINED = 0
)

/*

*/

type StringState struct {
	position int
	input    []byte
}

func (self *StringState) Next() (int, bool) {
	if self.position >= len(self.input) {
		return 0, false
	}
	self.position += 1
	return int(self.input[self.position-1]), true
}

func (self *StringState) SetInput(in string) {
	self.input = []byte(in)
}

func (self *StringState) GetInput() string {
	return string(self.input)
}

func (self *StringState) GetPosition() int {
	return self.position
}

func (self *StringState) SetPosition(position int) {
	self.position = position
}

/*

*/

/*
	Matches *
*/
func Many(match Parser) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		node := new(pt.ParseTree)
		for {
			out, ok := match(in)
			if !ok {
				break
			}

			appendChild(node, out)
		}
		flatten(node)
		return node, true
	}
}

/*
	Matches +
*/
func Many1(match Parser) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		node := new(pt.ParseTree)
		out, ok := match(in)
		if !ok {
			return nil, false
		}
		appendChild(node, out)

		for {
			out, ok := match(in)
			if !ok {
				break
			}
			appendChild(node, out)
		}
		flatten(node)
		return node, true
	}
}

/*
	Matches parsers in |
	Wrap parsers in Try(...) calls to preserve state
*/
func Any(matches ...Parser) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		for _, match := range matches {
			out, ok := match(in)
			if ok {
				return out, true
			}
		}
		return nil, false
	}
}

/*
	Matches concatenation
*/
func Concat(matches ...Parser) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		node := new(pt.ParseTree)
		for _, match := range matches {
			out, ok := match(in)
			if !ok {
				return nil, false
			}
			appendChild(node, out)
		}
		flatten(node)
		return node, true
	}
}

/*
	Tries to match, preserving initial state in case of fail
*/
func Try(match Parser) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		initialPosition := in.GetPosition()
		out, ok := match(in)
		if !ok {
			in.SetPosition(initialPosition)
		}
		return out, ok
	}
}

/*
	Matches a single character
*/
func Character(c int) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		target, ok := in.Next()
		if ok && c == int(target) {
			node := new(pt.ParseTree)
			node.Value = []byte{byte(c)}
			return node, true
		}
		return nil, false
	}
}

/*
	Matches [a-zA-Z]
*/
func Char() Parser {
	return func(in State) (*pt.ParseTree, bool) {
		target, ok := in.Next()
		if ok {
			match, _ := regexp.Match("[a-zA-Z]", []byte{byte(target)})
			if match {
				node := new(pt.ParseTree)
				node.Value = []byte{byte(target)}
				return node, true
			}
		}
		return nil, false
	}
}

/*
	Matches [0-9]
*/
func Number() Parser {
	return func(in State) (*pt.ParseTree, bool) {
		target, ok := in.Next()
		if ok {
			match, _ := regexp.Match("[0-9]", []byte{byte(target)})
			if match {
				node := new(pt.ParseTree)
				node.Value = []byte{byte(target)}
				return node, true
			}
		}
		return nil, false
	}
}

/*
	Matches [\s]
*/
func Whitespace() Parser {
	return func(in State) (*pt.ParseTree, bool) {
		target, ok := in.Next()
		if ok {
			match, _ := regexp.Match("\\s", []byte{byte(target)})
			if match {
				node := new(pt.ParseTree)
				node.Value = []byte{byte(target)}
				return node, true
			}
		}
		return nil, false
	}
}

/*
	Skips matching
*/
func Skip(match Parser) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		_, ok := match(in)
		return nil, ok
	}
}

/*
	Skips whitespace(s)
*/
func Whitespaces() Parser {
	return Skip(Many(Try(Whitespace())))
}

/*
	Trims matching
*/
func Trim(match Parser) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		Whitespaces()(in)
		out, ok := match(in)
		Whitespaces()(in)
		return out, ok
	}
}

/*
	Matches between parens
*/
func Parens(match Parser) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		_, okl := Character('(')(in)
		if !okl {
			return nil, false
		}
		Whitespaces()(in)
		out, okm := match(in)
		Whitespaces()(in)
		_, okr := Character(')')(in)
		if !okr {
			return nil, false
		}
		return out, okm
	}
}

/*
	Matches =
*/
func Equal() Parser {
	return Character('=')
}

/*
	Matches ;
*/
func Semi() Parser {
	return Character(';')
}

/*
	Matches exact string
*/
func String(s string) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		matched := make([]byte, 0)
		node := new(pt.ParseTree)
		for _, c := range s {
			out, ok := Character(int(c))(in)
			if !ok {
				node.Value = matched
				return node, false
			}
			matched = append(matched, out.Value[0])
		}
		node.Value = matched
		return node, true
	}
}

/*

*/

/*
	Specifies a Node Type
*/
func Specify(nodeType int, match Parser) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		out, ok := match(in)
		if !ok {
			return nil, false
		}
		out.Type = nodeType
		return out, true
	}

}

/*
	Utility
*/

func concatBytes(old1, old2 []byte) []byte {
	newslice := make([]byte, len(old1)+len(old2))
	copy(newslice, old1)
	copy(newslice[len(old1):], old2)
	return newslice
}

func appendChild(parent, child *pt.ParseTree) {
	if child == nil {
		return
	}
	if child.Type == TYPE_UNDEFINED && parent.Children != nil {
		last := parent.Children[len(parent.Children)-1]
		if last.Type == TYPE_UNDEFINED {
			// concat values
			last.Value = concatBytes(last.Value, child.Value)
			return
		}
	}
	// append as new child
	parent.Children = append(parent.Children, child)
}

func flatten(node *pt.ParseTree) {
	concat := []byte{}
	for _, child := range node.Children {
		if child.Type != TYPE_UNDEFINED {
			return
		}
		concat = concatBytes(concat, child.Value)
	}
	node.Children = []*pt.ParseTree{}
	node.Value = concat
}
