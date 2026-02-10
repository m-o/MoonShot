package main

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	NEWLINE

	// Identifiers and literals
	IDENT   // variable names
	INTEGER // 123
	FLOAT   // 123.45
	STRING  // "hello"

	// Keywords
	DEF
	FUN
	STRUCT
	EXTEND
	IF
	ELSE
	WHILE
	FOR
	IN
	RETURN
	MATCH
	SOME
	NONE
	OK
	ERROR
	IMPORT
	AND
	OR
	NOT
	IS
	BREAK
	CONTINUE
	MUTABLE
	TRUE
	FALSE

	// Operators
	ASSIGN     // =
	ASSIGN_MUT // ==
	PLUS       // +
	MINUS      // -
	MULTIPLY   // *
	DIVIDE     // /
	MODULO     // %
	GT         // >
	LT         // <
	GTE        // >=
	LTE        // <=
	ARROW      // ->

	// Delimiters
	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]
	COMMA    // ,
	COLON    // :
	DOT      // .
)

var tokenNames = map[TokenType]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	NEWLINE:    "NEWLINE",
	IDENT:      "IDENT",
	INTEGER:    "INTEGER",
	FLOAT:      "FLOAT",
	STRING:     "STRING",
	DEF:        "DEF",
	FUN:        "FUN",
	STRUCT:     "STRUCT",
	EXTEND:     "EXTEND",
	IF:         "IF",
	ELSE:       "ELSE",
	WHILE:      "WHILE",
	FOR:        "FOR",
	IN:         "IN",
	RETURN:     "RETURN",
	MATCH:      "MATCH",
	SOME:       "SOME",
	NONE:       "NONE",
	OK:         "OK",
	ERROR:      "ERROR",
	IMPORT:     "IMPORT",
	AND:        "AND",
	OR:         "OR",
	NOT:        "NOT",
	IS:         "IS",
	BREAK:      "BREAK",
	CONTINUE:   "CONTINUE",
	MUTABLE:    "MUTABLE",
	TRUE:       "TRUE",
	FALSE:      "FALSE",
	ASSIGN:     "=",
	ASSIGN_MUT: "==",
	PLUS:       "+",
	MINUS:      "-",
	MULTIPLY:   "*",
	DIVIDE:     "/",
	MODULO:     "%",
	GT:         ">",
	LT:         "<",
	GTE:        ">=",
	LTE:        "<=",
	ARROW:      "->",
	LPAREN:     "(",
	RPAREN:     ")",
	LBRACE:     "{",
	RBRACE:     "}",
	LBRACKET:   "[",
	RBRACKET:   "]",
	COMMA:      ",",
	COLON:      ":",
	DOT:        ".",
}

func (t TokenType) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return "UNKNOWN"
}

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Keywords maps keyword strings to token types
var keywords = map[string]TokenType{
	"def":      DEF,
	"fun":      FUN,
	"struct":   STRUCT,
	"extend":   EXTEND,
	"if":       IF,
	"else":     ELSE,
	"while":    WHILE,
	"for":      FOR,
	"in":       IN,
	"return":   RETURN,
	"match":    MATCH,
	"Some":     SOME,
	"None":     NONE,
	"Ok":       OK,
	"Error":    ERROR,
	"import":   IMPORT,
	"and":      AND,
	"or":       OR,
	"not":      NOT,
	"is":       IS,
	"break":    BREAK,
	"continue": CONTINUE,
	"Mutable":  MUTABLE,
	"true":     TRUE,
	"false":    FALSE,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
