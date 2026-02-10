package main

// Environment stores variable bindings
type Environment struct {
	store  map[string]Value
	parent *Environment
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	return &Environment{
		store:  make(map[string]Value),
		parent: nil,
	}
}

// NewEnclosedEnvironment creates a child environment
func NewEnclosedEnvironment(parent *Environment) *Environment {
	env := NewEnvironment()
	env.parent = parent
	return env
}

// Get retrieves a value from the environment
func (e *Environment) Get(name string) (Value, bool) {
	val, ok := e.store[name]
	if !ok && e.parent != nil {
		return e.parent.Get(name)
	}
	return val, ok
}

// Set defines a new variable in the current scope
func (e *Environment) Set(name string, val Value) Value {
	e.store[name] = val
	return val
}

// Update updates an existing variable in any scope
func (e *Environment) Update(name string, val Value) bool {
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return true
	}
	if e.parent != nil {
		return e.parent.Update(name, val)
	}
	return false
}

// GetDirect retrieves a value only from the current scope
func (e *Environment) GetDirect(name string) (Value, bool) {
	val, ok := e.store[name]
	return val, ok
}

// All returns all variable names in the current scope
func (e *Environment) All() []string {
	names := make([]string, 0, len(e.store))
	for name := range e.store {
		names = append(names, name)
	}
	return names
}

// Clone creates a shallow copy of the environment
func (e *Environment) Clone() *Environment {
	newStore := make(map[string]Value)
	for k, v := range e.store {
		newStore[k] = v
	}
	return &Environment{
		store:  newStore,
		parent: e.parent,
	}
}
