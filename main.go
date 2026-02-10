package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("MoonShot Language Interpreter")
		fmt.Println("Usage: moonshot <file.moon>")
		fmt.Println("       moonshot -e <expression>")
		os.Exit(0)
	}

	var source string
	var filename string

	if os.Args[1] == "-e" {
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: -e requires an expression")
			os.Exit(1)
		}
		source = os.Args[2]
		filename = "<eval>"
	} else {
		filename = os.Args[1]
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
			os.Exit(1)
		}
		source = string(content)
	}

	result := Run(source, filename)
	if result != nil {
		if errVal, ok := result.(*ErrorValue); ok {
			fmt.Fprintln(os.Stderr, errVal.String())
			os.Exit(1)
		}
	}
}

// Run executes MoonShot source code
func Run(source string, filename string) Value {
	lexer := NewLexer(source)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		for _, err := range parser.Errors() {
			fmt.Fprintf(os.Stderr, "Parse error: %s\n", err)
		}
		return &ErrorValue{Message: "Parse errors occurred"}
	}

	// Type check
	checker := NewTypeChecker()
	if err := checker.Check(program); err != nil {
		fmt.Fprintf(os.Stderr, "Type error: %s\n", err)
		return &ErrorValue{Message: err.Error()}
	}

	// Evaluate
	env := NewEnvironment()
	RegisterBuiltins(env)
	evaluator := NewEvaluator()

	return evaluator.Eval(program, env)
}
