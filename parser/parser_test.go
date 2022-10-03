package parser

import (
	"github.com/hollykbuck/muskmelon/ast"
	"github.com/hollykbuck/muskmelon/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
	tests := []struct {
		expectedLiteral string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedLiteral) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, literal string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		return false
	}
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", statement)
		return false
	}
	if letStatement.Name.Value != literal {
		t.Errorf("letStmt.Name.Value not '%s', got='%s'", literal, letStatement.Name.Value)
		return false
	}
	if letStatement.Name.TokenLiteral() != literal {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s', got='%s'", literal, letStatement.Name.TokenLiteral())
		return false
	}
	return true
}
