%{
package json
%}

// Defines the yySymType, which is used during parsing. For parsing
// JSON, all intermediate parse nodes can be represented as a Value.
// For more complex grammars, more than one field is needed.
%union {
	Value *Value
}

// Each of the following grammar productions has a result type, which
// must be one of the field names of the `%union` definition above.
%type <Value> value keyvalues values

// Tokens without additional information. These are part of the
// concrete syntax, but not of the abstract one.
%token TARRAYOPEN TARRAYCLOSE
%token TOBJECTOPEN TOBJECTCLOSE
%token TCOLON TCOMMA

// Tokens that carry additional information. The lexer has to fill in
// the additional information into the `Value` field.
%token <Value> TSTRING TNUMBER
%token <Value> TFALSE TTRUE TNULL

%%

// The tokens starting with TNUMBER don’t have explicit code, since
// they only copy the value, i.e. their code would be `{ $$ = $1 }`,
// which can be omitted.
value
	: TOBJECTOPEN keyvalues TOBJECTCLOSE { $$ = $2 }
	| TARRAYOPEN values TARRAYCLOSE { $$ = $2 }
	| TNUMBER
	| TSTRING
	| TTRUE
	| TFALSE
	| TNULL
	;

keyvalues
	: TSTRING TCOLON value {
		$$ = &Value{Type: VOBJECT, Properties: map[string]Value { $1.String: *$3 }}
	}
	| keyvalues TCOMMA TSTRING TCOLON value {
		$$.Properties[$3.String] = *$5
	}
	| /* empty */ {
		$$ = &Value{Type: VOBJECT, Properties: make(map[string]Value)}
	}
	;

values
	: value {
		$$ = &Value{Type: VARRAY, Elements: []Value{*$1}}
	}
	| values TCOMMA value {
		$$.Elements = append($$.Elements, *$3)
	}
	| /* empty */ {
		$$ = &Value{Type: VARRAY}
	}
	;

%%
