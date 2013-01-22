package main

import (
	"fmt"
	"parsego/parser"
	"parsego/parsetree"
)

const (
	TYPE_UNDEFINED = iota
	IDENTIFIER
	NUMBER_LITERAL
	STRING_LITERAL
	BOOL_LITERAL
	LITERAL
	ASSIGNMENT
	EXPRESSION
	SUM
	PRODUCT
	VALUE
	FOREACH
	BLOCK
)

var NODE_TYPES = map[int]string{
	TYPE_UNDEFINED: "|",

	IDENTIFIER:     "|IDENTIFIER",
	NUMBER_LITERAL: "|NUMBER_LITERAL",
	STRING_LITERAL: "|STRING_LITERAL",
	BOOL_LITERAL:   "|BOOL_LITERAL",
	ASSIGNMENT:     "|ASSIGNMENT",
	SUM:            "|SUM",
	PRODUCT:        "|PRODUCT",
	FOREACH:        "|FOREACH",
	BLOCK:          "|BLOCK",
}

/*
	Matches an IDENTIFIER
*/
func Identifier() pg.Parser {
	return pg.Specify(IDENTIFIER,
		pg.Concat(
			pg.Char(),
			pg.Many(
				pg.TryAny(
					pg.Char(),
					pg.Number()))))
}

/*
	Matches a NUMBER_LITERAL
*/
func NumberLiteral() pg.Parser {
	return pg.Specify(NUMBER_LITERAL,
		pg.Many1(
			pg.Try(
				pg.Number())))
}

/*
	Matches a STRING_LITERAL
*/
func StringLiteral() pg.Parser {
	return pg.Specify(STRING_LITERAL,
		pg.Concat(
			pg.Skip(
				pg.Character('"')),
			pg.Many(
				pg.TryAny(
					pg.Char(),
					pg.Number())),
			pg.Skip(
				pg.Character('"'))))
}

/*
	Matches a BOOL_LITERAL
*/
func BoolLiteral() pg.Parser {
	return pg.Specify(BOOL_LITERAL,
		pg.TryAny(
			pg.String("true"),
			pg.String("false")))
}

/*
	Matches a LITERAL
*/
func Literal() pg.Parser {
	return pg.TryAny(
		NumberLiteral(),
		StringLiteral(),
		BoolLiteral())
}

/*
	Matches an ASSIGNMENT
*/
func Assignment() pg.Parser {
	return pg.Specify(ASSIGNMENT,
		pg.Trim(
			pg.Concat(
				Identifier(),
				pg.Whitespaces(),
				pg.Skip(
					pg.Equal()),
				pg.Whitespaces(),
				Expression())))
}

/*
	Matches an EXPRESSION
	Expr ← Sum
*/
func Expression() pg.Parser {
	return pg.Trim(
		Sum())
}

/*
	Sum ← Product ('+' | '-') Product | Product
*/
func Sum() pg.Parser {
	return pg.TryAny(
		pg.Specify(SUM,
			pg.Concat(
				Product(),
				pg.Trim(
					pg.Concat(
						pg.Skip(
							SumOperator()),
						pg.Whitespaces(),
						Product())))),
		pg.Trim(
			Product()))
}

/*
	Product ← Value ('*' | '/') Value | Value
*/
func Product() pg.Parser {
	return pg.TryAny(
		pg.Specify(PRODUCT,
			pg.Concat(
				Value(),
				pg.Trim(
					pg.Concat(
						pg.Skip(
							ProductOperator()),
						pg.Whitespaces(),
						Value())))),
		pg.Trim(
			Value()))
}

/*
	Value   ← '(' Expr ')' | Identifier | Literal
*/
func Value() pg.Parser {
	return pg.TryAny(
		pg.Parens(
			pg.Recursive(
				Expression)),
		Identifier(),
		Literal())
}

func SumOperator() pg.Parser {
	return pg.TryAny(
		pg.Character('+'),
		pg.Character('-'))
}

func ProductOperator() pg.Parser {
	return pg.TryAny(
		pg.Character('*'),
		pg.Character('/'))
}

func Statement() pg.Parser {
	return pg.Trim(
		pg.TryAny(
			pg.Recursive(ControlStatement),
			Assignment()))
}

/*
	TODO
*/
func ControlStatement() pg.Parser {
	return pg.TryAny(
		Loop())
}

/*
	TODO
*/
func Loop() pg.Parser {
	return pg.TryAny(
		Foreach())
}

/*
	TODO
*/
func Foreach() pg.Parser {
	return pg.Specify(FOREACH,
		pg.Concat(
			pg.Skip(
				pg.String("for")),
			pg.Whitespaces(),
			Identifier(),
			pg.Whitespaces(),
			pg.Skip(
				pg.String("in")),
			pg.Whitespaces(),
			Identifier(),
			pg.Whitespaces(),
			Block()))
}

func Block() pg.Parser {
	return pg.Specify(BLOCK,
		pg.Trim(
			pg.Between(
				pg.Character('{'),
				pg.Many(Statement()),
				pg.Character('}'))))
}

/*

*/
func main() {
	in := new(pg.StringState)
	in.SetInput(`
		for man in men {
			test = 0 + (1 + 2 * 3) * 4
			varName = men
			for person in people {
				test = false
			}
		}
		`)
	out, ok := Statement()(in)

	fmt.Printf("Parse ok: %t\n", ok)
	fmt.Printf("Left: %d\n", len(in.GetInput())-in.GetPosition())
	fmt.Print("Parsed:\n")
	out.Walk(0, func(level int, node *pt.ParseTree) {
		for i := 0; i < level; i += 1 {
			fmt.Print("   ")
		}
		fmt.Printf("%s [%s]\n", NODE_TYPES[node.Type], node.Value)
	})
}
