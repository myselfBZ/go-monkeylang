package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/myselfBZ/interpreter/internal/ast"
)

type ObjType string

const (
    FUNCTION_OBJ = "FUNCTION"
	INTEGER_OBJ = "INTIGER_TYPE"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL        = "NULL"
    RETURN_VALUE = "RETURN_VALUE"
    ERROR_OBJ = "ERROR"
)

type Object interface {
	Type() ObjType
	Inspect() string
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: outer,
	}
}

func NewEnvironment() *Environment{
    return &Environment{
        store:make(map[string]Object),
		outer: nil,
    }
}

type Environment struct{
	outer *Environment
    store map[string]Object
}

func findIdent(env *Environment, name string) (Object, bool) {
	if env == nil {
		return nil, false
	}

	obj, ok := env.store[name]
	if !ok {
		return findIdent(env.outer, name)
	}
	return obj, true
}

func (e *Environment) Get(name string) (Object, bool) {
	return findIdent(e, name)
}

func (e *Environment) Set(name string, obj Object) Object {
    e.store[name] = obj
    return obj
}

type Integer struct {
	Value int
}

func (i *Integer) Type() ObjType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}
func (b *Boolean) Type() ObjType {
	return BOOLEAN_OBJ
}

type Null struct{}

func (n *Null) Type() ObjType {
	return NULL
}
func (n *Null) Inspect() string {
	return "NULL"
}


type ReturnValue struct{
    Value Object
}

func (r *ReturnValue) Type() ObjType{
    return RETURN_VALUE
}
func (r *ReturnValue) Inspect() string{
    return fmt.Sprintf("%s", r.Value.Inspect())
}


type Error struct{
    Message string
}

func (e *Error) Type() ObjType{
    return ERROR_OBJ
}
func (e *Error) Inspect() string{
    return e.Message
}


type Function struct{
    Params []*ast.Identifier
    Body   *ast.BlockStatement
    Env    *Environment
}

func (f *Function) Type() ObjType {
    return FUNCTION_OBJ
}
func (f *Function) Inspect() string{
    var out bytes.Buffer
    params := []string{}
    for _, p := range f.Params{
        params = append(params, p.String())
    }
    out.WriteString("fn")
    out.WriteString("(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(") {\n")
    out.WriteString(f.Body.String())
    out.WriteString("\n}")
    return out.String()
}


