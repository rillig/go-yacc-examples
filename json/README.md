# Example JSON parser using go-yacc

The file `json.y` contains a Yacc parser that parses JSON tokens into
JSON values.

To build it, run:

```
go generate
go test
go install
```

Note that this package only contains the parser, not the tokenizer.