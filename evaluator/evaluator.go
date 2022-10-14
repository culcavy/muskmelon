package evaluator

import (
	"fmt"
	"github.com/hollykbuck/muskmelon/ast"
	"github.com/hollykbuck/muskmelon/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// isError 判断 obj 类型是否是错误
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// Eval eval 传入的 ast 节点
func Eval(node ast.Node) object.Object {
	switch nodeActual := node.(type) {
	case *ast.Program:
		// Statements
		return evalProgram(nodeActual)
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
		// operand 出现错误应当返回错误
		if isError(right) {
			return right
		}
		return evalPrefixExpression(nodeActual.Operator, right)
	case *ast.InfixExpression:
		// Expressions
		// 任何一个 operand 出现错误都应返回错误
		left := Eval(nodeActual.Left)
		if isError(left) {
			return left
		}
		right := Eval(nodeActual.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(nodeActual.Operator, left, right)
	case *ast.BlockStatement:
		// Expression
		// 将实现委托给 evalStatements
		return evalBlockStatement(nodeActual.Statements)
	case *ast.IfExpression:
		// Expression
		return evalIfExpression(nodeActual)
	case *ast.ReturnStatement:
		// 计算表达式的值
		// 表达式的值作为返回值返回
		val := Eval(nodeActual.ReturnValue)
		// return 的 operand 出现错误应返回错误
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(nodeActual.Value)
		if isError(val) {
			return val
		}
	}
	return nil
}

// evalBlockStatement eval block statement
func evalBlockStatement(statements []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range statements {
		// 执行单条语句
		result = Eval(statement)

		if result != nil {
			resultType := result.Type()
			// 碰到 return 或者 error 了就打断流程
			// 这里我们不 unwrap return value
			if resultType == object.RETURN_VALUE_OBJ || resultType == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

// evalProgram eval Program 节点
func evalProgram(actual *ast.Program) object.Object {
	var result object.Object
	for _, statement := range actual.Statements {
		result = Eval(statement)
		switch resultActual := result.(type) {
		case *object.ReturnValue:
			// 碰到 return 了就打断流程
			// 直到 evalProgram 才 unwrap return value
			return resultActual.Value
		case *object.Error:
			// 如果运行出现了错误，选择不展开
			return resultActual
		}
	}
	return result
}

// evalIfExpression eval if 表达式
func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

// isTruthy 判断 Condition 的值是否为 if 意义上的 true
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// evalInfixExpression eval 中缀表达式
func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
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
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
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

// nativeBoolToBooleanObject 将 input 转换为布尔类型对象
func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

// newError Error 对象的构造函数
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
