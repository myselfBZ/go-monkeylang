package parser

import (
	"fmt"
	"strconv"

	"github.com/myselfBZ/interpreter/internal/ast"
	"github.com/myselfBZ/interpreter/internal/lexer"
	"github.com/myselfBZ/interpreter/internal/token"
)


const (
	_ int = iota
	LOWEST
	EQUALS
	// ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type (
	parsePrefix func() ast.Expression
	parseInfix  func(ast.Expression) ast.Expression
)

var precedences = map[token.TokenType]int{
	token.MINUS:          SUM,
	token.PLUS:           SUM,
	token.MULTIPLICATION: PRODUCT,
	token.DIVISION:       PRODUCT,
	token.LT:             LESSGREATER,
	token.GTOREQ:         LESSGREATER,
	token.LTOREQ:         LESSGREATER,
	token.GT:             LESSGREATER,
	token.NOT_EQ:         EQUALS,
	token.EQ:             EQUALS,
	token.LPAREN:         CALL,
}

type Parser struct {
	lexer     *lexer.Lexer
	curToken  *token.Token
	errors    []string
	peekToken *token.Token
	prefixFns map[token.TokenType]parsePrefix
	infixFns  map[token.TokenType]parseInfix
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:     l,
		curToken:  l.NextToken(),
		peekToken: l.NextToken(),
		prefixFns: make(map[token.TokenType]parsePrefix),
		infixFns:  make(map[token.TokenType]parseInfix),
	}
	//prefix
	p.registerPrefix(token.STRING, p.parseString)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpressions)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.IDENT, p.parseIdent)
	p.registerPrefix(token.INT, p.parseInt)
	p.registerPrefix(token.MINUS, p.parsePrefixOps)
	p.registerPrefix(token.BANG, p.parsePrefixOps)
	//infix
	p.registerInfix(token.LPAREN, p.parseCall)
	p.registerInfix(token.GTOREQ, p.parseInfixExpression)
	p.registerInfix(token.LTOREQ, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.DIVISION, p.parseInfixExpression)
	p.registerInfix(token.MULTIPLICATION, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for !p.currentTokenIs(token.EOF) {
		stmnt := p.parseStatement(p.curToken)
		if stmnt != nil {
			program.Statements = append(program.Statements, stmnt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement(t *token.Token) ast.Statement {
	switch t.Type {
	case token.LET:
		return p.parseLet()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	node := &ast.ExpressionStatement{Token: p.curToken}
	node.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return node
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixExpression(p.curToken.Type)
		return nil
	}
	left := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixFns[p.peekToken.Type]
		if infix == nil {
			return left
		}
		p.nextToken()
		left = infix(left)
	}
	return left
}

func (p *Parser) parseLet() ast.Statement {
	node := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeekToken(token.IDENT) {
		return nil
	}
	node.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeekToken(token.ASSIGN) {
		p.errors = append(p.errors, fmt.Sprintf("expected '=' got %s\n", p.peekToken.Literal))
		return nil
	}
	p.nextToken()
	node.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return node
}

func (p *Parser) parseReturnStatement() ast.Statement {
	node := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	node.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return node
}

func (p *Parser) parseInt() ast.Expression {
	number, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("%s is not a number", p.curToken.Literal))
		return nil
	}
	node := &ast.IntLiteral{Token: p.curToken, Value: int64(number)}
	return node
}

func (p *Parser) parseIdent() ast.Expression {
	node := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	return node
}

// parsing prefix operators
func (p *Parser) parsePrefixOps() ast.Expression {
	node := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}

	p.nextToken()

	node.Right = p.parseExpression(PREFIX)

	return node
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	node := &ast.InfixExperssion{Token: p.curToken, Left: left, Operator: p.curToken.Literal}
	precedence := p.curPrecedence()
	p.nextToken()
	node.Right = p.parseExpression(precedence)
	return node
}

func (p *Parser) parseGroupedExpressions() ast.Expression {
	p.nextToken()
	exprsn := p.parseExpression(LOWEST)
	if !p.peekTokenIs(token.RPAREN) {
		p.errors = append(p.errors, fmt.Sprintf("expected '(' . Got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken()
	return exprsn
}

func (p *Parser) parseString() ast.Expression {
	return &ast.String{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.currentTokenIs(token.TRUE)}
}

func (p *Parser) parseIfExpression() ast.Expression {
	node := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeekToken(token.LPAREN) {
		p.errors = append(p.errors, fmt.Sprintf("expected ) got %s", p.peekToken.Literal))
		return nil
	}
	p.nextToken()
	node.Condition = p.parseExpression(LOWEST)
	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}
	// skip the closing parentheses
	p.nextToken()
	node.Consequence = p.parseBlockStatements()
	if p.expectPeekToken(token.ELSE) {
		if !p.expectPeekToken(token.LBRACE) {
			return nil
		}
		node.Alternative = p.parseBlockStatements()
	}
	return node
}

func (p *Parser) parseBlockStatements() *ast.BlockStatement {
	node := &ast.BlockStatement{Token: p.curToken}
	node.Statements = []ast.Statement{}
	p.nextToken()
	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmnt := p.parseStatement(p.curToken)
		if stmnt != nil {
			node.Statements = append(node.Statements, stmnt)
		}
		p.nextToken()
	}
	return node
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	node := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeekToken(token.LPAREN) {
		return nil
	}
	node.Params = p.parseParams()
	if !p.expectPeekToken(token.LBRACE) {
		return nil
	}
	node.Body = p.parseBlockStatements()
	if !p.currentTokenIs(token.RBRACE) {
		return nil
	}
	return node
}

func (p *Parser) parseParams() []*ast.Identifier {
	idents := []*ast.Identifier{}
	if p.expectPeekToken(token.RPAREN) {
		return idents
	}
	p.nextToken()
	ident1 := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	idents = append(idents, ident1)
	for p.expectPeekToken(token.COMMA) {
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		idents = append(idents, ident)
	}
	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}
	return idents
}

func (p *Parser) parseCall(f ast.Expression) ast.Expression {
	node := &ast.Call{Token: p.curToken, Function: f}
	node.Arguments = p.parseCallArguements()
	return node
}

func (p *Parser) parseCallArguements() []ast.Expression {
	var arguments []ast.Expression
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return arguments
	}
	p.nextToken()
	arguments = append(arguments, p.parseExpression(LOWEST))
	// this is for (a,b,s,)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		node := p.parseExpression(LOWEST)
		arguments = append(arguments, node)
	}
	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}
	return arguments
}
