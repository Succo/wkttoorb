package wktToOrb

import (
	"strings"

	"github.com/paulmach/orb"
	"github.com/pkg/errors"
)

func Scan(s string) (orb.Geometry, error) {
	p := Parser{NewLexer(strings.NewReader(s)), 0}

	err := p.Scan()
	if err != nil {
		return nil, errors.Wrap(err, "p.Scan error")
	}

	geo, err := p.Parse()
	if err != nil {
		return nil, errors.Wrap(err, "p.Parse error")
	}

	return geo, nil
}