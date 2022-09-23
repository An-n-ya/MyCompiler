package ast

import (
	"MyCompiler/src/token"
	"bytes"
)

type Node interface {
	TokenLiteral() string
	String() string
}

// region statement

// Program statement 程序Node（顶级Node）
type Program struct {
	Statement []Statement
}

type Statement interface {
	Node
	statementNode()
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

// BlockStatement 语句块
type BlockStatement struct {
	Token      token.Token // 词法单元是 {
	Statements []Statement
}

// endregion

// region expression

type Expression interface {
	Node
	expressionNode()
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

// BooleanLiteral expression 布尔字面量表达式
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

// PrefixExpression expression 前缀表达式
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

// InfixExpression expression 中缀表达式
type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

// IfExpression expression if表达式
type IfExpression struct {
	Token       token.Token // 词法单元为 if
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

// endregion

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

func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }

func (b *BooleanLiteral) String() string { return b.Token.Literal }

func (b *BooleanLiteral) expressionNode() {}

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

func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }

func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

func (i *InfixExpression) expressionNode() {}

func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }

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

func (i *IfExpression) expressionNode() {}

func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }

func (b *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range b.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (b *BlockStatement) statementNode() {}

