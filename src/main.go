package main

import (
	"MyComiler/src/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
