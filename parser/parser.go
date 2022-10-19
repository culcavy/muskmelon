package parser

import (
	"fmt"
	"github.com/hollykbuck/muskmelon/ast"
	"github.com/hollykbuck/muskmelon/lexer"
	"github.com/hollykbuck/muskmelon/token"
	"strconv"
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

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
	// Identifier 和 Integer 是终止符。
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	// BANG 和 MINUS 是非终止符
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	// 中缀表达式没有一个是终结符
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.nextToken()
	p.nextToken()
	return p
}

// 解析括号表达式
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	// 从最低优先级开始重新解析
	expression := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expression
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
	p.nextToken()
	statement.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return statement
}

// expectPeek 检查下一个 token 是不是指定的类型。
// 如果是，让 curToken 指向下一个 token。
// 如果不是，在 Parser 上记录一个错误。
//
// mutable。
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
//
// const
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
	msg := fmt.Sprintf("expectedBool next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// parseReturnStatement 解析 return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: p.curToken,
	}
	p.nextToken()
	statement.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
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
	// statement 的末尾可以是 semicolon，也可以不是
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
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
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

func (p *Parser) parseBoolean() ast.Expression {
	lit := &ast.Boolean{
		Token: p.curToken,
	}
	if p.curToken.Literal == "true" {
		lit.Value = true
	} else if p.curToken.Literal == "false" {
		lit.Value = false
	}
	return lit
}

// noPrefixParseFnError 记录找不到 prefix 对应的解析函数错误
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// parsePrefixExpression 递归解析前缀表达式
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

// peekPrecedence 查询 peekToken 的优先级
func (p *Parser) peekPrecedence() int {
	precedence, ok := precedences[p.peekToken.Type]
	if ok {
		return precedence
	}
	//if !ok {
	//	p.errors = append(p.errors, fmt.Sprintf("找不到 token 对应的优先级"))
	//}
	return LOWEST
}

// curPrecedence 查询 curToken 的优先级
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// parseInfixExpression 实现中缀表达式的解析
func (p *Parser) parseInfixExpression(leftOperand ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Left:     leftOperand,
		Operator: p.curToken.Literal,
		Right:    nil,
		Token:    p.curToken,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

// parseIfExpression 解析 If 表达式
func (p *Parser) parseIfExpression() ast.Expression {
	// 先创建一个空白的 If 表达式
	expression := &ast.IfExpression{Token: p.curToken}
	// 下一个 token 应该是 `(`
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	//log.Println(p.curToken.Type)
	p.nextToken()
	//log.Println(p.curToken.Type)
	// 从最低优先级开始解析 Condition 表达式
	expression.Condition = p.parseExpression(LOWEST)
	// 解析完成后应该碰到的 token 是 `)`
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	// 再接下去是 `{`
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

// parseBlockStatement 解析块级表达式
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() (identifiers []*ast.Identifier) {
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return
	}
	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return
}

func (p *Parser) parseCallExpression(expression ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: expression}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	var list []ast.Expression
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}
