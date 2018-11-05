package wktToOrb

import (
	"strings"
	"testing"
)

func Test_scanToken(t *testing.T) {
	inputs := []string{
		"(",
		")",
		"empty",
		"z",
		"ZM",
		"M",
		"POINT",
		"3.14",
		"12",
		",",
		"-23",
		"5e-05",
	}

	outputs := []Token{
		{ttype: LeftParen, lexeme: "("},
		{ttype: RightParen, lexeme: ")"},
		{ttype: Empty, lexeme: "empty"},
		{ttype: Z, lexeme: "z"},
		{ttype: ZM, lexeme: "zm"},
		{ttype: M, lexeme: "m"},
		{ttype: Point, lexeme: "point"},
		{ttype: Float, lexeme: "3.14"},
		{ttype: Float, lexeme: "12"},
		{ttype: Comma, lexeme: ","},
		{ttype: Float, lexeme: "-23"},
		{ttype: Float, lexeme: "5e-05"},
	}

	for i, input := range inputs {
		l := NewLexer(strings.NewReader(input))
		ok, err := l.scanToken()
		if !ok {
			t.Error("scanToken reached unexpected eof")
		}
		if err != nil {
			t.Errorf("unexpected error %s", err)
		}
		if len(l.scanned) != 1 {
			t.Error("too many token scanned")
		}
		if l.scanned[0].ttype != outputs[i].ttype {
			t.Errorf("incorrect ttype for %s", input)
		}
		if l.scanned[0].lexeme != outputs[i].lexeme {
			t.Errorf("incorrect lexeme for %s", input)
		}
		if l.scanned[0].pos != 0 {
			t.Errorf("incorrect position for %s", input)
		}
	}
}

func Test_Scan(t *testing.T) {
	input := " ( ) empty z ZM M POINT 3.14 12   , "
	output := []Token{
		{ttype: LeftParen, lexeme: "(", pos: 1},
		{ttype: RightParen, lexeme: ")", pos: 3},
		{ttype: Empty, lexeme: "empty", pos: 5},
		{ttype: Z, lexeme: "z", pos: 11},
		{ttype: ZM, lexeme: "zm", pos: 13},
		{ttype: M, lexeme: "m", pos: 16},
		{ttype: Point, lexeme: "point", pos: 18},
		{ttype: Float, lexeme: "3.14", pos: 24},
		{ttype: Float, lexeme: "12", pos: 29},
		{ttype: Comma, lexeme: ",", pos: 34},
		{ttype: Eof, lexeme: "", pos: 36},
	}

	l := NewLexer(strings.NewReader(input))
	err := l.Scan()
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	if len(l.scanned) != len(output) {
		t.Error("incorrect number of tokens scanned")
	}

	for i, token := range l.scanned {
		if token.ttype != output[i].ttype {
			t.Errorf("incorrect ttype for token %d", i)
		}
		if token.lexeme != output[i].lexeme {
			t.Errorf("incorrect lexeme for token %d", i)
		}
		if token.pos != output[i].pos {
			t.Errorf("incorrect position for token %d", i)
		}
	}
}
