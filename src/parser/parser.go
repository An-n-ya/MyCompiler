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

// region 优先级相关

// 为每个符号规定优先级
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

// 查找下个符号的优先级
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	// 如果没找到， 就返回最低优先级
	return LOWEST
}

// 查找当前符号的优先级
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	// 如果没找到， 就返回最低优先级
	return LOWEST
}

// endregion

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// 读取两个词法单元, 此时curToken指向第一个词法单元
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	// 注册前缀解析函数(用来解析字面量表达式和！ -这样的前缀表达式)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFnExpression)

	// 注册中缀解析函数
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	// 用中缀解析左括号 用作解析调用函数
	p.registerInfix(token.LPAREN, p.parseCallExpression)

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
		stmt := p.parseStatement()
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

func (p *Parser) parseStatement() ast.Statement {
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

	// 跳过赋值号
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

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
	stmt.ReturnValue = p.parseExpression(LOWEST)

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

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// 往下找中缀表达式
		infix := p.infixParseFns[p.peekToken.Type]

		// 如果不是中缀表达式， 就直接返回之前的结果
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

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

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	// 这里需要用当前符号的优先级来解析右表达式
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	// 解析括号里的表达式
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		// 表达式应该被右括号括起来
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	// if 后应该是左括号
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// 跳过左括号
	p.nextToken()
	// 解析条件表达式
	expression.Condition = p.parseExpression(LOWEST)

	// 条件表达式后面应该是右括号
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	// 下面解析consequence
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

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	// 跳过左花括号
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFnExpression() ast.Expression {
	expression := &ast.FnExpression{Token: p.curToken}

	// fn后面应该是括号
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	// 解析参数
	expression.Parameters = p.parseFnParameters()

	// 参数后应该是花括号
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	// 解析函数体
	expression.Body = p.parseBlockStatement()
	return expression

}

func (p *Parser) parseFnParameters() []*ast.Identifier {
	var params []*ast.Identifier

	// 跳过左括号
	p.nextToken()

	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		ident := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		params = append(params, ident)
		p.nextToken()

		if p.curTokenIs(token.RPAREN) || p.curTokenIs(token.EOF) {
			break
		}

		// 当前应该是逗号
		if !p.curTokenIs(token.COMMA) {
			return nil
		}
		p.nextToken()
	}
	// 当前应该是右括号
	if !p.curTokenIs(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	// 中缀表达式传过来的左值就是函数
	expression := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	expression.Arguments = p.parseCallArguments()
	return expression
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
