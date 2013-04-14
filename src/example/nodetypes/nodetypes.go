package nt

const (
	TYPE_UNDEFINED = iota
	IDENTIFIER
	NUMBER_LITERAL
	STRING_LITERAL
	BOOL_LITERAL
	LITERAL
	ASSIGNMENT
	EXPRESSION
	ADD
	SUB
	MUL
	DIV
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
	N_E_COMPARISON
	MINUS
	NOT
	BREAK
	CONTINUE
	RETURN
	OR_EXPRESSION
	AND_EXPRESSION
	FUNCTION_CALL
	FUNCTION_DEFINITION
	VALUE_ACCESS
	PROGRAM
	EXPORT
)

var NODE_TYPES = map[int]string{
	TYPE_UNDEFINED: "?",

	IDENTIFIER:          "IDENTIFIER",
	NUMBER_LITERAL:      "NUMBER_LITERAL",
	STRING_LITERAL:      "STRING_LITERAL",
	BOOL_LITERAL:        "BOOL_LITERAL",
	ASSIGNMENT:          "ASSIGNMENT",
	EXPRESSION:          "EXPRESSION",
	ADD:                 "ADD",
	SUB:                 "SUB",
	MUL:                 "MUL",
	DIV:                 "DIV",
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
	N_E_COMPARISON:      "N_E_COMPARISON",
	MINUS:               "MINUS",
	NOT:                 "NOT",
	BREAK:               "BREAK",
	CONTINUE:            "CONTINUE",
	RETURN:              "RETURN",
	OR_EXPRESSION:       "OR_EXPRESSION",
	AND_EXPRESSION:      "AND_EXPRESSION",
	FUNCTION_CALL:       "FUNCTION_CALL",
	FUNCTION_DEFINITION: "FUNCTION_DEFINITION",
	PROGRAM:             "PROGRAM",
	VALUE_ACCESS:        "VALUE_ACCESS",
	EXPORT:              "EXPORT",
}
