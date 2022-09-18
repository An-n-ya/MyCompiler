package ast

import (
	"MyCompiler/src/token"
	"bytes"
)

type Node interface {
	TokenLiteral() string
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

// Program statement 程序Node（顶级Node）
type Program struct {
	Statement []Statement
}

// Identifier expression 标识符表达式
type Identifier struct {
	Token token.Token
	Value string
}

// IntegerLiteral expression 整数字面量表达式
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// PrefixExpression expression 前缀表达式
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

// LetStatement statement let 语句Node
// 由三部分组成：1.let 2.等号左边的标识符 3.等号右边的表达式
type LetStatement struct {
	Token token.Token // token.LET 词法单元
	Name  *Identifier
	Value Expression
}

// ReturnStatement statement return语句node
// 由两部分组成：return + 表达式
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

// ExpressionStatement statement 表达式语句
// 仅仅有一个表达式构成的语句，至此语言中的三种语句都定义完成
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (p *Program) TokenLiteral() string {
	if len(p.Statement) > 0 {
		return p.Statement[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) statementNode() {}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statement {
		out.WriteString(s.String())
	}
	return out.String()
}

func (l *LetStatement) TokenLiteral() string { return l.Token.Literal }

func (l *LetStatement) statementNode() {}

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

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) expressionNode() {}

func (i *Identifier) String() string { return i.Value }

func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }

func (i *IntegerLiteral) String() string { return i.Token.Literal }

func (i *IntegerLiteral) expressionNode() {}

func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }

func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

func (p *PrefixExpression) expressionNode() {}
