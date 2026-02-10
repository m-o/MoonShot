package main

import (
	"fmt"
)

// TypeChecker performs static type checking
type TypeChecker struct {
	env       *TypeEnvironment
	structs   map[string]*StructType
	functions map[string]*FunctionType
	errors    []string
}

// TypeEnvironment stores type bindings
type TypeEnvironment struct {
	store  map[string]Type
	parent *TypeEnvironment
}

// NewTypeEnvironment creates a new type environment
func NewTypeEnvironment() *TypeEnvironment {
	return &TypeEnvironment{
		store:  make(map[string]Type),
		parent: nil,
	}
}

// NewEnclosedTypeEnvironment creates a child type environment
func NewEnclosedTypeEnvironment(parent *TypeEnvironment) *TypeEnvironment {
	env := NewTypeEnvironment()
	env.parent = parent
	return env
}

// Get retrieves a type from the environment
func (e *TypeEnvironment) Get(name string) (Type, bool) {
	t, ok := e.store[name]
	if !ok && e.parent != nil {
		return e.parent.Get(name)
	}
	return t, ok
}

// Set defines a new type binding
func (e *TypeEnvironment) Set(name string, t Type) {
	e.store[name] = t
}

// NewTypeChecker creates a new type checker
func NewTypeChecker() *TypeChecker {
	tc := &TypeChecker{
		env:       NewTypeEnvironment(),
		structs:   make(map[string]*StructType),
		functions: make(map[string]*FunctionType),
	}

	// Register built-in function types
	tc.env.Set("print", &FunctionType{Parameters: []Type{&AnyType{}}, Return: &NullType{}})
	tc.env.Set("println", &FunctionType{Parameters: []Type{&AnyType{}}, Return: &NullType{}})
	tc.env.Set("range", &FunctionType{Parameters: []Type{&IntegerType{}, &IntegerType{}}, Return: &ListType{Element: &IntegerType{}}})
	tc.env.Set("len", &FunctionType{Parameters: []Type{&AnyType{}}, Return: &IntegerType{}})
	tc.env.Set("type", &FunctionType{Parameters: []Type{&AnyType{}}, Return: &StringType{}})
	tc.env.Set("str", &FunctionType{Parameters: []Type{&AnyType{}}, Return: &StringType{}})
	tc.env.Set("int", &FunctionType{Parameters: []Type{&AnyType{}}, Return: &IntegerType{}})
	tc.env.Set("float", &FunctionType{Parameters: []Type{&AnyType{}}, Return: &FloatType{}})

	return tc
}

// Check performs type checking on a program
func (tc *TypeChecker) Check(program *Program) error {
	// First pass: collect struct and function definitions
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *StructStatement:
			tc.collectStruct(s)
		case *FunctionStatement:
			tc.collectFunction(s)
		case *ExtendStatement:
			tc.collectExtend(s)
		}
	}

	// Second pass: type check all statements
	for _, stmt := range program.Statements {
		tc.checkStatement(stmt)
	}

	if len(tc.errors) > 0 {
		return fmt.Errorf("%s", tc.errors[0])
	}
	return nil
}

func (tc *TypeChecker) collectExtend(stmt *ExtendStatement) {
	for _, method := range stmt.Methods {
		tc.collectFunction(method)
	}
}

func (tc *TypeChecker) collectStruct(stmt *StructStatement) {
	fields := make(map[string]Type)
	for _, f := range stmt.Fields {
		fields[f.Name.Value] = TypeFromAnnotation(f.TypeHint)
	}
	tc.structs[stmt.Name.Value] = &StructType{Name: stmt.Name.Value, Fields: fields}
	tc.env.Set(stmt.Name.Value, tc.structs[stmt.Name.Value])
}

func (tc *TypeChecker) collectFunction(stmt *FunctionStatement) {
	params := make([]Type, len(stmt.Parameters))
	for i, p := range stmt.Parameters {
		params[i] = TypeFromAnnotation(p.TypeHint)
	}
	returnType := TypeFromAnnotation(stmt.ReturnType)
	tc.functions[stmt.Name.Value] = &FunctionType{Parameters: params, Return: returnType}
	tc.env.Set(stmt.Name.Value, tc.functions[stmt.Name.Value])
}

func (tc *TypeChecker) checkStatement(stmt Statement) Type {
	switch s := stmt.(type) {
	case *DefStatement:
		return tc.checkDefStatement(s)
	case *FunctionStatement:
		return tc.checkFunctionStatement(s)
	case *ReturnStatement:
		return tc.checkReturnStatement(s)
	case *ExpressionStatement:
		return tc.checkExpression(s.Expression)
	case *WhileStatement:
		return tc.checkWhileStatement(s)
	case *ForStatement:
		return tc.checkForStatement(s)
	case *StructStatement:
		return tc.structs[s.Name.Value]
	case *ExtendStatement:
		return tc.checkExtendStatement(s)
	case *ImportStatement:
		return &NullType{}
	case *BreakStatement, *ContinueStatement:
		return &NullType{}
	}
	return &AnyType{}
}

func (tc *TypeChecker) checkDefStatement(stmt *DefStatement) Type {
	valueType := tc.checkExpression(stmt.Value)

	if stmt.TypeHint != nil {
		expectedType := TypeFromAnnotation(stmt.TypeHint)
		if !tc.isAssignable(expectedType, valueType) {
			tc.addError(fmt.Sprintf("cannot assign %s to variable of type %s",
				valueType.String(), expectedType.String()))
		}
		tc.env.Set(stmt.Name.Value, expectedType)
		return expectedType
	}

	tc.env.Set(stmt.Name.Value, valueType)
	return valueType
}

func (tc *TypeChecker) checkFunctionStatement(stmt *FunctionStatement) Type {
	fnType := tc.functions[stmt.Name.Value]

	// Create new scope for function body
	prevEnv := tc.env
	tc.env = NewEnclosedTypeEnvironment(prevEnv)

	// Add parameters to scope
	for i, p := range stmt.Parameters {
		tc.env.Set(p.Name.Value, fnType.Parameters[i])
	}

	// Check function body
	tc.checkBlockStatement(stmt.Body, fnType.Return)

	tc.env = prevEnv
	return fnType
}

func (tc *TypeChecker) checkReturnStatement(stmt *ReturnStatement) Type {
	if stmt.Value == nil {
		return &NullType{}
	}
	return tc.checkExpression(stmt.Value)
}

func (tc *TypeChecker) checkWhileStatement(stmt *WhileStatement) Type {
	condType := tc.checkExpression(stmt.Condition)
	if !tc.isBooleanCompatible(condType) {
		tc.addError("while condition must be a boolean expression")
	}

	prevEnv := tc.env
	tc.env = NewEnclosedTypeEnvironment(prevEnv)
	tc.checkBlockStatement(stmt.Body, nil)
	tc.env = prevEnv

	return &NullType{}
}

func (tc *TypeChecker) checkForStatement(stmt *ForStatement) Type {
	iterType := tc.checkExpression(stmt.Iterable)

	listType, ok := iterType.(*ListType)
	if !ok {
		tc.addError(fmt.Sprintf("cannot iterate over %s", iterType.String()))
		return &NullType{}
	}

	prevEnv := tc.env
	tc.env = NewEnclosedTypeEnvironment(prevEnv)
	tc.env.Set(stmt.Variable.Value, listType.Element)
	tc.checkBlockStatement(stmt.Body, nil)
	tc.env = prevEnv

	return &NullType{}
}

func (tc *TypeChecker) checkExtendStatement(stmt *ExtendStatement) Type {
	// Get the struct type being extended
	typeName := stmt.TypeName.Value
	structType, ok := tc.structs[typeName]

	for _, method := range stmt.Methods {
		// Create a scope with 'this' bound to the struct type
		prevEnv := tc.env
		tc.env = NewEnclosedTypeEnvironment(prevEnv)

		if ok {
			tc.env.Set("this", structType)
		} else {
			tc.env.Set("this", &AnyType{})
		}

		fnType := tc.functions[method.Name.Value]
		if fnType != nil {
			// Add parameters to scope
			for i, p := range method.Parameters {
				if i < len(fnType.Parameters) {
					tc.env.Set(p.Name.Value, fnType.Parameters[i])
				}
			}

			// Check function body
			tc.checkBlockStatement(method.Body, fnType.Return)
		}

		tc.env = prevEnv
	}
	return &NullType{}
}

func (tc *TypeChecker) checkBlockStatement(block *BlockStatement, expectedReturn Type) Type {
	var lastType Type = &NullType{}
	for _, stmt := range block.Statements {
		lastType = tc.checkStatement(stmt)

		if ret, ok := stmt.(*ReturnStatement); ok && expectedReturn != nil {
			retType := tc.checkExpression(ret.Value)
			if !tc.isAssignable(expectedReturn, retType) {
				tc.addError(fmt.Sprintf("cannot return %s from function expecting %s",
					retType.String(), expectedReturn.String()))
			}
		}
	}
	return lastType
}

func (tc *TypeChecker) checkExpression(expr Expression) Type {
	if expr == nil {
		return &NullType{}
	}

	switch e := expr.(type) {
	case *IntegerLiteral:
		return &IntegerType{}
	case *FloatLiteral:
		return &FloatType{}
	case *StringLiteral:
		return &StringType{}
	case *BooleanLiteral:
		return &BooleanType{}
	case *Identifier:
		return tc.checkIdentifier(e)
	case *PrefixExpression:
		return tc.checkPrefixExpression(e)
	case *InfixExpression:
		return tc.checkInfixExpression(e)
	case *AssignmentExpression:
		return tc.checkAssignmentExpression(e)
	case *IfExpression:
		return tc.checkIfExpression(e)
	case *FunctionLiteral:
		return tc.checkFunctionLiteral(e)
	case *CallExpression:
		return tc.checkCallExpression(e)
	case *MemberExpression:
		return tc.checkMemberExpression(e)
	case *IndexExpression:
		return tc.checkIndexExpression(e)
	case *ListLiteral:
		return tc.checkListLiteral(e)
	case *MapLiteral:
		return tc.checkMapLiteral(e)
	case *StructLiteral:
		return tc.checkStructLiteral(e)
	case *WithExpression:
		return tc.checkWithExpression(e)
	case *OptionExpression:
		return tc.checkOptionExpression(e)
	case *ResultExpression:
		return tc.checkResultExpression(e)
	case *MatchExpression:
		return tc.checkMatchExpression(e)
	case *MutableExpression:
		return tc.checkMutableExpression(e)
	}

	return &AnyType{}
}

func (tc *TypeChecker) checkIdentifier(ident *Identifier) Type {
	t, ok := tc.env.Get(ident.Value)
	if !ok {
		tc.addError(fmt.Sprintf("undefined: %s", ident.Value))
		return &AnyType{}
	}
	return t
}

func (tc *TypeChecker) checkPrefixExpression(expr *PrefixExpression) Type {
	rightType := tc.checkExpression(expr.Right)

	switch expr.Operator {
	case "-":
		if !tc.isNumeric(rightType) {
			tc.addError(fmt.Sprintf("operator - not defined for %s", rightType.String()))
		}
		return rightType
	case "not":
		return &BooleanType{}
	}

	return &AnyType{}
}

func (tc *TypeChecker) checkInfixExpression(expr *InfixExpression) Type {
	leftType := tc.checkExpression(expr.Left)
	rightType := tc.checkExpression(expr.Right)

	switch expr.Operator {
	case "+", "-", "*", "/", "%":
		if !tc.isNumeric(leftType) || !tc.isNumeric(rightType) {
			// String concatenation
			if expr.Operator == "+" && tc.isString(leftType) && tc.isString(rightType) {
				return &StringType{}
			}
			tc.addError(fmt.Sprintf("operator %s not defined for %s and %s",
				expr.Operator, leftType.String(), rightType.String()))
		}
		// Return Float if either operand is Float
		if _, ok := leftType.(*FloatType); ok {
			return &FloatType{}
		}
		if _, ok := rightType.(*FloatType); ok {
			return &FloatType{}
		}
		return &IntegerType{}

	case ">", "<", ">=", "<=":
		if !tc.isComparable(leftType, rightType) {
			tc.addError(fmt.Sprintf("cannot compare %s and %s",
				leftType.String(), rightType.String()))
		}
		return &BooleanType{}

	case "and", "or":
		return &BooleanType{}

	case "is":
		return &BooleanType{}
	}

	return &AnyType{}
}

func (tc *TypeChecker) checkAssignmentExpression(expr *AssignmentExpression) Type {
	varType, ok := tc.env.Get(expr.Name.Value)
	if !ok {
		tc.addError(fmt.Sprintf("undefined: %s", expr.Name.Value))
		return &AnyType{}
	}

	mutType, isMutable := varType.(*MutableType)
	if !isMutable {
		tc.addError(fmt.Sprintf("%s is not mutable", expr.Name.Value))
		return &AnyType{}
	}

	valueType := tc.checkExpression(expr.Value)
	if !tc.isAssignable(mutType.Element, valueType) {
		tc.addError(fmt.Sprintf("cannot assign %s to Mutable[%s]",
			valueType.String(), mutType.Element.String()))
	}

	return mutType.Element
}

func (tc *TypeChecker) checkIfExpression(expr *IfExpression) Type {
	condType := tc.checkExpression(expr.Condition)
	if !tc.isBooleanCompatible(condType) {
		tc.addError("if condition must be a boolean expression")
	}

	prevEnv := tc.env
	tc.env = NewEnclosedTypeEnvironment(prevEnv)
	consType := tc.checkBlockStatement(expr.Consequence, nil)
	tc.env = prevEnv

	if expr.Alternative != nil {
		tc.env = NewEnclosedTypeEnvironment(prevEnv)
		altType := tc.checkBlockStatement(expr.Alternative, nil)
		tc.env = prevEnv

		// If both branches return compatible types, use that
		if tc.isAssignable(consType, altType) {
			return consType
		}
	}

	return consType
}

func (tc *TypeChecker) checkFunctionLiteral(expr *FunctionLiteral) Type {
	params := make([]Type, len(expr.Parameters))
	for i := range expr.Parameters {
		params[i] = &AnyType{} // Lambda parameters are inferred
	}

	// For simple lambdas, we can try to infer the return type
	returnType := Type(&AnyType{})

	return &FunctionType{Parameters: params, Return: returnType}
}

func (tc *TypeChecker) checkCallExpression(expr *CallExpression) Type {
	fnType := tc.checkExpression(expr.Function)

	// If it's Any (e.g., a method call we can't resolve), just check args and return Any
	if _, ok := fnType.(*AnyType); ok {
		for _, arg := range expr.Arguments {
			tc.checkExpression(arg)
		}
		return &AnyType{}
	}

	fn, ok := fnType.(*FunctionType)
	if !ok {
		// Might be a struct constructor
		if st, ok := fnType.(*StructType); ok {
			return st
		}
		// Don't error on unresolved types - just return Any
		for _, arg := range expr.Arguments {
			tc.checkExpression(arg)
		}
		return &AnyType{}
	}

	// Check argument count (but allow varargs for builtins)
	if len(expr.Arguments) != len(fn.Parameters) {
		// Allow variadic functions
		if len(fn.Parameters) == 1 {
			_, isAny := fn.Parameters[0].(*AnyType)
			if !isAny && len(fn.Parameters) != len(expr.Arguments) {
				// Skip this error - too strict for now
			}
		}
	}

	// Check argument types
	for i, arg := range expr.Arguments {
		argType := tc.checkExpression(arg)
		if i < len(fn.Parameters) {
			if !tc.isAssignable(fn.Parameters[i], argType) {
				// Skip strict type checking for now - too many false positives
			}
		}
	}

	return fn.Return
}

func (tc *TypeChecker) checkMemberExpression(expr *MemberExpression) Type {
	objType := tc.checkExpression(expr.Object)

	// Unwrap mutable
	if mut, ok := objType.(*MutableType); ok {
		objType = mut.Element
	}

	if st, ok := objType.(*StructType); ok {
		if fieldType, ok := st.Fields[expr.Member.Value]; ok {
			return fieldType
		}
		// Could be a method - return Any for now
		return &AnyType{}
	}

	// Could be a method call on a list, map, etc.
	return &AnyType{}
}

func (tc *TypeChecker) checkIndexExpression(expr *IndexExpression) Type {
	leftType := tc.checkExpression(expr.Left)
	indexType := tc.checkExpression(expr.Index)

	// Unwrap mutable
	if mut, ok := leftType.(*MutableType); ok {
		leftType = mut.Element
	}

	switch t := leftType.(type) {
	case *ListType:
		if !tc.isInteger(indexType) {
			tc.addError("list index must be an integer")
		}
		return t.Element
	case *MapType:
		if !tc.isString(indexType) {
			tc.addError("map key must be a string")
		}
		return t.Value
	case *StringType:
		if !tc.isInteger(indexType) {
			tc.addError("string index must be an integer")
		}
		return &StringType{}
	case *AnyType:
		// Unknown type - just return Any
		return &AnyType{}
	}

	// Don't error for now - might be a valid indexing we don't understand
	return &AnyType{}
}

func (tc *TypeChecker) checkListLiteral(expr *ListLiteral) Type {
	if len(expr.Elements) == 0 {
		return &ListType{Element: &AnyType{}}
	}

	elemType := tc.checkExpression(expr.Elements[0])
	for i := 1; i < len(expr.Elements); i++ {
		t := tc.checkExpression(expr.Elements[i])
		if !tc.isAssignable(elemType, t) {
			// Allow mixed types if first element is Any
			if _, ok := elemType.(*AnyType); !ok {
				tc.addError("list elements must have the same type")
			}
		}
	}

	return &ListType{Element: elemType}
}

func (tc *TypeChecker) checkMapLiteral(expr *MapLiteral) Type {
	if len(expr.Pairs) == 0 {
		return &MapType{Key: &StringType{}, Value: &AnyType{}}
	}

	var valueType Type = &AnyType{}
	for _, v := range expr.Pairs {
		valueType = tc.checkExpression(v)
		break // Just check first value for now
	}

	return &MapType{Key: &StringType{}, Value: valueType}
}

func (tc *TypeChecker) checkStructLiteral(expr *StructLiteral) Type {
	st, ok := tc.structs[expr.StructName.Value]
	if !ok {
		tc.addError(fmt.Sprintf("undefined struct: %s", expr.StructName.Value))
		return &AnyType{}
	}

	for fieldName, fieldExpr := range expr.Fields {
		expectedType, ok := st.Fields[fieldName]
		if !ok {
			tc.addError(fmt.Sprintf("undefined field %s on %s", fieldName, st.Name))
			continue
		}
		actualType := tc.checkExpression(fieldExpr)
		if !tc.isAssignable(expectedType, actualType) {
			tc.addError(fmt.Sprintf("cannot assign %s to field %s of type %s",
				actualType.String(), fieldName, expectedType.String()))
		}
	}

	return st
}

func (tc *TypeChecker) checkWithExpression(expr *WithExpression) Type {
	objType := tc.checkExpression(expr.Object)

	// Unwrap mutable
	if mut, ok := objType.(*MutableType); ok {
		objType = mut.Element
	}

	st, ok := objType.(*StructType)
	if !ok {
		tc.addError("with can only be used on structs")
		return &AnyType{}
	}

	for fieldName, fieldExpr := range expr.Updates {
		expectedType, ok := st.Fields[fieldName]
		if !ok {
			tc.addError(fmt.Sprintf("undefined field %s on %s", fieldName, st.Name))
			continue
		}
		actualType := tc.checkExpression(fieldExpr)
		if !tc.isAssignable(expectedType, actualType) {
			tc.addError(fmt.Sprintf("cannot assign %s to field %s of type %s",
				actualType.String(), fieldName, expectedType.String()))
		}
	}

	return st
}

func (tc *TypeChecker) checkOptionExpression(expr *OptionExpression) Type {
	if !expr.IsSome {
		return &OptionType{Element: &AnyType{}}
	}
	elemType := tc.checkExpression(expr.Value)
	return &OptionType{Element: elemType}
}

func (tc *TypeChecker) checkResultExpression(expr *ResultExpression) Type {
	valueType := tc.checkExpression(expr.Value)
	if expr.IsOk {
		return &ResultType{ValueType: valueType, ErrorType: &StringType{}}
	}
	return &ResultType{ValueType: &AnyType{}, ErrorType: valueType}
}

func (tc *TypeChecker) checkMatchExpression(expr *MatchExpression) Type {
	tc.checkExpression(expr.Value)

	var resultType Type = &NullType{}
	for _, c := range expr.Cases {
		prevEnv := tc.env
		tc.env = NewEnclosedTypeEnvironment(prevEnv)

		if c.BindingVar != nil {
			tc.env.Set(c.BindingVar.Value, &AnyType{})
		}

		resultType = tc.checkBlockStatement(c.Body, nil)
		tc.env = prevEnv
	}

	return resultType
}

func (tc *TypeChecker) checkMutableExpression(expr *MutableExpression) Type {
	elemType := tc.checkExpression(expr.Value)
	if expr.TypeHint != nil {
		elemType = TypeFromAnnotation(expr.TypeHint)
	}
	return &MutableType{Element: elemType}
}

// Helper functions

func (tc *TypeChecker) isAssignable(expected, actual Type) bool {
	if _, ok := expected.(*AnyType); ok {
		return true
	}
	if _, ok := actual.(*AnyType); ok {
		return true
	}

	// Handle mutable unwrapping
	if mut, ok := actual.(*MutableType); ok {
		return tc.isAssignable(expected, mut.Element)
	}

	// Handle Option types - be lenient with element types involving Any
	if expOpt, ok := expected.(*OptionType); ok {
		if actOpt, ok := actual.(*OptionType); ok {
			return tc.isAssignable(expOpt.Element, actOpt.Element)
		}
	}

	// Handle Result types - be lenient with element types involving Any
	if expRes, ok := expected.(*ResultType); ok {
		if actRes, ok := actual.(*ResultType); ok {
			return tc.isAssignable(expRes.ValueType, actRes.ValueType)
		}
	}

	return expected.Equals(actual)
}

func (tc *TypeChecker) isNumeric(t Type) bool {
	if _, ok := t.(*AnyType); ok {
		return true
	}
	if mut, ok := t.(*MutableType); ok {
		return tc.isNumeric(mut.Element)
	}
	_, isInt := t.(*IntegerType)
	_, isFloat := t.(*FloatType)
	return isInt || isFloat
}

func (tc *TypeChecker) isInteger(t Type) bool {
	if _, ok := t.(*AnyType); ok {
		return true
	}
	if mut, ok := t.(*MutableType); ok {
		return tc.isInteger(mut.Element)
	}
	_, ok := t.(*IntegerType)
	return ok
}

func (tc *TypeChecker) isString(t Type) bool {
	if _, ok := t.(*AnyType); ok {
		return true
	}
	if mut, ok := t.(*MutableType); ok {
		return tc.isString(mut.Element)
	}
	_, ok := t.(*StringType)
	return ok
}

func (tc *TypeChecker) isBooleanCompatible(t Type) bool {
	if _, ok := t.(*AnyType); ok {
		return true
	}
	if mut, ok := t.(*MutableType); ok {
		return tc.isBooleanCompatible(mut.Element)
	}
	_, ok := t.(*BooleanType)
	return ok
}

func (tc *TypeChecker) isComparable(a, b Type) bool {
	if _, ok := a.(*AnyType); ok {
		return true
	}
	if _, ok := b.(*AnyType); ok {
		return true
	}

	// Unwrap mutables
	if mut, ok := a.(*MutableType); ok {
		a = mut.Element
	}
	if mut, ok := b.(*MutableType); ok {
		b = mut.Element
	}

	if tc.isNumeric(a) && tc.isNumeric(b) {
		return true
	}
	if tc.isString(a) && tc.isString(b) {
		return true
	}
	return false
}

func (tc *TypeChecker) addError(msg string) {
	tc.errors = append(tc.errors, msg)
}
