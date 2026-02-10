// Structs and Extensions in MoonShot

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

// Create a user
def alice = User { name: "Alice", age: 30 }
println(alice)
println("Is adult: " + str(alice.isAdult()))
println(alice.greet())

// Create updated user with .with
def olderAlice = alice.with { age: 31 }
println("After birthday: " + str(olderAlice))
