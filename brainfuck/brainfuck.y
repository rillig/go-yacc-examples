%{
package brainfuck
%}

%union {
	Code Code
	Codes []Code
}

%token <Code> tkLoopStart tkLoopEnd
%token <Code> tkInstruction tkComment

%type <Codes> start instructions instructions
%type <Code> instruction

%%

start : instructions {
    yylex.(*lexer).result = Program{$1}
}

instructions : /* empty */ {
    $$ = []Code{}
}
instructions : instructions instruction {
    $$ = append($1, $2)
}

instruction : tkLoopStart instructions tkLoopEnd {
    $$ = Code{Loop: $2}
}
instruction : tkInstruction {
    $$ = $1
}
instruction : tkComment {
    $$ = $1
}

%%
