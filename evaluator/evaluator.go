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
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch nodeActual := node.(type) {
	case *ast.Program:
		// Statements
		return evalProgram(nodeActual, env)
	case *ast.ExpressionStatement:
		// Statements
		return Eval(nodeActual.Expression, env)
	case *ast.IntegerLiteral:
		// Expressions
		return &object.Integer{Value: nodeActual.Value}
	case *ast.Boolean:
		// Expressions
		return nativeBoolToBooleanObject(nodeActual.Value)
	case *ast.PrefixExpression:
		// Expressions
		right := Eval(nodeActual.Right, env)
		// operand 出现错误应当返回错误
		if isError(right) {
			return right
		}
		return evalPrefixExpression(nodeActual.Operator, right)
	case *ast.InfixExpression:
		// Expressions
		// 任何一个 operand 出现错误都应返回错误
		left := Eval(nodeActual.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(nodeActual.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(nodeActual.Operator, left, right)
	case *ast.BlockStatement:
		// Expression
		// 将实现委托给 evalStatements
		return evalBlockStatement(nodeActual.Statements, env)
	case *ast.IfExpression:
		// Expression
		return evalIfExpression(nodeActual, env)
	case *ast.ReturnStatement:
		// 计算表达式的值
		// 表达式的值作为返回值返回
		val := Eval(nodeActual.ReturnValue, env)
		// return 的 operand 出现错误应返回错误
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(nodeActual.Value, env)
		if isError(val) {
			return val
		}
		// 将等号右边的值存到 environment 中
		env.Set(nodeActual.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(nodeActual, env)
	case *ast.FunctionLiteral:
		params := nodeActual.Parameters
		body := nodeActual.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		function := Eval(nodeActual.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(nodeActual.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}
	return nil
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}
	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

// evalIdentifier 计算标识符的值
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}
	return val
}

// evalBlockStatement eval block statement
func evalBlockStatement(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range statements {
		// 执行单条语句
		result = Eval(statement, env)

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
func evalProgram(actual *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range actual.Statements {
		result = Eval(statement, env)
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
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
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
