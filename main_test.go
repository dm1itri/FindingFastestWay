package main

import (
	"reflect"
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func EqualMatrix(t *testing.T, actual, expected [][]int) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func EqualSlice[T comparable](t *testing.T, actual, expected []T) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func TestCorrectPoint(t *testing.T) {
	matrix := [][]int{{1, 2}, {0, 4}}
	tests := []struct {
		name   string
		point  Point
		matrix [][]int
		want   bool
	}{
		{
			name:   "Negative coords",
			point:  Point{-1, -2},
			matrix: matrix,
			want:   false,
		},
		{
			name:   "Large coordinates",
			point:  Point{1, 52},
			matrix: matrix,
			want:   false,
		},
		{
			name:   "Сoordinates in the wall",
			point:  Point{0, 1},
			matrix: matrix,
			want:   false,
		},
		{
			name:   "correct coordinates",
			point:  Point{1, 1},
			matrix: matrix,
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Equal(t, CorrectPoint(tt.point, tt.matrix), tt.want)
		})
	}
}

func TestReadInput(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		matrix [][]int
		start  Point
		finish Point
		err    string
	}{
		{
			name:   "Length == 0",
			input:  "0 10\n",
			matrix: nil,
			start:  Point{},
			finish: Point{},
			err:    "incorrect dimensions of the maze",
		},
		{
			name:   "Maze element > 9",
			input:  "2 2\n33 1\n",
			matrix: nil,
			start:  Point{},
			finish: Point{},
			err:    "incorrect maze structure",
		},
		{
			name:   "start == finish",
			input:  "2 2\n2 1\n5 3\n0 0 0 0\n",
			matrix: nil,
			start:  Point{},
			finish: Point{},
			err:    "incorrect data about the beginning or end of the way",
		},
		{
			name:   "correct input",
			input:  "2 2\n2 1\n5 3\n0 0 1 0\n",
			matrix: [][]int{{2, 1}, {5, 3}},
			start:  Point{0, 0},
			finish: Point{0, 1},
			err:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matrix, start, finish, err := ReadInput(strings.NewReader(tt.input))
			EqualMatrix(t, matrix, tt.matrix)
			Equal(t, start, tt.start)
			Equal(t, finish, tt.finish)
			if err != nil {
				Equal(t, err.Error(), tt.err)
			}
		})
	}
}

func TestFindingFastestWay(t *testing.T) {
	tests := []struct {
		name      string
		matrix    [][]int
		start     Point
		finish    Point
		wantWay   []Point
		wantError string
	}{
		{
			name:      "Easy way",
			matrix:    [][]int{{5, 2}, {3, 3}},
			start:     Point{0, 0},
			finish:    Point{1, 1},
			wantWay:   []Point{{1, 1}, {1, 0}, {0, 0}},
			wantError: "",
		},
		{
			name:      "Way with wall",
			matrix:    [][]int{{1, 1, 9}, {2, 0, 1}, {1, 4, 1}, {1, 1, 1}},
			start:     Point{0, 0},
			finish:    Point{2, 1},
			wantWay:   []Point{{2, 1}, {2, 2}, {2, 3}, {1, 3}, {0, 3}, {0, 2}, {0, 1}, {0, 0}},
			wantError: "",
		},
		{
			name:      "One way with walls",
			matrix:    [][]int{{1, 0, 1, 1, 1}, {1, 0, 1, 0, 1}, {1, 0, 1, 0, 1}, {1, 0, 0, 0, 1}, {1, 1, 1, 1, 1}},
			start:     Point{0, 0},
			finish:    Point{2, 2},
			wantWay:   []Point{{2, 2}, {2, 1}, {2, 0}, {3, 0}, {4, 0}, {4, 1}, {4, 2}, {4, 3}, {4, 4}, {3, 4}, {2, 4}, {1, 4}, {0, 4}, {0, 3}, {0, 2}, {0, 1}, {0, 0}},
			wantError: "",
		},
		{
			name:      "Way not found",
			matrix:    [][]int{{1, 1, 1}, {1, 0, 1}, {0, 9, 0}},
			start:     Point{0, 0},
			finish:    Point{1, 2},
			wantWay:   nil,
			wantError: "the way was not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			way, err := FindingFastestWay(tt.matrix, tt.start, tt.finish)
			EqualSlice(t, way, tt.wantWay)
			if err != nil {
				Equal(t, err.Error(), tt.wantError)
			}
		})
	}
}
