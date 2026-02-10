package main

// Type represents a type in the type system
type Type interface {
	typeNode()
	String() string
	Equals(Type) bool
}

// IntegerType represents the Integer type
type IntegerType struct{}

func (t *IntegerType) typeNode()        {}
func (t *IntegerType) String() string   { return "Integer" }
func (t *IntegerType) Equals(o Type) bool {
	_, ok := o.(*IntegerType)
	return ok
}

// FloatType represents the Float type
type FloatType struct{}

func (t *FloatType) typeNode()        {}
func (t *FloatType) String() string   { return "Float" }
func (t *FloatType) Equals(o Type) bool {
	_, ok := o.(*FloatType)
	return ok
}

// StringType represents the String type
type StringType struct{}

func (t *StringType) typeNode()        {}
func (t *StringType) String() string   { return "String" }
func (t *StringType) Equals(o Type) bool {
	_, ok := o.(*StringType)
	return ok
}

// BooleanType represents the Boolean type
type BooleanType struct{}

func (t *BooleanType) typeNode()        {}
func (t *BooleanType) String() string   { return "Boolean" }
func (t *BooleanType) Equals(o Type) bool {
	_, ok := o.(*BooleanType)
	return ok
}

// NullType represents the absence of a value
type NullType struct{}

func (t *NullType) typeNode()        {}
func (t *NullType) String() string   { return "Null" }
func (t *NullType) Equals(o Type) bool {
	_, ok := o.(*NullType)
	return ok
}

// ListType represents List[T]
type ListType struct {
	Element Type
}

func (t *ListType) typeNode()        {}
func (t *ListType) String() string   { return "List[" + t.Element.String() + "]" }
func (t *ListType) Equals(o Type) bool {
	if ot, ok := o.(*ListType); ok {
		return t.Element.Equals(ot.Element)
	}
	return false
}

// MapType represents Map[K, V]
type MapType struct {
	Key   Type
	Value Type
}

func (t *MapType) typeNode()        {}
func (t *MapType) String() string   { return "Map[" + t.Key.String() + ", " + t.Value.String() + "]" }
func (t *MapType) Equals(o Type) bool {
	if ot, ok := o.(*MapType); ok {
		return t.Key.Equals(ot.Key) && t.Value.Equals(ot.Value)
	}
	return false
}

// OptionType represents Option[T]
type OptionType struct {
	Element Type
}

func (t *OptionType) typeNode()        {}
func (t *OptionType) String() string   { return "Option[" + t.Element.String() + "]" }
func (t *OptionType) Equals(o Type) bool {
	if ot, ok := o.(*OptionType); ok {
		return t.Element.Equals(ot.Element)
	}
	return false
}

// ResultType represents Result[T, E]
type ResultType struct {
	ValueType Type
	ErrorType Type
}

func (t *ResultType) typeNode()        {}
func (t *ResultType) String() string {
	return "Result[" + t.ValueType.String() + ", " + t.ErrorType.String() + "]"
}
func (t *ResultType) Equals(o Type) bool {
	if ot, ok := o.(*ResultType); ok {
		return t.ValueType.Equals(ot.ValueType) && t.ErrorType.Equals(ot.ErrorType)
	}
	return false
}

// MutableType represents Mutable[T]
type MutableType struct {
	Element Type
}

func (t *MutableType) typeNode()        {}
func (t *MutableType) String() string   { return "Mutable[" + t.Element.String() + "]" }
func (t *MutableType) Equals(o Type) bool {
	if ot, ok := o.(*MutableType); ok {
		return t.Element.Equals(ot.Element)
	}
	return false
}

// FunctionType represents a function type
type FunctionType struct {
	Parameters []Type
	Return     Type
}

func (t *FunctionType) typeNode()        {}
func (t *FunctionType) String() string {
	params := ""
	for i, p := range t.Parameters {
		if i > 0 {
			params += ", "
		}
		params += p.String()
	}
	return "(" + params + ") -> " + t.Return.String()
}
func (t *FunctionType) Equals(o Type) bool {
	if ot, ok := o.(*FunctionType); ok {
		if len(t.Parameters) != len(ot.Parameters) {
			return false
		}
		for i, p := range t.Parameters {
			if !p.Equals(ot.Parameters[i]) {
				return false
			}
		}
		return t.Return.Equals(ot.Return)
	}
	return false
}

// StructType represents a struct type
type StructType struct {
	Name   string
	Fields map[string]Type
}

func (t *StructType) typeNode()        {}
func (t *StructType) String() string   { return t.Name }
func (t *StructType) Equals(o Type) bool {
	if ot, ok := o.(*StructType); ok {
		return t.Name == ot.Name
	}
	return false
}

// AnyType is a placeholder for unresolved types
type AnyType struct{}

func (t *AnyType) typeNode()           {}
func (t *AnyType) String() string      { return "Any" }
func (t *AnyType) Equals(o Type) bool { return true }

// TypeFromAnnotation converts a type annotation to a Type
func TypeFromAnnotation(ta *TypeAnnotation) Type {
	if ta == nil {
		return &AnyType{}
	}

	switch ta.Name {
	case "Integer":
		return &IntegerType{}
	case "Float":
		return &FloatType{}
	case "String":
		return &StringType{}
	case "Boolean":
		return &BooleanType{}
	case "List":
		if len(ta.TypeParams) > 0 {
			return &ListType{Element: TypeFromAnnotation(ta.TypeParams[0])}
		}
		return &ListType{Element: &AnyType{}}
	case "Map":
		keyType := &StringType{} // Default key type
		valueType := Type(&AnyType{})
		if len(ta.TypeParams) > 0 {
			keyType, _ = TypeFromAnnotation(ta.TypeParams[0]).(*StringType)
			if keyType == nil {
				keyType = &StringType{}
			}
		}
		if len(ta.TypeParams) > 1 {
			valueType = TypeFromAnnotation(ta.TypeParams[1])
		}
		return &MapType{Key: keyType, Value: valueType}
	case "Option":
		if len(ta.TypeParams) > 0 {
			return &OptionType{Element: TypeFromAnnotation(ta.TypeParams[0])}
		}
		return &OptionType{Element: &AnyType{}}
	case "Result":
		valueType := Type(&AnyType{})
		errorType := Type(&StringType{})
		if len(ta.TypeParams) > 0 {
			valueType = TypeFromAnnotation(ta.TypeParams[0])
		}
		if len(ta.TypeParams) > 1 {
			errorType = TypeFromAnnotation(ta.TypeParams[1])
		}
		return &ResultType{ValueType: valueType, ErrorType: errorType}
	case "Mutable":
		if len(ta.TypeParams) > 0 {
			return &MutableType{Element: TypeFromAnnotation(ta.TypeParams[0])}
		}
		return &MutableType{Element: &AnyType{}}
	default:
		return &StructType{Name: ta.Name, Fields: make(map[string]Type)}
	}
}
