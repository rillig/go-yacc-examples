package json

import "errors"

//go:generate goyacc json.y

// Value represents a JSON value. Depending on the `Type`, one of the
// other fields contains the actual value.
type Value struct {
	Type       ValueType
	Properties map[string]Value
	Elements   []Value
	String     string
	Number     float64
}

// Token is a basic building block of JSON text, e.g. an opening brace
// or a number.
type Token struct {
	TokenType int
	String    string
	Number    float64
}

type ValueType uint8

const (
	VOBJECT ValueType = iota
	VARRAY
	VSTRING
	VNUMBER
	VTRUE
	VFALSE
	VNULL
)

func (t ValueType) String() string {
	return [...]string{
		"VOBJECT", "VARRAY", "VSTRING", "VNUMBER", "VTRUE", "VFALSE", "VNULL",
	}[t]
}

// Parse converts a flat list of tokens into a tree of JSON values.
func Parse(tokens []Token) (Value, error) {
	p := yyNewParser()
	lexer := &lexer{tokens: tokens}
	_ = p.Parse(lexer)
	if lexer.error == "" {
		// See http://stackoverflow.com/q/36822702
		return *lexer.result, nil
	}
	return Value{}, errors.New(lexer.error)
}

type lexer struct {
	i      int
	tokens []Token
	error  string
	result *Value
}

func (l *lexer) Lex(lval *yySymType) int {
	if l.i == len(l.tokens) {
		return 0
	}

	token := l.tokens[l.i]
	l.i++

	// Those tokens that have an associated type in json.y (see the
	// %token definitions) must fill the corresponding field of the
	// SymType.
	switch token.TokenType {
	case TSTRING:
		lval.Value = &Value{Type: VSTRING, String: token.String}
	case TNUMBER:
		lval.Value = &Value{Type: VNUMBER, Number: token.Number}
	case TNULL:
		lval.Value = &Value{Type: VNULL}
	case TTRUE:
		lval.Value = &Value{Type: VTRUE}
	case TFALSE:
		lval.Value = &Value{Type: VFALSE}
	}
	return token.TokenType
}

func (l *lexer) Error(s string) {
	l.error = s
}
