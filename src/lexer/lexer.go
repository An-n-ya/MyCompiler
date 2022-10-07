package lexer

import (
	"MyCompiler/src/token"
	"bytes"
)

// TODO: 支持Unicode， 需要把ch改成rune，并且要更换readChar的逻辑
type Lexer struct {
	input        string // 输入字符串
	position     int    // 指向当前位置
	readPosition int    // 指向当前位置之后的一个字符
	ch           byte   // 当前字符 (只支持ASCII)
}

// 构造函数，使用input创建Lexer实例
func New(input string) *Lexer {
	l := &Lexer{input: input}
	// 初始化
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// 空白类字符不检测
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			// 如果下一个字符还是=，解析为==
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = tokenFactory(token.ASSIGN, l.ch)
		}
	case '+':
		tok = tokenFactory(token.PLUS, l.ch)
	case '-':
		tok = tokenFactory(token.MINUS, l.ch)
	case '*':
		tok = tokenFactory(token.ASTERISK, l.ch)
	case '/':
		tok = tokenFactory(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = tokenFactory(token.BANG, l.ch)
		}
	case '>':
		tok = tokenFactory(token.GT, l.ch)
	case '<':
		tok = tokenFactory(token.LT, l.ch)
	case ',':
		tok = tokenFactory(token.COMMA, l.ch)
	case ';':
		tok = tokenFactory(token.SEMICOLON, l.ch)
	case '(':
		tok = tokenFactory(token.LPAREN, l.ch)
	case ')':
		tok = tokenFactory(token.RPAREN, l.ch)
	case '{':
		tok = tokenFactory(token.LBRACE, l.ch)
	case '}':
		tok = tokenFactory(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = tokenFactory(token.LBRACKET, l.ch)
	case ']':
		tok = tokenFactory(token.RBRACKET, l.ch)
	case ':':
		tok = tokenFactory(token.COLON, l.ch)
	case 0:
		// 读到的字符为空 - 返回空字符串（这里需要特殊处理）
		tok = token.Token{Type: token.EOF}
	default:
		if isLetter(l.ch) {
			// 如果是字符
			tok.Literal = l.readIdentifier()
			// 根据字符值匹配关键词，从而决定Type
			tok.Type = token.LookupIdent(tok.Literal)
			// 这个时候不继续读下一个字符，直接返回
			return tok
		} else if isDigit(l.ch) {
			// 如果是数字
			tok.Type = token.INT
			tok.Literal = l.readNum()
			return tok
		} else {
			tok = tokenFactory(token.ILLEGAL, l.ch)
		}
	}

	// 指针移动到下一个字符
	l.readChar()
	return tok
}

// 跳过空白字符
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// 读取标识符
func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	endPos := l.position
	return l.input[startPos:endPos]
}

// 读取数字
func (l *Lexer) readNum() string {
	startPos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	endPos := l.position
	return l.input[startPos:endPos]
}

// 判断是否是数字
func isDigit(c byte) bool {
	// 只接受 a-zA-Z_ 作为标识符
	return '0' <= c && c <= '9'
}

// 判断是否是有效字母
func isLetter(c byte) bool {
	// 只接受 a-zA-Z_ 作为标识符
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

// 读取字符 移动指针
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// 如果指针到底了, 置ch为0
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	// 指向下一个位置
	l.position = l.readPosition
	l.readPosition += 1
}

// 查看下一个字符，不移动指针
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readString() string {
	// position := l.position + 1
	var out bytes.Buffer
	for {
		l.readChar()
		if l.ch == '\\' {
			// 如果当前是\,就往后读两次
			l.readChar()
			handleEscape(&out, l.ch)
			continue
		}
		if l.ch == '"' || l.ch == 0 {
			break
		}
		out.WriteByte(l.ch)
	}
	return out.String()
}

// 创建Token的工厂方法
func tokenFactory(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// 处理转义字符
func handleEscape(out *bytes.Buffer, ch byte) {
	if ch == 'n' {
		out.WriteByte('\n')
	} else if ch == '"' {
		out.WriteByte('"')
	} else if ch == '\\' {
		out.WriteByte('\\')
	} else if ch == 'b' {
		out.WriteByte('\b')
	} else if ch == '\r' {
		out.WriteByte('\r')
	} else if ch == '\t' {
		out.WriteByte('\t')
	}

}
