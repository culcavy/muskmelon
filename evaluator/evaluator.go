package evaluator

import (
	"github.com/hollykbuck/muskmelon/ast"
	"github.com/hollykbuck/muskmelon/object"
)

// Eval eval 传入的 ast 节点
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return nil
}
