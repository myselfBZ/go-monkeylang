package ast

import (
	"bytes"
	"strings"

	"github.com/myselfBZ/interpreter/internal/token"
)

type Node interface {
	String() string
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, v := range p.Statements {
		out.WriteString(v.String() + "\n")
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// ExpressionStatement
type ExpressionStatement struct {
	Token      *token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode() { return }

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

// Identifier
type Identifier struct {
	Token *token.Token `json:"token"`
	Value string       `json:"value"`
}

func (i *Identifier) expressionNode() { return }

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// infix operators
type InfixExperssion struct {
	Left     Expression
	Right    Expression
	Operator string
	Token    *token.Token
}

func (i *InfixExperssion) expressionNode() { return }
func (i *InfixExperssion) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}
func (i *InfixExperssion) TokenLiteral() string { return i.Token.Literal }

// int literals
type IntLiteral struct {
	Value int64
	Token *token.Token
}

func (i *IntLiteral) expressionNode() { return }
func (i *IntLiteral) String() string {
	return i.Token.Literal
}
func (i *IntLiteral) TokenLiteral() string {
	return i.Token.Literal
}

type LetStatement struct {
	Token *token.Token `json:"token"`
	Value Expression   `json:"value"`
	Name  *Identifier  `json:"name"`
}

func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

func (l *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String() + "=")
	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	return out.String()
}

func (l *LetStatement) statementNode() {
	return
}

// return statements
type ReturnStatement struct {
	Token       *token.Token `json:"token"`
	ReturnValue Expression   `json:"return_value"`
}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}
	return out.String()
}

func (r *ReturnStatement) statementNode() {
	return
}

// prefix expressions
type PrefixExpression struct {
	Token    *token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() { return }
func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

type Boolean struct {
	Token *token.Token
	Value bool
}

func (b *Boolean) expressionNode() { return }
func (b *Boolean) String() string {
	return b.TokenLiteral()
}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

type IfExpression struct {
	Token       *token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode()      { return }
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" " + i.Consequence.String())
	if i.Alternative != nil {
		out.WriteString(" else " + i.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      *token.Token // {
	Statements []Statement
}

func (b *BlockStatement) statementNode() { return }
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}
func (b *BlockStatement) String() string {
	var out bytes.Buffer
	for _, v := range b.Statements {
		out.WriteString(v.String() + "\n")
	}
	return out.String()
}

type FunctionLiteral struct {
	Token  *token.Token
	Params []*Identifier
	Body   *BlockStatement
}

func (f *FunctionLiteral) expressionNode() { return }
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(f.TokenLiteral() + " ")
	out.WriteString("(")
	params := []string{}
	for _, p := range f.Params {
		params = append(params, p.String())
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString("{\n")
	statements := []string{}
	for _, s := range f.Body.Statements {
		statements = append(statements, s.String())
	}
	out.WriteString(strings.Join(statements, "\n"))
	out.WriteString("\n}")
	return out.String()
}
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

type Call struct {
	Token     *token.Token
	Function  Expression
	Arguments []Expression
}

func (c *Call) expressionNode() {
	return
}
func (c *Call) TokenLiteral() string {
	return c.Token.Literal
}
func (c *Call) String() string {
	var out bytes.Buffer
	out.WriteString(c.Function.String())
	out.WriteString("(")
	var arguemnts []string
	for _, v := range c.Arguments {
		arguemnts = append(arguemnts, v.String())
	}
	out.WriteString(strings.Join(arguemnts, ", "))
	out.WriteString(")")
	return out.String()
}

type String struct {
	Token *token.Token
	Value string
}

func (s *String) expressionNode() {}  
func (s *String) TokenLiteral() string {
	return s.Value
}
func (s *String) String() string {
	return "\"" + s.Value + "\""
}
