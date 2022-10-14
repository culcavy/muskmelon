package object

import "fmt"

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
	ERROR_OBJ = "ERROR"
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
