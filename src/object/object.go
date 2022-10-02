package object

import (
	"MyCompiler/src/ast"
	"bytes"
	"fmt"
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
