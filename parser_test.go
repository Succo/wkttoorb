package wktToOrb

import (
	"reflect"
	"testing"

	"github.com/paulmach/orb"
)

func Test_parsePoint(t *testing.T) {
	geo, err := Scan("POINT ( 10.05  10.28 )")

	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if !reflect.DeepEqual(geo, orb.Point{10.05, 10.28}) {
		t.Error("incorrect value returned")
	}

	geo, err = Scan("POINT empty")

	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if !reflect.DeepEqual(geo, orb.Point{0, 0}) {
		t.Error("incorrect value returned")
	}
}
