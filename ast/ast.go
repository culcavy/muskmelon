package ast

import (
	"bytes"
	"fmt"
	"github.com/hollykbuck/muskmelon/token"
	"strings"
)

// Node AST 节点
type Node interface {
	// TokenLiteral 节点 Token 的字面
	TokenLiteral() string
	// String 输出节点的字符串表示。
	// 打印 AST 节点用于调试
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BlockStatement) String() string {
	var out bytes.Buffer
	for _, statement := range b.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}

func (b *BlockStatement) statementNode() {
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

// String 序列化 If 表达式
func (i *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		out.WriteString("else ")

		out.WriteString(i.Alternative.String())
	}
	return out.String()
}

func (i *IfExpression) expressionNode() {

}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionLiteral) String() string {
	var out bytes.Buffer
	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(f.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	return out.String()
}

func (f *FunctionLiteral) expressionNode() {
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

func (b *Boolean) expressionNode() {

}

// Program AST的根节点
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, statement := range p.Statements {
		// 将实现委托给 statement
		out.WriteString(statement.String())
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("(%s%s)", p.Operator, p.Right.String()))
	return out.String()
}

func (p *PrefixExpression) expressionNode() {
}

func (l *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")
	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
	Token    token.Token
}

func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Operator, i.Right.String())
}

func (i *InfixExpression) expressionNode() {

}

func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) expressionNode() {

}

func (r *ExpressionStatement) String() string {
	if r.Expression != nil {
		return r.Expression.String()
	}
	return ""
}

func (l *LetStatement) statementNode() {

}

func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) expressionNode() {

}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (r *ReturnStatement) statementNode() {

}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ExpressionStatement) statementNode() {

}

func (r *ExpressionStatement) TokenLiteral() string {
	return r.Token.Literal
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}

func (s *StringLiteral) String() string {
	return s.Token.Literal
}

func (s *StringLiteral) expressionNode() {

}
