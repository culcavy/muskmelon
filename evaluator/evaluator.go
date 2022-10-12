package evaluator

import (
	"github.com/hollykbuck/muskmelon/ast"
	"github.com/hollykbuck/muskmelon/object"
)

// Eval eval 传入的 ast 节点
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
		// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return nil
}

// evalStatements 批量解释 statement. 将实现委托给 Eval
func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement)
	}
	return result
}
