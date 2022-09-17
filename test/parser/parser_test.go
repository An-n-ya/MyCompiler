package parser

import (
	"MyCompiler/src/ast"
	"MyCompiler/src/lexer"
	"MyCompiler/src/parser"
	"testing"
)

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
