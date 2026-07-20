package evaluator

// testing comment
import (
	"testing"

	"github.com/myselfBZ/interpreter/internal/lexer"
	"github.com/myselfBZ/interpreter/internal/object"
	"github.com/myselfBZ/interpreter/internal/parser"
)

func TestEvaluatorInt(t *testing.T) {
	input := struct {
		input    string
		expected int
	}{
		"6",
		6,
	}
	l := lexer.New(input.input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	obj := Eval(program, env)
	i, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("expected int got %T\n", obj.(*object.Integer))
	}
	if i.Value != input.expected {
		t.Fatalf("expected %d got %d", input.expected, i.Value)
	}
}

func TestEvalBoolean(t *testing.T) {
	input := struct {
		input string
		expct bool
	}{
		"false;",
		false,
	}
	l := lexer.New(input.input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	obj := Eval(program, env)
	b, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("expected boolean got %T", obj)
	}
	if b.Value != input.expct {
		t.Fatalf("expected %v got %v", input.expct, b.Value)
	}
}

func TestBang(t *testing.T) {
	input := struct {
		input string
		expct bool
	}{
		"!true",
		false,
	}
	l := lexer.New(input.input)
	p := parser.New(l)
	env := object.NewEnvironment()
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("errors: %s", p.Errors()[0])
	}
	v := Eval(program, env)
	b, ok := v.(*object.Boolean)
	if !ok {
		t.Fatalf("expected boolean object got %T", v.(*object.Boolean))
	}
	if b.Value != input.expct {
		t.Fatalf("expected %v got %v", input.expct, b.Value)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
		result.Value, expected)
		return false
	}
	return true
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input string
		expected int
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
