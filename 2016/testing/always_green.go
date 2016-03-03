package main

import (
	"os"
	"testing"
)

// START OMIT
func zeroBadly(x int) {}

func TestThatNeverFails(t *testing.T) {
	x := 123
	zeroBadly(x)
	if x == 0 { // typo! // HL
		t.Errorf("zeroBadly(x) == %d, want 0", x)
	}
}

// END OMIT

func main() {
	os.Args = append(os.Args, "-test.v")
	m := testing.MainStart(all, []testing.InternalTest{{Name: "TestThatNeverFails", F: TestThatNeverFails}}, nil, nil)
	m.Run()
}

func all(p, s string) (bool, error) { return true, nil }
