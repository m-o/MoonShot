// Collections in MoonShot

// Lists
def numbers = [1, 2, 3, 4, 5]
println("Numbers: " + str(numbers))

// Map over list
def doubled = numbers.map({ x -> x * 2 })
println("Doubled: " + str(doubled))

// Filter list
def evens = numbers.filter({ x -> x % 2 is 0 })
println("Evens: " + str(evens))

// Reduce list
def sum = numbers.reduce({ acc, x -> acc + x }, 0)
println("Sum: " + str(sum))

// Find in list
def found = numbers.find({ x -> x > 3 })
match found {
    Some(val) -> { println("Found: " + str(val)) }
    None -> { println("Not found") }
}

// Maps
def person = {"name": "Alice", "city": "Paris"}
println("Person: " + str(person))

// Access map values
def cityOpt = person.get("city")
match cityOpt {
    Some(city) -> { println("City: " + city) }
    None -> { println("No city") }
}

// Insert into map (creates new map)
def personWithAge = person.insert("age", "30")
println("With age: " + str(personWithAge))
