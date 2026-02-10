// Option and Result types in MoonShot

fun divide(a: Integer, b: Integer) -> Result[Integer, String] {
    if b is 0 {
        return Error("Division by zero")
    }
    return Ok(a / b)
}

// Test divide function
def result1 = divide(10, 2)
match result1 {
    Ok(val) -> { println("10 / 2 = " + str(val)) }
    Error(msg) -> { println("Error: " + msg) }
}

def result2 = divide(10, 0)
match result2 {
    Ok(val) -> { println("Result: " + str(val)) }
    Error(msg) -> { println("Error: " + msg) }
}

// Working with Option
fun findUser(id: Integer) -> Option[String] {
    if id is 1 {
        return Some("Alice")
    }
    if id is 2 {
        return Some("Bob")
    }
    return None
}

def user1 = findUser(1)
match user1 {
    Some(name) -> { println("Found user: " + name) }
    None -> { println("User not found") }
}

def user3 = findUser(99)
match user3 {
    Some(name) -> { println("Found user: " + name) }
    None -> { println("User not found") }
}
