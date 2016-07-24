package brainfuck

import (
	"errors"
	"strings"
)

//go:generate go tool yacc brainfuck.y

type Program struct {
	Code []Code
}

type Code struct {
	Command rune
	Comment string
	Loop    []Code
}

func Parse(program string) (Program, error) {
	parser := yyNewParser()
	lexer := &lexer{input: program}
	_ = parser.Parse(lexer)
	if lexer.error == "" {
		return lexer.result, nil
	}
	return Program{}, errors.New(lexer.error)
}

type lexer struct {
	index  int
	input  string
	error  string
	result Program
}

func (lex *lexer) Lex(lval *yySymType) int {
	if lex.index == len(lex.input) {
		return 0
	}

	index := lex.index
	switch lex.input[index] {
	case '+', '-', '<', '>', ',', '.':
		lex.index++
		lval.Code = Code{Command: rune(lex.input[index])}
		return tkInstruction
	case '[':
		lex.index++
		return tkLoopStart
	case ']':
		lex.index++
		return tkLoopEnd
	default:
		rest := lex.input[index:]
		afterComment := strings.TrimFunc(rest, func(r rune) bool {
			return !strings.ContainsRune("+-<>.,[]", r)
		})
		commentLength := len(rest) - len(afterComment)
		lval.Code = Code{Comment: rest[:commentLength]}
		lex.index += commentLength
		return tkComment
	}
}

func (lex *lexer) Error(s string) {
	lex.error = s
}
