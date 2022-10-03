package parser

import (
	"fmt"
	"github.com/hollykbuck/muskmelon/ast"
	"github.com/hollykbuck/muskmelon/lexer"
	"github.com/hollykbuck/muskmelon/token"
	"strconv"
)

const (
	_ int = iota
	// 运算符优先级

	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	// errors 不记录 error 只记录报错信息（字符串）
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(expression ast.Expression) ast.Expression
)

// New 初始化 Parser 结构
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken 将 lexer 解析出来的 token 读入 parser。
// 该 parser 可以向前看一个 token。
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram 在该函数中实现 parser 的主要逻辑
func (p *Parser) ParseProgram() *ast.Program {
	// 创建一个 Program 结构体
	program := &ast.Program{}
	// Program 的孩子指向一个 Statement 数组
	program.Statements = []ast.Statement{}
	// 程序的解析过程可以概括为重复解析 Statement
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// parseStatement 解析 Statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		// 以 Let Token 为开头的 statement 是 let statement
		// 委托给 parseLetStatement 执行解析任务
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		// fallback 到表达式语句解析
		return p.parseExpressionStatement()
	}
}

// parseLetStatement 解析 let statement。
// 返回 nil 表示解析失败。
func (p *Parser) parseLetStatement() *ast.LetStatement {
	// 创建一个空的 let statement
	statement := &ast.LetStatement{Token: p.curToken}
	// let 关键字后必须跟标识符
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	// 配置 let statement 的标识符部分
	statement.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// TODO: 现在我们还没实现表达式解析，所以我们让解析跳到 semicolon 为止。
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return statement
}

// expectPeek 检查下一个 token 是不是指定的类型
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// curTokenIs 断言当前 token 的类型
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs 断言下一个 token 的类型
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) Errors() []string {
	return p.errors
}

// peekError 添加一个 token 错误到 parser 的错误列表中
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// parseReturnStatement 解析 return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: p.curToken,
	}
	p.nextToken()
	// TODO: 现在我们跳过一切
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return statement
}

// registerPrefix 注册前缀表达式解析
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix 注册中缀表达式解析
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// parseExpressionStatement 解析表达式语句
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: p.parseExpression(LOWEST),
	}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseExpression 解析表达式。
// 根据表达式开头的 token 寻找对应的解析函数。
// 如果查不到 prefix 对应的解析函数会记录错误。
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// 查找和当前 token 对应的前缀表达式解析
	prefix, ok := p.prefixParseFns[p.curToken.Type]
	if !ok {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	return leftExp
}

// parseIdentifier 解析标识符
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral 解析整型字面量表达式
func (p *Parser) parseIntegerLiteral() ast.Expression {
	// 创建一个空的字面量表达式
	lit := &ast.IntegerLiteral{
		Token: p.curToken,
	}
	// 使用 strconv 解析数字
	// base=0 表示进制根据字符串而定
	parseInt, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	// 如果出错，将错误添加到 parser 的错误列表中
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	// 将 value 设定为解析出的数字
	lit.Value = parseInt
	return lit
}

// noPrefixParseFnError 记录找不到 prefix 对应的解析函数错误
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
