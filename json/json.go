package json

//go:generate go tool yacc json.y

type VType uint8

const (
	VOBJECT VType = iota
	VARRAY
	VSTRING
	VNUMBER
	VTRUE
	VFALSE
	VNULL
)

func (t VType) String() string {
	return [...]string{
		"VOBJECT", "VARRAY", "VSTRING", "VNUMBER", "VTRUE", "VFALSE", "VNULL",
	}[t]
}

type Value struct {
	Type       VType
	Properties map[string]Value
	Elements   []Value
	String     string
	Number     float64
}

type Token struct {
	TokenType int
	String    string
	Number    float64
}

func Parse(tokens []Token) (Value, error) {
	p := &yyParserImpl{}
	lexer := &lexer{0, tokens, ""}
	_ = p.Parse(lexer)
	if lexer.error == "" {
		return *p.stack[1].Value, nil
	}
	return Value{}, jsonerror(lexer.error)
}

type lexer struct {
	i      int
	tokens []Token
	error  string
}
type jsonerror string

func (e jsonerror) Error() string {
	return string(e)
}

func (l *lexer) Lex(lval *yySymType) int {
	if l.i == len(l.tokens) {
		return 0
	}

	token := l.tokens[l.i]
	l.i++
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
