package parser

import (
	"MyCompiler/src/ast"
	"MyCompiler/src/lexer"
	"MyCompiler/src/token"
	"fmt"
	"strconv"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string // 错误信息

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression               // 前缀解析函数
	infixParseFn  func(ast.Expression) ast.Expression // 中缀解析函数
)

// 优先级
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunc(X)
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// 读取两个词法单元, 此时curToken指向第一个词法单元
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	// 注册前缀解析函数
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		// 如果是其他情况，按照表达式语句处理
		return p.parseExpressionStatement()
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: 跳过对表达式的处理，直接到分号
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	// 表达式分号可选
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt

}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// 把前缀解析函数找出来
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPreFixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	return leftExp
}

func (p *Parser) noPreFixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// 解析标识符
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// 解析整数字面量
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	// 将字符串转化为int64保存
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	// 解析 ！ 或 - 后面的内容
	// 递归解析右边的表达式
	expression.Right = p.parseExpression(PREFIX)

	return expression
}
