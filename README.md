# MoonShot

A statically-typed, expression-oriented programming language with immutable-by-default semantics.

## Building

Requires Go 1.21 or later.

```bash
# Build the interpreter
go build -o moonshot .

# Run a program
./moonshot examples/hello.moon

# Evaluate an expression directly
./moonshot -e 'println("Hello, World!")'
```

## Language Features

### Variables

Variables are immutable by default. Use `def` to declare:

```moonshot
def name = "Alice"
def age = 30
def pi = 3.14159
def active = true
```

Optional type annotations:

```moonshot
def count: Integer = 0
def message: String = "Hello"
```

### Mutable Variables

Use `Mutable[T]` for mutable state and `==` to update:

```moonshot
def counter = Mutable[Integer](0)
counter == counter + 1
counter == counter + 1
println(counter)  // 2
```

### Data Types

| Type | Example | Description |
|------|---------|-------------|
| `Integer` | `42`, `-17` | 64-bit signed integer |
| `Float` | `3.14`, `-0.5` | 64-bit floating point |
| `String` | `"hello"` | UTF-8 string |
| `Boolean` | `true`, `false` | Boolean value |
| `List[T]` | `[1, 2, 3]` | Immutable list |
| `Map[K, V]` | `{"key": "value"}` | Immutable map |
| `Option[T]` | `Some(x)`, `None` | Optional value |
| `Result[T, E]` | `Ok(x)`, `Error(e)` | Success or error |
| `Mutable[T]` | `Mutable[Integer](0)` | Mutable wrapper |

### Operators

```moonshot
// Arithmetic
def sum = 10 + 5      // 15
def diff = 10 - 5     // 5
def prod = 10 * 5     // 50
def quot = 10 / 5     // 2
def rem = 10 % 3      // 1

// Comparison
def gt = 10 > 5       // true
def lt = 10 < 5       // false
def gte = 10 >= 10    // true
def lte = 5 <= 10     // true

// Equality (use 'is')
def eq = 10 is 10     // true

// Logical
def both = true and false   // false
def either = true or false  // true
def negated = not true      // false

// String concatenation
def greeting = "Hello, " + "World!"
```

### Functions

```moonshot
fun add(a: Integer, b: Integer) -> Integer {
    return a + b
}

fun greet(name: String) -> String {
    return "Hello, " + name + "!"
}

println(add(5, 3))        // 8
println(greet("Alice"))   // Hello, Alice!
```

### Lambdas

Anonymous functions with concise syntax:

```moonshot
// Single parameter
def double = { x -> x * 2 }
println(double(5))  // 10

// Multiple parameters
def add = { a, b -> a + b }
println(add(3, 4))  // 7

// Used with higher-order functions
def numbers = [1, 2, 3, 4, 5]
def doubled = numbers.map({ x -> x * 2 })
println(doubled)  // [2, 4, 6, 8, 10]
```

### Control Flow

#### If/Else

```moonshot
def age = 20

if age >= 18 {
    println("Adult")
} else {
    println("Minor")
}

// If as expression
def status = if age >= 18 { "adult" } else { "minor" }
```

#### While Loop

```moonshot
def i = Mutable[Integer](0)
while i < 5 {
    println(i)
    i == i + 1
}
```

#### For Loop

```moonshot
for item in [1, 2, 3, 4, 5] {
    println(item)
}

// Using range
for i in range(5) {
    println(i)  // 0, 1, 2, 3, 4
}

for i in range(2, 5) {
    println(i)  // 2, 3, 4
}
```

#### Break and Continue

```moonshot
for i in range(10) {
    if i is 3 {
        continue  // Skip 3
    }
    if i is 7 {
        break     // Stop at 7
    }
    println(i)
}
```

### Lists

Lists are immutable. Operations return new lists.

```moonshot
def numbers = [1, 2, 3, 4, 5]

// Access by index
println(numbers[0])           // 1

// Methods
println(numbers.length())     // 5
println(numbers.append(6))    // [1, 2, 3, 4, 5, 6]
println(numbers.contains(3))  // true

// Higher-order functions
def doubled = numbers.map({ x -> x * 2 })
println(doubled)  // [2, 4, 6, 8, 10]

def evens = numbers.filter({ x -> x % 2 is 0 })
println(evens)  // [2, 4]

def sum = numbers.reduce({ acc, x -> acc + x }, 0)
println(sum)  // 15

def found = numbers.find({ x -> x > 3 })
// found is Some(4)
```

### Maps

Maps are immutable with string keys.

```moonshot
def person = {"name": "Alice", "city": "Paris"}

// Access (returns Option)
match person.get("name") {
    Some(name) -> { println(name) }
    None -> { println("Not found") }
}

// Methods
def updated = person.insert("age", "30")
println(updated)  // {"age": 30, "city": Paris, "name": Alice}

def removed = person.remove("city")
println(person.keys())      // ["city", "name"]
println(person.values())    // [Paris, Alice]
println(person.contains("name"))  // true
```

### Structs

Define custom data types:

```moonshot
struct User {
    name: String,
    age: Integer
}

// Create instance
def alice = User { name: "Alice", age: 30 }

// Access fields
println(alice.name)  // Alice
println(alice.age)   // 30

// Update with .with (returns new struct)
def older = alice.with { age: 31 }
println(older)  // User{age: 31, name: Alice}
```

### Extension Methods

Add methods to existing types:

```moonshot
struct User {
    name: String,
    age: Integer
}

extend User {
    fun isAdult() -> Boolean {
        return this.age >= 18
    }

    fun greet() -> String {
        return "Hello, I'm " + this.name
    }
}

def alice = User { name: "Alice", age: 30 }
println(alice.isAdult())  // true
println(alice.greet())    // Hello, I'm Alice
```

### Option Type

Represents optional values safely:

```moonshot
fun findUser(id: Integer) -> Option[String] {
    if id is 1 {
        return Some("Alice")
    }
    return None
}

def result = findUser(1)

match result {
    Some(name) -> { println("Found: " + name) }
    None -> { println("Not found") }
}

// Methods
println(result.isSome())        // true
println(result.isNone())        // false
println(result.unwrapOr("Unknown"))  // Alice
```

### Result Type

Handle errors explicitly:

```moonshot
fun divide(a: Integer, b: Integer) -> Result[Integer, String] {
    if b is 0 {
        return Error("Division by zero")
    }
    return Ok(a / b)
}

def result = divide(10, 2)

match result {
    Ok(value) -> { println("Result: " + str(value)) }
    Error(msg) -> { println("Error: " + msg) }
}

// Chaining with .then and .map
def chained = divide(10, 2)
    .then({ x -> divide(x, 2) })
    .map({ x -> x * 10 })
```

### Pattern Matching

Match on Option and Result types:

```moonshot
def value = Some(42)

match value {
    Some(x) -> { println("Got: " + str(x)) }
    None -> { println("Nothing") }
}

def result = Ok(100)

match result {
    Ok(x) -> { println("Success: " + str(x)) }
    Error(e) -> { println("Failed: " + e) }
}
```

### Comments

```moonshot
// This is a single-line comment

def x = 5  // Inline comment
```

### Built-in Functions

| Function | Description |
|----------|-------------|
| `print(args...)` | Print without newline |
| `println(args...)` | Print with newline |
| `range(end)` | Generate list `[0, 1, ..., end-1]` |
| `range(start, end)` | Generate list `[start, ..., end-1]` |
| `len(x)` | Length of string, list, or map |
| `type(x)` | Get type name as string |
| `str(x)` | Convert to string |
| `int(x)` | Convert to integer |
| `float(x)` | Convert to float |

### String Methods

```moonshot
def s = "Hello, World!"

println(s.length())          // 13
println(s.upper())           // HELLO, WORLD!
println(s.lower())           // hello, world!
println(s.trim())            // Removes whitespace
println(s.contains("World")) // true
println(s.split(", "))       // ["Hello", "World!"]
```

### Modules

Import other MoonShot files:

```moonshot
// utils.moon
fun helper() -> String {
    return "I'm helping!"
}

// main.moon
import utils
println(utils.helper())
```

## Complete Examples

### Fibonacci

```moonshot
fun fibonacci(n: Integer) -> Integer {
    if n <= 1 {
        return n
    }
    return fibonacci(n - 1) + fibonacci(n - 2)
}

for i in range(10) {
    println("fib(" + str(i) + ") = " + str(fibonacci(i)))
}
```

### FizzBuzz

```moonshot
for i in range(1, 101) {
    if i % 15 is 0 {
        println("FizzBuzz")
    } else {
        if i % 3 is 0 {
            println("Fizz")
        } else {
            if i % 5 is 0 {
                println("Buzz")
            } else {
                println(i)
            }
        }
    }
}
```

### User Management

```moonshot
struct User {
    id: Integer,
    name: String,
    email: String
}

extend User {
    fun validate() -> Result[User, String] {
        if this.name.length() is 0 {
            return Error("Name cannot be empty")
        }
        if not this.email.contains("@") {
            return Error("Invalid email")
        }
        return Ok(this)
    }
}

def user = User { id: 1, name: "Alice", email: "alice@example.com" }

match user.validate() {
    Ok(u) -> { println("Valid user: " + u.name) }
    Error(msg) -> { println("Validation failed: " + msg) }
}
```

### Data Processing

```moonshot
def data = [
    {"name": "Alice", "score": "95"},
    {"name": "Bob", "score": "87"},
    {"name": "Charlie", "score": "92"}
]

// Calculate average score
def scores = data.map({ item ->
    match item.get("score") {
        Some(s) -> { int(s) }
        None -> { 0 }
    }
})

def total = scores.reduce({ acc, x -> acc + x }, 0)
def average = total / scores.length()
println("Average score: " + str(average))

// Find top scorer
def topScore = scores.reduce({ a, b -> if a > b { a } else { b } }, 0)
println("Top score: " + str(topScore))
```

## Architecture

```
Source (.moon) -> Lexer -> Parser -> Type Checker -> Evaluator -> Output
```

| File | Purpose |
|------|---------|
| `token.go` | Token types and keywords |
| `lexer.go` | Tokenization |
| `ast.go` | AST node definitions |
| `parser.go` | Recursive descent parser with Pratt precedence |
| `types.go` | Type system |
| `checker.go` | Static type checking |
| `value.go` | Runtime values |
| `environment.go` | Variable scopes |
| `eval.go` | Tree-walking interpreter |
| `builtins.go` | Built-in functions and methods |
| `module.go` | Module loader |
| `errors.go` | Error handling |
| `main.go` | CLI entry point |

## Error Messages

MoonShot provides helpful error messages:

```
Error in divide
Input:
Reason: Division by zero
```

Type errors are caught before execution:

```
Type error: cannot assign String to variable of type Integer
```

## File Extension

MoonShot source files use the `.moon` extension.
