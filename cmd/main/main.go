package main

import (
	"fmt"
	"log"
	"os"

	"github.com/myselfBZ/interpreter/internal/evaluator"
	"github.com/myselfBZ/interpreter/internal/lexer"
	"github.com/myselfBZ/interpreter/internal/object"
	"github.com/myselfBZ/interpreter/internal/parser"
)

func open(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("couldn't find the file")
	}
	buff := make([]byte, 1024)
	size, err := file.Read(buff)
	if err != nil {
		log.Fatal("error reading from a file, ", err)
	}
	return string(buff[:size])
}

func main() {
    env := object.NewEnvironment()
	src := open("test.monkey")
	l := lexer.New(src)
    p := parser.New(l)
    program := p.ParseProgram()
    o := evaluator.Eval(program, env)
    if o != nil{
        fmt.Println(o.Inspect())
    }
}
