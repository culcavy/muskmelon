package evaluator

import (
	"github.com/hollykbuck/muskmelon/lexer"
	"github.com/hollykbuck/muskmelon/object"
	"github.com/hollykbuck/muskmelon/parser"
	"testing"
)

// TestEvalIntegerExpression 测试 eval 整型
func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// testEval 运行 input 代码
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

// testIntegerObject 检查 eval 的结果是否为整型对象, 检查值是否为 expected
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}
