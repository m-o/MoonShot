package main

import (
	"fmt"
	"strings"
)

// RegisterBuiltins registers all built-in functions
func RegisterBuiltins(env *Environment) {
	// I/O functions
	env.Set("print", &BuiltinFunction{
		Name: "print",
		Fn:   builtinPrint,
	})

	env.Set("println", &BuiltinFunction{
		Name: "println",
		Fn:   builtinPrintln,
	})

	// Collection functions
	env.Set("range", &BuiltinFunction{
		Name: "range",
		Fn:   builtinRange,
	})

	env.Set("len", &BuiltinFunction{
		Name: "len",
		Fn:   builtinLen,
	})

	env.Set("type", &BuiltinFunction{
		Name: "type",
		Fn:   builtinType,
	})

	env.Set("str", &BuiltinFunction{
		Name: "str",
		Fn:   builtinStr,
	})

	env.Set("int", &BuiltinFunction{
		Name: "int",
		Fn:   builtinInt,
	})

	env.Set("float", &BuiltinFunction{
		Name: "float",
		Fn:   builtinFloat,
	})
}

func builtinPrint(args ...Value) Value {
	var parts []string
	for _, arg := range args {
		parts = append(parts, UnwrapValue(arg).String())
	}
	fmt.Print(strings.Join(parts, " "))
	return &NullValue{}
}

func builtinPrintln(args ...Value) Value {
	var parts []string
	for _, arg := range args {
		parts = append(parts, UnwrapValue(arg).String())
	}
	fmt.Println(strings.Join(parts, " "))
	return &NullValue{}
}

func builtinRange(args ...Value) Value {
	if len(args) < 1 || len(args) > 2 {
		return &ErrorValue{Message: "range() requires 1 or 2 arguments"}
	}

	var start, end int64

	if len(args) == 1 {
		endVal, ok := UnwrapValue(args[0]).(*IntegerValue)
		if !ok {
			return &ErrorValue{Message: "range() argument must be an integer"}
		}
		start = 0
		end = endVal.Value
	} else {
		startVal, ok := UnwrapValue(args[0]).(*IntegerValue)
		if !ok {
			return &ErrorValue{Message: "range() start must be an integer"}
		}
		endVal, ok := UnwrapValue(args[1]).(*IntegerValue)
		if !ok {
			return &ErrorValue{Message: "range() end must be an integer"}
		}
		start = startVal.Value
		end = endVal.Value
	}

	elements := make([]Value, 0, end-start)
	for i := start; i < end; i++ {
		elements = append(elements, &IntegerValue{Value: i})
	}

	return &ListValue{Elements: elements}
}

func builtinLen(args ...Value) Value {
	if len(args) != 1 {
		return &ErrorValue{Message: "len() requires exactly 1 argument"}
	}

	arg := UnwrapValue(args[0])
	switch val := arg.(type) {
	case *StringValue:
		return &IntegerValue{Value: int64(len(val.Value))}
	case *ListValue:
		return &IntegerValue{Value: int64(len(val.Elements))}
	case *MapValue:
		return &IntegerValue{Value: int64(len(val.Pairs))}
	default:
		return &ErrorValue{Message: fmt.Sprintf("len() not supported for %s", arg.Type())}
	}
}

func builtinType(args ...Value) Value {
	if len(args) != 1 {
		return &ErrorValue{Message: "type() requires exactly 1 argument"}
	}
	return &StringValue{Value: UnwrapValue(args[0]).Type()}
}

func builtinStr(args ...Value) Value {
	if len(args) != 1 {
		return &ErrorValue{Message: "str() requires exactly 1 argument"}
	}
	return &StringValue{Value: UnwrapValue(args[0]).String()}
}

func builtinInt(args ...Value) Value {
	if len(args) != 1 {
		return &ErrorValue{Message: "int() requires exactly 1 argument"}
	}

	arg := UnwrapValue(args[0])
	switch val := arg.(type) {
	case *IntegerValue:
		return val
	case *FloatValue:
		return &IntegerValue{Value: int64(val.Value)}
	case *StringValue:
		var i int64
		_, err := fmt.Sscanf(val.Value, "%d", &i)
		if err != nil {
			return &ErrorValue{Message: fmt.Sprintf("cannot convert %q to integer", val.Value)}
		}
		return &IntegerValue{Value: i}
	case *BooleanValue:
		if val.Value {
			return &IntegerValue{Value: 1}
		}
		return &IntegerValue{Value: 0}
	default:
		return &ErrorValue{Message: fmt.Sprintf("cannot convert %s to integer", arg.Type())}
	}
}

func builtinFloat(args ...Value) Value {
	if len(args) != 1 {
		return &ErrorValue{Message: "float() requires exactly 1 argument"}
	}

	arg := UnwrapValue(args[0])
	switch val := arg.(type) {
	case *FloatValue:
		return val
	case *IntegerValue:
		return &FloatValue{Value: float64(val.Value)}
	case *StringValue:
		var f float64
		_, err := fmt.Sscanf(val.Value, "%f", &f)
		if err != nil {
			return &ErrorValue{Message: fmt.Sprintf("cannot convert %q to float", val.Value)}
		}
		return &FloatValue{Value: f}
	default:
		return &ErrorValue{Message: fmt.Sprintf("cannot convert %s to float", arg.Type())}
	}
}

// List methods

func listLength(list *ListValue) Value {
	return &IntegerValue{Value: int64(len(list.Elements))}
}

func listGet(list *ListValue, index int64) Value {
	if index < 0 || index >= int64(len(list.Elements)) {
		return &OptionValue{IsSome: false}
	}
	return &OptionValue{IsSome: true, Value: list.Elements[index]}
}

func listAppend(list *ListValue, val Value) *ListValue {
	return list.Append(val)
}

func listMap(list *ListValue, fn *FunctionValue, eval *Evaluator, env *Environment) *ListValue {
	newElements := make([]Value, len(list.Elements))
	for i, elem := range list.Elements {
		result := eval.applyFunction(fn, []Value{elem}, env)
		newElements[i] = result
	}
	return &ListValue{Elements: newElements}
}

func listFilter(list *ListValue, fn *FunctionValue, eval *Evaluator, env *Environment) *ListValue {
	var newElements []Value
	for _, elem := range list.Elements {
		result := eval.applyFunction(fn, []Value{elem}, env)
		if IsTruthy(result) {
			newElements = append(newElements, elem)
		}
	}
	return &ListValue{Elements: newElements}
}

func listReduce(list *ListValue, fn *FunctionValue, initial Value, eval *Evaluator, env *Environment) Value {
	acc := initial
	for _, elem := range list.Elements {
		acc = eval.applyFunction(fn, []Value{acc, elem}, env)
	}
	return acc
}

func listFind(list *ListValue, fn *FunctionValue, eval *Evaluator, env *Environment) *OptionValue {
	for _, elem := range list.Elements {
		result := eval.applyFunction(fn, []Value{elem}, env)
		if IsTruthy(result) {
			return &OptionValue{IsSome: true, Value: elem}
		}
	}
	return &OptionValue{IsSome: false}
}

func listContains(list *ListValue, val Value) bool {
	for _, elem := range list.Elements {
		if valuesEqual(elem, val) {
			return true
		}
	}
	return false
}

// Map methods

func mapGet(m *MapValue, key string) *OptionValue {
	if val, ok := m.Pairs[key]; ok {
		return &OptionValue{IsSome: true, Value: val}
	}
	return &OptionValue{IsSome: false}
}

func mapInsert(m *MapValue, key string, val Value) *MapValue {
	newPairs := make(map[string]Value)
	for k, v := range m.Pairs {
		newPairs[k] = v
	}
	newPairs[key] = val
	return &MapValue{Pairs: newPairs}
}

func mapRemove(m *MapValue, key string) *MapValue {
	newPairs := make(map[string]Value)
	for k, v := range m.Pairs {
		if k != key {
			newPairs[k] = v
		}
	}
	return &MapValue{Pairs: newPairs}
}

func mapKeys(m *MapValue) *ListValue {
	keys := make([]Value, 0, len(m.Pairs))
	for k := range m.Pairs {
		keys = append(keys, &StringValue{Value: k})
	}
	return &ListValue{Elements: keys}
}

func mapValues(m *MapValue) *ListValue {
	values := make([]Value, 0, len(m.Pairs))
	for _, v := range m.Pairs {
		values = append(values, v)
	}
	return &ListValue{Elements: values}
}

func mapContains(m *MapValue, key string) bool {
	_, ok := m.Pairs[key]
	return ok
}

// String methods

func stringLength(s *StringValue) Value {
	return &IntegerValue{Value: int64(len(s.Value))}
}

func stringSplit(s *StringValue, sep string) *ListValue {
	parts := strings.Split(s.Value, sep)
	elements := make([]Value, len(parts))
	for i, p := range parts {
		elements[i] = &StringValue{Value: p}
	}
	return &ListValue{Elements: elements}
}

func stringContains(s *StringValue, substr string) bool {
	return strings.Contains(s.Value, substr)
}

func stringTrim(s *StringValue) *StringValue {
	return &StringValue{Value: strings.TrimSpace(s.Value)}
}

func stringUpper(s *StringValue) *StringValue {
	return &StringValue{Value: strings.ToUpper(s.Value)}
}

func stringLower(s *StringValue) *StringValue {
	return &StringValue{Value: strings.ToLower(s.Value)}
}

// Helper function to compare values
func valuesEqual(a, b Value) bool {
	a = UnwrapValue(a)
	b = UnwrapValue(b)

	switch av := a.(type) {
	case *IntegerValue:
		if bv, ok := b.(*IntegerValue); ok {
			return av.Value == bv.Value
		}
	case *FloatValue:
		if bv, ok := b.(*FloatValue); ok {
			return av.Value == bv.Value
		}
	case *StringValue:
		if bv, ok := b.(*StringValue); ok {
			return av.Value == bv.Value
		}
	case *BooleanValue:
		if bv, ok := b.(*BooleanValue); ok {
			return av.Value == bv.Value
		}
	case *NullValue:
		_, ok := b.(*NullValue)
		return ok
	}
	return false
}
