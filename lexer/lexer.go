package lexer

import (
	"github.com/hollykbuck/muskmelon/token"
	"unicode/utf8"
)

// Lexer Lexer 记录
type Lexer struct {
	input string
	// 当前 Position
	position int
	// 下一个 Position
	readPosition int
	// 当前的字符
	ch rune
}

// New Lexer 的构造函数
func New(input string) *Lexer {
	l := &Lexer{input: input}
	// 初始化 lexer
	l.readChar()
	return l
}

// readChar 读下一个字符，将读到的字符串写到 context 中
// 如果 EOF 那么字符设定为 EOF
func (l *Lexer) readChar() {
	step := 1
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
		step = size
		l.ch = r
		//l.ch =
	}
	// 如果成功读取，移动 position
	l.position = l.readPosition
	// 如果成功读取，移动 readPosition
	l.readPosition += step
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
