package object

import (
	"MyCompiler/src/ast"
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

// 值类型对象
// 包含三种数据类型：空值 布尔值 整数

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// region Integer

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// endregion

// region Boolean

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// endregion

// region Null

type Null struct{}

func (n Null) Type() ObjectType { return NULL_OBJ }

func (n Null) Inspect() string { return "null" }

// endregion

// region Return

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

func (r *ReturnValue) Inspect() string { return r.Value.Inspect() }

// endregion

// region Function

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

func (f *Function) Inspect() string {
	var out bytes.Buffer

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}

// endregion

// region Error

type Error struct {
	Message string
}

func (e Error) Type() ObjectType { return ERROR_OBJ }

func (e Error) Inspect() string { return "ERROR: " + e.Message }

// endregion

// region String

type String struct {
	Value string
}

func (s String) Type() ObjectType { return STRING_OBJ }

func (s String) Inspect() string { return s.Value }

// endregion

// region Builtin function

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b Builtin) Type() ObjectType { return BUILTIN_OBJ }

func (b Builtin) Inspect() string { return "builtin function" }

// endregion

// region Array

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }

func (a *Array) Inspect() string {
	var out bytes.Buffer

	var elements []string
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteByte('[')
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// endregion

// region Hash

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{
		Type:  b.Type(),
		Value: value,
	}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(),
	}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	var pairs []string
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// endregion
