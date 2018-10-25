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

func Test_parsePolygon(t *testing.T) {
	inputs := []string{
		"polygon empty",
		"polygon z empty",
		"polygon m empty",
		"polygon zm empty",
		"polygon (( 10 10, 10 20, 20 20, 20 15, 10 10))",
		"polygon z ((10 10 3, 10 20 3, 20 20 3, 20 15 4, 10 10 3))",
		"polygon m (( 10 10 8, 10 20 9, 20 20 9, 20 15 9, 10 10 8 ))",
		"polygon zm (( 10 10 3 8, 10 20 3 9, 20 20 3 9, 20 15 4 9, 10 10 3 8 ))",
		"polygon (( 10 10, 10 20, 20 20, 20 15, 10 10),( 10 10, 10 20, 20 20, 20 15, 10 10))",
	}
	outputs := []orb.Polygon{
		orb.Polygon{},
		orb.Polygon{},
		orb.Polygon{},
		orb.Polygon{},
		orb.Polygon{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10.0, 10.0}}},
		orb.Polygon{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10.0, 10.0}}},
		orb.Polygon{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10.0, 10.0}}},
		orb.Polygon{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10.0, 10.0}}},
		orb.Polygon{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10.0, 10.0}},
			{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10.0, 10.0}}},
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
