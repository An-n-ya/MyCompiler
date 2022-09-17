package lexer

import (
	"MyCompiler/src/lexer"
	"MyCompiler/src/token"
	"testing"
)

type expectStruct []struct {
	expectedType    token.TokenType
	expectedLiteral string
}

type testSet struct {
	input   string
	expects expectStruct
	name    string
}

func TestNextToken(t *testing.T) {

	basicToken := testSet{
		"=+(){},;*/!<> true false == !=",
		expectStruct{
			{token.ASSIGN, "="},
			{token.PLUS, "+"},
			{token.LPAREN, "("},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RBRACE, "}"},
			{token.COMMA, ","},
			{token.SEMICOLON, ";"},
			{token.ASTERISK, "*"},
			{token.SLASH, "/"},
			{token.BANG, "!"},
			{token.LT, "<"},
			{token.GT, ">"},
			{token.TRUE, "true"},
			{token.FALSE, "false"},
			{token.EQ, "=="},
			{token.NOT_EQ, "!="},
			{token.EOF, ""},
		},
		"basicToken",
	}

	expAndFunc := testSet{
		`let five = 5;
let ten = 10;
let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);

if (ten < five) {
	return true;
} else {
	return false;
}
`,
		expectStruct{
			{token.LET, "let"},
			{token.IDENT, "five"},
			{token.ASSIGN, "="},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "ten"},
			{token.ASSIGN, "="},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "add"},
			{token.ASSIGN, "="},
			{token.FUNCTION, "fn"},
			{token.LPAREN, "("},
			{token.IDENT, "x"},
			{token.COMMA, ","},
			{token.IDENT, "y"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.IDENT, "x"},
			{token.PLUS, "+"},
			{token.IDENT, "y"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "result"},
			{token.ASSIGN, "="},
			{token.IDENT, "add"},
			{token.LPAREN, "("},
			{token.IDENT, "five"},
			{token.COMMA, ","},
			{token.IDENT, "ten"},
			{token.RPAREN, ")"},
			{token.SEMICOLON, ";"},
			{token.IF, "if"},
			{token.LPAREN, "("},
			{token.IDENT, "ten"},
			{token.LT, "<"},
			{token.IDENT, "five"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.TRUE, "true"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.ELSE, "else"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.FALSE, "false"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.EOF, ""},
		},
		"expAndFunc",
	}

	tests := []testSet{
		basicToken,
		expAndFunc,
	}
	for _, test := range tests {
		input, expects, name := test.input, test.expects, test.name
		l := lexer.New(input)

		for i, tt := range expects {
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Errorf("%d of %s - tokentype wrong. expected=%q, got=%q",
					i, name, tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("%d of %s - literal wrong. expected=%q, got=%q",
					i, name, tt.expectedLiteral, tok.Literal)
			}
		}

	}

}
