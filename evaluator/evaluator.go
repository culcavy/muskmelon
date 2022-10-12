package evaluator

import (
	"github.com/hollykbuck/muskmelon/ast"
	"github.com/hollykbuck/muskmelon/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval eval 传入的 ast 节点
func Eval(node ast.Node) object.Object {
	switch nodeType := node.(type) {
	case *ast.Program:
		// Statements
		return evalStatements(nodeType.Statements)
	case *ast.ExpressionStatement:
		// Statements
		return Eval(nodeType.Expression)
	case *ast.IntegerLiteral:
		// Expressions
		return &object.Integer{Value: nodeType.Value}
	case *ast.Boolean:
		// Expressions
		return nativeBoolToBooleanObject(nodeType.Value)
	case *ast.PrefixExpression:
		// Expressions
		right := Eval(nodeType.Right)
		return evalPrefixExpression(nodeType.Operator, right)
	}
	return nil
}

// evalPrefixExpression eval 前缀表达式. 根据 operator 类型将实现委托给具体的 eval 函数.
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// evalBangOperatorExpression eval 取反表达式
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

// evalStatements 批量解释 statement. 将实现委托给 Eval
func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement)
	}
	return result
}
