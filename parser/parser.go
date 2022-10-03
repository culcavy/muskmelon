package parser

import (
	"github.com/hollykbuck/muskmelon/ast"
	"github.com/hollykbuck/muskmelon/lexer"
	"github.com/hollykbuck/muskmelon/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

// New 初始化 Parser 结构
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
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
	default:
		return nil
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
