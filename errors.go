package main

import (
	"fmt"
)

// MoonShotError represents a rich error with context
type MoonShotError struct {
	Type    string
	Method  string
	Input   string
	Message string
	Line    int
	Column  int
}

func (e *MoonShotError) Error() string {
	if e.Method != "" {
		return fmt.Sprintf("Error in %s\nInput: %s\nReason: %s", e.Method, e.Input, e.Message)
	}
	if e.Line > 0 {
		return fmt.Sprintf("Line %d, Column %d: %s", e.Line, e.Column, e.Message)
	}
	return e.Message
}

// NewParseError creates a parse error
func NewParseError(line, col int, msg string) *MoonShotError {
	return &MoonShotError{
		Type:    "ParseError",
		Line:    line,
		Column:  col,
		Message: msg,
	}
}

// NewTypeError creates a type error
func NewTypeError(msg string) *MoonShotError {
	return &MoonShotError{
		Type:    "TypeError",
		Message: msg,
	}
}

// NewRuntimeError creates a runtime error
func NewRuntimeError(method, input, msg string) *MoonShotError {
	return &MoonShotError{
		Type:    "RuntimeError",
		Method:  method,
		Input:   input,
		Message: msg,
	}
}

// EnrichError adds context to an error value
func EnrichError(err *ErrorValue, method string, input Value) *ErrorValue {
	if err.Method == "" {
		err.Method = method
	}
	if err.Input == "" && input != nil {
		err.Input = input.String()
	}
	return err
}

// FormatError formats an error for display
func FormatError(err *ErrorValue) string {
	if err.Method != "" {
		result := fmt.Sprintf("Error in %s", err.Method)
		if err.Input != "" {
			result += fmt.Sprintf("\nInput: %s", err.Input)
		}
		result += fmt.Sprintf("\nReason: %s", err.Message)
		return result
	}
	return err.Message
}
