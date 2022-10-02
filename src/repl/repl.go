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
)

const PROMPT = ">> "

// TODO: 解析命令行参数

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

func Start(in io.Reader, out io.Writer) {
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
