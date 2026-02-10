package main

import (
	"fmt"
)

// Evaluator evaluates AST nodes
type Evaluator struct {
	structs    map[string]*StructDefinition
	extensions map[string]map[string]*FunctionValue
	modules    map[string]*ModuleValue
	loader     *ModuleLoader
	currentFn  string // current function name for error context
}

// NewEvaluator creates a new Evaluator
func NewEvaluator() *Evaluator {
	return &Evaluator{
		structs:    make(map[string]*StructDefinition),
		extensions: make(map[string]map[string]*FunctionValue),
		modules:    make(map[string]*ModuleValue),
		loader:     NewModuleLoader(),
	}
}

// Eval evaluates an AST node
func (e *Evaluator) Eval(node Node, env *Environment) Value {
	switch node := node.(type) {
	// Statements
	case *Program:
		return e.evalProgram(node, env)
	case *DefStatement:
		return e.evalDefStatement(node, env)
	case *ReturnStatement:
		return e.evalReturnStatement(node, env)
	case *ExpressionStatement:
		return e.Eval(node.Expression, env)
	case *BlockStatement:
		return e.evalBlockStatement(node, env)
	case *FunctionStatement:
		return e.evalFunctionStatement(node, env)
	case *WhileStatement:
		return e.evalWhileStatement(node, env)
	case *ForStatement:
		return e.evalForStatement(node, env)
	case *BreakStatement:
		return &BreakValue{}
	case *ContinueStatement:
		return &ContinueValue{}
	case *StructStatement:
		return e.evalStructStatement(node, env)
	case *ExtendStatement:
		return e.evalExtendStatement(node, env)
	case *ImportStatement:
		return e.evalImportStatement(node, env)

	// Expressions
	case *IntegerLiteral:
		return &IntegerValue{Value: node.Value}
	case *FloatLiteral:
		return &FloatValue{Value: node.Value}
	case *StringLiteral:
		return &StringValue{Value: node.Value}
	case *BooleanLiteral:
		return &BooleanValue{Value: node.Value}
	case *Identifier:
		return e.evalIdentifier(node, env)
	case *PrefixExpression:
		return e.evalPrefixExpression(node, env)
	case *InfixExpression:
		return e.evalInfixExpression(node, env)
	case *AssignmentExpression:
		return e.evalAssignmentExpression(node, env)
	case *IfExpression:
		return e.evalIfExpression(node, env)
	case *FunctionLiteral:
		return e.evalFunctionLiteral(node, env)
	case *CallExpression:
		return e.evalCallExpression(node, env)
	case *MemberExpression:
		return e.evalMemberExpression(node, env)
	case *IndexExpression:
		return e.evalIndexExpression(node, env)
	case *ListLiteral:
		return e.evalListLiteral(node, env)
	case *MapLiteral:
		return e.evalMapLiteral(node, env)
	case *StructLiteral:
		return e.evalStructLiteral(node, env)
	case *WithExpression:
		return e.evalWithExpression(node, env)
	case *OptionExpression:
		return e.evalOptionExpression(node, env)
	case *ResultExpression:
		return e.evalResultExpression(node, env)
	case *MatchExpression:
		return e.evalMatchExpression(node, env)
	case *MutableExpression:
		return e.evalMutableExpression(node, env)
	}

	return &NullValue{}
}

func (e *Evaluator) evalProgram(program *Program, env *Environment) Value {
	var result Value = &NullValue{}

	for _, stmt := range program.Statements {
		result = e.Eval(stmt, env)

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *ErrorValue:
			return result
		}
	}

	return result
}

func (e *Evaluator) evalDefStatement(stmt *DefStatement, env *Environment) Value {
	val := e.Eval(stmt.Value, env)
	// Note: ErrorValue is a valid value to assign, so don't propagate it as an error
	env.Set(stmt.Name.Value, val)
	return val
}

func (e *Evaluator) evalReturnStatement(stmt *ReturnStatement, env *Environment) Value {
	if stmt.Value == nil {
		return &ReturnValue{Value: &NullValue{}}
	}
	val := e.Eval(stmt.Value, env)
	if isError(val) {
		return val
	}
	return &ReturnValue{Value: val}
}

func (e *Evaluator) evalBlockStatement(block *BlockStatement, env *Environment) Value {
	var result Value = &NullValue{}

	for _, stmt := range block.Statements {
		result = e.Eval(stmt, env)

		if result != nil {
			switch result.(type) {
			case *ReturnValue, *BreakValue, *ContinueValue:
				return result
			}
		}
	}

	return result
}

func (e *Evaluator) evalFunctionStatement(stmt *FunctionStatement, env *Environment) Value {
	fn := &FunctionValue{
		Name:       stmt.Name.Value,
		Parameters: stmt.Parameters,
		Body:       stmt.Body,
		Env:        env,
	}
	env.Set(stmt.Name.Value, fn)
	return fn
}

func (e *Evaluator) evalWhileStatement(stmt *WhileStatement, env *Environment) Value {
	for {
		condition := e.Eval(stmt.Condition, env)
		if isError(condition) {
			return condition
		}

		if !IsTruthy(condition) {
			break
		}

		result := e.Eval(stmt.Body, NewEnclosedEnvironment(env))

		switch result.(type) {
		case *BreakValue:
			return &NullValue{}
		case *ContinueValue:
			continue
		case *ReturnValue, *ErrorValue:
			return result
		}
	}

	return &NullValue{}
}

func (e *Evaluator) evalForStatement(stmt *ForStatement, env *Environment) Value {
	iterable := e.Eval(stmt.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	list, ok := UnwrapValue(iterable).(*ListValue)
	if !ok {
		return &ErrorValue{Message: fmt.Sprintf("cannot iterate over %s", iterable.Type())}
	}

	for _, elem := range list.Elements {
		loopEnv := NewEnclosedEnvironment(env)
		loopEnv.Set(stmt.Variable.Value, elem)

		result := e.Eval(stmt.Body, loopEnv)

		switch result.(type) {
		case *BreakValue:
			return &NullValue{}
		case *ContinueValue:
			continue
		case *ReturnValue, *ErrorValue:
			return result
		}
	}

	return &NullValue{}
}

func (e *Evaluator) evalStructStatement(stmt *StructStatement, env *Environment) Value {
	def := &StructDefinition{
		Name:   stmt.Name.Value,
		Fields: stmt.Fields,
	}
	e.structs[stmt.Name.Value] = def
	env.Set(stmt.Name.Value, def)
	return def
}

func (e *Evaluator) evalExtendStatement(stmt *ExtendStatement, env *Environment) Value {
	typeName := stmt.TypeName.Value

	if _, ok := e.extensions[typeName]; !ok {
		e.extensions[typeName] = make(map[string]*FunctionValue)
	}

	for _, method := range stmt.Methods {
		fn := &FunctionValue{
			Name:       method.Name.Value,
			Parameters: method.Parameters,
			Body:       method.Body,
			Env:        env,
		}
		e.extensions[typeName][method.Name.Value] = fn
	}

	return &NullValue{}
}

func (e *Evaluator) evalImportStatement(stmt *ImportStatement, env *Environment) Value {
	moduleName := stmt.Path[0]

	if mod, ok := e.modules[moduleName]; ok {
		env.Set(moduleName, mod)
		return mod
	}

	program, err := e.loader.Load(moduleName)
	if err != nil {
		return &ErrorValue{Message: err.Error()}
	}

	modEnv := NewEnvironment()
	RegisterBuiltins(modEnv)

	result := e.Eval(program, modEnv)
	if isError(result) {
		return result
	}

	mod := &ModuleValue{
		Name:    moduleName,
		Exports: modEnv,
	}
	e.modules[moduleName] = mod
	env.Set(moduleName, mod)

	return mod
}

func (e *Evaluator) evalIdentifier(node *Identifier, env *Environment) Value {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return &ErrorValue{Message: fmt.Sprintf("undefined: %s", node.Value)}
}

func (e *Evaluator) evalPrefixExpression(node *PrefixExpression, env *Environment) Value {
	right := e.Eval(node.Right, env)
	if isError(right) {
		return right
	}

	right = UnwrapValue(right)

	switch node.Operator {
	case "-":
		return e.evalMinusPrefixExpression(right)
	case "not":
		return e.evalNotPrefixExpression(right)
	default:
		return &ErrorValue{Message: fmt.Sprintf("unknown operator: %s%s", node.Operator, right.Type())}
	}
}

func (e *Evaluator) evalMinusPrefixExpression(right Value) Value {
	switch val := right.(type) {
	case *IntegerValue:
		return &IntegerValue{Value: -val.Value}
	case *FloatValue:
		return &FloatValue{Value: -val.Value}
	default:
		return &ErrorValue{Message: fmt.Sprintf("unknown operator: -%s", right.Type())}
	}
}

func (e *Evaluator) evalNotPrefixExpression(right Value) Value {
	return &BooleanValue{Value: !IsTruthy(right)}
}

func (e *Evaluator) evalInfixExpression(node *InfixExpression, env *Environment) Value {
	left := e.Eval(node.Left, env)
	if isError(left) {
		return left
	}

	right := e.Eval(node.Right, env)
	if isError(right) {
		return right
	}

	left = UnwrapValue(left)
	right = UnwrapValue(right)

	switch {
	case node.Operator == "and":
		return &BooleanValue{Value: IsTruthy(left) && IsTruthy(right)}
	case node.Operator == "or":
		return &BooleanValue{Value: IsTruthy(left) || IsTruthy(right)}
	case node.Operator == "is":
		return &BooleanValue{Value: valuesEqual(left, right)}
	}

	leftInt, leftIsInt := left.(*IntegerValue)
	rightInt, rightIsInt := right.(*IntegerValue)
	if leftIsInt && rightIsInt {
		return e.evalIntegerInfixExpression(node.Operator, leftInt.Value, rightInt.Value)
	}

	leftFloat, leftIsFloat := left.(*FloatValue)
	rightFloat, rightIsFloat := right.(*FloatValue)
	if leftIsFloat && rightIsFloat {
		return e.evalFloatInfixExpression(node.Operator, leftFloat.Value, rightFloat.Value)
	}
	if leftIsInt && rightIsFloat {
		return e.evalFloatInfixExpression(node.Operator, float64(leftInt.Value), rightFloat.Value)
	}
	if leftIsFloat && rightIsInt {
		return e.evalFloatInfixExpression(node.Operator, leftFloat.Value, float64(rightInt.Value))
	}

	leftStr, leftIsStr := left.(*StringValue)
	rightStr, rightIsStr := right.(*StringValue)
	if leftIsStr && rightIsStr {
		return e.evalStringInfixExpression(node.Operator, leftStr.Value, rightStr.Value)
	}

	return &ErrorValue{Message: fmt.Sprintf("type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())}
}

func (e *Evaluator) evalIntegerInfixExpression(op string, left, right int64) Value {
	switch op {
	case "+":
		return &IntegerValue{Value: left + right}
	case "-":
		return &IntegerValue{Value: left - right}
	case "*":
		return &IntegerValue{Value: left * right}
	case "/":
		if right == 0 {
			return &ErrorValue{Message: "division by zero"}
		}
		return &IntegerValue{Value: left / right}
	case "%":
		if right == 0 {
			return &ErrorValue{Message: "division by zero"}
		}
		return &IntegerValue{Value: left % right}
	case ">":
		return &BooleanValue{Value: left > right}
	case "<":
		return &BooleanValue{Value: left < right}
	case ">=":
		return &BooleanValue{Value: left >= right}
	case "<=":
		return &BooleanValue{Value: left <= right}
	default:
		return &ErrorValue{Message: fmt.Sprintf("unknown operator: Integer %s Integer", op)}
	}
}

func (e *Evaluator) evalFloatInfixExpression(op string, left, right float64) Value {
	switch op {
	case "+":
		return &FloatValue{Value: left + right}
	case "-":
		return &FloatValue{Value: left - right}
	case "*":
		return &FloatValue{Value: left * right}
	case "/":
		if right == 0 {
			return &ErrorValue{Message: "division by zero"}
		}
		return &FloatValue{Value: left / right}
	case ">":
		return &BooleanValue{Value: left > right}
	case "<":
		return &BooleanValue{Value: left < right}
	case ">=":
		return &BooleanValue{Value: left >= right}
	case "<=":
		return &BooleanValue{Value: left <= right}
	default:
		return &ErrorValue{Message: fmt.Sprintf("unknown operator: Float %s Float", op)}
	}
}

func (e *Evaluator) evalStringInfixExpression(op string, left, right string) Value {
	switch op {
	case "+":
		return &StringValue{Value: left + right}
	case ">":
		return &BooleanValue{Value: left > right}
	case "<":
		return &BooleanValue{Value: left < right}
	case ">=":
		return &BooleanValue{Value: left >= right}
	case "<=":
		return &BooleanValue{Value: left <= right}
	default:
		return &ErrorValue{Message: fmt.Sprintf("unknown operator: String %s String", op)}
	}
}

func (e *Evaluator) evalAssignmentExpression(node *AssignmentExpression, env *Environment) Value {
	val := e.Eval(node.Value, env)
	if isError(val) {
		return val
	}

	existing, ok := env.Get(node.Name.Value)
	if !ok {
		return &ErrorValue{Message: fmt.Sprintf("undefined: %s", node.Name.Value)}
	}

	mut, isMutable := existing.(*MutableValue)
	if !isMutable {
		return &ErrorValue{Message: fmt.Sprintf("%s is not mutable", node.Name.Value)}
	}

	mut.Value = UnwrapValue(val)
	return mut.Value
}

func (e *Evaluator) evalIfExpression(node *IfExpression, env *Environment) Value {
	condition := e.Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if IsTruthy(condition) {
		return e.Eval(node.Consequence, NewEnclosedEnvironment(env))
	} else if node.Alternative != nil {
		return e.Eval(node.Alternative, NewEnclosedEnvironment(env))
	}

	return &NullValue{}
}

func (e *Evaluator) evalFunctionLiteral(node *FunctionLiteral, env *Environment) Value {
	params := make([]*FunctionParameter, len(node.Parameters))
	for i, p := range node.Parameters {
		params[i] = &FunctionParameter{Name: p}
	}

	return &FunctionValue{
		Parameters: params,
		Body:       nil,
		LambdaBody: node.Body,
		Env:        env,
		IsLambda:   true,
	}
}

func (e *Evaluator) evalCallExpression(node *CallExpression, env *Environment) Value {
	// Check if it's a method call
	if member, ok := node.Function.(*MemberExpression); ok {
		return e.evalMethodCall(member, node.Arguments, env)
	}

	function := e.Eval(node.Function, env)

	args := e.evalExpressions(node.Arguments, env)

	return e.applyFunction(function, args, env)
}

func (e *Evaluator) evalMethodCall(member *MemberExpression, args []Expression, env *Environment) Value {
	obj := e.Eval(member.Object, env)

	methodName := member.Member.Value
	argValues := e.evalExpressions(args, env)

	// Check for built-in methods
	result := e.evalBuiltinMethod(obj, methodName, argValues, env)
	if result != nil {
		return result
	}

	// Check for extension methods
	typeName := obj.Type()
	if extMethods, ok := e.extensions[typeName]; ok {
		if method, ok := extMethods[methodName]; ok {
			// Create new environment with 'this' bound to the object
			extEnv := NewEnclosedEnvironment(method.Env)
			extEnv.Set("this", obj)

			// Bind parameters
			for i, param := range method.Parameters {
				if i < len(argValues) {
					extEnv.Set(param.Name.Value, argValues[i])
				}
			}

			// Evaluate the method body directly
			result := e.Eval(method.Body, extEnv)
			return e.unwrapReturnValue(result)
		}
	}

	return &ErrorValue{Message: fmt.Sprintf("undefined method %s on %s", methodName, typeName)}
}

func (e *Evaluator) evalBuiltinMethod(obj Value, method string, args []Value, env *Environment) Value {
	obj = UnwrapValue(obj)

	switch val := obj.(type) {
	case *ListValue:
		return e.evalListMethod(val, method, args, env)
	case *MapValue:
		return e.evalMapMethod(val, method, args, env)
	case *StringValue:
		return e.evalStringMethod(val, method, args)
	case *ResultValue:
		return e.evalResultMethod(val, method, args, env)
	case *OptionValue:
		return e.evalOptionMethod(val, method, args, env)
	case *ModuleValue:
		if member, ok := val.Exports.Get(method); ok {
			return member
		}
		return nil
	}

	return nil
}

func (e *Evaluator) evalListMethod(list *ListValue, method string, args []Value, env *Environment) Value {
	switch method {
	case "length":
		return listLength(list)
	case "get":
		if len(args) != 1 {
			return &ErrorValue{Message: "get() requires 1 argument"}
		}
		idx, ok := UnwrapValue(args[0]).(*IntegerValue)
		if !ok {
			return &ErrorValue{Message: "get() argument must be an integer"}
		}
		return listGet(list, idx.Value)
	case "append":
		if len(args) != 1 {
			return &ErrorValue{Message: "append() requires 1 argument"}
		}
		return listAppend(list, args[0])
	case "map":
		if len(args) != 1 {
			return &ErrorValue{Message: "map() requires 1 argument"}
		}
		fn, ok := args[0].(*FunctionValue)
		if !ok {
			return &ErrorValue{Message: "map() argument must be a function"}
		}
		return listMap(list, fn, e, env)
	case "filter":
		if len(args) != 1 {
			return &ErrorValue{Message: "filter() requires 1 argument"}
		}
		fn, ok := args[0].(*FunctionValue)
		if !ok {
			return &ErrorValue{Message: "filter() argument must be a function"}
		}
		return listFilter(list, fn, e, env)
	case "reduce":
		if len(args) != 2 {
			return &ErrorValue{Message: "reduce() requires 2 arguments"}
		}
		fn, ok := args[0].(*FunctionValue)
		if !ok {
			return &ErrorValue{Message: "reduce() first argument must be a function"}
		}
		return listReduce(list, fn, args[1], e, env)
	case "find":
		if len(args) != 1 {
			return &ErrorValue{Message: "find() requires 1 argument"}
		}
		fn, ok := args[0].(*FunctionValue)
		if !ok {
			return &ErrorValue{Message: "find() argument must be a function"}
		}
		return listFind(list, fn, e, env)
	case "contains":
		if len(args) != 1 {
			return &ErrorValue{Message: "contains() requires 1 argument"}
		}
		return &BooleanValue{Value: listContains(list, args[0])}
	}
	return nil
}

func (e *Evaluator) evalMapMethod(m *MapValue, method string, args []Value, env *Environment) Value {
	switch method {
	case "get":
		if len(args) != 1 {
			return &ErrorValue{Message: "get() requires 1 argument"}
		}
		key, ok := UnwrapValue(args[0]).(*StringValue)
		if !ok {
			return &ErrorValue{Message: "get() argument must be a string"}
		}
		return mapGet(m, key.Value)
	case "insert":
		if len(args) != 2 {
			return &ErrorValue{Message: "insert() requires 2 arguments"}
		}
		key, ok := UnwrapValue(args[0]).(*StringValue)
		if !ok {
			return &ErrorValue{Message: "insert() first argument must be a string"}
		}
		return mapInsert(m, key.Value, args[1])
	case "remove":
		if len(args) != 1 {
			return &ErrorValue{Message: "remove() requires 1 argument"}
		}
		key, ok := UnwrapValue(args[0]).(*StringValue)
		if !ok {
			return &ErrorValue{Message: "remove() argument must be a string"}
		}
		return mapRemove(m, key.Value)
	case "keys":
		return mapKeys(m)
	case "values":
		return mapValues(m)
	case "contains":
		if len(args) != 1 {
			return &ErrorValue{Message: "contains() requires 1 argument"}
		}
		key, ok := UnwrapValue(args[0]).(*StringValue)
		if !ok {
			return &ErrorValue{Message: "contains() argument must be a string"}
		}
		return &BooleanValue{Value: mapContains(m, key.Value)}
	}
	return nil
}

func (e *Evaluator) evalStringMethod(s *StringValue, method string, args []Value) Value {
	switch method {
	case "length":
		return stringLength(s)
	case "split":
		if len(args) != 1 {
			return &ErrorValue{Message: "split() requires 1 argument"}
		}
		sep, ok := UnwrapValue(args[0]).(*StringValue)
		if !ok {
			return &ErrorValue{Message: "split() argument must be a string"}
		}
		return stringSplit(s, sep.Value)
	case "contains":
		if len(args) != 1 {
			return &ErrorValue{Message: "contains() requires 1 argument"}
		}
		substr, ok := UnwrapValue(args[0]).(*StringValue)
		if !ok {
			return &ErrorValue{Message: "contains() argument must be a string"}
		}
		return &BooleanValue{Value: stringContains(s, substr.Value)}
	case "trim":
		return stringTrim(s)
	case "upper":
		return stringUpper(s)
	case "lower":
		return stringLower(s)
	}
	return nil
}

func (e *Evaluator) evalResultMethod(r *ResultValue, method string, args []Value, env *Environment) Value {
	switch method {
	case "then":
		if len(args) != 1 {
			return &ErrorValue{Message: "then() requires 1 argument"}
		}
		if !r.IsOk {
			return r // Short-circuit on error
		}
		fn, ok := args[0].(*FunctionValue)
		if !ok {
			return &ErrorValue{Message: "then() argument must be a function"}
		}
		result := e.applyFunction(fn, []Value{r.Value}, env)
		// If the function returns a Result, return it; otherwise wrap in Ok
		if res, ok := result.(*ResultValue); ok {
			return res
		}
		return &ResultValue{IsOk: true, Value: result}
	case "map":
		if len(args) != 1 {
			return &ErrorValue{Message: "map() requires 1 argument"}
		}
		if !r.IsOk {
			return r // Short-circuit on error
		}
		fn, ok := args[0].(*FunctionValue)
		if !ok {
			return &ErrorValue{Message: "map() argument must be a function"}
		}
		result := e.applyFunction(fn, []Value{r.Value}, env)
		return &ResultValue{IsOk: true, Value: result}
	case "unwrap":
		if !r.IsOk {
			return r.Error
		}
		return r.Value
	case "unwrapOr":
		if len(args) != 1 {
			return &ErrorValue{Message: "unwrapOr() requires 1 argument"}
		}
		if !r.IsOk {
			return args[0]
		}
		return r.Value
	}
	return nil
}

func (e *Evaluator) evalOptionMethod(o *OptionValue, method string, args []Value, env *Environment) Value {
	switch method {
	case "unwrap":
		if !o.IsSome {
			return &ErrorValue{Message: "called unwrap on None"}
		}
		return o.Value
	case "unwrapOr":
		if len(args) != 1 {
			return &ErrorValue{Message: "unwrapOr() requires 1 argument"}
		}
		if !o.IsSome {
			return args[0]
		}
		return o.Value
	case "map":
		if len(args) != 1 {
			return &ErrorValue{Message: "map() requires 1 argument"}
		}
		if !o.IsSome {
			return o // Return None
		}
		fn, ok := args[0].(*FunctionValue)
		if !ok {
			return &ErrorValue{Message: "map() argument must be a function"}
		}
		result := e.applyFunction(fn, []Value{o.Value}, env)
		return &OptionValue{IsSome: true, Value: result}
	case "isSome":
		return &BooleanValue{Value: o.IsSome}
	case "isNone":
		return &BooleanValue{Value: !o.IsSome}
	}
	return nil
}

func (e *Evaluator) evalExpressions(exprs []Expression, env *Environment) []Value {
	result := make([]Value, len(exprs))
	for i, expr := range exprs {
		evaluated := e.Eval(expr, env)
		result[i] = evaluated
	}
	return result
}

func (e *Evaluator) applyFunction(fn Value, args []Value, callerEnv *Environment) Value {
	switch function := fn.(type) {
	case *FunctionValue:
		oldFn := e.currentFn
		e.currentFn = function.Name

		extendedEnv := e.extendFunctionEnv(function, args)
		var evaluated Value

		if function.IsLambda && function.LambdaBody != nil {
			evaluated = e.Eval(function.LambdaBody, extendedEnv)
		} else {
			evaluated = e.Eval(function.Body, extendedEnv)
		}

		e.currentFn = oldFn
		return e.unwrapReturnValue(evaluated)

	case *BuiltinFunction:
		return function.Fn(args...)

	case *StructDefinition:
		// Struct instantiation like User { ... } is handled elsewhere
		// This is for when a struct is called like a function (which shouldn't happen)
		return &ErrorValue{Message: fmt.Sprintf("%s is not callable", fn.Type())}

	default:
		return &ErrorValue{Message: fmt.Sprintf("not a function: %s", fn.Type())}
	}
}

func (e *Evaluator) extendFunctionEnv(fn *FunctionValue, args []Value) *Environment {
	env := NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		if i < len(args) {
			env.Set(param.Name.Value, args[i])
		}
	}
	return env
}

func (e *Evaluator) unwrapReturnValue(val Value) Value {
	if returnValue, ok := val.(*ReturnValue); ok {
		return returnValue.Value
	}
	return val
}

func (e *Evaluator) evalMemberExpression(node *MemberExpression, env *Environment) Value {
	obj := e.Eval(node.Object, env)
	if isError(obj) {
		return obj
	}

	// Handle struct field access
	if structVal, ok := UnwrapValue(obj).(*StructValue); ok {
		if val, ok := structVal.Fields[node.Member.Value]; ok {
			return val
		}
		return &ErrorValue{Message: fmt.Sprintf("undefined field %s on %s", node.Member.Value, structVal.Type())}
	}

	// Handle module access
	if mod, ok := obj.(*ModuleValue); ok {
		if val, ok := mod.Exports.Get(node.Member.Value); ok {
			return val
		}
		return &ErrorValue{Message: fmt.Sprintf("undefined export %s in module %s", node.Member.Value, mod.Name)}
	}

	return &ErrorValue{Message: fmt.Sprintf("cannot access member of %s", obj.Type())}
}

func (e *Evaluator) evalIndexExpression(node *IndexExpression, env *Environment) Value {
	left := e.Eval(node.Left, env)
	if isError(left) {
		return left
	}

	index := e.Eval(node.Index, env)
	if isError(index) {
		return index
	}

	left = UnwrapValue(left)
	index = UnwrapValue(index)

	switch obj := left.(type) {
	case *ListValue:
		idx, ok := index.(*IntegerValue)
		if !ok {
			return &ErrorValue{Message: "list index must be an integer"}
		}
		if idx.Value < 0 || idx.Value >= int64(len(obj.Elements)) {
			return &ErrorValue{Message: "index out of bounds"}
		}
		return obj.Elements[idx.Value]

	case *MapValue:
		key, ok := index.(*StringValue)
		if !ok {
			return &ErrorValue{Message: "map key must be a string"}
		}
		if val, ok := obj.Pairs[key.Value]; ok {
			return val
		}
		return &OptionValue{IsSome: false}

	case *StringValue:
		idx, ok := index.(*IntegerValue)
		if !ok {
			return &ErrorValue{Message: "string index must be an integer"}
		}
		if idx.Value < 0 || idx.Value >= int64(len(obj.Value)) {
			return &ErrorValue{Message: "index out of bounds"}
		}
		return &StringValue{Value: string(obj.Value[idx.Value])}

	default:
		return &ErrorValue{Message: fmt.Sprintf("cannot index %s", left.Type())}
	}
}

func (e *Evaluator) evalListLiteral(node *ListLiteral, env *Environment) Value {
	elements := e.evalExpressions(node.Elements, env)
	if len(elements) == 1 && isError(elements[0]) {
		return elements[0]
	}
	return &ListValue{Elements: elements}
}

func (e *Evaluator) evalMapLiteral(node *MapLiteral, env *Environment) Value {
	pairs := make(map[string]Value)

	for keyNode, valueNode := range node.Pairs {
		key := e.Eval(keyNode, env)
		if isError(key) {
			return key
		}

		keyStr, ok := UnwrapValue(key).(*StringValue)
		if !ok {
			return &ErrorValue{Message: "map key must be a string"}
		}

		value := e.Eval(valueNode, env)
		if isError(value) {
			return value
		}

		pairs[keyStr.Value] = value
	}

	return &MapValue{Pairs: pairs}
}

func (e *Evaluator) evalStructLiteral(node *StructLiteral, env *Environment) Value {
	def, ok := e.structs[node.StructName.Value]
	if !ok {
		return &ErrorValue{Message: fmt.Sprintf("undefined struct: %s", node.StructName.Value)}
	}

	fields := make(map[string]Value)
	for name, valueNode := range node.Fields {
		value := e.Eval(valueNode, env)
		if isError(value) {
			return value
		}
		fields[name] = value
	}

	return &StructValue{
		Definition: def,
		Fields:     fields,
	}
}

func (e *Evaluator) evalWithExpression(node *WithExpression, env *Environment) Value {
	obj := e.Eval(node.Object, env)
	if isError(obj) {
		return obj
	}

	structVal, ok := UnwrapValue(obj).(*StructValue)
	if !ok {
		return &ErrorValue{Message: fmt.Sprintf("with can only be used on structs, got %s", obj.Type())}
	}

	updates := make(map[string]Value)
	for name, valueNode := range node.Updates {
		value := e.Eval(valueNode, env)
		if isError(value) {
			return value
		}
		updates[name] = value
	}

	return structVal.With(updates)
}

func (e *Evaluator) evalOptionExpression(node *OptionExpression, env *Environment) Value {
	if !node.IsSome {
		return &OptionValue{IsSome: false}
	}

	value := e.Eval(node.Value, env)
	if isError(value) {
		return value
	}

	return &OptionValue{IsSome: true, Value: value}
}

func (e *Evaluator) evalResultExpression(node *ResultExpression, env *Environment) Value {
	value := e.Eval(node.Value, env)
	if isError(value) {
		return value
	}

	if node.IsOk {
		return &ResultValue{IsOk: true, Value: value}
	}

	errVal, ok := value.(*StringValue)
	if ok {
		return &ResultValue{IsOk: false, Error: &ErrorValue{
			Method:  e.currentFn,
			Message: errVal.Value,
		}}
	}

	return &ResultValue{IsOk: false, Error: &ErrorValue{
		Method:  e.currentFn,
		Message: value.String(),
	}}
}

func (e *Evaluator) evalMatchExpression(node *MatchExpression, env *Environment) Value {
	value := e.Eval(node.Value, env)
	if isError(value) {
		return value
	}

	for _, matchCase := range node.Cases {
		if matched, bindings := e.matchPattern(value, matchCase, env); matched {
			caseEnv := NewEnclosedEnvironment(env)
			for name, val := range bindings {
				caseEnv.Set(name, val)
			}
			return e.Eval(matchCase.Body, caseEnv)
		}
	}

	return &NullValue{}
}

func (e *Evaluator) matchPattern(value Value, matchCase *MatchCase, env *Environment) (bool, map[string]Value) {
	bindings := make(map[string]Value)

	switch pat := matchCase.Pattern.(type) {
	case *OptionExpression:
		opt, ok := value.(*OptionValue)
		if !ok {
			return false, nil
		}
		if pat.IsSome != opt.IsSome {
			return false, nil
		}
		if pat.IsSome && matchCase.BindingVar != nil {
			bindings[matchCase.BindingVar.Value] = opt.Value
		}
		return true, bindings

	case *ResultExpression:
		res, ok := value.(*ResultValue)
		if !ok {
			return false, nil
		}
		if pat.IsOk != res.IsOk {
			return false, nil
		}
		if matchCase.BindingVar != nil {
			if res.IsOk {
				bindings[matchCase.BindingVar.Value] = res.Value
			} else {
				bindings[matchCase.BindingVar.Value] = res.Error
			}
		}
		return true, bindings

	case *Identifier:
		// Wildcard pattern - matches anything
		bindings[pat.Value] = value
		return true, bindings
	}

	return false, nil
}

func (e *Evaluator) evalMutableExpression(node *MutableExpression, env *Environment) Value {
	value := e.Eval(node.Value, env)
	if isError(value) {
		return value
	}
	return &MutableValue{Value: UnwrapValue(value)}
}

func isError(val Value) bool {
	if val == nil {
		return false
	}
	_, ok := val.(*ErrorValue)
	return ok
}
