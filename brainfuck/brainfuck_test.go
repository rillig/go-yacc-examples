package brainfuck

import (
	"testing"
	"strings"
)

func TestParse(t *testing.T) {
	program, err := Parse("+++[>,Comment.<-]")

	if err != nil {
		t.Errorf("Parse error: %q", err.Error())
	} else {
		expected := "+++\n[\n  >, Comment .<-\n]\n"
		indenter := &Indenter{Indentation:"  "}
		actual := indenter.Indent(program)
		if expected != actual {
			t.Errorf("Wrong program: expected %q, got %q", expected, actual)
		}
	}
}

type Indenter struct {
	Indentation string
	depth int
	output string
}

func (ind *Indenter) Indent(program Program) string {
	ind.output = ""
	ind.indentCodes(program.Code)
	return ind.output
}

func (ind *Indenter) indentCodes(codes []Code) {
	for _, code := range codes {
		ind.indentCode(code)
	}
}

func (ind *Indenter) indentCode(code Code) {
	switch {
	case code.Command != 0:
		ind.output += string([]rune{code.Command})
	case code.Comment != "":
		ind.output += " " + strings.TrimSpace(code.Comment) + " "
	case code.Loop != nil:
		ind.output += "\n"
		ind.output += ind.indent() + "[\n"
		ind.depth++
		ind.output += ind.indent()
		ind.indentCodes(code.Loop)
		ind.depth--
		ind.output += "\n" + ind.indent() + "]\n"
	}
}

func (ind *Indenter) indent() string {
	return strings.Repeat(ind.Indentation, ind.depth)
}
