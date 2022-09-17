package parser

import (
	"MyCompiler/src/ast"
	"MyCompiler/src/lexer"
	"MyCompiler/src/token"
	"fmt"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string // 错误信息
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// 读取两个词法单元, 此时curToken指向第一个词法单元
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram 解析返回ast
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statement: []ast.Statement{}}

	for p.curToken.Type != token.EOF {
		stmt := p.parseProgram()
		if stmt != nil {
			program.Statement = append(program.Statement, stmt)
		}
		p.nextToken()
	}

	return program
}

// Error 返回错误信息
func (p *Parser) Error() []string {
	return p.errors
}

func (p *Parser) parseProgram() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// 解析标识符
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// 解析赋值号
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: 跳过表达式的处理，直接找分号
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 语法分析的断言函数
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		// 如果断言失败，就往errors里加错误信息
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekTokenIs(t token.TokenType) bool { return p.peekToken.Type == t }
func (p *Parser) curTokenIs(t token.TokenType) bool  { return p.curToken.Type == t }

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
