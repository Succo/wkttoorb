package main

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
		{ttype: Eof, lexeme: ""},
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
	}
}
