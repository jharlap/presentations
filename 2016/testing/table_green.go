package main

import (
	"os"
	"testing"
)

func TestTable(t *testing.T) {
	// START TC OMIT
	testCases := []struct {
		x, y int
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 4},
		{10, 512},
	}

	for _, tc := range testCases {
		y := Exp2Less1(tc.x)
		if y != tc.y {
			t.Errorf("fn(%d) == %d, want %d", tc.x, y, tc.y)
		}
	}
	// END TC OMIT
}

// START OMIT
func Exp2Less1(x int) int {
	if x < 2 {
		return 1
	}

	e := 1
	for i := x - 1; i > -1; i-- { // FIXME: -1?
		e = e * 2
	}
	return e
}

// END OMIT

func main() {
	os.Args = append(os.Args, "-test.v")
	m := testing.MainStart(all, []testing.InternalTest{{Name: "TestTable", F: TestTable}}, nil, nil)
	m.Run()
}

func all(p, s string) (bool, error) { return true, nil }
