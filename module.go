package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ModuleLoader handles loading and caching of modules
type ModuleLoader struct {
	basePath string
	cache    map[string]*Program
}

// NewModuleLoader creates a new module loader
func NewModuleLoader() *ModuleLoader {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	return &ModuleLoader{
		basePath: cwd,
		cache:    make(map[string]*Program),
	}
}

// SetBasePath sets the base path for module resolution
func (ml *ModuleLoader) SetBasePath(path string) {
	ml.basePath = path
}

// Load loads a module by name
func (ml *ModuleLoader) Load(modulePath string) (*Program, error) {
	// Check cache first
	if program, ok := ml.cache[modulePath]; ok {
		return program, nil
	}

	// Resolve module path
	filePath := ml.resolvePath(modulePath)

	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot load module %s: %v", modulePath, err)
	}

	// Parse module
	lexer := NewLexer(string(content))
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		return nil, fmt.Errorf("parse errors in module %s:\n%s",
			modulePath, strings.Join(parser.Errors(), "\n"))
	}

	// Cache the parsed module
	ml.cache[modulePath] = program

	return program, nil
}

// resolvePath converts a module path to a file path
func (ml *ModuleLoader) resolvePath(modulePath string) string {
	// Convert dot notation to path separators
	// e.g., "utils.math" -> "utils/math.moon"
	parts := strings.Split(modulePath, ".")
	relativePath := filepath.Join(parts...) + ".moon"

	return filepath.Join(ml.basePath, relativePath)
}

// ResolveImport resolves an import statement and returns the module path
func (ml *ModuleLoader) ResolveImport(importPath []string) (string, string) {
	modulePath := importPath[0]
	var itemName string

	if len(importPath) > 1 {
		// import user.User -> module "user", item "User"
		itemName = importPath[len(importPath)-1]
		// Check if this is a submodule or an item import
		// For now, assume single level: import module.Item
	}

	return modulePath, itemName
}

// CreateModuleEnvironment creates an environment for a module
func (ml *ModuleLoader) CreateModuleEnvironment(program *Program, eval *Evaluator) (*Environment, error) {
	env := NewEnvironment()
	RegisterBuiltins(env)

	result := eval.Eval(program, env)
	if errVal, ok := result.(*ErrorValue); ok {
		return nil, fmt.Errorf(errVal.String())
	}

	return env, nil
}

// GetExports returns the public exports of a module
func (ml *ModuleLoader) GetExports(env *Environment) map[string]Value {
	exports := make(map[string]Value)

	for _, name := range env.All() {
		// For now, all top-level definitions are exported
		// Could add convention: _ prefix means private
		if !strings.HasPrefix(name, "_") {
			if val, ok := env.GetDirect(name); ok {
				exports[name] = val
			}
		}
	}

	return exports
}
