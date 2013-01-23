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
	FOR
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
	FOR:            "|FOR",
	BLOCK:          "|BLOCK",
}

/*
	Identifier  ←  [a-zA-Z][a-zA-Z0-9]*
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
	NumberLiteral  ←  [0-9]+
*/
func NumberLiteral() pg.Parser {
	return pg.Specify(NUMBER_LITERAL,
		pg.Many1(
			pg.Try(
				pg.Number())))
}

/*
	StringLiteral  ←  '"' [a-zA-Z0-9]* '"'
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
	BoolLiteral  ←  'true' | 'false'
*/
func BoolLiteral() pg.Parser {
	return pg.Specify(BOOL_LITERAL,
		pg.TryAny(
			pg.String("true"),
			pg.String("false")))
}

func Literal() pg.Parser {
	return pg.TryAny(
		NumberLiteral(),
		StringLiteral(),
		BoolLiteral())
}

/*
	Assignment  ←  Identifier '=' Expression
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

/*
	Statement  ←  ControlStatement | Assignment
*/
func Statement() pg.Parser {
	return pg.Trim(
		pg.TryAny(
			pg.Recursive(ControlStatement),
			Assignment()))
}

/*
	ControlStatement  ←  Loop
*/
func ControlStatement() pg.Parser {
	return pg.TryAny(
		Loop())
}

/*
	Loop  ←  Foreach | For
*/
func Loop() pg.Parser {
	return pg.TryAny(
		Foreach(),
		For())
}

/*
	Foreach  ←  'for' Identifier 'in' Identifier Block
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

/*
	For  ←  'for' '(' Assignment* ';' Expression ';'
		Assignment* ')' Block
*/
func For() pg.Parser {
	return pg.Specify(FOR,
		pg.Concat(
			pg.Skip(
				pg.String("for")),
			pg.Whitespaces(),
			pg.Parens(
				pg.Concat(
					pg.Many(
						Assignment()),
					pg.Whitespaces(),
					pg.SkipChar(';'),
					pg.Whitespaces(),
					Expression(),
					pg.Whitespaces(),
					pg.SkipChar(';'),
					pg.Whitespaces(),
					pg.Many(
						Assignment()))),
			pg.Whitespaces(),
			Block()))
}

/*
	Block  ←  '{' Statement* '}'
*/
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
		test = true
		for ( i = 0; test; i = i + 1 ) {
			test = i + (1 + 2 * 3) * 4
			varName = men
			for person in people {
				test = false
			}
		}
		`)
	out, ok := pg.Many(Statement())(in)

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
