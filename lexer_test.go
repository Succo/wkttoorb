package main

import (
	"strings"
	"testing"
)

func Test_scan(t *testing.T) {
	inputs := []string{
		"(",
		")",
		"empty",
		"z",
		"ZM",
		"M",
		"POINT",
	}

	outputs := []Token{
		{ttype: LeftParen, lexeme: "("},
		{ttype: RightParen, lexeme: ")"},
		{ttype: Empty, lexeme: "empty"},
		{ttype: Z, lexeme: "z"},
		{ttype: ZM, lexeme: "zm"},
		{ttype: M, lexeme: "m"},
		{ttype: Point, lexeme: "point"},
	}

	for i, input := range inputs {
		l := NewLexer(strings.NewReader(input))
		ok, err := l.scan()
		if !ok {
			t.Error("scan reached unexpected eof")
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
	}
}
