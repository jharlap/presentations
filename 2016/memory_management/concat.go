package main

import "fmt"

// START OMIT
func Concat0(a, b string) int {
	s := a + b
	return len(s)
}

func Concat1(a, b string) int {
	a += b
	return len(a)
}

func Concat2(a, b string) int {
	return len(fmt.Sprintf("%s%s", a, b))
}

func Concat3(a, b string) int {
	return len(a) + len(b)
}

// END OMIT

var l int

func main() {
	var h, w = "hello", "world"
	l = Concat0(h, w)
	l = Concat1(h, w)
	l = Concat2(h, w)
	l = Concat3(h, w)
}
