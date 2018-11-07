package wktToOrb

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
	l.pos += len(lexeme)
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
	for isFloatRune(r) {
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
		return l.scanToken()
	case r == '(':
		l.addToken(LeftParen, "(")
	case r == ')':
		l.addToken(RightParen, ")")
	case r == ',':
		l.addToken(Comma, ",")
	case unicode.IsLetter(r):
		w := l.scanToLowerWord(r)
		switch w {
		case "empty":
			l.addToken(Empty, "empty")
		case "z":
			l.addToken(Z, "z")
		case "m":
			l.addToken(M, "m")
		case "zm":
			l.addToken(ZM, "zm")
		case "point":
			l.addToken(Point, "point")
		case "linestring":
			l.addToken(Linestring, "linestring")
		case "polygon":
			l.addToken(Polygon, "polygon")
		case "multipoint":
			l.addToken(Multipoint, "multipoint")
		case "multilinestring":
			l.addToken(MultilineString, "multilinestring")
		case "multipolygon":
			l.addToken(MultiPolygon, "multipolygon")
		default:
			return false, fmt.Errorf("Unexpected word %s on character %d", w, l.pos)
		}
	case beginFloat(r):
		w := l.scanFloat(r)
		l.addToken(Float, w)
	case r == eof:
		l.addToken(Eof, "")
		return false, nil
	default:
		return false, fmt.Errorf("Unexpected rune %s on character %d", string(r), l.pos)
	}
	return true, nil
}

func beginFloat(r rune) bool {
	return r == '-' || r == '.' || unicode.IsNumber(r)
}

func isFloatRune(r rune) bool {
	return beginFloat(r) || r == 'e'
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
