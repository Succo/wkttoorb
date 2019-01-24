package wkttoorb

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
		"POINT ZM ( 5e-05 21.28 12.6 34.6 )",
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
		orb.Point{5e-05, 21.28},
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

func Test_parseMultipoint(t *testing.T) {
	inputs := []string{
		"MULTIPOINT EMPTY",
		"MULTIPOINT M EMPTY",
		"MULTIPOINT Z EMPTY",
		"MULTIPOINT ZM EMPTY",
		" multipoint ( 10.05  10.28 , 20.95  20.89 )",
		"multipoint z ( 10.05 10.28 3.09, 20.95 31.98 4.72, 21.98 29.80 3.51 )",
		"multipoint m ( 10.05 10.28 5.84, 20.95 31.98 9.01, 21.98 29.80 12.84 )",
		"multipoint zm (10.05 10.28 3.09 5.84, 20.95 31.98 4.72 9.01, 21.98 29.80 3.51 12.84)",
	}
	outputs := []orb.MultiPoint{
		orb.MultiPoint{},
		orb.MultiPoint{},
		orb.MultiPoint{},
		orb.MultiPoint{},
		orb.MultiPoint{{10.05, 10.28}, {20.95, 20.89}},
		orb.MultiPoint{{10.05, 10.28}, {20.95, 31.98}, {21.98, 29.80}},
		orb.MultiPoint{{10.05, 10.28}, {20.95, 31.98}, {21.98, 29.80}},
		orb.MultiPoint{{10.05, 10.28}, {20.95, 31.98}, {21.98, 29.80}},
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

func Test_parseMultiLineString(t *testing.T) {
	inputs := []string{
		"multilinestring empty",
		"multilinestring z empty",
		"multilinestring m empty",
		"multilinestring zm empty",
		"multilinestring (( 10 10, 10 20, 20 20, 20 15, 10 10))",
		"multilinestring z ((10 10 3, 10 20 3, 20 20 3, 20 15 4, 10 10 3))",
		"multilinestring m (( 10 10 8, 10 20 9, 20 20 9, 20 15 9, 10 10 8 ))",
		"multilinestring zm (( 10 10 3 8, 10 20 3 9, 20 20 3 9, 20 15 4 9, 10 10 3 8 ))",
		"multilinestring (( 10 10, 10 20, 20 20, 20 15, 10 10),( 10 10, 10 20, 20 20, 20 15, 10 10))",
	}
	outputs := []orb.MultiLineString{
		orb.MultiLineString{},
		orb.MultiLineString{},
		orb.MultiLineString{},
		orb.MultiLineString{},
		orb.MultiLineString{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10.0, 10.0}}},
		orb.MultiLineString{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10.0, 10.0}}},
		orb.MultiLineString{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10.0, 10.0}}},
		orb.MultiLineString{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10.0, 10.0}}},
		orb.MultiLineString{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10.0, 10.0}},
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

func Test_parseMultiPolygon(t *testing.T) {
	inputs := []string{
		"multipolygon empty",
		"multipolygon z empty",
		"multipolygon m empty",
		"multipolygon zm empty",
		"multipolygon (((10 10, 10 20, 20 20, 20 15 , 10 10), (50 40, 50 50, 60 50, 60 40, 50 40)))",
		"multipolygon z (((10 10 7, 10 20 8, 20 20 7, 20 15 5, 10 10 7), (50 40 6, 50 50 6, 60 50 5, 60 40 6, 50 40 7)))",
		"multipolygon m (((10 10 2, 10 20 3, 20 20 4, 20 15 5, 10 10 2), (50 40 7, 50 50 3, 60 50 4, 60 40 5, 50 40 7)))",
		"multipolygon zm (((10 10 7 2, 10 20 8 3, 20 20 7 4, 20 15 5 5, 10 10 7 2), (50 40 6 7, 50 50 6 3, 60 50 5 4, 60 40 6 5, 50 40 7 7)))",
	}
	outputs := []orb.MultiPolygon{
		orb.MultiPolygon{},
		orb.MultiPolygon{},
		orb.MultiPolygon{},
		orb.MultiPolygon{},
		orb.MultiPolygon{{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10, 10}},
			{{50.0, 40.0}, {50.0, 50.0}, {60.0, 50.0}, {60.0, 40.0}, {50.0, 40.0}}}},
		orb.MultiPolygon{{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10, 10}},
			{{50.0, 40.0}, {50.0, 50.0}, {60.0, 50.0}, {60.0, 40.0}, {50.0, 40.0}}}},
		orb.MultiPolygon{{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10, 10}},
			{{50.0, 40.0}, {50.0, 50.0}, {60.0, 50.0}, {60.0, 40.0}, {50.0, 40.0}}}},
		orb.MultiPolygon{{{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 15.0}, {10, 10}},
			{{50.0, 40.0}, {50.0, 50.0}, {60.0, 50.0}, {60.0, 40.0}, {50.0, 40.0}}}},
	}

	for i, str := range inputs {
		geo, err := Scan(str)

		if err != nil {
			t.Errorf("unexpected error %s on test %d", err, i)
		}
		if !reflect.DeepEqual(geo, outputs[i]) {
			t.Errorf("incorrect value returned on test %d", i)
			fmt.Println(geo)
			fmt.Println(outputs[i])
		}
	}
}
