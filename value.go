package main

import (
	"fmt"
	"sort"
	"strings"
)

// Value represents a runtime value
type Value interface {
	Type() string
	String() string
}

// IntegerValue represents an integer
type IntegerValue struct {
	Value int64
}

func (iv *IntegerValue) Type() string   { return "Integer" }
func (iv *IntegerValue) String() string { return fmt.Sprintf("%d", iv.Value) }

// FloatValue represents a float
type FloatValue struct {
	Value float64
}

func (fv *FloatValue) Type() string   { return "Float" }
func (fv *FloatValue) String() string { return fmt.Sprintf("%g", fv.Value) }

// StringValue represents a string
type StringValue struct {
	Value string
}

func (sv *StringValue) Type() string   { return "String" }
func (sv *StringValue) String() string { return sv.Value }

// BooleanValue represents a boolean
type BooleanValue struct {
	Value bool
}

func (bv *BooleanValue) Type() string { return "Boolean" }
func (bv *BooleanValue) String() string {
	if bv.Value {
		return "true"
	}
	return "false"
}

// NullValue represents the absence of a value
type NullValue struct{}

func (nv *NullValue) Type() string   { return "Null" }
func (nv *NullValue) String() string { return "null" }

// ListValue represents a list
type ListValue struct {
	Elements []Value
}

func (lv *ListValue) Type() string { return "List" }
func (lv *ListValue) String() string {
	var elements []string
	for _, e := range lv.Elements {
		elements = append(elements, e.String())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

// Append creates a new list with the element appended (immutable)
func (lv *ListValue) Append(v Value) *ListValue {
	newElements := make([]Value, len(lv.Elements)+1)
	copy(newElements, lv.Elements)
	newElements[len(lv.Elements)] = v
	return &ListValue{Elements: newElements}
}

// MapValue represents a map
type MapValue struct {
	Pairs map[string]Value
}

func (mv *MapValue) Type() string { return "Map" }
func (mv *MapValue) String() string {
	var pairs []string
	// Sort keys for consistent output
	keys := make([]string, 0, len(mv.Pairs))
	for k := range mv.Pairs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%q: %s", k, mv.Pairs[k].String()))
	}
	return "{" + strings.Join(pairs, ", ") + "}"
}

// FunctionValue represents a function
type FunctionValue struct {
	Name       string
	Parameters []*FunctionParameter
	Body       *BlockStatement
	Env        *Environment
	IsLambda   bool
	LambdaBody Expression // for single-expression lambdas
}

func (fv *FunctionValue) Type() string { return "Function" }
func (fv *FunctionValue) String() string {
	if fv.Name != "" {
		return fmt.Sprintf("<function %s>", fv.Name)
	}
	return "<lambda>"
}

// BuiltinFunction represents a built-in function
type BuiltinFunction struct {
	Name string
	Fn   func(args ...Value) Value
}

func (bf *BuiltinFunction) Type() string   { return "Builtin" }
func (bf *BuiltinFunction) String() string { return fmt.Sprintf("<builtin %s>", bf.Name) }

// StructDefinition represents a struct type definition
type StructDefinition struct {
	Name   string
	Fields []*StructField
}

func (sd *StructDefinition) Type() string   { return "StructDef" }
func (sd *StructDefinition) String() string { return fmt.Sprintf("<struct %s>", sd.Name) }

// StructValue represents an instance of a struct
type StructValue struct {
	Definition *StructDefinition
	Fields     map[string]Value
}

func (sv *StructValue) Type() string { return sv.Definition.Name }
func (sv *StructValue) String() string {
	var fields []string
	// Sort keys for consistent output
	keys := make([]string, 0, len(sv.Fields))
	for k := range sv.Fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fields = append(fields, fmt.Sprintf("%s: %s", k, sv.Fields[k].String()))
	}
	return sv.Definition.Name + "{" + strings.Join(fields, ", ") + "}"
}

// With creates a new struct with updated fields (immutable)
func (sv *StructValue) With(updates map[string]Value) *StructValue {
	newFields := make(map[string]Value)
	for k, v := range sv.Fields {
		newFields[k] = v
	}
	for k, v := range updates {
		newFields[k] = v
	}
	return &StructValue{Definition: sv.Definition, Fields: newFields}
}

// OptionValue represents Some(x) or None
type OptionValue struct {
	IsSome bool
	Value  Value
}

func (ov *OptionValue) Type() string { return "Option" }
func (ov *OptionValue) String() string {
	if ov.IsSome {
		return fmt.Sprintf("Some(%s)", ov.Value.String())
	}
	return "None"
}

// ResultValue represents Ok(x) or Error(x)
type ResultValue struct {
	IsOk  bool
	Value Value
	Error *ErrorValue
}

func (rv *ResultValue) Type() string { return "Result" }
func (rv *ResultValue) String() string {
	if rv.IsOk {
		return fmt.Sprintf("Ok(%s)", rv.Value.String())
	}
	return fmt.Sprintf("Error(%s)", rv.Error.String())
}

// MutableValue wraps a value to make it mutable
type MutableValue struct {
	Value Value
}

func (mv *MutableValue) Type() string { return "Mutable" }
func (mv *MutableValue) String() string {
	return mv.Value.String()
}

// Unwrap returns the inner value
func (mv *MutableValue) Unwrap() Value {
	return mv.Value
}

// ErrorValue represents an error with context
type ErrorValue struct {
	Method  string
	Input   string
	Message string
}

func (ev *ErrorValue) Type() string { return "Error" }
func (ev *ErrorValue) String() string {
	if ev.Method != "" {
		return fmt.Sprintf("Error in %s\nInput: %s\nReason: %s", ev.Method, ev.Input, ev.Message)
	}
	return ev.Message
}

// ReturnValue signals a return from a function
type ReturnValue struct {
	Value Value
}

func (rv *ReturnValue) Type() string   { return "Return" }
func (rv *ReturnValue) String() string { return rv.Value.String() }

// BreakValue signals a break from a loop
type BreakValue struct{}

func (bv *BreakValue) Type() string   { return "Break" }
func (bv *BreakValue) String() string { return "break" }

// ContinueValue signals a continue in a loop
type ContinueValue struct{}

func (cv *ContinueValue) Type() string   { return "Continue" }
func (cv *ContinueValue) String() string { return "continue" }

// ModuleValue represents an imported module
type ModuleValue struct {
	Name    string
	Exports *Environment
}

func (mv *ModuleValue) Type() string   { return "Module" }
func (mv *ModuleValue) String() string { return fmt.Sprintf("<module %s>", mv.Name) }

// Helper functions for unwrapping mutable values
func UnwrapValue(v Value) Value {
	if mv, ok := v.(*MutableValue); ok {
		return mv.Value
	}
	return v
}

// IsTruthy returns whether a value is truthy
func IsTruthy(v Value) bool {
	switch val := v.(type) {
	case *BooleanValue:
		return val.Value
	case *NullValue:
		return false
	case *IntegerValue:
		return val.Value != 0
	case *StringValue:
		return val.Value != ""
	case *ListValue:
		return len(val.Elements) > 0
	case *OptionValue:
		return val.IsSome
	case *MutableValue:
		return IsTruthy(val.Value)
	default:
		return true
	}
}
