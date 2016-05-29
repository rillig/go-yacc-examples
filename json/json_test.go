package json

import (
	"gopkg.in/check.v1"
	"testing"
)

type Suite struct{}

var _ = check.Suite(new(Suite))

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *Suite) TestParse(c *check.C) {
	var currentTokens []Token
	given := func(tokens ...Token) {
		currentTokens = tokens
	}
	expect := func(expected Value) {
		actual, error := Parse(currentTokens)
		if c.Check(error, check.IsNil) {
			c.Check(actual, check.DeepEquals, expected)
		}
	}

	tstring := func(s string) Token {
		return Token{TokenType: TSTRING, String: s}
	}
	tnumber := func(n float64) Token {
		return Token{TokenType: TNUMBER, Number: n}
	}
	tobjectopen := Token{TokenType: TOBJECTOPEN}
	tobjectclose := Token{TokenType: TOBJECTCLOSE}
	tarrayopen := Token{TokenType: TARRAYOPEN}
	tarrayclose := Token{TokenType: TARRAYCLOSE}
	tcolon := Token{TokenType: TCOLON}
	tcomma := Token{TokenType: TCOMMA}
	ttrue := Token{TokenType: TTRUE}
	tfalse := Token{TokenType: TFALSE}
	tnull := Token{TokenType: TNULL}

	type keyvalue struct {
		key   string
		value Value
	}
	jobject := func(keyvalues ...keyvalue) Value {
		value := Value{Type: VOBJECT, Properties: make(map[string]Value)}
		for _, keyvalue := range keyvalues {
			value.Properties[keyvalue.key] = keyvalue.value
		}
		return value
	}
	jkeyvalue := func(key string, value Value) keyvalue {
		return keyvalue{key, value}
	}
	jarray := func(elements ...Value) Value {
		return Value{Type: VARRAY, Elements: elements}
	}
	jnumber := func(n float64) Value {
		return Value{Type: VNUMBER, Number: n}
	}
	jstring := func(s string) Value {
		return Value{Type: VSTRING, String: s}
	}
	jtrue := Value{Type: VTRUE}
	jfalse := Value{Type: VFALSE}
	jnull := Value{Type: VNULL}

	given(tnull)
	expect(jnull)

	given(tfalse)
	expect(jfalse)

	given(ttrue)
	expect(jtrue)

	given(tnumber(1.2))
	expect(jnumber(1.2))

	given(tstring("str"))
	expect(jstring("str"))

	given(tobjectopen, tobjectclose)
	expect(jobject())

	given(tarrayopen, tarrayclose)
	expect(jarray())

	given(tobjectopen,
		tstring("key"), tcolon, tstring("value"),
		tobjectclose)
	expect(jobject(jkeyvalue("key", jstring("value"))))

	given(tobjectopen,
		tstring("key1"), tcolon, tstring("value1"), tcomma,
		tstring("key2"), tcolon, tstring("value2"),
		tobjectclose)
	expect(jobject(
		jkeyvalue("key1", jstring("value1")),
		jkeyvalue("key2", jstring("value2"))))

	given(tarrayopen,
		tstring("element"),
		tarrayclose)
	expect(jarray(jstring("element")))

	given(tarrayopen,
		tstring("element0"), tcomma,
		tnumber(1.0), tcomma,
		ttrue, tcomma,
		tfalse, tcomma,
		tnull,
		tarrayclose)
	expect(jarray(
		jstring("element0"),
		jnumber(1.0),
		jtrue,
		jfalse,
		jnull))
}
