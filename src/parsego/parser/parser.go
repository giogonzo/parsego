package pg

import (
	"fmt"
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
		fmt.Print("\n")
		return 0, false
	}
	self.position += 1
	fmt.Printf("%c", self.input[self.position-1])
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
		return new(pt.ParseTree), false
	}
}

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

			node = concat(node, out)
		}
		return node, true
	}
}

/*
	Matches {n,n}
*/
func ManyN(match Parser, n int) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		node := new(pt.ParseTree)
		for i := 0; i < n; i += 1 {
			out, ok := match(in)
			if !ok {
				return new(pt.ParseTree), false
			}

			node = concat(node, out)
		}
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
			return new(pt.ParseTree), false
		}
		node = concat(node, out)

		for {
			out, ok := match(in)
			if !ok {
				break
			}
			node = concat(node, out)
		}
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
		return new(pt.ParseTree), false
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
		return new(pt.ParseTree), false
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
		return new(pt.ParseTree), false
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
				return node, false
			}
			node = concat(node, out)
		}
		return node, true
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
		return new(pt.ParseTree), false
	}
}

/*
	Skips matching
*/
func Skip(match Parser) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		_, ok := match(in)
		if ok {
			in.SetPosition(in.GetPosition() - 1)
		}
		return new(pt.ParseTree), ok
	}
}

/*
	Skips whitespace(s)
*/
func Whitespaces() Parser {
	return Skip(Many(Whitespace()))
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
			return new(pt.ParseTree), false
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

func concat(o1, o2 *pt.ParseTree) *pt.ParseTree {
	if o1.Type == TYPE_UNDEFINED && o2.Type == TYPE_UNDEFINED {
		// concat values
		fmt.Print("\n->concat values\n")
		o1.Value = concatBytes(o1.Value, o2.Value)
		return o1
	} else if o1.Type == TYPE_UNDEFINED && o2.Type != TYPE_UNDEFINED {
		// append children and return parent
		fmt.Print("\n->append children\n")
		parent := new(pt.ParseTree)
		parent.Children = []*pt.ParseTree{o1, o2}
		return parent
	}

	// else append new child
	fmt.Print("\n->append new child\n")
	if o1.Children == nil {
		o1.Children = []*pt.ParseTree{}
	}
	o1.Children = append(o1.Children, o2)
	return o1
}
