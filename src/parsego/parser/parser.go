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
	SetLineCount(lineCount int)
	GetLineCount() int
	GetProbeCount() int
}

type Parser func(in State) ([]*pt.ParseTree, bool)

type Cache interface {
	Get(id string) Parser
	Set(id string, match Parser)
}

const (
	TYPE_UNDEFINED = 0
)

/*

*/

type ParseState struct {
	input      []byte
	position   int
	lineCount  int
	probeCount int
}

func (self *ParseState) Next() (int, bool) {
	if self.position >= len(self.input) {
		return 0, false
	}

	next := int(self.input[self.position])
	self.position += 1
	self.probeCount += 1
	if next == '\n' {
		self.lineCount += 1
	}
	return next, true
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

func (self *ParseState) GetLineCount() int {
	return self.lineCount
}

func (self *ParseState) SetLineCount(lineCount int) {
	self.lineCount = lineCount
}

func (self *ParseState) GetProbeCount() int {
	return self.probeCount
}

func InitParser() *ParseState {
	state := new(ParseState)
	state.SetPosition(0)
	state.SetLineCount(1)
	state.SetInput("")
	return state
}

type ParserCache struct {
	parsers map[string]Parser
}

func (self *ParserCache) Get(id string) Parser {
	return self.parsers[id]
}

func (self *ParserCache) Set(id string, match Parser) {
	self.parsers[id] = match
}

var cache Cache = initParserCache()

func initParserCache() Cache {
	cache := new(ParserCache)
	cache.parsers = make(map[string]Parser)
	return cache
}

/*

*/

/*
	Matches *
*/
func Many(match Parser) Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
		nodes := []*pt.ParseTree{}
		for {
			out, ok := Try(match)(in)
			if !ok {
				break
			}
			nodes = concat(nodes, out)
		}
		return nodes, true
	}
}

/*
	Matches +
*/
func Many1(match Parser) Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
		nodes := []*pt.ParseTree{}
		out, ok := match(in)
		if !ok {
			return nil, false
		}
		nodes = concat(nodes, out)

		for {
			out, ok := Try(match)(in)
			if !ok {
				break
			}
			nodes = concat(nodes, out)
		}
		return nodes, true
	}
}

/*
	Matches disjunction
	Wrap parsers in Try(...) calls to preserve state
*/
func Any(matches ...Parser) Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
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
	return Any(matches...)
}

/*
	Matches concatenation
*/
func Concat(matches ...Parser) Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
		nodes := []*pt.ParseTree{}
		for _, match := range matches {
			out, ok := match(in)
			if !ok {
				return nil, false
			}
			nodes = concat(nodes, out)
		}
		return nodes, true
	}
}

/*
	Tries to match, preserving initial state in case of fail
*/
func Try(match Parser) Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
		initialPosition := in.GetPosition()
		initialLineCount := in.GetLineCount()
		out, ok := match(in)
		if !ok {
			in.SetPosition(initialPosition)
			in.SetLineCount(initialLineCount)
		}
		return out, ok
	}
}

/*
	Matches a single character
*/
func Character(c int) Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
		target, ok := in.Next()
		if ok && c == int(target) {
			node := new(pt.ParseTree)
			node.Value = []byte{byte(c)}
			return []*pt.ParseTree{node}, true
		}
		return nil, false
	}
}

/*
	Matches [a-zA-Z]
*/
func Char() Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
		target, ok := in.Next()
		if ok {
			match, _ := regexp.Match("[a-zA-Z]", []byte{byte(target)})
			if match {
				node := new(pt.ParseTree)
				node.Value = []byte{byte(target)}
				return []*pt.ParseTree{node}, true
			}
		}
		return nil, false
	}
}

/*
	Matches [^c]
*/
func AnyCharBut(c int) Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
		target, ok := in.Next()
		if ok {
			match, _ := regexp.Match(fmt.Sprintf("[^%c]", c), []byte{byte(target)})
			if match {
				node := new(pt.ParseTree)
				node.Value = []byte{byte(target)}
				return []*pt.ParseTree{node}, true
			}
		}
		return nil, false
	}
}

/*
	Matches [0-9]
*/
func Number() Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
		target, ok := in.Next()
		if ok {
			match, _ := regexp.Match("[0-9]", []byte{byte(target)})
			if match {
				node := new(pt.ParseTree)
				node.Value = []byte{byte(target)}
				return []*pt.ParseTree{node}, true
			}
		}
		return nil, false
	}
}

/*
	Matches [\s]
*/
func Whitespace() Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
		target, ok := in.Next()
		if ok {
			match, _ := regexp.Match("\\s", []byte{byte(target)})
			if match {
				node := new(pt.ParseTree)
				node.Value = []byte{byte(target)}
				return []*pt.ParseTree{node}, true
			}
		}
		return nil, false
	}
}

/*
	Skips matching
*/
func Skip(match Parser) Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
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
	return Between(
		Whitespaces(),
		match,
		Whitespaces())
}

/*
	Matches between two parsers
*/
func Between(left, match, right Parser) Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
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
	return func(in State) ([]*pt.ParseTree, bool) {
		matched := make([]byte, 0)
		node := new(pt.ParseTree)
		for _, c := range s {
			out, ok := Character(int(c))(in)
			if !ok {
				node.Value = matched
				return []*pt.ParseTree{node}, false
			}
			if out != nil {
				matched = append(matched, out[0].Value[0])
			}
		}
		node.Value = matched
		return []*pt.ParseTree{node}, true
	}
}

/*
	Matches emptiness
*/
func Empty() Parser {
	return func(in State) ([]*pt.ParseTree, bool) {
		return nil, true
	}
}

/*

*/

/*
	Specifies a Node Type
*/
func Specify(nodeType int, match Parser) Parser {
	specId := fmt.Sprintf("_SPEC_%d", nodeType)
	cached := cache.Get(specId)
	if cached == nil {
		cache.Set(specId, func(in State) ([]*pt.ParseTree, bool) {
			pos := new(pt.InputPosition)
			pos.StartPosition = in.GetPosition()
			pos.StartLine = in.GetLineCount()
			out, ok := match(in)
			if !ok {
				return nil, false
			}
			pos.EndPosition = in.GetPosition()
			pos.EndLine = in.GetLineCount()

			nodes := []*pt.ParseTree{new(pt.ParseTree)}
			nodes[0].Type = nodeType
			nodes[0].Position = *pos
			if out != nil && len(out) > 0 && out[0].Type == TYPE_UNDEFINED {
				nodes[0].Value = out[0].Value
			} else {
				appendChildren(nodes[0], out)
			}
			return nodes, true
		})
	}
	return cache.Get(specId)
}

/*
	Helper for recursive rules
*/
func Recursive(id string, matchMaker func() Parser) Parser {
	recId := "_REC_" + id
	cachedRec := cache.Get(recId)
	if cachedRec == nil {
		cache.Set(recId, func(in State) ([]*pt.ParseTree, bool) {
			cached := cache.Get(id)
			if cached == nil {
				cache.Set(id, matchMaker())
			}
			return cache.Get(id)(in)
		})
	}
	return cache.Get(recId)
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

func concat(a, b []*pt.ParseTree) []*pt.ParseTree {
	if b == nil {
		return a
	}

	if len(a) == 1 && len(b) == 1 && a[0].Type == TYPE_UNDEFINED {
		a[0].Value = concatBytes(a[0].Value, b[0].Value)
		return a
	}

	for _, e := range b {
		a = append(a, e)
	}
	return a
}

func appendChildren(node *pt.ParseTree, children []*pt.ParseTree) {
	for _, child := range children {
		node.Children = append(node.Children, child)
	}
}
