%{
package json
%}

%union {
	Value *Value
}

%type <Value> value keyvalues keyvalue values

%token TARRAYOPEN TARRAYCLOSE
%token TOBJECTOPEN TOBJECTCLOSE
%token TCOLON TCOMMA
%token <Value> TSTRING TNUMBER
%token <Value> TFALSE TTRUE TNULL

%%

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
	: keyvalue
	| keyvalues TCOMMA keyvalue {
		for k, v := range $3.Properties {
			$$.Properties[k] = v
		}
	}
	| /* empty */ {
		$$ = &Value{Type: VOBJECT, Properties: make(map[string]Value)}
	}
	;

keyvalue
	: TSTRING TCOLON value {
		$$ = &Value{Type: VOBJECT, Properties: map[string]Value { $1.String: *$3 }}
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
