package ast

import (
	"MyCompiler/src/token"
	"bytes"
	"strings"
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

// StringLiteral expression 字符串字面量
type StringLiteral struct {
	Token token.Token
	Value string
}

// BooleanLiteral expression 布尔字面量表达式
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

// ArrayLiteral expression 数组字面量
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

// IndexExpression expression 数组索引表达式
type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
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

// FnExpression expression 函数表达式
type FnExpression struct {
	Token      token.Token // 词法单元是 fn
	Parameters []*Identifier
	Body       *BlockStatement
}

// CallExpression expression 调用函数表达式
type CallExpression struct {
	Token     token.Token // 词法单元是 '('
	Function  Expression  // 标识符或者是函数表达式
	Arguments []Expression
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

func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }

func (s *StringLiteral) String() string { return s.Token.Literal }

func (s *StringLiteral) expressionNode() {}

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

func (f *FnExpression) TokenLiteral() string { return f.Token.Literal }

func (f *FnExpression) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range f.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(f.Body.String())

	return out.String()
}

func (f *FnExpression) expressionNode() {}

func (c *CallExpression) TokenLiteral() string { return c.Token.Literal }

func (c *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")

	return out.String()
}

func (c *CallExpression) expressionNode() {}

func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }

func (a *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, ele := range a.Elements {
		elements = append(elements, ele.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (a *ArrayLiteral) expressionNode() {}

func (i *IndexExpression) TokenLiteral() string { return i.Token.Literal }

func (i *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("])")

	return out.String()
}

func (i *IndexExpression) expressionNode() {}
