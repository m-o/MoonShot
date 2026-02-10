# MoonShot Implementation Specification (Go)

## Overview

This document describes how to implement the MoonShot interpreter in Go.

### Architecture

```
┌─────────────┐
│   Source    │
│   Code      │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Lexer     │  Tokenize source code
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Parser    │  Build AST
└──────┬──────┘
       │
       ▼
┌─────────────┐
│Type Checker │  Validate types
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Interpreter │  Execute AST
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Output    │
└─────────────┘
```

---

## 1. Project Structure

```
moonshot/
  cmd/
    moonshot/
      main.go           # Entry point
  pkg/
    lexer/
      lexer.go          # Tokenizer
      token.go          # Token types
    parser/
      parser.go         # Parser
      ast.go            # AST node definitions
    types/
      types.go          # Type system
      checker.go        # Type checker
    runtime/
      value.go          # Runtime value representation
      environment.go    # Variable scopes
      builtins.go       # Built-in functions
    evaluator/
      evaluator.go      # Interpreter
      errors.go         # Error enrichment
    module/
      loader.go         # Module system
  tests/
    lexer_test.go
    parser_test.go
    evaluator_test.go
  examples/
    hello.moon
    user.moon
```

---

## 2. Token Types

```go
// pkg/lexer/token.go
package lexer

type TokenType int

const (
    // Literals
    INTEGER TokenType = iota
    FLOAT
    STRING
    TRUE
    FALSE
    
    // Identifiers
    IDENT
    
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
    BREAK
    CONTINUE
    RETURN
    MATCH
    SOME
    NONE
    OK
    ERROR
    IMPORT
    
    // Operators
    ASSIGN        // =
    ASSIGN_MUT    // ==
    IS            // is
    IS_NOT        // is not
    PLUS          // +
    MINUS         // -
    MULTIPLY      // *
    DIVIDE        // /
    MODULO        // %
    GT            // >
    LT            // <
    GTE           // >=
    LTE           // <=
    AND           // and
    OR            // or
    NOT           // not
    ARROW         // ->
    
    // Delimiters
    LPAREN        // (
    RPAREN        // )
    LBRACE        // {
    RBRACE        // }
    LBRACKET      // [
    RBRACKET      // ]
    COMMA         // ,
    COLON         // :
    DOT           // .
    
    // Special
    NEWLINE
    EOF
)

type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Column  int
}
```

---

## 3. Lexer Implementation

```go
// pkg/lexer/lexer.go
package lexer

type Lexer struct {
    input        string
    position     int    // current position
    readPosition int    // next position
    ch           byte   // current char
    line         int
    column       int
}

func New(input string) *Lexer {
    l := &Lexer{input: input, line: 1, column: 0}
    l.readChar()
    return l
}

func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }
    l.position = l.readPosition
    l.readPosition++
    l.column++
}

func (l *Lexer) NextToken() Token {
    var tok Token
    
    l.skipWhitespace()
    
    tok.Line = l.line
    tok.Column = l.column
    
    switch l.ch {
    case '=':
        if l.peekChar() == '=' {
            l.readChar()
            tok = Token{Type: ASSIGN_MUT, Literal: "=="}
        } else {
            tok = Token{Type: ASSIGN, Literal: "="}
        }
    case '+':
        tok = Token{Type: PLUS, Literal: "+"}
    case '-':
        if l.peekChar() == '>' {
            l.readChar()
            tok = Token{Type: ARROW, Literal: "->"}
        } else {
            tok = Token{Type: MINUS, Literal: "-"}
        }
    case '{':
        tok = Token{Type: LBRACE, Literal: "{"}
    case '}':
        tok = Token{Type: RBRACE, Literal: "}"}
    case '[':
        tok = Token{Type: LBRACKET, Literal: "["}
    case ']':
        tok = Token{Type: RBRACKET, Literal: "]"}
    case '(':
        tok = Token{Type: LPAREN, Literal: "("}
    case ')':
        tok = Token{Type: RPAREN, Literal: ")"}
    case ',':
        tok = Token{Type: COMMA, Literal: ","}
    case ':':
        tok = Token{Type: COLON, Literal: ":"}
    case '.':
        tok = Token{Type: DOT, Literal: "."}
    case '"':
        tok.Type = STRING
        tok.Literal = l.readString()
    case 0:
        tok = Token{Type: EOF, Literal: ""}
    default:
        if isLetter(l.ch) {
            tok.Literal = l.readIdentifier()
            tok.Type = lookupKeyword(tok.Literal)
            return tok
        } else if isDigit(l.ch) {
            return l.readNumber()
        }
    }
    
    l.readChar()
    return tok
}

func (l *Lexer) readIdentifier() string {
    position := l.position
    for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
        l.readChar()
    }
    return l.input[position:l.position]
}

func (l *Lexer) readNumber() Token {
    position := l.position
    isFloat := false
    
    for isDigit(l.ch) {
        l.readChar()
    }
    
    if l.ch == '.' && isDigit(l.peekChar()) {
        isFloat = true
        l.readChar() // consume '.'
        for isDigit(l.ch) {
            l.readChar()
        }
    }
    
    literal := l.input[position:l.position]
    if isFloat {
        return Token{Type: FLOAT, Literal: literal}
    }
    return Token{Type: INTEGER, Literal: literal}
}

func lookupKeyword(ident string) TokenType {
    keywords := map[string]TokenType{
        "def":      DEF,
        "fun":      FUN,
        "struct":   STRUCT,
        "extend":   EXTEND,
        "if":       IF,
        "else":     ELSE,
        "while":    WHILE,
        "for":      FOR,
        "in":       IN,
        "break":    BREAK,
        "continue": CONTINUE,
        "return":   RETURN,
        "match":    MATCH,
        "Some":     SOME,
        "None":     NONE,
        "Ok":       OK,
        "Error":    ERROR,
        "true":     TRUE,
        "false":    FALSE,
        "and":      AND,
        "or":       OR,
        "not":      NOT,
        "is":       IS,
        "import":   IMPORT,
    }
    
    if tok, ok := keywords[ident]; ok {
        return tok
    }
    return IDENT
}
```

---

## 4. AST Node Definitions

```go
// pkg/parser/ast.go
package parser

import "moonshot/pkg/lexer"

// Base node interface
type Node interface {
    TokenLiteral() string
    String() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

// Program root
type Program struct {
    Statements []Statement
}

// Statements
type DefStatement struct {
    Token lexer.Token  // 'def' token
    Name  *Identifier
    Value Expression
    Type  TypeAnnotation  // optional
}

type FunctionStatement struct {
    Token      lexer.Token  // 'fun' token
    Name       *Identifier
    Parameters []*Parameter
    ReturnType TypeAnnotation  // optional
    Body       *BlockStatement
}

type StructStatement struct {
    Token  lexer.Token  // 'struct' token
    Name   *Identifier
    Fields []*StructField
}

type ExtendStatement struct {
    Token   lexer.Token  // 'extend' token
    Type    *Identifier
    Methods []*FunctionStatement
}

type ReturnStatement struct {
    Token       lexer.Token  // 'return' token
    ReturnValue Expression
}

type IfStatement struct {
    Token       lexer.Token  // 'if' token
    Condition   Expression
    Consequence *BlockStatement
    Alternative *BlockStatement  // optional else
}

type WhileStatement struct {
    Token     lexer.Token  // 'while' token
    Condition Expression
    Body      *BlockStatement
}

type ForStatement struct {
    Token    lexer.Token  // 'for' token
    Variable *Identifier  // or multiple for map iteration
    Index    *Identifier  // optional
    Iterable Expression
    Body     *BlockStatement
}

type BreakStatement struct {
    Token lexer.Token  // 'break' token
}

type ContinueStatement struct {
    Token lexer.Token  // 'continue' token
}

type MatchStatement struct {
    Token    lexer.Token  // 'match' token
    Value    Expression
    Cases    []*MatchCase
}

type MatchCase struct {
    Pattern Expression
    Body    *BlockStatement
}

type BlockStatement struct {
    Token      lexer.Token  // '{' token
    Statements []Statement
}

type ExpressionStatement struct {
    Token      lexer.Token
    Expression Expression
}

// Expressions
type Identifier struct {
    Token lexer.Token
    Value string
}

type IntegerLiteral struct {
    Token lexer.Token
    Value int64
}

type FloatLiteral struct {
    Token lexer.Token
    Value float64
}

type StringLiteral struct {
    Token lexer.Token
    Value string
}

type BooleanLiteral struct {
    Token lexer.Token
    Value bool
}

type ListLiteral struct {
    Token    lexer.Token  // '[' token
    Elements []Expression
}

type MapLiteral struct {
    Token lexer.Token  // '{' token
    Pairs map[Expression]Expression
}

type StructLiteral struct {
    Token  lexer.Token  // struct name
    Type   *Identifier
    Fields map[string]Expression
}

type FunctionLiteral struct {
    Token      lexer.Token  // '{' for lambda
    Parameters []*Parameter
    Body       *BlockStatement
}

type CallExpression struct {
    Token     lexer.Token  // '(' token
    Function  Expression   // identifier or function literal
    Arguments []Expression
}

type MethodCallExpression struct {
    Token    lexer.Token  // '.' token
    Object   Expression
    Method   *Identifier
    Arguments []Expression
}

type MemberExpression struct {
    Token  lexer.Token  // '.' token
    Object Expression
    Member *Identifier
}

type IndexExpression struct {
    Token lexer.Token  // '[' token
    Left  Expression
    Index Expression
}

type InfixExpression struct {
    Token    lexer.Token
    Left     Expression
    Operator string
    Right    Expression
}

type PrefixExpression struct {
    Token    lexer.Token
    Operator string
    Right    Expression
}

type OptionExpression struct {
    Token lexer.Token  // 'Some' or 'None'
    Value Expression   // nil for None
}

type ResultExpression struct {
    Token lexer.Token  // 'Ok' or 'Error'
    Value Expression
}

type MutableExpression struct {
    Token lexer.Token  // 'Mutable'
    Type  TypeAnnotation
    Value Expression
}

type WithExpression struct {
    Token   lexer.Token  // '.with' token
    Object  Expression
    Updates map[string]Expression
}

// Type annotations
type TypeAnnotation interface {
    Node
    typeAnnotation()
}

type SimpleType struct {
    Token lexer.Token
    Name  string
}

type ListType struct {
    Token       lexer.Token  // '[' token
    ElementType TypeAnnotation
}

type MapType struct {
    Token     lexer.Token  // '{' token
    KeyType   TypeAnnotation
    ValueType TypeAnnotation
}

type OptionType struct {
    Token       lexer.Token  // 'Option' token
    ElementType TypeAnnotation
}

type ResultType struct {
    Token     lexer.Token  // 'Result' token
    ValueType TypeAnnotation
    ErrorType TypeAnnotation
}

type MutableType struct {
    Token       lexer.Token  // 'Mutable' token
    ElementType TypeAnnotation
}

// Helper structs
type Parameter struct {
    Name *Identifier
    Type TypeAnnotation
}

type StructField struct {
    Name *Identifier
    Type TypeAnnotation
}
```

---

## 5. Parser Implementation

```go
// pkg/parser/parser.go
package parser

import (
    "moonshot/pkg/lexer"
)

type Parser struct {
    l         *lexer.Lexer
    curToken  lexer.Token
    peekToken lexer.Token
    errors    []string
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{l: l}
    
    // Read two tokens to initialize curToken and peekToken
    p.nextToken()
    p.nextToken()
    
    return p
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
    program := &Program{}
    program.Statements = []Statement{}
    
    for p.curToken.Type != lexer.EOF {
        stmt := p.parseStatement()
        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.nextToken()
    }
    
    return program
}

func (p *Parser) parseStatement() Statement {
    switch p.curToken.Type {
    case lexer.DEF:
        return p.parseDefStatement()
    case lexer.FUN:
        return p.parseFunctionStatement()
    case lexer.STRUCT:
        return p.parseStructStatement()
    case lexer.EXTEND:
        return p.parseExtendStatement()
    case lexer.RETURN:
        return p.parseReturnStatement()
    case lexer.IF:
        return p.parseIfStatement()
    case lexer.WHILE:
        return p.parseWhileStatement()
    case lexer.FOR:
        return p.parseForStatement()
    case lexer.BREAK:
        return &BreakStatement{Token: p.curToken}
    case lexer.CONTINUE:
        return &ContinueStatement{Token: p.curToken}
    case lexer.MATCH:
        return p.parseMatchStatement()
    default:
        return p.parseExpressionStatement()
    }
}

func (p *Parser) parseDefStatement() *DefStatement {
    stmt := &DefStatement{Token: p.curToken}
    
    if !p.expectPeek(lexer.IDENT) {
        return nil
    }
    
    stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
    
    // Optional type annotation
    if p.peekTokenIs(lexer.COLON) {
        p.nextToken()
        p.nextToken()
        stmt.Type = p.parseTypeAnnotation()
    }
    
    if !p.expectPeek(lexer.ASSIGN) {
        return nil
    }
    
    p.nextToken()
    stmt.Value = p.parseExpression(LOWEST)
    
    return stmt
}

// ... implement other parse methods
```

---

## 6. Runtime Value Representation

```go
// pkg/runtime/value.go
package runtime

type ValueType int

const (
    INTEGER_VAL ValueType = iota
    FLOAT_VAL
    STRING_VAL
    BOOLEAN_VAL
    LIST_VAL
    MAP_VAL
    STRUCT_VAL
    FUNCTION_VAL
    OPTION_VAL
    RESULT_VAL
    MUTABLE_VAL
    NONE_VAL
    ERROR_VAL
)

type Value interface {
    Type() ValueType
    String() string
}

type IntegerValue struct {
    Value int64
}

type FloatValue struct {
    Value float64
}

type StringValue struct {
    Value string
}

type BooleanValue struct {
    Value bool
}

type ListValue struct {
    Elements []Value
}

type MapValue struct {
    Pairs map[string]Value  // keys are hashed to strings
}

type StructValue struct {
    Name   string
    Fields map[string]Value
}

type FunctionValue struct {
    Parameters []*parser.Parameter
    Body       *parser.BlockStatement
    Env        *Environment
}

type OptionValue struct {
    HasValue bool
    Value    Value  // nil if None
}

type ResultValue struct {
    IsOk  bool
    Value Value  // actual value if Ok, error message if Error
}

type MutableValue struct {
    Value Value  // the wrapped value that can be mutated
}

type ErrorValue struct {
    Method  string
    Input   string
    Message string
}

func (i *IntegerValue) Type() ValueType { return INTEGER_VAL }
func (f *FloatValue) Type() ValueType   { return FLOAT_VAL }
func (s *StringValue) Type() ValueType  { return STRING_VAL }
func (b *BooleanValue) Type() ValueType { return BOOLEAN_VAL }
func (l *ListValue) Type() ValueType    { return LIST_VAL }
func (m *MapValue) Type() ValueType     { return MAP_VAL }
func (s *StructValue) Type() ValueType  { return STRUCT_VAL }
func (f *FunctionValue) Type() ValueType { return FUNCTION_VAL }
func (o *OptionValue) Type() ValueType  { return OPTION_VAL }
func (r *ResultValue) Type() ValueType  { return RESULT_VAL }
func (m *MutableValue) Type() ValueType { return MUTABLE_VAL }
func (e *ErrorValue) Type() ValueType   { return ERROR_VAL }

// Auto-generated toString implementations
func (s *StructValue) String() string {
    // Generate: "StructName{field1: value1, field2: value2}"
    result := s.Name + "{"
    first := true
    for k, v := range s.Fields {
        if !first {
            result += ", "
        }
        result += k + ": " + v.String()
        first = false
    }
    result += "}"
    return result
}

func (e *ErrorValue) String() string {
    return "Error in " + e.Method + "\nInput: " + e.Input + "\nReason: " + e.Message
}
```

---

## 7. Environment (Scoping)

```go
// pkg/runtime/environment.go
package runtime

type Environment struct {
    store map[string]Value
    outer *Environment  // parent scope
}

func NewEnvironment() *Environment {
    return &Environment{
        store: make(map[string]Value),
        outer: nil,
    }
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
    env := NewEnvironment()
    env.outer = outer
    return env
}

func (e *Environment) Get(name string) (Value, bool) {
    val, ok := e.store[name]
    if !ok && e.outer != nil {
        return e.outer.Get(name)
    }
    return val, ok
}

func (e *Environment) Set(name string, val Value) Value {
    e.store[name] = val
    return val
}

func (e *Environment) Update(name string, val Value) bool {
    if _, ok := e.store[name]; ok {
        e.store[name] = val
        return true
    }
    if e.outer != nil {
        return e.outer.Update(name, val)
    }
    return false
}
```

---

## 8. Evaluator Implementation

```go
// pkg/evaluator/evaluator.go
package evaluator

import (
    "moonshot/pkg/parser"
    "moonshot/pkg/runtime"
)

func Eval(node parser.Node, env *runtime.Environment) runtime.Value {
    switch node := node.(type) {
    
    // Program
    case *parser.Program:
        return evalProgram(node, env)
    
    // Statements
    case *parser.DefStatement:
        val := Eval(node.Value, env)
        if isError(val) {
            return val
        }
        env.Set(node.Name.Value, val)
        return val
    
    case *parser.ExpressionStatement:
        return Eval(node.Expression, env)
    
    case *parser.ReturnStatement:
        val := Eval(node.ReturnValue, env)
        if isError(val) {
            return val
        }
        return &runtime.ReturnValue{Value: val}
    
    case *parser.BlockStatement:
        return evalBlockStatement(node, env)
    
    case *parser.IfStatement:
        return evalIfStatement(node, env)
    
    case *parser.WhileStatement:
        return evalWhileStatement(node, env)
    
    case *parser.ForStatement:
        return evalForStatement(node, env)
    
    // Expressions
    case *parser.IntegerLiteral:
        return &runtime.IntegerValue{Value: node.Value}
    
    case *parser.FloatLiteral:
        return &runtime.FloatValue{Value: node.Value}
    
    case *parser.StringLiteral:
        return &runtime.StringValue{Value: node.Value}
    
    case *parser.BooleanLiteral:
        return &runtime.BooleanValue{Value: node.Value}
    
    case *parser.Identifier:
        return evalIdentifier(node, env)
    
    case *parser.ListLiteral:
        return evalListLiteral(node, env)
    
    case *parser.MapLiteral:
        return evalMapLiteral(node, env)
    
    case *parser.StructLiteral:
        return evalStructLiteral(node, env)
    
    case *parser.InfixExpression:
        left := Eval(node.Left, env)
        if isError(left) {
            return left
        }
        right := Eval(node.Right, env)
        if isError(right) {
            return right
        }
        return evalInfixExpression(node.Operator, left, right)
    
    case *parser.CallExpression:
        return evalCallExpression(node, env)
    
    case *parser.MethodCallExpression:
        return evalMethodCallExpression(node, env)
    
    case *parser.MemberExpression:
        return evalMemberExpression(node, env)
    
    case *parser.MutableExpression:
        val := Eval(node.Value, env)
        if isError(val) {
            return val
        }
        return &runtime.MutableValue{Value: val}
    
    case *parser.OptionExpression:
        if node.Token.Type == lexer.NONE {
            return &runtime.OptionValue{HasValue: false}
        }
        val := Eval(node.Value, env)
        if isError(val) {
            return val
        }
        return &runtime.OptionValue{HasValue: true, Value: val}
    
    case *parser.ResultExpression:
        val := Eval(node.Value, env)
        if isError(val) {
            return val
        }
        isOk := node.Token.Type == lexer.OK
        return &runtime.ResultValue{IsOk: isOk, Value: val}
    
    case *parser.MatchStatement:
        return evalMatchStatement(node, env)
    }
    
    return nil
}

func evalMethodCallExpression(node *parser.MethodCallExpression, env *runtime.Environment) runtime.Value {
    object := Eval(node.Object, env)
    if isError(object) {
        return object
    }
    
    // Handle Result auto-unwrapping and short-circuiting
    if result, ok := object.(*runtime.ResultValue); ok {
        if !result.IsOk {
            // Propagate error
            return result
        }
        // Unwrap and call method on the value
        object = result.Value
    }
    
    // Handle Mutable implicit get
    if mut, ok := object.(*runtime.MutableValue); ok {
        object = mut.Value
    }
    
    // Look up method (from extensions or built-ins)
    method := lookupMethod(object, node.Method.Value)
    if method == nil {
        return newError("method not found: " + node.Method.Value)
    }
    
    // Evaluate arguments
    args := []runtime.Value{object}  // receiver is first arg
    for _, arg := range node.Arguments {
        val := Eval(arg, env)
        if isError(val) {
            return val
        }
        args = append(args, val)
    }
    
    // Call method
    result := applyFunction(method, args, env)
    
    // Enrich error if Result type
    if errVal, ok := result.(*runtime.ResultValue); ok && !errVal.IsOk {
        return enrichError(errVal, node.Method.Value, object)
    }
    
    return result
}

func enrichError(result *runtime.ResultValue, methodName string, input runtime.Value) runtime.Value {
    // Extract error message
    errorMsg := ""
    if strVal, ok := result.Value.(*runtime.StringValue); ok {
        errorMsg = strVal.Value
    }
    
    // Create enriched error
    return &runtime.ErrorValue{
        Method:  methodName,
        Input:   input.String(),  // uses auto-generated toString
        Message: errorMsg,
    }
}
```

---

## 9. Built-in Functions and Methods

```go
// pkg/runtime/builtins.go
package runtime

var Builtins = map[string]*FunctionValue{
    "print": {
        Builtin: func(args ...Value) Value {
            for _, arg := range args {
                fmt.Print(arg.String())
            }
            return nil
        },
    },
    "println": {
        Builtin: func(args ...Value) Value {
            for _, arg := range args {
                fmt.Print(arg.String())
            }
            fmt.Println()
            return nil
        },
    },
    "range": {
        Builtin: func(args ...Value) Value {
            if len(args) != 2 {
                return newError("range requires 2 arguments")
            }
            start := args[0].(*IntegerValue).Value
            end := args[1].(*IntegerValue).Value
            
            elements := []Value{}
            for i := start; i < end; i++ {
                elements = append(elements, &IntegerValue{Value: i})
            }
            return &ListValue{Elements: elements}
        },
    },
}

// List methods
var ListMethods = map[string]func(Value, ...Value) Value{
    "length": func(list Value, args ...Value) Value {
        l := list.(*ListValue)
        return &IntegerValue{Value: int64(len(l.Elements))}
    },
    "get": func(list Value, args ...Value) Value {
        l := list.(*ListValue)
        index := args[0].(*IntegerValue).Value
        if index < 0 || index >= int64(len(l.Elements)) {
            return &OptionValue{HasValue: false}
        }
        return l.Elements[index]
    },
    "append": func(list Value, args ...Value) Value {
        l := list.(*ListValue)
        newElements := make([]Value, len(l.Elements)+1)
        copy(newElements, l.Elements)
        newElements[len(l.Elements)] = args[0]
        return &ListValue{Elements: newElements}
    },
    "map": func(list Value, args ...Value) Value {
        l := list.(*ListValue)
        fn := args[0].(*FunctionValue)
        
        newElements := make([]Value, len(l.Elements))
        for i, elem := range l.Elements {
            newElements[i] = applyFunction(fn, []Value{elem}, nil)
        }
        return &ListValue{Elements: newElements}
    },
    "filter": func(list Value, args ...Value) Value {
        l := list.(*ListValue)
        fn := args[0].(*FunctionValue)
        
        newElements := []Value{}
        for _, elem := range l.Elements {
            result := applyFunction(fn, []Value{elem}, nil)
            if isTruthy(result) {
                newElements = append(newElements, elem)
            }
        }
        return &ListValue{Elements: newElements}
    },
    // ... implement other list methods
}

// Map methods
var MapMethods = map[string]func(Value, ...Value) Value{
    "get": func(m Value, args ...Value) Value {
        mapVal := m.(*MapValue)
        key := args[0].String()  // hash key to string
        if val, ok := mapVal.Pairs[key]; ok {
            return &OptionValue{HasValue: true, Value: val}
        }
        return &OptionValue{HasValue: false}
    },
    "insert": func(m Value, args ...Value) Value {
        mapVal := m.(*MapValue)
        key := args[0].String()
        val := args[1]
        
        newPairs := make(map[string]Value)
        for k, v := range mapVal.Pairs {
            newPairs[k] = v
        }
        newPairs[key] = val
        return &MapValue{Pairs: newPairs}
    },
    // ... implement other map methods
}
```

---

## 10. Type Checker

```go
// pkg/types/checker.go
package types

import "moonshot/pkg/parser"

type TypeChecker struct {
    env *TypeEnvironment
}

type TypeEnvironment struct {
    store map[string]Type
    outer *TypeEnvironment
}

func (tc *TypeChecker) Check(program *parser.Program) error {
    for _, stmt := range program.Statements {
        if err := tc.checkStatement(stmt); err != nil {
            return err
        }
    }
    return nil
}

func (tc *TypeChecker) checkStatement(stmt parser.Statement) error {
    switch stmt := stmt.(type) {
    case *parser.DefStatement:
        return tc.checkDefStatement(stmt)
    case *parser.FunctionStatement:
        return tc.checkFunctionStatement(stmt)
    // ... check other statements
    }
    return nil
}

func (tc *TypeChecker) checkMutableReturn(fn *parser.FunctionStatement) error {
    // Ensure function doesn't return Mutable type
    if isMutableType(fn.ReturnType) {
        return fmt.Errorf("function '%s' cannot return Mutable type", fn.Name.Value)
    }
    return nil
}

func isMutableType(t parser.TypeAnnotation) bool {
    if mt, ok := t.(*parser.MutableType); ok {
        return true
    }
    return false
}
```

---

## 11. Module Loader

```go
// pkg/module/loader.go
package module

import (
    "os"
    "path/filepath"
)

type ModuleLoader struct {
    cache map[string]*parser.Program
    basePath string
}

func NewLoader(basePath string) *ModuleLoader {
    return &ModuleLoader{
        cache: make(map[string]*parser.Program),
        basePath: basePath,
    }
}

func (ml *ModuleLoader) Load(moduleName string) (*parser.Program, error) {
    // Check cache
    if prog, ok := ml.cache[moduleName]; ok {
        return prog, nil
    }
    
    // Construct file path
    filePath := filepath.Join(ml.basePath, moduleName + ".moon")
    
    // Read file
    content, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    
    // Lex and parse
    l := lexer.New(string(content))
    p := parser.New(l)
    prog := p.ParseProgram()
    
    // Cache
    ml.cache[moduleName] = prog
    
    return prog, nil
}
```

---

## 12. Main Entry Point

```go
// cmd/moonshot/main.go
package main

import (
    "fmt"
    "os"
    "moonshot/pkg/lexer"
    "moonshot/pkg/parser"
    "moonshot/pkg/evaluator"
    "moonshot/pkg/runtime"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: moonshot <file.moon>")
        os.Exit(1)
    }
    
    filename := os.Args[1]
    
    // Read file
    content, err := os.ReadFile(filename)
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }
    
    // Lex
    l := lexer.New(string(content))
    
    // Parse
    p := parser.New(l)
    program := p.ParseProgram()
    
    if len(p.Errors()) > 0 {
        for _, err := range p.Errors() {
            fmt.Printf("Parse error: %s\n", err)
        }
        os.Exit(1)
    }
    
    // Type check
    tc := types.NewTypeChecker()
    if err := tc.Check(program); err != nil {
        fmt.Printf("Type error: %v\n", err)
        os.Exit(1)
    }
    
    // Evaluate
    env := runtime.NewEnvironment()
    result := evaluator.Eval(program, env)
    
    if result != nil && evaluator.IsError(result) {
        fmt.Printf("Runtime error: %s\n", result.String())
        os.Exit(1)
    }
}
```

---

## 13. Testing Strategy

### Lexer Tests
```go
// tests/lexer_test.go
func TestLexer(t *testing.T) {
    input := `def x = 5`
    
    tests := []struct {
        expectedType    lexer.TokenType
        expectedLiteral string
    }{
        {lexer.DEF, "def"},
        {lexer.IDENT, "x"},
        {lexer.ASSIGN, "="},
        {lexer.INTEGER, "5"},
    }
    
    l := lexer.New(input)
    for i, tt := range tests {
        tok := l.NextToken()
        if tok.Type != tt.expectedType {
            t.Fatalf("tests[%d] - wrong token type. expected=%q, got=%q",
                i, tt.expectedType, tok.Type)
        }
    }
}
```

### Parser Tests
```go
func TestDefStatement(t *testing.T) {
    input := `def x = 5`
    
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()
    
    if len(program.Statements) != 1 {
        t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
    }
    
    stmt, ok := program.Statements[0].(*parser.DefStatement)
    if !ok {
        t.Fatalf("statement is not DefStatement. got=%T", program.Statements[0])
    }
    
    if stmt.Name.Value != "x" {
        t.Errorf("name not 'x'. got=%s", stmt.Name.Value)
    }
}
```

### Evaluator Tests
```go
func TestEvalIntegerExpression(t *testing.T) {
    tests := []struct {
        input    string
        expected int64
    }{
        {"5", 5},
        {"10", 10},
        {"5 + 5", 10},
        {"5 * 2", 10},
    }
    
    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testIntegerObject(t, evaluated, tt.expected)
    }
}
```

---

## 14. Performance Optimizations

### Phase 1: Basic Interpreter
- Direct AST interpretation
- Simple environment lookups
- No optimizations

### Phase 2: Optimizations
1. **Constant folding** - Evaluate constant expressions at parse time
2. **Tail call optimization** - Reuse stack frames for tail recursion
3. **Inline caching** - Cache method lookups
4. **Copy-on-write** - For immutable collections

### Phase 3: JIT Compilation
1. **Bytecode compilation** - Compile AST to bytecode
2. **Bytecode interpreter** - Faster than AST walking
3. **JIT for hot paths** - Compile frequently-executed code to machine code

---

## 15. Implementation Phases

### Phase 1: Core Language (MVP)
- [x] Lexer
- [x] Parser (basic)
- [x] AST
- [x] Evaluator (basic)
- [x] Variables and functions
- [x] Primitives (Integer, String, Boolean)

### Phase 2: Collections
- [ ] List type and methods
- [ ] Map type and methods
- [ ] Immutability

### Phase 3: Type System
- [ ] Type annotations
- [ ] Type checker
- [ ] Option type
- [ ] Result type
- [ ] Mutable type validation

### Phase 4: Control Flow
- [ ] If/else
- [ ] While loops
- [ ] For loops
- [ ] Break/continue
- [ ] Match expressions

### Phase 5: Structs and Extensions
- [ ] Struct definitions
- [ ] Struct literals
- [ ] Extension functions
- [ ] Auto toString

### Phase 6: Error Enrichment
- [ ] Auto error context
- [ ] Method name tracking
- [ ] Input value capture

### Phase 7: Module System
- [ ] Import statements
- [ ] Module loader
- [ ] File-based modules

### Phase 8: Optimizations
- [ ] Performance profiling
- [ ] Optimize hot paths
- [ ] JIT compilation (future)

---

## 16. Build and Run

```bash
# Build
go build -o moonshot cmd/moonshot/main.go

# Run
./moonshot examples/hello.moon

# Test
go test ./...

# Benchmark
go test -bench=. ./...
```

---

## Summary

This implementation spec provides:
1. **Complete architecture** - Lexer → Parser → Type Checker → Evaluator
2. **Data structures** - AST nodes, runtime values, environments
3. **Core algorithms** - Parsing, type checking, evaluation
4. **Error handling** - Auto-enrichment with context
5. **Module system** - File-based imports
6. **Testing strategy** - Unit tests for each component
7. **Optimization path** - From interpreter to JIT

Start with Phase 1 (MVP), then incrementally add features!

**MoonShot Implementation - Ready to Build!** 🚀
