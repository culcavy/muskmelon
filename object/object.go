package object

import (
	"bytes"
	"fmt"
	"github.com/hollykbuck/muskmelon/ast"
	"strings"
)

type ObjectType string

const (
	// INTEGER_OBJ 整型类型
	INTEGER_OBJ = "INTEGER"
	//BOOLEAN_OBJ Bool 类型
	BOOLEAN_OBJ = "BOOLEAN"
	//NULL_OBJ Null 类型
	NULL_OBJ = "NULL"
	//RETURN_VALUE_OBJ 返回值类型
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	//ERROR_OBJ 错误类型
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	STRING_OBJ   = "STRING"
)

// Object 所有的对象的父类型
type Object interface {
	// Type 对象的类型
	Type() ObjectType
	// Inspect 查看对象的信息
	Inspect() string
}

// Integer 整型类型对象
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Boolean Bool 类型对象
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue 返回值类型
type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

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
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// String 字符串字面量
type String struct {
	Value string
}

// Type 字符串字面量的类型
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
