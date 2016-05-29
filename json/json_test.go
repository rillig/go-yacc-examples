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

	test := func(expected Value, tokens ...Token) {
		actual, error := Parse(tokens)
		if c.Check(error, check.IsNil) {
			c.Check(actual, check.DeepEquals, expected)
		}
	}

	test(jnull, tnull)
	test(jfalse, tfalse)
	test(jtrue, ttrue)
	test(jnumber(1.2), tnumber(1.2))
	test(jstring("str"), tstring("str"))
	test(jobject(), tobjectopen, tobjectclose)
	test(jarray(), tarrayopen, tarrayclose)

	test(jobject(jkeyvalue("key", jstring("value"))),
		tobjectopen,
		tstring("key"), tcolon, tstring("value"),
		tobjectclose)

	test(
		jobject(
			jkeyvalue("key1", jstring("value1")),
			jkeyvalue("key2", jstring("value2"))),

		tobjectopen,
		tstring("key1"), tcolon, tstring("value1"),
		tcomma,
		tstring("key2"), tcolon, tstring("value2"),
		tobjectclose)

	test(
		jarray(jstring("element")),

		tarrayopen,
		tstring("element"),
		tarrayclose)

	test(
		jarray(
			jstring("element0"),
			jnumber(1.0),
			jtrue,
			jfalse,
			jnull),

		tarrayopen,
		tstring("element0"),
		tcomma,
		tnumber(1.0),
		tcomma,
		ttrue,
		tcomma,
		tfalse,
		tcomma,
		tnull,
		tarrayclose)
}
