package ast

import "MyCompiler/src/token"

type Node interface {
	TokenLiteral() string
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

// Identifier expression 标识符Node
type Identifier struct {
	Token token.Token
	Value string
}

// LetStatement statement let 语句Node
// 由三部分组成：1.let 2.等号左边的标识符 3.等号右边的表达式
type LetStatement struct {
	Token token.Token // token.LET 词法单元
	Name  *Identifier
	Value Expression
}

func (p *Program) TokenLiteral() string {
	if len(p.Statement) > 0 {
		return p.Statement[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) statementNode() {}

func (l *LetStatement) TokenLiteral() string { return l.Token.Literal }

func (l *LetStatement) statementNode() {}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) expressionNode() {}
