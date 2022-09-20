package parser

import (
	"MyCompiler/src/ast"
	"MyCompiler/src/lexer"
	"MyCompiler/src/parser"
	"MyCompiler/src/token"
	"fmt"
	"testing"
)

// 构建"let myVar = anotherVar;"语句的ast，测试字符串输出
func TestString(t *testing.T) {
	program := &ast.Program{
		Statement: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got = %q", program.String())
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993;
`
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statement) != 3 {
		t.Fatalf("program.Statement does contain 3 statements. got = %d",
			len(program.Statement))
	}

	// TODO: 添加对值的检查
	for _, stmt := range program.Statement {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not ReturnStatement. got = %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return'. got = %q", returnStmt.TokenLiteral())
		}
	}
}

func TestLetStatement(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 78787878;
`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statement) != 3 {
		t.Fatalf("program.Statement does contain 3 statements. got = %d",
			len(program.Statement))
	}

	// TODO: 添加对值的检查
	expects := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range expects {
		stmt := program.Statement[i]
		testLetStatement(t, stmt, tt.expectedIdentifier)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statement) != 1 {
		t.Fatalf("IntegerLiteralExpression does contain 1 statements. got = %d",
			len(program.Statement))
	}
	stmt, ok := program.Statement[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ExpressionStatement. got = %T",
			program.Statement[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not IntegerLiteral. got = %T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value is not %s. got %d", input, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral is not %s. got %s", input, literal.TokenLiteral())
	}

}

// 仅有两种前缀运算符： ！ 和 -
func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statement) != 1 {
			t.Fatalf("program.Statement should contain 1 statements. got = %d",
				len(program.Statement))
		}
		stmt, ok := program.Statement[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not ExpressionStatement. got = %T",
				program.Statement[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not prefixExpression. got = %T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not %s. got %s", tt.operator, exp.Operator)
		}
		// 测试右值是否是整数字面量，并检查值是否正确
		testIntegerLiteral(t, exp.Right, tt.integerValue)

	}

}

func TestParsingInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		leftVal  int64
		operator string
		rightVal int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statement) != 1 {
			t.Fatalf("program.Statement should have 1 statement. got %d\n",
				len(program.Statement))
		}

		stmt, ok := program.Statement[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not ast.ExpressionStatement. got %T\n",
				program.Statement[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not InfixExpression. got %T\n",
				stmt.Expression)
		}

		testIntegerLiteral(t, exp.Left, tt.leftVal)

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s. got %s", tt.operator, exp.Operator)
		}

		testIntegerLiteral(t, exp.Right, tt.rightVal)
	}

}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b",
			"((-a) * b)",
		}, {
			"!-a", "(!(-a))",
		}, {
			"a + b + c",
			"((a + b) + c)",
		}, {
			"a + b - c",
			"((a + b) - c)",
		}, {
			"a * b * c",
			"((a * b) * c)",
		}, {
			"a * b / c",
			"((a * b) / c)",
		}, {
			"a + b / c",
			"(a + (b / c))",
		}, {
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		}, {
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		}, {
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		}, {
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		}, {
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foo"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statement) != 1 {
		t.Fatalf("IdentifierExpression does contain 1 statements. got = %d",
			len(program.Statement))
	}
	stmt, ok := program.Statement[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ExpressionStatement. got = %T",
			program.Statement[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not Indentifier. got = %T", stmt.Expression)
	}
	if ident.Value != "foo" {
		t.Errorf("ident.Value is not %s. got %s", input, ident.Value)
	}
	if ident.TokenLiteral() != "foo" {
		t.Errorf("ident.TokenLiteral is not %s. got %s", input, ident.TokenLiteral())
	}

}

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not 'let'. got = %q", s.TokenLiteral())
		return
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("the statement is not LetStatement type. got = %T", s)
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got = '%s'", name, letStmt.Name.Value)
		return
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got = '%s'",
			name, letStmt.Name.TokenLiteral())
		return
	}
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Error()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	integ, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp is not IntegerLiteral. got %T", exp)
	}

	if integ.Value != value {
		t.Fatalf("integ.Value is not %d. got %d", value, integ.Value)
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Fatalf("integ.TokenLiteral is not %d. got %s", value, integ.TokenLiteral())
	}
}
