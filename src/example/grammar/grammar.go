package grammar

import (
	"el/nodetypes"
	"parsego/parser"
)

/*
	Identifier  ←  [a-zA-Z][a-zA-Z0-9]*
*/
func Identifier() pg.Parser {
	return pg.Specify(nt.IDENTIFIER,
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
	return pg.Specify(nt.NUMBER_LITERAL,
		pg.Many1(
			pg.Number()))
}

/*
	StringLiteral  ←  '"' [a-zA-Z0-9]* '"'
*/
func StringLiteral() pg.Parser {
	return pg.Specify(nt.STRING_LITERAL,
		pg.Concat(
			pg.Skip(
				pg.Character('"')),
			pg.Many(
				pg.AnyCharBut('"')),
			pg.Skip(
				pg.Character('"'))))
}

/*
	BoolLiteral  ←  'true' | 'false'
*/
func BoolLiteral() pg.Parser {
	return pg.Specify(nt.BOOL_LITERAL,
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
	return pg.Specify(nt.ASSIGNMENT,
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
	return pg.Specify(nt.EXPRESSION,
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
		pg.Specify(nt.OR_EXPRESSION,
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
		pg.Specify(nt.AND_EXPRESSION,
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
		EComparison(),
		NEComparison(),
		LComparison(),
		LEComparison(),
		GComparison(),
		GEComparison(),
		Sum())
}

/*
	LComparison  ←  Sum '<' Sum
*/
func LComparison() pg.Parser {
	return pg.Specify(nt.L_COMPARISON,
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
	return pg.Specify(nt.L_E_COMPARISON,
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
	return pg.Specify(nt.G_COMPARISON,
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
	return pg.Specify(nt.G_E_COMPARISON,
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
	return pg.Specify(nt.E_COMPARISON,
		pg.Concat(
			Sum(),
			pg.Whitespaces(),
			pg.Skip(
				pg.String("==")),
			pg.Whitespaces(),
			Sum()))
}

/*
	NEComparison  ←  Sum '!=' Sum
*/
func NEComparison() pg.Parser {
	return pg.Specify(nt.N_E_COMPARISON,
		pg.Concat(
			Sum(),
			pg.Whitespaces(),
			pg.Skip(
				pg.String("!=")),
			pg.Whitespaces(),
			Sum()))
}

/*
	Sum ← Product ('+' | '-') Product | Product
*/
func Sum() pg.Parser {
	return pg.TryAny(
		Add(),
		Sub(),
		pg.Trim(
			Product()))
}

func Add() pg.Parser {
	return pg.Specify(nt.ADD,
		pg.Concat(
			Product(),
			pg.Trim(
				pg.Concat(
					pg.SkipChar('+'),
					pg.Whitespaces(),
					Product()))))
}

func Sub() pg.Parser {
	return pg.Specify(nt.SUB,
		pg.Concat(
			Product(),
			pg.Trim(
				pg.Concat(
					pg.SkipChar('-'),
					pg.Whitespaces(),
					Product()))))
}

/*
	Product ← Value ('*' | '/') Value | Value
*/
func Product() pg.Parser {
	return pg.TryAny(
		Mul(),
		Div(),
		pg.Trim(
			Value()))
}

func Mul() pg.Parser {
	return pg.Specify(nt.MUL,
		pg.Concat(
			Value(),
			pg.Trim(
				pg.Concat(
					pg.SkipChar('*'),
					pg.Whitespaces(),
					Value()))))
}

func Div() pg.Parser {
	return pg.Specify(nt.DIV,
		pg.Concat(
			Value(),
			pg.Trim(
				pg.Concat(
					pg.SkipChar('/'),
					pg.Whitespaces(),
					Value()))))
}

/*
	Value  ←
		FunctionCall | Literal | ValueAccess | Identifier | '(' Expr ')'
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
	return pg.Specify(nt.FOREACH,
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
	For  ←  'for' ForInit ';' ForCondition ';'
		ForStep Block
*/
func For() pg.Parser {
	return pg.Specify(nt.FOR,
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
	return pg.Specify(nt.FOR_INIT,
		AssignmentList())
}

/*
	ForCondition  ←  Expression
*/
func ForCondition() pg.Parser {
	return pg.Specify(nt.FOR_CONDITION,
		Expression())
}

/*
	ForStep  ←  AssignmentList
*/
func ForStep() pg.Parser {
	return pg.Specify(nt.FOR_STEP,
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
	return pg.Specify(nt.IFTHEN,
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
	return pg.Specify(nt.IFTHENELSE,
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
	return pg.Specify(nt.BREAK,
		pg.Skip(
			pg.String("break")))
}

/*
	Continue  ←  'continue'
*/
func Continue() pg.Parser {
	return pg.Specify(nt.CONTINUE,
		pg.Skip(
			pg.String("continue")))
}

/*
	Return  ←  'return' Expression | 'return' ';'
*/
func Return() pg.Parser {
	return pg.Specify(nt.RETURN,
		pg.Concat(
			pg.Skip(
				pg.String("return")),
			pg.TryAny(
				pg.Concat(
					pg.Whitespaces(),
					pg.SkipChar(';')),
				pg.Concat(
					pg.Whitespaces(),
					Expression()))))
}

/*
	Switch  ←  'switch' Expression SwitchBlock
*/
func Switch() pg.Parser {
	return pg.Specify(nt.SWITCH,
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
	return pg.Specify(nt.BLOCK,
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
	return pg.Specify(nt.CASE,
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
	return pg.Specify(nt.CASE_ELSE,
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
	return pg.Specify(nt.FUNCTION_CALL,
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
	return pg.Specify(nt.FUNCTION_DEFINITION,
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
	return pg.Specify(nt.PROGRAM,
		pg.Many(
			Statement()))
}
