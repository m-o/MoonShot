// Mutable variables in MoonShot

// Create a mutable counter
def counter = Mutable[Integer](0)
println("Initial: " + str(counter))

// Update the counter using ==
counter == counter + 1
println("After increment: " + str(counter))

// While loop with mutable
def i = Mutable[Integer](0)
while i < 5 {
    println("i = " + str(i))
    i == i + 1
}

// Accumulator example
def total = Mutable[Integer](0)
for n in [1, 2, 3, 4, 5] {
    total == total + n
}
println("Total: " + str(total))
