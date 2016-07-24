%{
package json
%}

// Defines the yySymType, which is used during parsing. For parsing
// JSON, all intermediate parse nodes can be represented as a Value.
// For more complex grammars, more than one field is needed.
%union {
	Value *Value
}

// By convention, tokens are written in UPPERCASE while grammar
// productions are written in snake_case.

// Tokens without additional information. These are part of the
// concrete syntax, but not of the abstract one.
%token TARRAYOPEN TARRAYCLOSE
%token TOBJECTOPEN TOBJECTCLOSE
%token TCOLON TCOMMA

// Tokens that carry additional information. The lexer has to fill in
// the additional information into the `Value` field.
%token <Value> TSTRING TNUMBER
%token <Value> TFALSE TTRUE TNULL

// Each of the following grammar productions has a result type, which
// must be one of the field names of the `%union` definition above.
%type <Value> start value objectvalues arrayvalues

%%

start : value {
	yylex.(*lexer).result = $1
}

// The tokens TNUMBER and below don’t have explicit code, since they
// only copy the value, i.e. their code would be `{ $$ = $1 }`, which
// can be omitted.
value
	: TOBJECTOPEN objectvalues TOBJECTCLOSE { $$ = $2 }
	| TARRAYOPEN arrayvalues TARRAYCLOSE { $$ = $2 }
	| TNUMBER
	| TSTRING
	| TTRUE
	| TFALSE
	| TNULL
	;

objectvalues
	: TSTRING TCOLON value {
		$$ = &Value{Type: VOBJECT, Properties: map[string]Value { $1.String: *$3 }}
	}
	| objectvalues TCOMMA TSTRING TCOLON value {
		$$.Properties[$3.String] = *$5
	}
	| /* empty */ {
		$$ = &Value{Type: VOBJECT, Properties: make(map[string]Value)}
	}
	;

arrayvalues
	: value {
		$$ = &Value{Type: VARRAY, Elements: []Value{*$1}}
	}
	| arrayvalues TCOMMA value {
		$$.Elements = append($$.Elements, *$3)
	}
	| /* empty */ {
		$$ = &Value{Type: VARRAY}
	}
	;

%%
