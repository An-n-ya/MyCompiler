package main

import (
	"MyCompiler/src/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
