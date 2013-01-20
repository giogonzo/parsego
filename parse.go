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
)

var NODE_TYPES = map[int]string{
	TYPE_UNDEFINED: "|",

	IDENTIFIER:     "|IDENTIFIER",
	NUMBER_LITERAL: "|NUMBER_LITERAL",
	STRING_LITERAL: "|STRING_LITERAL",
	BOOL_LITERAL:   "|BOOL_LITERAL",
	LITERAL:        "|LITERAL",
	ASSIGNMENT:     "|ASSIGNMENT",
	EXPRESSION:     "|EXPRESSION",
}

/*
	Matches an IDENTIFIER
*/
func Identifier() pg.Parser {
	return pg.Specify(IDENTIFIER,
		pg.Concat(
			pg.Char(),
			pg.Many(
				pg.Any(
					pg.Try(
						pg.Char()),
					pg.Try(
						pg.Number())))))
}

/*
	Matches a NUMBER_LITERAL
*/
func NumberLiteral() pg.Parser {
	return pg.Specify(NUMBER_LITERAL,
		pg.Many1(
			pg.Number()))
}

/*
	Matches a STRING_LITERAL
*/
func StringLiteral() pg.Parser {
	return pg.Specify(STRING_LITERAL,
		pg.Concat(
			pg.Character('"'),
			pg.Many(
				pg.Any(
					pg.Try(
						pg.Char()),
					pg.Try(
						pg.Number()))),
			pg.Character('"')))
}

/*
	Matches a BOOL_LITERAL
*/
func BoolLiteral() pg.Parser {
	return pg.Specify(BOOL_LITERAL,
		pg.Any(
			pg.Try(
				pg.String("true")),
			pg.Try(
				pg.String("false"))))
}

/*
	Matches a LITERAL
*/
func Literal() pg.Parser {
	return pg.Specify(LITERAL,
		pg.Any(
			pg.Try(
				NumberLiteral()),
			pg.Try(
				StringLiteral()),
			pg.Try(
				BoolLiteral())))
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
				pg.Equal(),
				pg.Whitespaces(),
				Expression())))
}

/*
	TODO
	Matches an EXPRESSION
*/
func Expression() pg.Parser {
	return pg.Specify(EXPRESSION,
		pg.Any(
			pg.Try(
				Literal()),
			pg.Try(
				Identifier())))
}

/*

*/
func main() {
	in := new(pg.StringState)

	in.SetInput(`
		varName = 10
		test = "string"
		`)
	fmt.Printf("%s\n", in.GetInput())
	out, ok := pg.Many(Assignment())(in)

	fmt.Printf("-------\n")
	fmt.Printf("Parse ok: %t\n", ok)
	fmt.Printf("Last position: %d\n", in.GetPosition())
	fmt.Print("Parsed:\n")
	out.Walk(0, func(level int, node *pt.ParseTree) {
		for i := 0; i < level; i += 1 {
			fmt.Print("  ")
		}
		fmt.Printf("%s [%s]\n", NODE_TYPES[node.Type], node.Value)
	})
}
