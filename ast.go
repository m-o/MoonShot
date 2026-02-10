package main

import (
	"bytes"
	"strings"
)

// Node is the base interface for all AST nodes
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents a statement node
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// DefStatement represents a variable definition: def x = 5
type DefStatement struct {
	Token    Token      // the DEF token
	Name     *Identifier
	TypeHint *TypeAnnotation // optional type hint
	Value    Expression
}

func (ds *DefStatement) statementNode()       {}
func (ds *DefStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DefStatement) String() string {
	var out bytes.Buffer
	out.WriteString("def ")
	out.WriteString(ds.Name.String())
	if ds.TypeHint != nil {
		out.WriteString(": ")
		out.WriteString(ds.TypeHint.String())
	}
	out.WriteString(" = ")
	if ds.Value != nil {
		out.WriteString(ds.Value.String())
	}
	return out.String()
}

// TypeAnnotation represents a type hint
type TypeAnnotation struct {
	Token      Token // the type name token
	Name       string
	TypeParams []*TypeAnnotation // for generics like List[Integer]
}

func (ta *TypeAnnotation) String() string {
	if len(ta.TypeParams) == 0 {
		return ta.Name
	}
	var params []string
	for _, p := range ta.TypeParams {
		params = append(params, p.String())
	}
	return ta.Name + "[" + strings.Join(params, ", ") + "]"
}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Token Token // the RETURN token
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString("return ")
	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}
	return out.String()
}

// ExpressionStatement wraps an expression as a statement
type ExpressionStatement struct {
	Token      Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement represents a block of statements
type BlockStatement struct {
	Token      Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")
	return out.String()
}

// Identifier represents a variable name
type Identifier struct {
	Token Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents an integer value
type IntegerLiteral struct {
	Token Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral represents a floating-point value
type FloatLiteral struct {
	Token Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral represents a string value
type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// BooleanLiteral represents true or false
type BooleanLiteral struct {
	Token Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// PrefixExpression represents a prefix operation like -5 or not true
type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

// InfixExpression represents a binary operation like 5 + 3
type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

// AssignmentExpression represents mutable assignment: counter == counter + 1
type AssignmentExpression struct {
	Token Token
	Name  *Identifier
	Value Expression
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	return ae.Name.String() + " == " + ae.Value.String()
}

// IfExpression represents an if-else expression
type IfExpression struct {
	Token       Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// WhileStatement represents a while loop
type WhileStatement struct {
	Token     Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())
	return out.String()
}

// ForStatement represents a for-in loop
type ForStatement struct {
	Token    Token
	Variable *Identifier
	Iterable Expression
	Body     *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	out.WriteString(fs.Variable.String())
	out.WriteString(" in ")
	out.WriteString(fs.Iterable.String())
	out.WriteString(" ")
	out.WriteString(fs.Body.String())
	return out.String()
}

// BreakStatement represents a break statement
type BreakStatement struct {
	Token Token
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string       { return "break" }

// ContinueStatement represents a continue statement
type ContinueStatement struct {
	Token Token
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string       { return "continue" }

// FunctionStatement represents a function definition
type FunctionStatement struct {
	Token      Token
	Name       *Identifier
	Parameters []*FunctionParameter
	ReturnType *TypeAnnotation
	Body       *BlockStatement
}

type FunctionParameter struct {
	Name     *Identifier
	TypeHint *TypeAnnotation
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer
	out.WriteString("fun ")
	out.WriteString(fs.Name.String())
	out.WriteString("(")
	var params []string
	for _, p := range fs.Parameters {
		param := p.Name.String()
		if p.TypeHint != nil {
			param += ": " + p.TypeHint.String()
		}
		params = append(params, param)
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	if fs.ReturnType != nil {
		out.WriteString(" -> ")
		out.WriteString(fs.ReturnType.String())
	}
	out.WriteString(" ")
	out.WriteString(fs.Body.String())
	return out.String()
}

// FunctionLiteral represents an anonymous function (lambda)
type FunctionLiteral struct {
	Token      Token
	Parameters []*Identifier
	Body       Expression // single expression for lambdas
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(" -> ")
	out.WriteString(fl.Body.String())
	out.WriteString(" }")
	return out.String()
}

// CallExpression represents a function call
type CallExpression struct {
	Token     Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// MemberExpression represents member access: obj.field
type MemberExpression struct {
	Token  Token
	Object Expression
	Member *Identifier
}

func (me *MemberExpression) expressionNode()      {}
func (me *MemberExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MemberExpression) String() string {
	return me.Object.String() + "." + me.Member.String()
}

// IndexExpression represents index access: list[0]
type IndexExpression struct {
	Token Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	return "(" + ie.Left.String() + "[" + ie.Index.String() + "])"
}

// ListLiteral represents a list: [1, 2, 3]
type ListLiteral struct {
	Token    Token
	Elements []Expression
}

func (ll *ListLiteral) expressionNode()      {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }
func (ll *ListLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	var elements []string
	for _, e := range ll.Elements {
		elements = append(elements, e.String())
	}
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// MapLiteral represents a map: {"key": value}
type MapLiteral struct {
	Token Token
	Pairs map[Expression]Expression
}

func (ml *MapLiteral) expressionNode()      {}
func (ml *MapLiteral) TokenLiteral() string { return ml.Token.Literal }
func (ml *MapLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	var pairs []string
	for k, v := range ml.Pairs {
		pairs = append(pairs, k.String()+": "+v.String())
	}
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// StructStatement represents a struct definition
type StructStatement struct {
	Token  Token
	Name   *Identifier
	Fields []*StructField
}

type StructField struct {
	Name     *Identifier
	TypeHint *TypeAnnotation
}

func (ss *StructStatement) statementNode()       {}
func (ss *StructStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *StructStatement) String() string {
	var out bytes.Buffer
	out.WriteString("struct ")
	out.WriteString(ss.Name.String())
	out.WriteString(" { ")
	var fields []string
	for _, f := range ss.Fields {
		field := f.Name.String()
		if f.TypeHint != nil {
			field += ": " + f.TypeHint.String()
		}
		fields = append(fields, field)
	}
	out.WriteString(strings.Join(fields, ", "))
	out.WriteString(" }")
	return out.String()
}

// StructLiteral represents a struct instantiation: User { name: "Alice" }
type StructLiteral struct {
	Token      Token
	StructName *Identifier
	Fields     map[string]Expression
}

func (sl *StructLiteral) expressionNode()      {}
func (sl *StructLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StructLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(sl.StructName.String())
	out.WriteString(" { ")
	var fields []string
	for k, v := range sl.Fields {
		fields = append(fields, k+": "+v.String())
	}
	out.WriteString(strings.Join(fields, ", "))
	out.WriteString(" }")
	return out.String()
}

// WithExpression represents struct update: user.with { age: 31 }
type WithExpression struct {
	Token   Token
	Object  Expression
	Updates map[string]Expression
}

func (we *WithExpression) expressionNode()      {}
func (we *WithExpression) TokenLiteral() string { return we.Token.Literal }
func (we *WithExpression) String() string {
	var out bytes.Buffer
	out.WriteString(we.Object.String())
	out.WriteString(".with { ")
	var updates []string
	for k, v := range we.Updates {
		updates = append(updates, k+": "+v.String())
	}
	out.WriteString(strings.Join(updates, ", "))
	out.WriteString(" }")
	return out.String()
}

// OptionExpression represents Some(x) or None
type OptionExpression struct {
	Token   Token
	IsSome  bool
	Value   Expression // nil if None
}

func (oe *OptionExpression) expressionNode()      {}
func (oe *OptionExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *OptionExpression) String() string {
	if oe.IsSome {
		return "Some(" + oe.Value.String() + ")"
	}
	return "None"
}

// ResultExpression represents Ok(x) or Error(x)
type ResultExpression struct {
	Token Token
	IsOk  bool
	Value Expression
}

func (re *ResultExpression) expressionNode()      {}
func (re *ResultExpression) TokenLiteral() string { return re.Token.Literal }
func (re *ResultExpression) String() string {
	if re.IsOk {
		return "Ok(" + re.Value.String() + ")"
	}
	return "Error(" + re.Value.String() + ")"
}

// MatchExpression represents pattern matching
type MatchExpression struct {
	Token Token
	Value Expression
	Cases []*MatchCase
}

type MatchCase struct {
	Pattern    Expression
	BindingVar *Identifier // the variable in Some(x) or Ok(x)
	Body       *BlockStatement
}

func (me *MatchExpression) expressionNode()      {}
func (me *MatchExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MatchExpression) String() string {
	var out bytes.Buffer
	out.WriteString("match ")
	out.WriteString(me.Value.String())
	out.WriteString(" { ")
	for _, c := range me.Cases {
		out.WriteString(c.Pattern.String())
		out.WriteString(" -> ")
		out.WriteString(c.Body.String())
		out.WriteString(" ")
	}
	out.WriteString("}")
	return out.String()
}

// MutableExpression represents Mutable[T](value)
type MutableExpression struct {
	Token    Token
	TypeHint *TypeAnnotation
	Value    Expression
}

func (me *MutableExpression) expressionNode()      {}
func (me *MutableExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MutableExpression) String() string {
	var out bytes.Buffer
	out.WriteString("Mutable")
	if me.TypeHint != nil {
		out.WriteString("[")
		out.WriteString(me.TypeHint.String())
		out.WriteString("]")
	}
	out.WriteString("(")
	out.WriteString(me.Value.String())
	out.WriteString(")")
	return out.String()
}

// ExtendStatement represents extension methods
type ExtendStatement struct {
	Token    Token
	TypeName *Identifier
	Methods  []*FunctionStatement
}

func (es *ExtendStatement) statementNode()       {}
func (es *ExtendStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExtendStatement) String() string {
	var out bytes.Buffer
	out.WriteString("extend ")
	out.WriteString(es.TypeName.String())
	out.WriteString(" { ")
	for _, m := range es.Methods {
		out.WriteString(m.String())
		out.WriteString(" ")
	}
	out.WriteString("}")
	return out.String()
}

// ImportStatement represents an import
type ImportStatement struct {
	Token Token
	Path  []string // e.g., ["user", "User"]
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImportStatement) String() string {
	return "import " + strings.Join(is.Path, ".")
}
