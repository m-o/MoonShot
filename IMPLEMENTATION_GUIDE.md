# MoonShot Implementation Guide

## Quick Start

You're building a **statically-typed interpreter** for MoonShot in Go.

**What you're building:**
```
.moon file → Type check → Execute
```

No complex build system. Just: check types, then run.

---

## File Structure

```
moonshot/
  main.go              # Entry point
  lexer.go             # Tokenize source
  parser.go            # Build AST
  ast.go               # AST node types
  types.go             # Type system
  checker.go           # Type checker
  eval.go              # Interpreter
  value.go             # Runtime values
  builtins.go          # Built-in functions
  examples/
    hello.moon         # Test files
    user.moon
```

Keep it simple - one file per major component.

---

## Implementation Order

Build in this order. Each step produces something runnable.

### Step 1: Lexer (Tokenizer)
**Goal:** Turn source code into tokens

**Input:**
```moonshot
def x = 5
```

**Output:**
```
[DEF] [IDENT "x"] [ASSIGN "="] [INTEGER "5"]
```

**What to build:**
- `Token` struct (type, literal, line, column)
- `Lexer` struct (input string, current position)
- `NextToken()` method

**Test:**
```go
func main() {
    input := "def x = 5"
    lexer := NewLexer(input)
    
    for tok := lexer.NextToken(); tok.Type != EOF; tok = lexer.NextToken() {
        fmt.Printf("%v\n", tok)
    }
}
```

**Keywords to recognize:**
```
def, fun, struct, extend, if, else, while, for, in, 
return, match, Some, None, Ok, Error, true, false,
and, or, not, is, break, continue, import
```

**Operators:**
```
= == + - * / > < >= <= ( ) { } [ ] , : . ->
```

---

### Step 2: Parser (AST Builder)
**Goal:** Turn tokens into Abstract Syntax Tree

**Input:** Tokens from lexer

**Output:** AST nodes

**What to build:**
- AST node interfaces (`Node`, `Statement`, `Expression`)
- Concrete node types:
  - `Program` (list of statements)
  - `DefStatement` (def x = value)
  - `IntegerLiteral`, `StringLiteral`, `BooleanLiteral`
  - `Identifier`
  - `InfixExpression` (a + b)
- `Parser` struct with recursive descent parsing

**Start simple - just parse:**
1. Variable declarations: `def x = 5`
2. Integers and strings
3. Binary operations: `1 + 2`

**Test:**
```go
func main() {
    input := "def x = 5"
    lexer := NewLexer(input)
    parser := NewParser(lexer)
    program := parser.ParseProgram()
    
    fmt.Printf("%#v\n", program)
}
```

---

### Step 3: Basic Evaluator
**Goal:** Execute the simplest programs

**What to build:**
- `Value` interface (all runtime values)
- `IntegerValue`, `StringValue`, `BooleanValue`
- `Environment` (variable storage - map of name → value)
- `Eval()` function

**Can now run:**
```moonshot
def x = 5
def y = 10
def z = x + y
```

**Test:**
```go
func main() {
    input := `
        def x = 5
        def y = 10
        def z = x + y
    `
    
    lexer := NewLexer(input)
    parser := NewParser(lexer)
    program := parser.ParseProgram()
    
    env := NewEnvironment()
    result := Eval(program, env)
    
    fmt.Println(env.Get("z")) // 15
}
```

---

### Step 4: Functions
**Goal:** Define and call functions

**Add to parser:**
- `FunctionStatement` (fun name(params) -> Type { body })
- `CallExpression` (functionName(args))

**Add to evaluator:**
- `FunctionValue` (stores params, body, closure environment)
- Function call evaluation

**Can now run:**
```moonshot
fun add(a: Integer, b: Integer) -> Integer {
    return a + b
}

def result = add(5, 3)
```

---

### Step 5: Control Flow
**Goal:** If/else, loops

**Add to parser:**
- `IfStatement`
- `WhileStatement`
- `ForStatement`

**Add to evaluator:**
- Evaluate conditions (truthiness)
- Execute blocks conditionally
- Loop execution

**Can now run:**
```moonshot
def counter = Mutable[Integer](0)

while counter < 10 {
    println(counter)
    counter == counter + 1
}

for i in range(0, 5) {
    println(i)
}
```

---

### Step 6: Collections
**Goal:** Lists and Maps

**Add to parser:**
- `ListLiteral` → `[1, 2, 3]`
- `MapLiteral` → `{ "a": 1 }`
- `IndexExpression` → `list[0]`

**Add to evaluator:**
- `ListValue` (slice of Values)
- `MapValue` (map[string]Value)
- Built-in methods (append, get, map, filter, etc.)

**Can now run:**
```moonshot
def numbers = [1, 2, 3, 4, 5]
def doubled = numbers.map({ x -> x * 2 })
def ages = { "Alice": 30, "Bob": 25 }
```

---

### Step 7: Structs
**Goal:** User-defined types

**Add to parser:**
- `StructStatement` (struct definition)
- `StructLiteral` (creating instances)
- `MemberExpression` (accessing fields)

**Add to evaluator:**
- `StructValue`
- Field access
- Auto-generated `toString()`

**Can now run:**
```moonshot
struct User {
    name: String
    age: Integer
}

def alice = User { name: "Alice", age: 30 }
println(alice.name)
println(alice)  // "User{name: Alice, age: 30}"
```

---

### Step 8: Type Checker
**Goal:** Validate types before execution

**What to build:**
- `Type` interface
- `TypeChecker` with environment
- Check each statement/expression
- Validate:
  - Variable types match values
  - Function parameters match arguments
  - Return types match
  - No returning `Mutable[T]` from functions

**Can now catch:**
```moonshot
def x: Integer = "hello"  // ERROR: type mismatch

fun bad() -> Mutable[Integer] {  // ERROR: can't return Mutable
    return Mutable[Integer](0)
}
```

---

### Step 9: Option & Result Types
**Goal:** No null, errors as values

**Add to parser:**
- `OptionExpression` → `Some(value)`, `None`
- `ResultExpression` → `Ok(value)`, `Error("msg")`
- `MatchStatement`

**Add to evaluator:**
- `OptionValue` (hasValue bool, value)
- `ResultValue` (isOk bool, value)
- Match evaluation
- Auto short-circuit for Result chains

**Can now run:**
```moonshot
fun findUser(id: Integer) -> Option[User] {
    if id is 42 {
        return Some(User { name: "Alice", age: 30 })
    }
    return None
}

fun divide(a: Integer, b: Integer) -> Result[Integer, String] {
    if b is 0 {
        return Error("Division by zero")
    }
    return Ok(a / b)
}

def result = divide(10, 2)
    .then({ x -> divide(x, 2) })
    .map({ y -> y * 2 })
```

---

### Step 10: Mutable Type
**Goal:** Explicit mutability

**Add to evaluator:**
- `MutableValue` (wraps any value)
- `==` operator for mutable assignment
- Implicit `.get()` when reading

**Add to type checker:**
- Validate functions don't return `Mutable[T]`
- Validate `Mutable[T]` only created at declaration

**Can now run:**
```moonshot
def counter = Mutable[Integer](0)
counter == 5
counter == counter + 1

println(counter)  // implicit get → "6"
```

---

### Step 11: Error Enrichment
**Goal:** Rich error messages

**When a method returns `Error("msg")`, automatically add:**
- Method name
- Input value (using auto toString)
- Error message

**Add to evaluator:**
- Track current method name during execution
- When Error detected, create `ErrorValue` with context
- Format nicely for output

**Output:**
```
Error in User.validate
Input: User{id: 1, name: Bob, age: 16}
Reason: Must be 18+
```

---

### Step 12: Extensions
**Goal:** Add methods to any type

**Add to parser:**
- `ExtendStatement` (single or block syntax)

**Add to evaluator:**
- Extension registry (map type → methods)
- Method lookup (check extensions, then built-ins)

**Can now run:**
```moonshot
extend User {
    fun isAdult() -> Boolean {
        return this.age >= 18
    }
}

def alice = User { name: "Alice", age: 30 }
println(alice.isAdult())  // true
```

---

### Step 13: Modules
**Goal:** Import other files

**Add to parser:**
- `ImportStatement`

**What to build:**
- `ModuleLoader` (cache loaded modules)
- Load `.moon` files
- Merge module environment into current environment

**Can now run:**
```moonshot
// user.moon
struct User {
    name: String
    age: Integer
}

// main.moon
import user

def alice = user.User { name: "Alice", age: 30 }
```

---

## Testing Strategy

### Test each component independently:

```go
// lexer_test.go
func TestLexer(t *testing.T) {
    input := "def x = 5"
    expected := []Token{
        {Type: DEF, Literal: "def"},
        {Type: IDENT, Literal: "x"},
        {Type: ASSIGN, Literal: "="},
        {Type: INTEGER, Literal: "5"},
    }
    // ... test
}

// parser_test.go
func TestDefStatement(t *testing.T) {
    input := "def x = 5"
    program := parseProgram(input)
    // ... test AST structure
}

// eval_test.go
func TestIntegerArithmetic(t *testing.T) {
    tests := []struct{
        input    string
        expected int64
    }{
        {"5", 5},
        {"10", 10},
        {"5 + 5", 10},
        {"2 * 3", 6},
    }
    // ... test evaluation
}
```

---

## Main Entry Point

```go
// main.go
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: moonshot <file.moon>")
        os.Exit(1)
    }
    
    // Read file
    content, err := os.ReadFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    
    // Lex
    lexer := NewLexer(string(content))
    
    // Parse
    parser := NewParser(lexer)
    program := parser.ParseProgram()
    
    if len(parser.Errors()) > 0 {
        for _, err := range parser.Errors() {
            fmt.Printf("Parse error: %s\n", err)
        }
        os.Exit(1)
    }
    
    // Type check
    checker := NewTypeChecker()
    if err := checker.Check(program); err != nil {
        fmt.Printf("Type error: %v\n", err)
        os.Exit(1)
    }
    
    // Evaluate
    env := NewEnvironment()
    result := Eval(program, env)
    
    if result != nil && IsError(result) {
        fmt.Printf("%s\n", result.String())
        os.Exit(1)
    }
}
```

---

## Key Data Structures

### Token
```go
type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Column  int
}
```

### AST Nodes
```go
type Node interface {
    String() string
}

type Program struct {
    Statements []Statement
}

type DefStatement struct {
    Name  string
    Value Expression
}

type IntegerLiteral struct {
    Value int64
}

type InfixExpression struct {
    Left     Expression
    Operator string
    Right    Expression
}
```

### Runtime Values
```go
type Value interface {
    Type() string
    String() string
}

type IntegerValue struct {
    Value int64
}

type ListValue struct {
    Elements []Value
}

type StructValue struct {
    Name   string
    Fields map[string]Value
}

type MutableValue struct {
    Value Value  // wrapped value
}
```

### Environment
```go
type Environment struct {
    store map[string]Value
    outer *Environment  // parent scope
}

func (e *Environment) Get(name string) (Value, bool)
func (e *Environment) Set(name string, val Value)
```

---

## Build & Run

```bash
# Build
go build -o moonshot main.go

# Run example
./moonshot examples/hello.moon

# Test
go test ./...
```

---

## Debugging Tips

### Print AST
```go
program := parser.ParseProgram()
fmt.Printf("%#v\n", program)
```

### Print tokens
```go
for tok := lexer.NextToken(); tok.Type != EOF; tok = lexer.NextToken() {
    fmt.Printf("%v\n", tok)
}
```

### Trace evaluation
```go
func Eval(node Node, env *Environment) Value {
    fmt.Printf("Evaluating: %T\n", node)
    // ... rest of eval
}
```

---

## Common Patterns

### Recursive Descent Parsing
```go
func (p *Parser) parseExpression(precedence int) Expression {
    // Get prefix parser for current token
    prefix := p.prefixParsers[p.curToken.Type]
    if prefix == nil {
        return nil
    }
    
    leftExp := prefix()
    
    // Get infix parsers while precedence allows
    for precedence < p.peekPrecedence() {
        infix := p.infixParsers[p.peekToken.Type]
        if infix == nil {
            return leftExp
        }
        p.nextToken()
        leftExp = infix(leftExp)
    }
    
    return leftExp
}
```

### Environment Scoping
```go
// Create new scope
innerEnv := NewEnclosedEnvironment(outerEnv)

// Look up in current scope, then parent
func (e *Environment) Get(name string) (Value, bool) {
    val, ok := e.store[name]
    if !ok && e.outer != nil {
        return e.outer.Get(name)
    }
    return val, ok
}
```

### Type Checking
```go
func (tc *TypeChecker) checkExpression(expr Expression) (Type, error) {
    switch e := expr.(type) {
    case *IntegerLiteral:
        return IntegerType, nil
    case *InfixExpression:
        leftType, err := tc.checkExpression(e.Left)
        if err != nil {
            return nil, err
        }
        rightType, err := tc.checkExpression(e.Right)
        if err != nil {
            return nil, err
        }
        if leftType != rightType {
            return nil, fmt.Errorf("type mismatch")
        }
        return leftType, nil
    }
}
```

---

## Start Here

1. Create `main.go` with the entry point
2. Implement `lexer.go` - get tokens working
3. Implement `parser.go` - parse simple statements
4. Implement `eval.go` - evaluate integers and arithmetic
5. Test with `examples/hello.moon`
6. Iterate through the 13 steps above

**Keep it simple. Make it work. Then make it better.**

---

## Example Test Program

Create `examples/hello.moon`:
```moonshot
def greeting = "Hello, MoonShot!"
println(greeting)

fun add(a: Integer, b: Integer) -> Integer {
    return a + b
}

def result = add(5, 3)
println(result)

def numbers = [1, 2, 3, 4, 5]
def doubled = numbers.map({ x -> x * 2 })
println(doubled)
```

**Goal:** Get this running!

---

## Remember

- **Start simple** - integers and addition first
- **Test constantly** - each feature should have tests
- **Iterate** - don't try to build everything at once
- **Debug with prints** - see what's happening
- **Types help** - static types make the interpreter simpler

**You're building a real programming language - that's awesome!** 🚀

Good luck!
