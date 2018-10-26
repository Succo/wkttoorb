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

func (p *Parser) Parse() (orb.Geometry, error) {
	t := p.pop()
	switch t.ttype {
	case Point:
		return p.parsePoint()
	case Linestring:
		return p.parseLineString()
	case Polygon:
		return p.parsePolygon()
	case Multipoint:
		line, err := p.parseLineString()
		return orb.MultiPoint(line), err
	case MultilineString:
		poly, err := p.parsePolygon()
		multiline := make(orb.MultiLineString, 0, len(poly))
		for _, ring := range poly {
			multiline = append(multiline, orb.LineString(ring))
		}
		return multiline, err
	case MultiPolygon:
		return p.parseMultiPolygon()
	default:
		return nil, fmt.Errorf("unexpected token %s on pos %d", t.lexeme, t.pos)
	}
}

func (p *Parser) parsePoint() (point orb.Point, err error) {
	t := p.pop()
	switch t.ttype {
	case Empty:
		point = orb.Point{0, 0}
	case Z, M, ZM:
		t1 := p.pop()
		if t1.ttype == Empty {
			point = orb.Point{0, 0}
			break
		}
		if t1.ttype != LeftParen {
			return point, fmt.Errorf("parse point unexpected token on pos %d", t.pos)
		}
		fallthrough
	case LeftParen:
		switch t.ttype { // reswitch on the type because of the fallthrough
		case Z:
			point, err = p.parseZCoord()
		case M:
			point, err = p.parseMCoord()
		case ZM:
			point, err = p.parseZMCoord()
		default:
			point, err = p.parseCoord()
		}
		if err != nil {
			return point, err
		}

		t = p.pop()
		if t.ttype != RightParen {
			return point, fmt.Errorf("parse point unexpected token on pos %d", t.pos)
		}
	default:
		return point, fmt.Errorf("parse point unexpected token on pos %d", t.pos)
	}

	t = p.pop()
	if t.ttype != Eof {
		return point, fmt.Errorf("parse point unexpected token on pos %d", t.pos)
	}

	return point, nil
}

func (p *Parser) parseLineString() (line orb.LineString, err error) {
	line = make([]orb.Point, 0)
	t := p.pop()
	switch t.ttype {
	case Empty:
	case Z, M, ZM:
		t1 := p.pop()
		if t1.ttype == Empty {
			break
		}
		if t1.ttype != LeftParen {
			return line, fmt.Errorf("unexpected token %s on pos %d expected '('", t.lexeme, t.pos)
		}
		fallthrough
	case LeftParen:
		line, err = p.parseLineStringText(t.ttype)
		if err != nil {
			return line, err
		}
	default:
		return line, fmt.Errorf("unexpected token %s on pos %d", t.lexeme, t.pos)
	}

	t = p.pop()
	if t.ttype != Eof {
		return line, fmt.Errorf("unexpected token %s on pos %d, expected Eof", t.lexeme, t.pos)
	}

	return line, nil
}

func (p *Parser) parseLineStringText(ttype tokenType) (line orb.LineString, err error) {
	line = make([]orb.Point, 0)
	for {
		var point orb.Point
		switch ttype {
		case Z:
			point, err = p.parseZCoord()
		case M:
			point, err = p.parseMCoord()
		case ZM:
			point, err = p.parseZMCoord()
		default:
			point, err = p.parseCoord()
		}
		if err != nil {
			return line, err
		}
		line = append(line, point)
		t := p.pop()
		if t.ttype == RightParen {
			break
		} else if t.ttype != Comma {
			return line, fmt.Errorf("unexpected token %s on pos %d expected ','", t.lexeme, t.pos)
		}
	}
	return line, nil
}

func (p *Parser) parsePolygon() (poly orb.Polygon, err error) {
	poly = make([]orb.Ring, 0)
	t := p.pop()
	switch t.ttype {
	case Empty:
	case Z, M, ZM:
		t1 := p.pop()
		if t1.ttype == Empty {
			break
		}
		if t1.ttype != LeftParen {
			return poly, fmt.Errorf("unexpected token %s on pos %d expected '('", t.lexeme, t.pos)
		}
		fallthrough
	case LeftParen:
		poly, err = p.parsePolygonText(t.ttype)
		if err != nil {
			return poly, err
		}
	default:
		return poly, fmt.Errorf("unexpected token %s on pos %d", t.lexeme, t.pos)
	}

	t = p.pop()
	if t.ttype != Eof {
		return poly, fmt.Errorf("unexpected token %s on pos %d, expected Eof", t.lexeme, t.pos)
	}

	return poly, nil
}

func (p *Parser) parsePolygonText(ttype tokenType) (poly orb.Polygon, err error) {
	poly = make([]orb.Ring, 0)
	for {
		var line orb.LineString
		t := p.pop()
		if t.ttype != LeftParen {
			return poly, fmt.Errorf("unexpected token %s on pos %d expected '('", t.lexeme, t.pos)
		}
		line, err = p.parseLineStringText(ttype)
		if err != nil {
			return poly, err
		}
		poly = append(poly, orb.Ring(line))
		t = p.pop()
		if t.ttype == RightParen {
			break
		} else if t.ttype != Comma {
			return poly, fmt.Errorf("unexpected token %s on pos %d expected ','", t.lexeme, t.pos)
		}
	}
	return poly, nil
}

func (p *Parser) parseMultiPolygon() (multi orb.MultiPolygon, err error) {
	multi = make([]orb.Polygon, 0)
	t := p.pop()
	switch t.ttype {
	case Empty:
	case Z, M, ZM:
		t1 := p.pop()
		if t1.ttype == Empty {
			break
		}
		if t1.ttype != LeftParen {
			return multi, fmt.Errorf("unexpected token %s on pos %d expected '('", t.lexeme, t.pos)
		}
		fallthrough
	case LeftParen:
		multi, err = p.parseMultiPolygonText(t.ttype)
		if err != nil {
			return multi, err
		}
	default:
		return multi, fmt.Errorf("unexpected token %s on pos %d", t.lexeme, t.pos)
	}

	t = p.pop()
	if t.ttype != Eof {
		return multi, fmt.Errorf("unexpected token %s on pos %d, expected Eof", t.lexeme, t.pos)
	}

	return multi, nil
}

func (p *Parser) parseMultiPolygonText(ttype tokenType) (multi orb.MultiPolygon, err error) {
	multi = make([]orb.Polygon, 0)
	for {
		var poly orb.Polygon
		t := p.pop()
		if t.ttype != LeftParen {
			return multi, fmt.Errorf("unexpected token %s on pos %d expected '('", t.lexeme, t.pos)
		}
		poly, err = p.parsePolygonText(ttype)
		if err != nil {
			return multi, err
		}
		multi = append(multi, poly)
		t = p.pop()
		if t.ttype == RightParen {
			break
		} else if t.ttype != Comma {
			return multi, fmt.Errorf("unexpected token %s on pos %d expected ','", t.lexeme, t.pos)
		}
	}
	return multi, nil
}

func (p *Parser) parseCoord() (point orb.Point, err error) {
	t1 := p.pop()
	if t1.ttype != Float {
		return point, fmt.Errorf("parse coordinates unexpected token %s on pos %d", t1.lexeme, t1.pos)
	}
	t2 := p.pop()
	if t2.ttype != Float {
		return point, fmt.Errorf("parse coordinates unexpected token %s on pos %d", t1.lexeme, t2.pos)
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
	t := p.pop()
	if t.ttype != Float {
		return point, fmt.Errorf("parseZCoord unexpected token %s on pos %d expected Float", t.lexeme, t.pos)
	}

	return point, nil
}

func (p *Parser) parseMCoord() (point orb.Point, err error) {
	point, err = p.parseCoord()
	if err != nil {
		return point, err
	}

	// drop the last value M values are not really supported
	t := p.pop()
	if t.ttype != Float {
		return point, fmt.Errorf("parseZCoord unexpected token %s on pos %d expected Float", t.lexeme, t.pos)
	}

	return point, nil
}

func (p *Parser) parseZMCoord() (point orb.Point, err error) {
	point, err = p.parseCoord()
	if err != nil {
		return point, err
	}

	// drop the last value M values
	// and Z coordinates are not really supported
	for i := 0; i < 2; i++ {
		t := p.pop()
		if t.ttype != Float {
			return point, fmt.Errorf("parseZCoord unexpected token %s on pos %d expected Float", t.lexeme, t.pos)
		}
	}

	return point, nil
}
