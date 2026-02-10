// Fibonacci sequence in MoonShot

fun fibonacci(n: Integer) -> Integer {
    if n <= 1 {
        return n
    }
    return fibonacci(n - 1) + fibonacci(n - 2)
}

// Calculate first 10 fibonacci numbers
for i in range(10) {
    def result = fibonacci(i)
    println("fibonacci(" + str(i) + ") = " + str(result))
}
