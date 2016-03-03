package main

import (
	"strconv"
	"testing"
)

func TestExp2Less1_0(t *testing.T) {
	y := Exp2Less1(0)
	if y != 1 {
		t.Errorf("Exp2Less1(0) == %d, want 1", y)
	}
}

func TestExp2Less1_1(t *testing.T) {
	y := Exp2Less1(1)
	if y != 1 {
		t.Errorf("Exp2Less1(1) == %d, want 1", y)
	}
}

func TestExp2Less1_2(t *testing.T) {
	y := Exp2Less1(2)
	if y != 2 {
		t.Errorf("Exp2Less1(2) == %d, want 2", y)
	}
}

func TestExp2Less1_3(t *testing.T) {
	y := Exp2Less1(3)
	if y != 4 {
		t.Errorf("Exp2Less1(3) == %d, want 4", y)
	}
}

func TestExp2Less1_10(t *testing.T) {
	y := Exp2Less1(10)
	if y != 512 {
		t.Errorf("Exp2Less1(10) == %d, want 512", y)
	}
}

// START OMIT
func Exp2Less1(x int) int {
	// NOT: (int)(math.Ceil(math.Exp2((float64)(x - 1))))
	if x < 2 {
		return 1
	}
	return 1 << (uint)(x-2) // FIXME: x-1 // HL
}

// END OMIT

func main() {
	var tests []testing.InternalTest
	for i, f := range []func(*testing.T){
		TestExp2Less1_0,
		TestExp2Less1_1,
		TestExp2Less1_2,
		TestExp2Less1_3,
		TestExp2Less1_10,
	} {
		tests = append(tests, testing.InternalTest{
			Name: "TestExp2Less1_" + strconv.Itoa(i),
			F:    f,
		})
	}
	m := testing.MainStart(all, tests, nil, nil)
	m.Run()
}

func all(p, s string) (bool, error) { return true, nil }
