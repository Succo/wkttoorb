package wktToOrb

import (
	"fmt"
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

func Test_parseLine(t *testing.T) {
	inputs := []string{
		"LINESTRING EMPTY",
		"LINESTRING M EMPTY",
		"LINESTRING Z EMPTY",
		"LINESTRING ZM EMPTY",
		" linestring ( 10.05  10.28 , 20.95  20.89 )",
		"linestring z ( 10.05 10.28 3.09, 20.95 31.98 4.72, 21.98 29.80 3.51 )",
		"linestring m ( 10.05 10.28 5.84, 20.95 31.98 9.01, 21.98 29.80 12.84 )",
		"linestring zm (10.05 10.28 3.09 5.84, 20.95 31.98 4.72 9.01, 21.98 29.80 3.51 12.84)",
	}
	outputs := []orb.LineString{
		orb.LineString{},
		orb.LineString{},
		orb.LineString{},
		orb.LineString{},
		orb.LineString{{10.05, 10.28}, {20.95, 20.89}},
		orb.LineString{{10.05, 10.28}, {20.95, 31.98}, {21.98, 29.80}},
		orb.LineString{{10.05, 10.28}, {20.95, 31.98}, {21.98, 29.80}},
		orb.LineString{{10.05, 10.28}, {20.95, 31.98}, {21.98, 29.80}},
	}

	for i, str := range inputs {
		geo, err := Scan(str)

		if err != nil {
			t.Errorf("unexpected error %s on test %d", err, i)
		}
		if !reflect.DeepEqual(geo, outputs[i]) {
			t.Errorf("incorrect value returned on test %d", i)
			fmt.Println(geo)
		}
	}
}
