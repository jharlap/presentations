package main

import "fmt"

// START OMIT
// BenchmarkConcat0-4	 30000000     45.5 ns/op	       0 B/op	       0 allocs/op
func Concat0(a, b string) int {
	s := a + b
	return len(s)
}

// BenchmarkConcat1-4	 20000000     69.6 ns/op	      16 B/op	       1 allocs/op
func Concat1(a, b string) int {
	a += b // HLplusequal
	return len(a)
}

// BenchmarkConcat2-4	  5000000      358 ns/op	      48 B/op	       3 allocs/op
func Concat2(a, b string) int {
	return len(fmt.Sprintf("%s%s", a, b)) // HLsprintf
}

// BenchmarkConcat3-4	500000000     3.10 ns/op	       0 B/op	       0 allocs/op
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
