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
	switch nodeActual := node.(type) {
	case *ast.Program:
		// Statements
		return evalStatements(nodeActual.Statements)
	case *ast.ExpressionStatement:
		// Statements
		return Eval(nodeActual.Expression)
	case *ast.IntegerLiteral:
		// Expressions
		return &object.Integer{Value: nodeActual.Value}
	case *ast.Boolean:
		// Expressions
		return nativeBoolToBooleanObject(nodeActual.Value)
	case *ast.PrefixExpression:
		// Expressions
		right := Eval(nodeActual.Right)
		return evalPrefixExpression(nodeActual.Operator, right)
	case *ast.InfixExpression:
		// Expressions
		left := Eval(nodeActual.Left)
		right := Eval(nodeActual.Right)
		return evalInfixExpression(nodeActual.Operator, left, right)
	}
	return nil
}

// evalInfixExpression eval 中缀表达式
func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	default:
		return NULL
	}
}

// evalIntegerInfixExpression 整型中缀表达式计算
func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	default:
		return NULL
	}
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
