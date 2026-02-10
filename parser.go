package main

import (
	"fmt"
	"strconv"
)

// Operator precedence levels
const (
	_ int = iota
	LOWEST
	ASSIGN_PREC  // ==
	OR_PREC      // or
	AND_PREC     // and
	IS_PREC      // is
	COMPARE_PREC // >, <, >=, <=
	SUM_PREC     // +, -
	PRODUCT_PREC // *, /, %
	PREFIX_PREC  // not, -
	CALL_PREC    // .
	INDEX_PREC   // [
)

var precedences = map[TokenType]int{
	ASSIGN_MUT: ASSIGN_PREC,
	OR:         OR_PREC,
	AND:        AND_PREC,
	IS:         IS_PREC,
	GT:         COMPARE_PREC,
	LT:         COMPARE_PREC,
	GTE:        COMPARE_PREC,
	LTE:        COMPARE_PREC,
	PLUS:       SUM_PREC,
	MINUS:      SUM_PREC,
	MULTIPLY:   PRODUCT_PREC,
	DIVIDE:     PRODUCT_PREC,
	MODULO:     PRODUCT_PREC,
	LPAREN:     CALL_PREC,
	DOT:        CALL_PREC,
	LBRACKET:   INDEX_PREC,
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

// Parser parses MoonShot source code into an AST
type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	errors    []string

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

// NewParser creates a new Parser
func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(INTEGER, p.parseIntegerLiteral)
	p.registerPrefix(FLOAT, p.parseFloatLiteral)
	p.registerPrefix(STRING, p.parseStringLiteral)
	p.registerPrefix(TRUE, p.parseBooleanLiteral)
	p.registerPrefix(FALSE, p.parseBooleanLiteral)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(NOT, p.parsePrefixExpression)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)
	p.registerPrefix(LBRACKET, p.parseListLiteral)
	p.registerPrefix(LBRACE, p.parseBraceExpression)
	p.registerPrefix(IF, p.parseIfExpression)
	p.registerPrefix(SOME, p.parseOptionExpression)
	p.registerPrefix(NONE, p.parseOptionExpression)
	p.registerPrefix(OK, p.parseResultExpression)
	p.registerPrefix(ERROR, p.parseResultExpression)
	p.registerPrefix(MATCH, p.parseMatchExpression)
	p.registerPrefix(MUTABLE, p.parseMutableExpression)

	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(MULTIPLY, p.parseInfixExpression)
	p.registerInfix(DIVIDE, p.parseInfixExpression)
	p.registerInfix(MODULO, p.parseInfixExpression)
	p.registerInfix(GT, p.parseInfixExpression)
	p.registerInfix(LT, p.parseInfixExpression)
	p.registerInfix(GTE, p.parseInfixExpression)
	p.registerInfix(LTE, p.parseInfixExpression)
	p.registerInfix(AND, p.parseInfixExpression)
	p.registerInfix(OR, p.parseInfixExpression)
	p.registerInfix(IS, p.parseInfixExpression)
	p.registerInfix(LPAREN, p.parseCallExpression)
	p.registerInfix(DOT, p.parseMemberExpression)
	p.registerInfix(LBRACKET, p.parseIndexExpression)
	p.registerInfix(ASSIGN_MUT, p.parseAssignmentExpression)

	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("line %d: expected next token to be %s, got %s instead",
		p.peekToken.Line, t.String(), p.peekToken.Type.String())
	p.errors = append(p.errors, msg)
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() *Program {
	program := &Program{Statements: []Statement{}}

	for !p.curTokenIs(EOF) {
		p.skipNewlines()
		if p.curTokenIs(EOF) {
			break
		}
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) skipNewlines() {
	for p.curTokenIs(NEWLINE) {
		p.nextToken()
	}
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case DEF:
		return p.parseDefStatement()
	case FUN:
		return p.parseFunctionStatement()
	case RETURN:
		return p.parseReturnStatement()
	case IF:
		return p.parseIfStatement()
	case WHILE:
		return p.parseWhileStatement()
	case FOR:
		return p.parseForStatement()
	case BREAK:
		return &BreakStatement{Token: p.curToken}
	case CONTINUE:
		return &ContinueStatement{Token: p.curToken}
	case STRUCT:
		return p.parseStructStatement()
	case EXTEND:
		return p.parseExtendStatement()
	case IMPORT:
		return p.parseImportStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDefStatement() *DefStatement {
	stmt := &DefStatement{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Optional type hint
	if p.peekTokenIs(COLON) {
		p.nextToken() // consume ':'
		p.nextToken() // move to type
		stmt.TypeHint = p.parseTypeAnnotation()
	}

	if !p.expectPeek(ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseTypeAnnotation() *TypeAnnotation {
	ta := &TypeAnnotation{Token: p.curToken, Name: p.curToken.Literal}

	// Check for type parameters like List[Integer]
	if p.peekTokenIs(LBRACKET) {
		p.nextToken() // consume '['
		p.nextToken() // move to first type param

		for !p.curTokenIs(RBRACKET) {
			param := p.parseTypeAnnotation()
			ta.TypeParams = append(ta.TypeParams, param)

			if p.peekTokenIs(COMMA) {
				p.nextToken()
				p.nextToken()
			} else if p.peekTokenIs(RBRACKET) {
				p.nextToken()
			} else {
				break
			}
		}
	}

	return ta
}

func (p *Parser) parseFunctionStatement() *FunctionStatement {
	stmt := &FunctionStatement{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	// Optional return type
	if p.peekTokenIs(ARROW) {
		p.nextToken() // consume '->'
		p.nextToken() // move to type
		stmt.ReturnType = p.parseTypeAnnotation()
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFunctionParameters() []*FunctionParameter {
	params := []*FunctionParameter{}

	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	param := &FunctionParameter{
		Name: &Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}

	// Optional type hint
	if p.peekTokenIs(COLON) {
		p.nextToken()
		p.nextToken()
		param.TypeHint = p.parseTypeAnnotation()
	}

	params = append(params, param)

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()

		param := &FunctionParameter{
			Name: &Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}

		if p.peekTokenIs(COLON) {
			p.nextToken()
			p.nextToken()
			param.TypeHint = p.parseTypeAnnotation()
		}

		params = append(params, param)
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}

	p.nextToken()

	if !p.curTokenIs(NEWLINE) && !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		stmt.Value = p.parseExpression(LOWEST)
	}

	return stmt
}

func (p *Parser) parseIfStatement() Statement {
	expr := p.parseIfExpression()
	return &ExpressionStatement{Token: expr.(*IfExpression).Token, Expression: expr}
}

func (p *Parser) parseWhileStatement() *WhileStatement {
	stmt := &WhileStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatement() *ForStatement {
	stmt := &ForStatement{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Variable = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(IN) {
		return nil
	}

	p.nextToken()
	stmt.Iterable = p.parseExpression(LOWEST)

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseStructStatement() *StructStatement {
	stmt := &StructStatement{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Fields = p.parseStructFields()

	return stmt
}

func (p *Parser) parseStructFields() []*StructField {
	fields := []*StructField{}

	p.nextToken()
	p.skipNewlines()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		field := &StructField{
			Name: &Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}

		if p.peekTokenIs(COLON) {
			p.nextToken()
			p.nextToken()
			field.TypeHint = p.parseTypeAnnotation()
		}

		fields = append(fields, field)

		p.nextToken()
		if p.curTokenIs(COMMA) || p.curTokenIs(NEWLINE) {
			p.nextToken()
		}
		p.skipNewlines()
	}

	return fields
}

func (p *Parser) parseExtendStatement() *ExtendStatement {
	stmt := &ExtendStatement{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.TypeName = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	p.nextToken()
	p.skipNewlines()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		if p.curTokenIs(FUN) {
			method := p.parseFunctionStatement()
			if method != nil {
				stmt.Methods = append(stmt.Methods, method)
			}
		}
		p.nextToken()
		p.skipNewlines()
	}

	return stmt
}

func (p *Parser) parseImportStatement() *ImportStatement {
	stmt := &ImportStatement{Token: p.curToken}

	p.nextToken()

	stmt.Path = append(stmt.Path, p.curToken.Literal)

	for p.peekTokenIs(DOT) {
		p.nextToken()
		p.nextToken()
		stmt.Path = append(stmt.Path, p.curToken.Literal)
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("line %d: no prefix parse function for %s found",
			p.curToken.Line, p.curToken.Type.String()))
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(NEWLINE) && !p.peekTokenIs(EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() Expression {
	ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check if this is a struct literal: StructName { field: value }
	if p.peekTokenIs(LBRACE) {
		// Need to look ahead to see if it's a struct literal
		return p.maybeParseStructLiteral(ident)
	}

	return ident
}

func (p *Parser) maybeParseStructLiteral(name *Identifier) Expression {
	// Peek to see if the { is followed by IDENT :
	// Save current position - we need to commit or backtrack
	// For simplicity, we'll assume uppercase identifiers followed by { are struct literals
	if len(name.Value) > 0 && name.Value[0] >= 'A' && name.Value[0] <= 'Z' && p.peekTokenIs(LBRACE) {
		p.nextToken() // consume '{'
		return p.parseStructLiteralBody(name)
	}
	return name
}

func (p *Parser) parseStructLiteralBody(name *Identifier) Expression {
	lit := &StructLiteral{
		Token:      p.curToken,
		StructName: name,
		Fields:     make(map[string]Expression),
	}

	p.nextToken()
	p.skipNewlines()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		fieldName := p.curToken.Literal

		if !p.expectPeek(COLON) {
			return nil
		}

		p.nextToken()
		lit.Fields[fieldName] = p.parseExpression(LOWEST)

		p.nextToken()
		if p.curTokenIs(COMMA) || p.curTokenIs(NEWLINE) {
			p.nextToken()
		}
		p.skipNewlines()
	}

	return lit
}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("line %d: could not parse %q as integer",
			p.curToken.Line, p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() Expression {
	lit := &FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("line %d: could not parse %q as float",
			p.curToken.Line, p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() Expression {
	return &BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(TRUE)}
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX_PREC)

	return expression
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseAssignmentExpression(left Expression) Expression {
	ident, ok := left.(*Identifier)
	if !ok {
		p.errors = append(p.errors, fmt.Sprintf("line %d: left side of == must be an identifier",
			p.curToken.Line))
		return nil
	}

	expression := &AssignmentExpression{
		Token: p.curToken,
		Name:  ident,
	}

	p.nextToken()
	expression.Value = p.parseExpression(LOWEST)

	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseListLiteral() Expression {
	list := &ListLiteral{Token: p.curToken}
	list.Elements = p.parseExpressionList(RBRACKET)
	return list
}

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	list := []Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// parseBraceExpression handles { which could be:
// - Lambda: { x -> x * 2 }
// - Map literal: { "key": value }
// - Block statement (in certain contexts)
func (p *Parser) parseBraceExpression() Expression {
	token := p.curToken

	// Peek ahead to determine what kind of expression this is
	p.nextToken()
	p.skipNewlines()

	// Empty map/block
	if p.curTokenIs(RBRACE) {
		return &MapLiteral{Token: token, Pairs: make(map[Expression]Expression)}
	}

	// Check for lambda: identifier followed by ->
	if p.curTokenIs(IDENT) && p.peekTokenIs(ARROW) {
		return p.parseLambdaWithFirstParam(token)
	}

	// Check for lambda with multiple params: identifier followed by comma
	if p.curTokenIs(IDENT) && p.peekTokenIs(COMMA) {
		return p.parseLambdaMultiParam(token)
	}

	// Otherwise it's a map literal
	return p.parseMapLiteralBody(token)
}

func (p *Parser) parseLambdaWithFirstParam(token Token) Expression {
	lambda := &FunctionLiteral{Token: token}

	param := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	lambda.Parameters = []*Identifier{param}

	p.nextToken() // consume ->
	p.nextToken() // move to body

	lambda.Body = p.parseExpression(LOWEST)

	if !p.expectPeek(RBRACE) {
		return nil
	}

	return lambda
}

func (p *Parser) parseLambdaMultiParam(token Token) Expression {
	lambda := &FunctionLiteral{Token: token}

	// First parameter
	param := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	lambda.Parameters = []*Identifier{param}

	// Additional parameters
	for p.peekTokenIs(COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next param

		param := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		lambda.Parameters = append(lambda.Parameters, param)
	}

	if !p.expectPeek(ARROW) {
		return nil
	}

	p.nextToken()
	lambda.Body = p.parseExpression(LOWEST)

	if !p.expectPeek(RBRACE) {
		return nil
	}

	return lambda
}

func (p *Parser) parseMapLiteralBody(token Token) Expression {
	ml := &MapLiteral{Token: token, Pairs: make(map[Expression]Expression)}

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		ml.Pairs[key] = value

		p.nextToken()
		if p.curTokenIs(COMMA) || p.curTokenIs(NEWLINE) {
			p.nextToken()
		}
		p.skipNewlines()
	}

	return ml
}

func (p *Parser) parseIfExpression() Expression {
	expression := &IfExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(ELSE) {
		p.nextToken()

		if !p.expectPeek(LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	p.nextToken()
	p.skipNewlines()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
		p.skipNewlines()
	}

	return block
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(RPAREN)
	return exp
}

func (p *Parser) parseMemberExpression(object Expression) Expression {
	exp := &MemberExpression{Token: p.curToken, Object: object}

	if !p.expectPeek(IDENT) {
		return nil
	}

	exp.Member = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check for .with { ... } syntax
	if exp.Member.Value == "with" && p.peekTokenIs(LBRACE) {
		return p.parseWithExpression(object)
	}

	return exp
}

func (p *Parser) parseWithExpression(object Expression) Expression {
	we := &WithExpression{Token: p.curToken, Object: object}
	we.Updates = make(map[string]Expression)

	p.nextToken() // consume 'with'
	p.nextToken() // consume '{'
	p.skipNewlines()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		fieldName := p.curToken.Literal

		if !p.expectPeek(COLON) {
			return nil
		}

		p.nextToken()
		we.Updates[fieldName] = p.parseExpression(LOWEST)

		p.nextToken()
		if p.curTokenIs(COMMA) || p.curTokenIs(NEWLINE) {
			p.nextToken()
		}
		p.skipNewlines()
	}

	return we
}

func (p *Parser) parseIndexExpression(left Expression) Expression {
	exp := &IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseOptionExpression() Expression {
	exp := &OptionExpression{Token: p.curToken}

	if p.curTokenIs(NONE) {
		exp.IsSome = false
		return exp
	}

	exp.IsSome = true

	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseResultExpression() Expression {
	exp := &ResultExpression{Token: p.curToken}
	exp.IsOk = p.curTokenIs(OK)

	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseMatchExpression() Expression {
	exp := &MatchExpression{Token: p.curToken}

	p.nextToken()
	exp.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(LBRACE) {
		return nil
	}

	p.nextToken()
	p.skipNewlines()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		matchCase := p.parseMatchCase()
		if matchCase != nil {
			exp.Cases = append(exp.Cases, matchCase)
		}
		p.nextToken()
		p.skipNewlines()
	}

	return exp
}

func (p *Parser) parseMatchCase() *MatchCase {
	mc := &MatchCase{}

	// Parse pattern: Some(x), None, Ok(x), Error(x)
	mc.Pattern = p.parseExpression(LOWEST)

	// Extract binding variable from pattern
	switch pat := mc.Pattern.(type) {
	case *OptionExpression:
		if pat.IsSome {
			if ident, ok := pat.Value.(*Identifier); ok {
				mc.BindingVar = ident
			}
		}
	case *ResultExpression:
		if ident, ok := pat.Value.(*Identifier); ok {
			mc.BindingVar = ident
		}
	}

	if !p.expectPeek(ARROW) {
		return nil
	}

	if !p.expectPeek(LBRACE) {
		// Single expression form
		p.nextToken()
		expr := p.parseExpression(LOWEST)
		mc.Body = &BlockStatement{
			Statements: []Statement{
				&ExpressionStatement{Expression: expr},
			},
		}
		return mc
	}

	mc.Body = p.parseBlockStatement()

	return mc
}

func (p *Parser) parseMutableExpression() Expression {
	exp := &MutableExpression{Token: p.curToken}

	// Optional type parameter: Mutable[Integer]
	if p.peekTokenIs(LBRACKET) {
		p.nextToken() // consume '['
		p.nextToken() // move to type
		exp.TypeHint = p.parseTypeAnnotation()
		if !p.expectPeek(RBRACKET) {
			return nil
		}
	}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return exp
}
