package main

import (
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
	return 0
}

// END OMIT

func main() {
	m := testing.MainStart(all, []testing.InternalTest{{Name: "TestTable", F: TestTable}}, nil, nil)
	m.Run()
}

func all(p, s string) (bool, error) { return true, nil }
