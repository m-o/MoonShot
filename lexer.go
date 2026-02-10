package main

// Lexer tokenizes MoonShot source code
type Lexer struct {
	input   string
	pos     int  // current position in input
	readPos int  // current reading position (after current char)
	ch      byte // current char under examination
	line    int  // current line number
	column  int  // current column number
}

// NewLexer creates a new Lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	l.column++
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	l.skipWhitespaceExceptNewline()

	var tok Token
	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '\n':
		tok = l.newToken(NEWLINE, string(l.ch))
		l.line++
		l.column = 0
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: ASSIGN_MUT, Literal: "==", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(ASSIGN, string(l.ch))
		}
	case '+':
		tok = l.newToken(PLUS, string(l.ch))
	case '-':
		if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Type: ARROW, Literal: "->", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(MINUS, string(l.ch))
		}
	case '*':
		tok = l.newToken(MULTIPLY, string(l.ch))
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken()
		}
		tok = l.newToken(DIVIDE, string(l.ch))
	case '%':
		tok = l.newToken(MODULO, string(l.ch))
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: GTE, Literal: ">=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(GT, string(l.ch))
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: LTE, Literal: "<=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(LT, string(l.ch))
		}
	case '(':
		tok = l.newToken(LPAREN, string(l.ch))
	case ')':
		tok = l.newToken(RPAREN, string(l.ch))
	case '{':
		tok = l.newToken(LBRACE, string(l.ch))
	case '}':
		tok = l.newToken(RBRACE, string(l.ch))
	case '[':
		tok = l.newToken(LBRACKET, string(l.ch))
	case ']':
		tok = l.newToken(RBRACKET, string(l.ch))
	case ',':
		tok = l.newToken(COMMA, string(l.ch))
	case ':':
		tok = l.newToken(COLON, string(l.ch))
	case '.':
		tok = l.newToken(DOT, string(l.ch))
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Type = EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal, tok.Type = l.readNumber()
			return tok
		} else {
			tok = l.newToken(ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) newToken(tokenType TokenType, literal string) Token {
	return Token{Type: tokenType, Literal: literal, Line: l.line, Column: l.column}
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() (string, TokenType) {
	pos := l.pos
	tokenType := INTEGER

	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for float
	if l.ch == '.' && isDigit(l.peekChar()) {
		tokenType = FLOAT
		l.readChar() // consume the '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[pos:l.pos], tokenType
}

func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	pos := l.pos

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' && l.peekChar() != 0 {
			l.readChar() // skip escape char
		}
		l.readChar()
	}

	return l.input[pos:l.pos]
}

func (l *Lexer) skipWhitespaceExceptNewline() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
