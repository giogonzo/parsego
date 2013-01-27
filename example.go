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
	FOR_INIT
	FOR_CONDITION
	FOR_STEP
	BLOCK
	IFTHEN
	IFTHENELSE
	SWITCH
	CASE
	CASE_ELSE
	L_COMPARISON
	L_E_COMPARISON
	G_COMPARISON
	G_E_COMPARISON
	E_COMPARISON
	BREAK
	CONTINUE
	RETURN
	OR_EXPRESSION
	AND_EXPRESSION
	FUNCTION_CALL
	FUNCTION_DEFINITION
	PROGRAM
)

var NODE_TYPES = map[int]string{
	TYPE_UNDEFINED: "?",

	IDENTIFIER:          "IDENTIFIER",
	NUMBER_LITERAL:      "NUMBER_LITERAL",
	STRING_LITERAL:      "STRING_LITERAL",
	BOOL_LITERAL:        "BOOL_LITERAL",
	ASSIGNMENT:          "ASSIGNMENT",
	EXPRESSION:          "EXPRESSION",
	SUM:                 "SUM",
	PRODUCT:             "PRODUCT",
	FOREACH:             "FOREACH",
	FOR:                 "FOR",
	FOR_INIT:            "FOR_INIT",
	FOR_CONDITION:       "FOR_CONDITION",
	FOR_STEP:            "FOR_STEP",
	BLOCK:               "BLOCK",
	IFTHEN:              "IFTHEN",
	IFTHENELSE:          "IFTHENELSE",
	SWITCH:              "SWITCH",
	CASE:                "CASE",
	CASE_ELSE:           "CASE_ELSE",
	L_COMPARISON:        "L_COMPARISON",
	L_E_COMPARISON:      "L_E_COMPARISON",
	G_COMPARISON:        "G_COMPARISON",
	G_E_COMPARISON:      "G_E_COMPARISON",
	E_COMPARISON:        "E_COMPARISON",
	BREAK:               "BREAK",
	CONTINUE:            "CONTINUE",
	RETURN:              "RETURN",
	OR_EXPRESSION:       "OR_EXPRESSION",
	AND_EXPRESSION:      "AND_EXPRESSION",
	FUNCTION_CALL:       "FUNCTION_CALL",
	FUNCTION_DEFINITION: "FUNCTION_DEFINITION",
	PROGRAM:             "PROGRAM",
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
	Expr  ←  BoolExpression
*/
func Expression() pg.Parser {
	return pg.Specify(EXPRESSION,
		pg.Trim(
			BoolExpression()))
}

/*
	BoolExpression  ←
		  OrExpression
*/
func BoolExpression() pg.Parser {
	return OrExpression()
}

/*
	OrExpression  ←
		  AndExpression '||' AndExpression
		| AndExpression
*/
func OrExpression() pg.Parser {
	return pg.TryAny(
		pg.Specify(OR_EXPRESSION,
			pg.Concat(
				AndExpression(),
				pg.Whitespaces(),
				pg.Skip(
					pg.String("||")),
				pg.Whitespaces(),
				AndExpression())),
		AndExpression())
}

/*
	AndExpression  ←
		  Comparison '&&' Comparison
		| Comparison
*/
func AndExpression() pg.Parser {
	return pg.TryAny(
		pg.Specify(AND_EXPRESSION,
			pg.Concat(
				Comparison(),
				pg.Whitespaces(),
				pg.Skip(
					pg.String("&&")),
				pg.Whitespaces(),
				Comparison())),
		Comparison())
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
		Literal(),
		Identifier(),
		pg.Parens(
			pg.Recursive(
				"Expression",
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
			FunctionCall(),
			pg.Recursive(
				"ControlStatement",
				ControlStatement),
			Assignment()))
}

/*
	ControlStatement  ←  Loop | If | Break | Continue
*/
func ControlStatement() pg.Parser {
	return pg.TryAny(
		Break(),
		Continue(),
		Return(),
		Loop(),
		If(),
		Switch())
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
	For  ←  'for' '(' ForInit ';' ForCondition ';'
		ForStep ')' Block
*/
func For() pg.Parser {
	return pg.Specify(FOR,
		pg.Concat(
			pg.Skip(
				pg.String("for")),
			pg.Whitespaces(),
			ForInit(),
			pg.Whitespaces(),
			pg.SkipChar(';'),
			pg.Whitespaces(),
			ForCondition(),
			pg.Whitespaces(),
			pg.SkipChar(';'),
			pg.Whitespaces(),
			ForStep(),
			pg.Whitespaces(),
			Block()))
}

/*
	ForInit  ←  AssignmentList
*/
func ForInit() pg.Parser {
	return pg.Specify(FOR_INIT,
		AssignmentList())
}

/*
	ForCondition  ←  Expression
*/
func ForCondition() pg.Parser {
	return pg.Specify(FOR_CONDITION,
		Expression())
}

/*
	ForStep  ←  AssignmentList
*/
func ForStep() pg.Parser {
	return pg.Specify(FOR_STEP,
		AssignmentList())
}

/*
	AssignmentList  ←
		  AssignmentList1
		| 'empty'
*/
func AssignmentList() pg.Parser {
	return pg.Trim(
		pg.TryAny(
			AssignmentList1(),
			pg.Empty()))
}

/*
	AssignmentList1  ←
		  Assignment ',' AssignmentList1
		| Assignment
*/
func AssignmentList1() pg.Parser {
	return pg.Trim(
		pg.TryAny(
			pg.Concat(
				Assignment(),
				pg.Whitespaces(),
				pg.SkipChar(','),
				pg.Recursive(
					"AssignmentList1",
					AssignmentList1)),
			Assignment()))
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
	Return  ←  'return'
*/
func Return() pg.Parser {
	return pg.Specify(RETURN,
		pg.Concat(
			pg.Skip(
				pg.String("return")),
			pg.Whitespaces(),
			Expression()))
}

/*
	Switch  ←  'switch' Expression SwitchBlock
*/
func Switch() pg.Parser {
	return pg.Specify(SWITCH,
		pg.Concat(
			pg.Skip(
				pg.String("switch")),
			pg.Whitespaces(),
			pg.Recursive(
				"Expression",
				Expression),
			pg.Whitespaces(),
			SwitchBlock()))
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
	SwitchBlock  ←  '{' Case* CaseElse? '}'
*/
func SwitchBlock() pg.Parser {
	return pg.Trim(
		pg.Between(
			pg.Character('{'),
			pg.Trim(
				pg.Concat(
					pg.Many(
						Case()),
					pg.Try(
						CaseElse()))),
			pg.Character('}')))
}

/*
	Case  ←  'case' Expression ':' Block
*/
func Case() pg.Parser {
	return pg.Specify(CASE,
		pg.Trim(
			pg.Concat(
				pg.Skip(
					pg.String("case")),
				pg.Whitespaces(),
				pg.Recursive(
					"Expression",
					Expression),
				pg.Whitespaces(),
				pg.SkipChar(':'),
				pg.Whitespaces(),
				pg.Recursive(
					"Block",
					Block))))
}

/*
	CaseElse  ←  'else' ':' Block
*/
func CaseElse() pg.Parser {
	return pg.Specify(CASE_ELSE,
		pg.Trim(
			pg.Concat(
				pg.Skip(
					pg.String("else")),
				pg.Whitespaces(),
				pg.SkipChar(':'),
				pg.Whitespaces(),
				pg.Recursive(
					"Block",
					Block))))
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
		  ParamsList1
		| 'empty'
*/
func ParamsList() pg.Parser {
	return pg.Trim(
		pg.TryAny(
			ParamsList1(),
			pg.Empty()))
}

/*
	ParamsList1  ←
		  Expression ',' ParamsList1
		| Expression
*/
func ParamsList1() pg.Parser {
	return pg.Trim(
		pg.TryAny(
			pg.Concat(
				pg.Recursive(
					"Expression",
					Expression),
				pg.Whitespaces(),
				pg.SkipChar(','),
				pg.Whitespaces(),
				pg.Recursive(
					"ParamsList1",
					ParamsList1)),
			pg.Recursive(
				"Expression",
				Expression)))
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
			pg.Recursive(
				"Block",
				Block)))
}

/*
	NamedParamsList  ←
		  NamedParamsList1
		| 'empty'
*/
func NamedParamsList() pg.Parser {
	return pg.Trim(
		pg.TryAny(
			NamedParamsList1(),
			pg.Empty()))
}

/*
	NamedParamsList1  ←
		  Identifier ',' NamedParamsList1
		| Identifier
*/
func NamedParamsList1() pg.Parser {
	return pg.Trim(
		pg.TryAny(
			pg.Concat(
				Identifier(),
				pg.Whitespaces(),
				pg.SkipChar(','),
				pg.Whitespaces(),
				pg.Recursive(
					"NamedParamsList1",
					NamedParamsList1)),
			Identifier()))
}

/*
	Program  ←  Statement*
*/
func Program() pg.Parser {
	return pg.Specify(PROGRAM,
		pg.Many(
			Statement()))
}

/*

*/
func main() {
	in := pg.InitParser()
	in.SetInput(`
		func callMe(a, b) {
			return a == b
		}
		if test == false {
			test = callMe(0, 0)
		}
		for i = 0, k = 2; test; i = i + 1 {
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
		switch kind {
			case "3": {
				kind = "0"
			}
			case "4": {
				kind = "1"
			}
			else: {
				kind = "3"
			}
		}
		`)

	start := time.Now()
	out, ok := Program()(in)
	end := time.Now()

	for _, o := range out {
		o.Walk(0, printNode, none, nil)
	}
	fmt.Printf("Input length: %d, probe count: %d, total: %s\n", len(in.GetInput()), in.GetProbeCount(), end.Sub(start).String())
	fmt.Printf("Parse ok: %t\n", ok)
	if len(in.GetInput())-in.GetPosition() > 0 {
		fmt.Printf("Early stop at line: %d\n", in.GetLineCount())
	}
}

func printNode(level int, node *pt.ParseTree, e interface{}) {
	for i := 0; i < level; i += 1 {
		fmt.Print("|  ")
	}
	fmt.Printf("%s [%s]\n", NODE_TYPES[node.Type], node.Value)
}

func none(level int, node *pt.ParseTree, e interface{}) {}
