package token

// TokenType Token的类型
// 使用string是为了方便调试，如果使用int之类的可读性会变差，但是性能会更高
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// TODO: 支持字符串字面量 支持浮点数
// 所有的类型
const (
	// 特殊类型
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 标识符 字面量
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// 运算符
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"
	GT       = ">"
	LT       = "<"

	EQ     = "=="
	NOT_EQ = "!="

	// 分隔符
	COMMA     = ","
	SEMICOLON = ";"

	// 括号
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// 关键字
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "if"
	ELSE     = "else"
	RETURN   = "return"
	TRUE     = "true"
	FALSE    = "false"
)

// 所有的关键字
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

// 关键字匹配
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
