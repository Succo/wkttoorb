package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

type tokenType int

const (
	// Separator
	LeftParen tokenType = iota
	RightParen
	Comma

	// Keyword
	Empty
	Z
	M
	ZM

	// Geometry type
	Point
	Linestring
	Polygon
	Multipoint
	MultilineString
	MultiPolygon

	// Values
	Float

	// Eof
	Eof
)

// eof is used to simplify treatment of file end
const eof = rune(0)

type Token struct {
	ttype  tokenType
	lexeme string
	pos    int
}

type Lexer struct {
	reader  *bufio.Reader
	scanned []Token

	pos int
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader:  bufio.NewReader(reader),
		scanned: make([]Token, 0),
	}
}

// addToken add a parsed token to the token list
func (l *Lexer) addToken(ttype tokenType, lexeme string) {
	t := Token{ttype, lexeme, l.pos}
	l.scanned = append(l.scanned, t)
}

func (l *Lexer) read() rune {
	ch, _, err := l.reader.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (l *Lexer) unread() {
	_ = l.reader.UnreadRune()
}

func (l *Lexer) peek() rune {
	ch, _, err := l.reader.ReadRune()
	if err != nil {
		ch = eof
	}
	l.reader.UnreadRune()
	return ch
}

// scanToLowerWord scan a word and returns its value in lower letters
func (l *Lexer) scanToLowerWord(r rune) string {
	var buf bytes.Buffer
	buf.WriteRune(unicode.ToLower(r))
	r = l.read()
	for unicode.IsLetter(r) {
		buf.WriteRune(unicode.ToLower(r))
		r = l.read()
	}
	l.unread()
	return buf.String()
}

// scanFloat scan a string representing a float
func (l *Lexer) scanFloat(r rune) string {
	var buf bytes.Buffer
	buf.WriteRune(r)
	r = l.read()
	for unicode.IsDigit(r) || r == '.' {
		buf.WriteRune(r)
		r = l.read()
	}
	l.unread()
	return buf.String()
}

// scanToken scans the next lexeme
// return false is eof is reached true otherwise
// error is non nil only in case of unexpected character or word
func (l *Lexer) scanToken() (bool, error) {
	r := l.read()
	switch {
	case unicode.IsSpace(r):
		l.pos++
	case r == '(':
		l.pos++
		l.addToken(LeftParen, "(")
	case r == ')':
		l.pos++
		l.addToken(RightParen, ")")
	case r == ',':
		l.pos++
		l.addToken(Comma, ",")
	case unicode.IsLetter(r):
		w := l.scanToLowerWord(r)
		switch w {
		case "empty":
			l.pos = l.pos + 5
			l.addToken(Empty, "empty")
		case "z":
			l.pos++
			l.addToken(Z, "z")
		case "m":
			l.pos++
			l.addToken(M, "m")
		case "zm":
			l.pos = l.pos + 2
			l.addToken(ZM, "zm")
		case "point":
			l.pos = l.pos + 5
			l.addToken(Point, "point")
		case "linestring":
			l.pos = l.pos + 10
			l.addToken(Linestring, "linestring")
		case "polygon":
			l.pos = l.pos + 7
			l.addToken(Polygon, "polygon")
		case "multipoint":
			l.pos = l.pos + 10
			l.addToken(Multipoint, "multipoint")
		case "multilinestring":
			l.pos = l.pos + 15
			l.addToken(MultilineString, "multilinestring")
		case "multipolygon":
			l.pos = l.pos + 12
			l.addToken(MultiPolygon, "multipolygon")
		default:
			return false, fmt.Errorf("Unexpected word %s on character %d", w, l.pos)
		}
	case unicode.IsNumber(r):
		w := l.scanFloat(r)
		l.pos = l.pos + len(w)
		l.addToken(Float, w)
	case r == eof:
		l.addToken(Eof, "")
		return false, nil
	default:
		return false, fmt.Errorf("Unexpected rune %s on character %d", string(r), l.pos)
	}
	return true, nil
}

func (l *Lexer) Scan() error {
	for {
		ok, err := l.scanToken()
		switch {
		case ok:
		case err != nil:
			return err
		default:
			return nil
		}
	}
}
