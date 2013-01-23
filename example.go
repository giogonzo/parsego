package main

import (
	"fmt"
	"parsego/parser"
	"parsego/parsetree"
	"time"
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
	IFTHEN
	IFTHENELSE
	L_COMPARISON
	L_E_COMPARISON
	G_COMPARISON
	G_E_COMPARISON
	E_COMPARISON
	BREAK
	CONTINUE
	OR_EXPRESSION
	AND_EXPRESSION
	FUNCTION_CALL
	FUNCTION_DEFINITION
)

var NODE_TYPES = map[int]string{
	TYPE_UNDEFINED: "?",

	IDENTIFIER:          "IDENTIFIER",
	NUMBER_LITERAL:      "NUMBER_LITERAL",
	STRING_LITERAL:      "STRING_LITERAL",
	BOOL_LITERAL:        "BOOL_LITERAL",
	ASSIGNMENT:          "ASSIGNMENT",
	SUM:                 "SUM",
	PRODUCT:             "PRODUCT",
	FOREACH:             "FOREACH",
	FOR:                 "FOR",
	BLOCK:               "BLOCK",
	IFTHEN:              "IFTHEN",
	IFTHENELSE:          "IFTHENELSE",
	L_COMPARISON:        "L_COMPARISON",
	L_E_COMPARISON:      "L_E_COMPARISON",
	G_COMPARISON:        "G_COMPARISON",
	G_E_COMPARISON:      "G_E_COMPARISON",
	E_COMPARISON:        "E_COMPARISON",
	BREAK:               "BREAK",
	CONTINUE:            "CONTINUE",
	OR_EXPRESSION:       "OR_EXPRESSION",
	AND_EXPRESSION:      "AND_EXPRESSION",
	FUNCTION_CALL:       "FUNCTION_CALL",
	FUNCTION_DEFINITION: "FUNCTION_DEFINITION",
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
			pg.Number()))
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
	Expr ← Comparison
*/
func Expression() pg.Parser {
	return pg.Trim(
		BoolExpression())
}

/*
	BoolExpression ←
		  OrExpression
		| AndExpression
		| Comparison
*/
func BoolExpression() pg.Parser {
	return pg.TryAny(
		OrExpression(),
		AndExpression(),
		Comparison())
}

/*
	OrExpression ← 
*/
func OrExpression() pg.Parser {
	return pg.Specify(OR_EXPRESSION,
		pg.Concat(
			Comparison(),
			pg.Whitespaces(),
			pg.Skip(
				pg.String("||")),
			pg.Whitespaces(),
			Comparison()))
}

/*
	AndExpression ← 
*/
func AndExpression() pg.Parser {
	return pg.Specify(AND_EXPRESSION,
		pg.Concat(
			Comparison(),
			pg.Whitespaces(),
			pg.Skip(
				pg.String("&&")),
			pg.Whitespaces(),
			Comparison()))
}

/*
	Comparison ← 
		  LComparison
		| LEComparison
		| GComparison
		| GEComparison
		| EComparison
		| Sum
*/
func Comparison() pg.Parser {
	return pg.TryAny(
		LComparison(),
		LEComparison(),
		GComparison(),
		GEComparison(),
		EComparison(),
		Sum())
}

/*
	LComparison  ←  Sum '<' Sum
*/
func LComparison() pg.Parser {
	return pg.Specify(L_COMPARISON,
		pg.Concat(
			Sum(),
			pg.Whitespaces(),
			pg.SkipChar('<'),
			pg.Whitespaces(),
			Sum()))
}

/*
	LEComparison  ←  Sum '<=' Sum
*/
func LEComparison() pg.Parser {
	return pg.Specify(L_E_COMPARISON,
		pg.Concat(
			Sum(),
			pg.Whitespaces(),
			pg.Skip(
				pg.String("<=")),
			pg.Whitespaces(),
			Sum()))
}

/*
	GComparison  ←  Sum '>' Sum
*/
func GComparison() pg.Parser {
	return pg.Specify(G_COMPARISON,
		pg.Concat(
			Sum(),
			pg.Whitespaces(),
			pg.SkipChar('>'),
			pg.Whitespaces(),
			Sum()))
}

/*
	GEComparison  ←  Sum '>=' Sum
*/
func GEComparison() pg.Parser {
	return pg.Specify(G_E_COMPARISON,
		pg.Concat(
			Sum(),
			pg.Whitespaces(),
			pg.Skip(
				pg.String(">=")),
			pg.Whitespaces(),
			Sum()))
}

/*
	EComparison  ←  Sum '==' Sum
*/
func EComparison() pg.Parser {
	return pg.Specify(E_COMPARISON,
		pg.Concat(
			Sum(),
			pg.Whitespaces(),
			pg.Skip(
				pg.String("==")),
			pg.Whitespaces(),
			Sum()))
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
		FunctionCall(),
		Identifier(),
		Literal(),
		pg.Parens(
			pg.Recursive(
				Expression)))
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
	Statement  ←
		  FunctionDefinition
		| ControlStatement
		| Assignment
*/
func Statement() pg.Parser {
	return pg.Trim(
		pg.TryAny(
			FunctionDefinition(),
			pg.Recursive(ControlStatement),
			Assignment()))
}

/*
	ControlStatement  ←  Loop | If | Break | Continue
*/
func ControlStatement() pg.Parser {
	return pg.TryAny(
		Loop(),
		If(),
		Break(),
		Continue())
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
	TODO
	For  ←  'for' '(' Assignment* ';' Expression ';'
		Assignment* ')' Block
*/
func For() pg.Parser {
	return pg.Specify(FOR,
		pg.Concat(
			pg.Skip(
				pg.String("for")),
			pg.Whitespaces(),
			//pg.Many(
			Assignment(), //),
			pg.Whitespaces(),
			pg.SkipChar(';'),
			pg.Whitespaces(),
			Expression(),
			pg.Whitespaces(),
			pg.SkipChar(';'),
			pg.Whitespaces(),
			//pg.Many(
			Assignment(), //),
			pg.Whitespaces(),
			Block()))
}

/*
	If  ←  IfThen | IfThenElse
*/
func If() pg.Parser {
	return pg.TryAny(
		IfThenElse(),
		IfThen())
}

/*
	IfThen  ←  'if' Expression Block
*/
func IfThen() pg.Parser {
	return pg.Specify(IFTHEN,
		pg.Concat(
			pg.Skip(
				pg.String("if")),
			pg.Whitespaces(),
			Expression(),
			pg.Whitespaces(),
			Block()))
}

/*
	IfThenElse  ←  'if' Expression Block
		'else' Block
*/
func IfThenElse() pg.Parser {
	return pg.Specify(IFTHENELSE,
		pg.Concat(
			pg.Skip(
				pg.String("if")),
			pg.Whitespaces(),
			Expression(),
			pg.Whitespaces(),
			Block(),
			pg.Whitespaces(),
			pg.Skip(
				pg.String("else")),
			pg.Whitespaces(),
			Block()))
}

/*
	Break  ←  'break'
*/
func Break() pg.Parser {
	return pg.Specify(BREAK,
		pg.Skip(
			pg.String("break")))
}

/*
	Continue  ←  'continue'
*/
func Continue() pg.Parser {
	return pg.Specify(CONTINUE,
		pg.Skip(
			pg.String("continue")))
}

/*
	Block  ←  '{' Statement* '}'
*/
func Block() pg.Parser {
	return pg.Specify(BLOCK,
		pg.Trim(
			pg.Between(
				pg.Character('{'),
				pg.Trim(
					pg.Many(
						Statement())),
				pg.Character('}'))))
}

/*
	FunctionCall  ←  Identifier '(' ParamsList ')'
*/
func FunctionCall() pg.Parser {
	return pg.Specify(FUNCTION_CALL,
		pg.Concat(
			Identifier(),
			pg.Parens(
				ParamsList())))
}

/*
	ParamsList  ←
		  Expression ',' ParamsList
		| Expression
		| 'empty'
*/
func ParamsList() pg.Parser {
	return pg.Trim(
		pg.TryAny(
			pg.Concat(
				pg.Recursive(Expression),
				pg.Whitespaces(),
				pg.SkipChar(','),
				pg.Whitespaces(),
				pg.Recursive(ParamsList)),
			pg.Recursive(Expression),
			pg.Empty()))
}

/*
	FunctionDefinition  ←
		Identifier '(' NamedParamsList ')' Block
*/
func FunctionDefinition() pg.Parser {
	return pg.Specify(FUNCTION_DEFINITION,
		pg.Concat(
			pg.Skip(
				pg.String("func")),
			pg.Whitespaces(),
			Identifier(),
			pg.Parens(
				NamedParamsList()),
			pg.Whitespaces(),
			pg.Recursive(Block)))
}

/*
	NamedParamsList  ←
		  Identifier ',' NamedParamsList
		| Identifier
		| 'empty'
*/
func NamedParamsList() pg.Parser {
	return pg.Trim(
		pg.TryAny(
			pg.Concat(
				Identifier(),
				pg.Whitespaces(),
				pg.SkipChar(','),
				pg.Whitespaces(),
				pg.Recursive(ParamsList)),
			Identifier(),
			pg.Empty()))
}

/*

*/
func main() {
	in := new(pg.ParseState)
	in.SetInput(`
		func callMe() {
		}
		if test == false {
			test = callMe()
		}
		for i = 0; test; i = i + 1 {
			test = test || i + (1 + 2 * 3) * 4 >= 20
			varName = man
			for person in people {
				if test {
					test = false
					continue
				} else {
					test = true
					break
				}
			}
		}
	`)
	start := time.Now()
	out, ok := pg.Many(Statement())(in)
	end := time.Now()

	fmt.Printf("Input length: %d, probe count: %d, total: %s\n", len(in.GetInput()), in.GetProbeCount(), end.Sub(start).String())
	fmt.Printf("Parse ok: %t\n", ok)
	fmt.Printf("Left: %d\n", len(in.GetInput())-in.GetPosition())
	fmt.Print("Parsed:\n")
	out.Walk(0, func(level int, node *pt.ParseTree) {
		for i := 0; i < level; i += 1 {
			fmt.Print("|  ")
		}
		fmt.Printf("%s [%s]\n", NODE_TYPES[node.Type], node.Value)
	})
}
