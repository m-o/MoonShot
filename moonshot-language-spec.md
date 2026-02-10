# MoonShot Language Specification v1.0

## Introduction

**MoonShot** is a statically-typed, functional programming language with:
- **Immutability by default** - Safe and predictable code
- **No null values** - Use `Option[T]` instead
- **Errors as values** - Use `Result[T, E]` instead of exceptions
- **Explicit mutability** - `Mutable[T]` wrapper type
- **Auto error context** - Interpreter provides rich debugging info
- **World-class interpreter** - JIT compilation for performance

### Design Philosophy

1. **Safety first** - Prevent errors at compile time
2. **Clarity over cleverness** - Readable, full words
3. **Immutable by default** - Opt-in to mutation
4. **Explicit over implicit** - Types show intent
5. **Simple and practical** - Small, consistent language

---

## Syntax Overview

### Variables

```moonshot
// Immutable binding (default)
def x = 5
def name = "Alice"
def items = [1, 2, 3]

// Mutable - explicit wrapper type
def counter = Mutable[Integer](0)
def scores = Mutable[[Integer]]([])

// Mutable assignment
counter == 10
counter == counter + 1
scores == [1, 2, 3]
```

**Rules:**
- `def` declares a variable
- `=` binds immutable values (cannot be reassigned)
- `==` assigns to mutable values (override/update)
- `is` checks equality (comparison)
- `is not` checks inequality

### Primitive Types

```moonshot
Integer    // 42, -10, 0
Float      // 3.14, -0.5, 2.0
String     // "hello", "world"
Boolean    // true, false
```

### Collection Types

```moonshot
// List - dynamic array
[Integer]              // type
[1, 2, 3, 4, 5]        // literal

// Nested lists
[[String]]             // list of list of strings
[[[Integer]]]          // list of list of list

// Map - key-value pairs
{String: Integer}      // type
{ "Alice": 30, "Bob": 25 }  // literal

// Nested maps and lists
{String: [Integer]}              // map to list
[{String: Integer}]              // list of maps
{String: {Integer: User}}        // map of maps
```

### Wrapper Types

```moonshot
// Option - represents value that might not exist
Option[User]
Some(user)   // has value
None         // no value

// Result - represents operation that might fail
Result[Integer, String]
Ok(42)              // success
Error("failed")     // failure

// Mutable - explicit mutability wrapper
Mutable[Integer]
Mutable[[String]]
Mutable[{String: Integer}]
```

---

## Functions

### Function Declaration

```moonshot
// Basic function
fun greet(name: String) -> String {
    return "Hello, " + name
}

// Multiple parameters
fun add(a: Integer, b: Integer) -> Integer {
    return a + b
}

// No return type (returns Nothing/unit)
fun printMessage(msg: String) {
    println(msg)
}

// Functions are first-class values
def double = { x -> x * 2 }

// Higher-order functions
fun map(list: [Integer], transform: fun(Integer) -> Integer) -> [Integer] {
    def result = Mutable[[Integer]]([])
    for item in list {
        result.append(transform(item))
    }
    return result.get()
}
```

### Lambda Syntax

```moonshot
// Single parameter
{ x -> x * 2 }

// Multiple parameters
{ x, y -> x + y }

// Multi-line
{ x ->
    def doubled = x * 2
    return doubled + 1
}

// No parameters
{ -> println("hello") }

// Type annotations (optional)
{ x: Integer -> x * 2 }
{ x: Integer, y: Integer -> x + y }

// Usage
[1, 2, 3].map({ x -> x * 2 })
[1, 2, 3].filter({ x -> x > 1 })
[1, 2, 3].forEach({ x -> println(x) })
```

---

## Structs

### Definition

```moonshot
// Define a struct
struct User {
    id: Integer
    name: String
    email: String
    age: Integer
}

// Create instances
def alice = User {
    name: "Alice"
    email: "alice@example.com"
    age: 30
    id: 1
}

// Access fields
println(alice.name)  // "Alice"

// Auto-generated toString()
println(alice)  // "User{id: 1, name: Alice, email: alice@example.com, age: 30}"
```

### Immutability

```moonshot
// Structs are immutable
alice.age = 31  // ERROR!

// Create new version with changes
def older = alice.with { age: 31 }
// alice unchanged, older is new instance

// Update multiple fields
def updated = alice.with {
    age: 31
    email: "newalice@example.com"
}
```

---

## Extension Functions

### Single Extension

```moonshot
// Dot syntax for single function
extend Point.distance() -> Float {
    return sqrt(this.x * this.x + this.y * this.y)
}

extend String.isPalindrome() -> Boolean {
    return this is this.reverse()
}
```

### Multiple Extensions

```moonshot
// Block syntax for multiple functions
extend User {
    fun isAdult() -> Boolean {
        return this.age >= 18
    }
    
    fun greet() -> String {
        return "Hello, " + this.name + "!"
    }
    
    fun celebrateBirthday() -> User {
        return this.with { age: this.age + 1 }
    }
}

// Can extend any type, including built-ins
extend [Integer] {
    fun sum() -> Integer {
        def total = Mutable[Integer](0)
        for item in this {
            total == total + item
        }
        return total.get()
    }
    
    fun average() -> Float {
        return this.sum() / this.length()
    }
}
```

---

## Control Flow

### If/Else

```moonshot
// Basic if
if condition {
    // do something
}

// If/else
if condition {
    // do something
} else {
    // do something else
}

// If/else if/else
if score >= 90 {
    println("A")
} else if score >= 80 {
    println("B")
} else {
    println("C")
}

// If is an expression (returns value)
def max = if a > b { a } else { b }

def status = if user.age >= 18 {
    "adult"
} else {
    "minor"
}
```

### While Loop

```moonshot
// Basic while
while condition {
    // do something
}

// Example with mutable counter
def counter = Mutable[Integer](0)

while counter < 10 {
    println(counter)
    counter == counter + 1
}

// Infinite loop
while true {
    // do something
    if shouldStop {
        break
    }
}
```

### For Loop

```moonshot
// For over list
for item in items {
    println(item)
}

// For over range
for i in range(0, 10) {
    println(i)  // 0 to 9
}

// For over map
for key, value in ages {
    println(key + ": " + value)
}

// With index
for index, item in items {
    println(index + ": " + item)
}

// Nested loops
for row in matrix {
    for cell in row {
        println(cell)
    }
}
```

### Break and Continue

```moonshot
// Break - exit loop
for i in range(0, 100) {
    if i is 50 {
        break
    }
    println(i)
}

// Continue - skip to next iteration
for i in range(0, 10) {
    if i is 5 {
        continue
    }
    println(i)  // skips 5
}

// Works in while loops too
def counter = Mutable[Integer](0)
while true {
    counter == counter + 1
    if counter > 100 {
        break
    }
    if counter is 50 {
        continue
    }
    println(counter)
}
```

---

## Type System

### No Null - Option Type

```moonshot
// Function that might not find something
fun findUser(id: Integer) -> Option[User] {
    if found {
        return Some(user)
    } else {
        return None
    }
}

// Using Options - must handle both cases
def result = findUser(42)

match result {
    Some(user) -> println(user.name)
    None -> println("Not found")
}

// Convenience methods
def email = findUser(42)
    .map({ u -> u.email })
    .orElse("unknown@example.com")

// Chaining Options
def managerEmail = findEmployee(123)
    .andThen({ emp -> emp.managerId })
    .andThen({ id -> findEmployee(id) })
    .map({ manager -> manager.email })
```

### Error Handling - Result Type

```moonshot
// Function that can fail
fun divide(a: Integer, b: Integer) -> Result[Integer, String] {
    if b is 0 {
        return Error("Cannot divide by zero")
    }
    return Ok(a / b)
}

// Must handle errors
def result = divide(10, 2)

match result {
    Ok(value) -> println("Result: " + value)
    Error(msg) -> println("Error: " + msg)
}

// Auto short-circuit chaining
fun calculate(a: Integer, b: Integer, c: Integer) -> Result[Integer, String] {
    return divide(a, b)
        .then({ x -> divide(x, c) })
        .map({ y -> y * 2 })
}
// If any step fails, chain stops and returns Error
// Interpreter adds method name, input value, and error message automatically
```

### Auto Error Context

When a function returns an Error, the interpreter automatically enriches it with:
- **Method name** - which function/method failed
- **Input value** - what data was passed in (using auto-generated toString)
- **Error message** - the error description

```moonshot
fun validate() -> Result[User, String] {
    if this.age < 18 {
        return Error("Must be 18+")
    }
    return Ok(this)
}

// When error occurs, interpreter produces:
// Error in User.validate
// Input: User{id: 1, name: Bob, age: 16}
// Reason: Must be 18+
```

### Pattern Matching

```moonshot
// Match on Option
match findUser(42) {
    Some(user) -> println(user.name)
    None -> println("Not found")
}

// Match on Result
match divide(10, 0) {
    Ok(value) -> println(value)
    Error(msg) -> println(msg)
}

// Match with binding and logic
def result = processUser(123)

match result {
    Ok(user) -> {
        println("Success!")
        saveToDatabase(user)
    }
    Error(msg) -> {
        println("Failed: " + msg)
        logError(msg)
    }
}
```

### Mutable Type

`Mutable[T]` is a wrapper type that allows mutation. It has special rules:

**Rules:**
1. Can only be created via direct declaration (`def x = Mutable[T](value)`)
2. Functions can accept `Mutable[T]` parameters (for side effects)
3. Functions **cannot** return `Mutable[T]` (prevents leaking mutable state)
4. Cannot be shared across threads (when concurrency is added)

```moonshot
// ✅ Create mutable
def counter = Mutable[Integer](0)
def items = Mutable[[String]]([])

// ✅ Assignment with ==
counter == 5
counter == counter + 1

// ✅ Read (implicit .get())
println(counter)        // prints "5"
val x = counter + 10    // x is 15

// ✅ Function can accept Mutable
fun increment(counter: Mutable[Integer]) {
    counter == counter + 1
}

// ❌ Function cannot return Mutable
fun createCounter() -> Mutable[Integer] {  // COMPILE ERROR!
    return Mutable[Integer](0)
}

// ✅ Pattern: local mutation, return immutable
fun sum(numbers: [Integer]) -> Integer {
    def total = Mutable[Integer](0)
    for n in numbers {
        total == total + n
    }
    return total.get()  // or just: return total
}
```

---

## Collections API

### List Methods

```moonshot
// Creation
def empty = [Integer]()
def numbers = [1, 2, 3, 4, 5]

// Size
numbers.length()        // 5
numbers.isEmpty()       // false

// Access
numbers.get(0)          // 1
numbers.first()         // Option[Integer] - Some(1)
numbers.last()          // Option[Integer] - Some(5)

// Add/Remove (returns new list - immutable)
numbers.append(6)           // [1, 2, 3, 4, 5, 6]
numbers.prepend(0)          // [0, 1, 2, 3, 4, 5]
numbers.concat([6, 7])      // [1, 2, 3, 4, 5, 6, 7]

// Transform
numbers.map({ x -> x * 2 })              // [2, 4, 6, 8, 10]
numbers.filter({ x -> x > 2 })           // [3, 4, 5]
numbers.reduce(0, { acc, x -> acc + x }) // 15

// Find
numbers.find({ x -> x > 3 })    // Option[Integer] - Some(4)
numbers.contains(3)             // true
numbers.indexOf(3)              // Option[Integer] - Some(2)

// Slice
numbers.slice(1, 3)     // [2, 3] - from index 1 to 3 (exclusive)
numbers.take(3)         // [1, 2, 3]
numbers.drop(2)         // [3, 4, 5]

// Sort
numbers.sort()                      // [1, 2, 3, 4, 5]
numbers.sortBy({ x -> -x })         // [5, 4, 3, 2, 1]
numbers.reverse()                   // [5, 4, 3, 2, 1]

// Check
numbers.any({ x -> x > 10 })        // false
numbers.all({ x -> x > 0 })         // true

// Iteration
numbers.forEach({ x -> println(x) })
```

### Map Methods

```moonshot
// Creation
def empty = {String: Integer}()
def ages = { "Alice": 30, "Bob": 25 }

// Size
ages.size()             // 2
ages.isEmpty()          // false

// Access
ages.get("Alice")               // Option[Integer] - Some(30)
ages.getOrElse("Charlie", 0)    // 0

// Add/Remove (returns new map - immutable)
ages.insert("Charlie", 35)      // new map with Charlie
ages.remove("Bob")              // new map without Bob

// Check
ages.contains("Alice")          // true

// Keys and Values
ages.keys()         // ["Alice", "Bob"]
ages.values()       // [30, 25]
ages.entries()      // [("Alice", 30), ("Bob", 25)]

// Transform
ages.map({ k, v -> v + 1 })             // { "Alice": 31, "Bob": 26 }
ages.filter({ k, v -> v > 25 })         // { "Alice": 30 }

// Iteration
ages.forEach({ k, v -> println(k + ": " + v) })
```

### Mutable Collection Helpers

For performance, `Mutable[List]` and `Mutable[Map]` have convenience methods:

```moonshot
def items = Mutable[[Integer]]([1, 2, 3])

// Mutable operations (mutate in place)
items.append(4)         // [1, 2, 3, 4]
items.prepend(0)        // [0, 1, 2, 3, 4]
items.clear()           // []

def cache = Mutable[{String: Integer}]({})
cache.insert("key", 42)
cache.remove("key")
cache.clear()
```

---

## Operators

### Comparison

```moonshot
// Equality
x is 5
x is not 5

// Comparison
x > 5
x < 10
x >= 5
x <= 10

// Logical
condition1 and condition2
condition1 or condition2
not condition
```

### Arithmetic

```moonshot
// Basic
a + b       // addition
a - b       // subtraction
a * b       // multiplication
a / b       // division
a % b       // modulo

// Assignment (only for Mutable)
counter == counter + 1
counter == counter * 2
```

### String Operations

```moonshot
// Concatenation
"Hello" + " " + "World"     // "Hello World"

// Comparison
"abc" is "abc"              // true
"abc" < "xyz"               // true (lexicographic)
```

---

## Module System

### File Structure

```
project/
  main.moon
  user.moon
  database.moon
  utils/
    math.moon
    string.moon
```

### Exports

Everything in a file is automatically exported:

```moonshot
// user.moon
struct User {
    name: String
    age: Integer
}

extend User {
    fun isAdult() -> Boolean {
        return this.age >= 18
    }
}

fun createUser(name: String, age: Integer) -> User {
    return User { name: name, age: age }
}
```

### Imports

```moonshot
// main.moon

// Import entire module
import user

fun main() {
    def alice = user.User { name: "Alice", age: 30 }
    println(alice.isAdult())
}

// Import specific items
import user.User
import user.createUser

fun main() {
    def alice = createUser("Alice", 30)
    println(alice.isAdult())
}

// Import from subdirectory
import utils.math
import utils.string.capitalize

fun main() {
    println(math.abs(-5))
    println(capitalize("hello"))
}
```

---

## Built-in Functions

### Print

```moonshot
// Print without newline
print("Hello")
print(42)

// Print with newline
println("Hello")
println(42)
println(user)  // uses auto-generated toString()

// Automatically converts any type to string
println([1, 2, 3])                    // "[1, 2, 3]"
println({ "a": 1 })                   // "{a: 1}"
println(User { name: "Alice", age: 30 })  // "User{name: Alice, age: 30}"
```

### Utility Functions

```moonshot
// Range
range(0, 10)        // [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
range(5, 10)        // [5, 6, 7, 8, 9]

// Math (future - not in v1.0)
// abs, sqrt, pow, min, max, etc.
```

---

## Complete Example

```moonshot
// user.moon
struct User {
    id: Integer
    name: String
    email: String
    age: Integer
}

extend User {
    fun validate() -> Result[User, String] {
        if this.age < 18 {
            return Error("Must be 18+")
        }
        if not this.email.contains("@") {
            return Error("Invalid email")
        }
        return Ok(this)
    }
    
    fun isAdult() -> Boolean {
        return this.age >= 18
    }
}

// database.moon
import user.User

struct Database {
    users: [User]
    index: {Integer: User}
}

extend Database {
    fun findUser(id: Integer) -> Option[User] {
        return this.index.get(id)
    }
    
    fun addUser(user: User) -> Result[Database, String] {
        return user.validate()
            .then({ validated ->
                if this.index.contains(validated.id) {
                    return Error("User already exists")
                }
                return Ok(this.with {
                    users: this.users.append(validated)
                    index: this.index.insert(validated.id, validated)
                })
            })
    }
    
    fun listAdults() -> [User] {
        return this.users.filter({ u -> u.isAdult() })
    }
}

// main.moon
import user.User
import database.Database

fun main() {
    def db = Mutable[Database](Database {
        users: []
        index: {}
    })
    
    def users = [
        User { id: 1, name: "Alice", email: "alice@test.com", age: 30 },
        User { id: 2, name: "Bob", email: "invalid", age: 16 },
        User { id: 3, name: "Charlie", email: "charlie@test.com", age: 25 }
    ]
    
    for user in users {
        match db.get().addUser(user) {
            Ok(newDb) -> {
                db == newDb
                println("Added: " + user.name)
            }
            Error(msg) -> {
                println("Failed for " + user.name + ": " + msg)
            }
        }
    }
    
    println("\nAdults:")
    def adults = db.get().listAdults()
    for adult in adults {
        println("  " + adult.name + " (" + adult.age + ")")
    }
}

// Output:
// Added: Alice
// Failed for Bob: Invalid email
// Added: Charlie
//
// Adults:
//   Alice (30)
//   Charlie (25)
```

---

## Quick Reference

| Feature | Syntax |
|---------|--------|
| **Variables** | `def x = 5`, `def y = Mutable[Integer](0)` |
| **Assignment** | `=` immutable, `==` mutable |
| **Equality** | `is`, `is not` |
| **Functions** | `fun name(params) -> Type { }` |
| **Structs** | `struct Name { fields }` |
| **Single Extension** | `extend Type.method() -> T { }` |
| **Multi Extension** | `extend Type { fun method() { } }` |
| **If/Else** | `if cond { } else { }` |
| **While** | `while cond { }` |
| **For** | `for item in list { }` |
| **Break/Continue** | `break`, `continue` |
| **Lambda** | `{ x -> x * 2 }`, `{ x, y -> x + y }` |
| **List Type** | `[Integer]` |
| **Map Type** | `{String: Integer}` |
| **List Literal** | `[1, 2, 3]` |
| **Map Literal** | `{ "a": 1, "b": 2 }` |
| **Option** | `Some(value)`, `None` |
| **Result** | `Ok(value)`, `Error("msg")` |
| **Mutable** | `Mutable[T](value)` |
| **Match** | `match value { Ok(x) -> ... }` |
| **Update Struct** | `.with { field: value }` |
| **Import** | `import module` |
| **Print** | `println(value)` |

---

## Language Guarantees

1. **No null pointer exceptions** - Option type forces handling
2. **No unhandled errors** - Result type forces handling
3. **No hidden mutation** - Mutable type is explicit and visible
4. **No race conditions** - Mutable cannot cross threads (future)
5. **Rich error context** - Interpreter adds method, input, and message
6. **Type safety** - Static typing catches errors at compile time
7. **Immutable by default** - Safe, predictable code

---

## Future Features (Not in v1.0)

- Concurrency (async/await)
- Generics (type parameters)
- Traits/Interfaces
- More pattern matching (destructuring, guards)
- Standard library (math, string utilities, etc.)
- Package manager
- REPL

---

**MoonShot v1.0 - Safe, Fast, Beautiful** 🚀
