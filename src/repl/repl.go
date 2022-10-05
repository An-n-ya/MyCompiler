package repl

import (
	"MyCompiler/src/evaluator"
	"MyCompiler/src/lexer"
	"MyCompiler/src/object"
	"MyCompiler/src/parser"
	"MyCompiler/src/token"
	"bufio"
	"fmt"
	"io"
	"os"
)

const PROMPT = ">> "

// TODO: 解析命令行参数
func Start(in io.Reader, out io.Writer) {
	if len(os.Args) == 1 {
		EvaluateStart(in, out)
		return
	}
	helpMsg := `
Usage:

	liu <command>

The commands are:

	lexer/lex       show the lexer structure
	parser/ast      show the ast structure
	[default]       evaluate the expression
	
`
	switch os.Args[1] {
	case "lexer", "lex":
		LexerStart(in, out)
	case "parser", "ast":
		ParserStart(in, out)
	case "help":
		fmt.Println(helpMsg)
	default:
		fmt.Println(os.Args[1], ": unknown command\nRun 'help' for usage")
	}
}

func LexerStart(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}

func ParserStart(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Error()) != 0 {
			printParserErrors(out, p.Error())
			continue
		}
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func EvaluateStart(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	// 创建变量表
	env := object.NewEnvironment()
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Error()) != 0 {
			printParserErrors(out, p.Error())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

// region 帮助函数

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

// endregion
