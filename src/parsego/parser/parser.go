package pg

import (
	"fmt"
	"regexp"
)

type State interface {
	Next() (int, bool)
	SetInput(in string)
	GetInput() string
	GetPosition() int
	SetPosition(position int)
}

type Parser func(in State) (Output, bool)

type Output []byte

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
	return func(in State) (Output, bool) {
		target, ok := in.Next()
		if ok && c == int(target) {
			return []byte{byte(c)}, true
		}
		return nil, false
	}
}

/*
	Matches *
*/
func Many(match Parser) Parser {
	return func(in State) (Output, bool) {
		matched := make([]byte, 0)
		for {
			out, ok := match(in)
			if !ok {
				break
			}

			if out != nil {
				matched = concat(matched, out)
			}
		}
		return matched, true
	}
}

/*
	Matches {n,n}
*/
func ManyN(match Parser, n int) Parser {
	return func(in State) (Output, bool) {
		matched := make([]byte, 0)
		for i := 0; i < n; i += 1 {
			out, ok := match(in)
			if !ok {
				return nil, false
			}

			if out != nil {
				matched = concat(matched, out)
			}
		}
		return matched, true
	}
}

/*
	Matches +
*/
func Many1(match Parser) Parser {
	return func(in State) (Output, bool) {
		out, ok := match(in)
		if !ok {
			return nil, false
		}
		matched := []byte(out)

		for {
			out, ok := match(in)
			if !ok {
				break
			}

			if out != nil {
				matched = concat(matched, out)
			}
		}
		return matched, true
	}
}

/*
	Matches parsers in |
	Wrap parsers in Try(...) calls to preserve state
*/
func Any(matches ...Parser) Parser {
	return func(in State) (Output, bool) {
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
	Tries to match, preserving initial state in case of fail
*/
func Try(match Parser) Parser {
	return func(in State) (Output, bool) {
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
	return func(in State) (Output, bool) {
		target, ok := in.Next()
		if ok {
			match, _ := regexp.Match("[a-zA-Z]", []byte{byte(target)})
			if match {
				return []byte{byte(target)}, true
			}
		}
		return nil, false
	}
}

/*
	Matches [0-9]
*/
func Number() Parser {
	return func(in State) (Output, bool) {
		target, ok := in.Next()
		if ok {
			match, _ := regexp.Match("[0-9]", []byte{byte(target)})
			if match {
				return []byte{byte(target)}, true
			}
		}
		return nil, false
	}
}

/*
	Matches concatenation
*/
func Concat(matches ...Parser) Parser {
	return func(in State) (Output, bool) {
		matched := make([]byte, 0)
		for _, match := range matches {
			out, ok := match(in)
			if !ok {
				return matched, false
			}
			matched = concat(matched, out)
		}
		return matched, true
	}
}

/*
	Matches [\s]
*/
func Whitespace() Parser {
	return func(in State) (Output, bool) {
		target, ok := in.Next()
		if ok {
			match, _ := regexp.Match("\\s", []byte{byte(target)})
			if match {
				return []byte{byte(target)}, true
			}
		}
		return nil, false
	}
}

/*
	Skips matching
*/
func Skip(match Parser) Parser {
	return func(in State) (Output, bool) {
		_, ok := match(in)
		if ok {
			in.SetPosition(in.GetPosition() - 1)
		}
		return nil, ok
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
	return func(in State) (Output, bool) {
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
	return func(in State) (Output, bool) {
		matched := make([]byte, 0)
		for _, c := range s {
			out, ok := Character(int(c))(in)
			if !ok {
				return matched, false
			}
			matched = append(matched, out[0])
		}
		return matched, true
	}
}

/*
	Utility
*/

func concat(old1, old2 []byte) []byte {
	newslice := make([]byte, len(old1)+len(old2))
	copy(newslice, old1)
	copy(newslice[len(old1):], old2)
	return newslice
}
