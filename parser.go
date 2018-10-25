package wktToOrb

import (
	"fmt"
	"strconv"

	"github.com/paulmach/orb"
	"github.com/pkg/errors"
)

type Parser struct {
	*Lexer
	tPos int
}

func (p *Parser) pop() Token {
	t := p.scanned[p.tPos]
	p.tPos++
	return t
}

func (p *Parser) rewind() {
	p.tPos--
}

func (p *Parser) Parse() (orb.Geometry, error) {
	t := p.pop()
	if t.ttype < Point || t.ttype > MultiPolygon {
		return nil, fmt.Errorf("unexpected token on pos %d", t.pos)
	}

	switch t.ttype {
	case Point:
		return p.parsePoint()
	//case Linestring:
	//	return p.parseLineString()
	//case Polygon:
	//	return p.parsePolygon()
	//case Multipoint:
	//	return p.parseMultipoint()
	//case MultilineString:
	//	return p.parseMultilinestring()
	//case MultiPolygon:
	//	return p.parseMultiPolygon()
	default:
		return nil, fmt.Errorf("unexpected token %s on pos %d", t.lexeme, t.pos)
	}
}

func (p *Parser) parsePoint() (point orb.Point, err error) {
	t := p.pop()
	switch t.ttype {
	case Empty:
		point = orb.Point{0, 0}
		t = p.pop()
		if t.ttype != Eof {
			return point, fmt.Errorf("unexpected token on pos %d", t.pos)
		}
		return point, nil
	case Z:
		t := p.pop()
		if t.ttype != LeftParen {
			return point, fmt.Errorf("unexpected token on pos %d", t.pos)
		}
		point, err = p.parseZCoord()
	case M:
		t := p.pop()
		if t.ttype != LeftParen {
			return point, fmt.Errorf("unexpected token on pos %d", t.pos)
		}
		point, err = p.parseMCoord()
	case ZM:
		t := p.pop()
		if t.ttype != LeftParen {
			return point, fmt.Errorf("unexpected token on pos %d", t.pos)
		}
		point, err = p.parseZMCoord()
	case LeftParen:
		point, err = p.parseCoord()
	default:
		return point, fmt.Errorf("unexpected token on pos %d", t.pos)
	}

	if err != nil {
		return point, err
	}

	t = p.pop()
	if t.ttype != RightParen {
		return point, fmt.Errorf("unexpected token on pos %d", t.pos)
	}
	t = p.pop()
	if t.ttype != Eof {
		return point, fmt.Errorf("unexpected token on pos %d", t.pos)
	}

	return point, nil
}

//func (p *Parser) parseLineString() (line orb.LineString, err error) {
//	line = make([]orb.Point, 0)
//	t := p.pop()
//	switch t.ttype {
//	case Empty:
//	case Z:
//		point, err = p.parseZCoord()
//	case M:
//		point, err = p.parseMCoord()
//	case ZM:
//		point, err = p.parseZMCoord()
//	case Float:
//		p.rewind()
//		point, err = p.parseCoord(t)
//	}
//
//	if err != nil {
//		return line, err
//	}
//
//	t := p.pop()
//	if t.ttype != Eof {
//		return line, fmt.Errorf("unexpected token on pos %d", t.pos)
//	}
//
//	return line, nil
//}

func (p *Parser) parseCoord() (point orb.Point, err error) {
	t1 := p.pop()
	if t1.ttype != Float {
		return point, fmt.Errorf("unexpected token on pos %d", t1.pos)
	}
	t2 := p.pop()
	if t2.ttype != Float {
		return point, fmt.Errorf("unexpected token on pos %d", t2.pos)
	}

	c1, err := strconv.ParseFloat(t1.lexeme, 64)
	if err != nil {
		return point, errors.Wrap(err, fmt.Sprintf("invalid lexeme for token on pos %d", t1.pos))
	}
	c2, err := strconv.ParseFloat(t2.lexeme, 64)
	if err != nil {
		return point, errors.Wrap(err, fmt.Sprintf("invalid lexeme for token on pos %d", t2.pos))
	}

	return orb.Point{c1, c2}, nil
}

func (p *Parser) parseZCoord() (point orb.Point, err error) {
	point, err = p.parseCoord()
	if err != nil {
		return point, err
	}

	// drop the last value Z coordinates are not really supported
	p.pop()

	return point, nil
}

func (p *Parser) parseMCoord() (point orb.Point, err error) {
	point, err = p.parseCoord()
	if err != nil {
		return point, err
	}

	// drop the last value M values are not really supported
	p.pop()

	return point, nil
}

func (p *Parser) parseZMCoord() (point orb.Point, err error) {
	point, err = p.parseCoord()
	if err != nil {
		return point, err
	}

	// drop the last value M values
	// and Z coordinates are not really supported
	p.pop()
	p.pop()

	return point, nil
}
