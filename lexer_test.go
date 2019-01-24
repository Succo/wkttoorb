package wkttoorb

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
		token, err := l.scanToken()
		if err != nil {
			t.Errorf("unexpected error %s", err)
		}
		if token.ttype != outputs[i].ttype {
			t.Errorf("incorrect ttype for %s", input)
		}
		if token.lexeme != outputs[i].lexeme {
			t.Errorf("incorrect lexeme for %s", input)
		}
		if token.pos != 0 {
			t.Errorf("incorrect position for %s", input)
		}
	}
}
