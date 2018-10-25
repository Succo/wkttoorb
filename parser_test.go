package wktToOrb

import (
	"reflect"
	"testing"

	"github.com/paulmach/orb"
)

func Test_parsePoint(t *testing.T) {
	inputs := []string{
		"POINT ( 10.05  10.28 )",
		"POINT empty",
		"POINT Z ( 79.1 21.28 12.6 )",
		"POINT M ( 79.1 21.28 12.6 )",
		"POINT ZM ( 79.1 21.28 12.6 34.6 )",
		"POINT Z EmPty ",
		"POINT M EMPTY",
		"POINT ZM Empty",
	}
	outputs := []orb.Point{
		orb.Point{10.05, 10.28},
		orb.Point{0, 0},
		orb.Point{79.1, 21.28},
		orb.Point{79.1, 21.28},
		orb.Point{79.1, 21.28},
		orb.Point{0, 0},
		orb.Point{0, 0},
		orb.Point{0, 0},
	}

	for i, str := range inputs {
		geo, err := Scan(str)

		if err != nil {
			t.Errorf("unexpected error %s", err)
		}
		if !reflect.DeepEqual(geo, outputs[i]) {
			t.Error("incorrect value returned")
		}
	}
}
