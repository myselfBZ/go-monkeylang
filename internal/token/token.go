package token

type TokenType string

type Token struct {
	Type    TokenType `json:"type"`
	Literal string    `json:"literal"`
}

func NewToken(kind TokenType, literal string) Token {
	return Token{Type: kind, Literal: literal}
}

var Keywords = map[string]string{
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"false":  FALSE,
	"true":   TRUE,
	"let":    LET,
}

const (
	ILLEGAL        = "ILLEGAL"
	EOF            = "EOF"
	IDENT          = "IDENT"
	INT            = "INT"
	ASSIGN         = "="
	PLUS           = "+"
	COMMA          = ","
	SEMICOLON      = ";"
	LPAREN         = "("
	RPAREN         = ")"
	LBRACE         = "{"
	RBRACE         = "}"
	FUNCTION       = "FUNCTION"
	LET            = "LET"
	MINUS          = "-"
	DIVISION       = "/"
	MULTIPLICATION = "*"
	LT             = "<"
	GT             = ">"
	RETURN         = "RETURN"
	IF             = "IF"
	ELSE           = "ELSE"
	TRUE           = "TRUE"
	FALSE          = "FALSE"
	EQ             = "=="
	NOT_EQ         = "!="
	GTOREQ         = ">="
	LTOREQ         = "<="
	BANG           = "!"
	STRING		   = "STRING"
)
