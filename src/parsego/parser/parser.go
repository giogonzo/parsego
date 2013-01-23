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
	GetProbeCount() int
}

type Parser func(in State) (*pt.ParseTree, bool)

const (
	TYPE_UNDEFINED = 0
)

/*

*/

type ParseState struct {
	position   int
	input      []byte
	probeCount int
}

func (self *ParseState) Next() (int, bool) {
	if self.position >= len(self.input) {
		return 0, false
	}
	self.position += 1
	self.probeCount += 1
	return int(self.input[self.position-1]), true
}

func (self *ParseState) SetInput(in string) {
	self.input = []byte(in)
}

func (self *ParseState) GetInput() string {
	return string(self.input)
}

func (self *ParseState) GetPosition() int {
	return self.position
}

func (self *ParseState) SetPosition(position int) {
	self.position = position
}

func (self *ParseState) GetProbeCount() int {
	return self.probeCount
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
			initialPosition := in.GetPosition()
			out, ok := match(in)
			if !ok {
				in.SetPosition(initialPosition)
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
			initialPosition := in.GetPosition()
			out, ok := match(in)
			if !ok {
				in.SetPosition(initialPosition)
				break
			}
			appendChild(node, out)
		}
		flatten(node)
		return node, true
	}
}

/*
	Matches disjunction
	Wrap parsers in Try(...) calls to preserve state
*/
func Any(matches ...Parser) Parser {
	return any(matches)
}

func any(matches []Parser) Parser {
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
	Matches disjunction,
	preserving state in case of failure
*/
func TryAny(matches ...Parser) Parser {
	for i, match := range matches {
		matches[i] = Try(match)
	}
	return any(matches)
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
	Skips a character
*/
func SkipChar(c int) Parser {
	return Skip(Character(c))
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
	Matches between two parsers
*/
func Between(left, match, right Parser) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		_, okl := left(in)
		if !okl {
			return nil, false
		}

		out, ok := match(in)

		_, okr := right(in)
		if !okr {
			return nil, false
		}

		return out, ok
	}
}

/*
	Matches between parens, skipping internal whitespaces
*/
func Parens(match Parser) Parser {
	return Between(
		Concat(
			Skip(
				Character('(')),
			Whitespaces()),
		match,
		Concat(
			Skip(
				Whitespaces()),
			Character(')')))
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
	Matches emptiness
*/
func Empty() Parser {
	return func(in State) (*pt.ParseTree, bool) {
		return nil, true
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
		if out == nil {
			out = new(pt.ParseTree)
		}
		out.Type = nodeType
		return out, true
	}
}

/*
	Helper for recursive rules
*/
func Recursive(match func() Parser) Parser {
	return func(in State) (*pt.ParseTree, bool) {
		return match()(in)
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
		} else {
			for _, c := range child.Children {
				if len(c.Value) > 0 || c.Type != TYPE_UNDEFINED {
					appendChild(parent, c)
				}
			}
			child.Children = []*pt.ParseTree{}
		}
	}
	// append as new child
	if len(child.Value) > 0 || child.Type != TYPE_UNDEFINED {
		parent.Children = append(parent.Children, child)
	}
}

func flatten(node *pt.ParseTree) {
	if node.Type != TYPE_UNDEFINED {
		return
	}
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
