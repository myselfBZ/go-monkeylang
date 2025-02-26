package main

import (
	"fmt"
	"os"

	"github.com/myselfBZ/interpreter/internal/evaluator"
	"github.com/myselfBZ/interpreter/internal/lexer"
	"github.com/myselfBZ/interpreter/internal/object"
	"github.com/myselfBZ/interpreter/internal/parser"
	"github.com/peterh/liner"
)

func Start() {
	env := object.NewEnviroment()
	l := liner.NewLiner()
	defer l.Close()
	l.SetCtrlCAborts(true)
	his := ".history"
	if f, err := os.Open(his); err == nil {
		l.ReadHistory(f)
		f.Close()
	}

	fmt.Println("Welcome to the Monkey programming language! By Boburmirzo")
	for {
		input, err := l.Prompt(">>> ")
		if err != nil {
			if err == liner.ErrPromptAborted {
				fmt.Println("byeeee")
				break
			}
			fmt.Println("error: ", err)
			break
		}
		if input == "" {
			continue
		}
		l.AppendHistory(input)
		if input == "exit" || input == "quit" {
			break
		}
		lex := lexer.New(input)
		p := parser.New(lex)
		program := p.ParseProgram()
		e := evaluator.Eval(program, env)
		if e != nil {
			fmt.Println(e.Inspect())
		}

	}
}

func main() {
	Start()
}
