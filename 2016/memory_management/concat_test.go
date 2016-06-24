package main

import "testing"

type concatter func(string, string) int

var h, w = "hello", "world"

func benchmark(b *testing.B, c concatter) {
	b.ReportAllocs()
	var r int
	for i := 0; i < b.N; i++ {
		r = c(h, w)
	}
	l = r
}
func BenchmarkConcat0(b *testing.B) { benchmark(b, Concat0) }
func BenchmarkConcat1(b *testing.B) { benchmark(b, Concat1) }
func BenchmarkConcat2(b *testing.B) { benchmark(b, Concat2) }
func BenchmarkConcat3(b *testing.B) { benchmark(b, Concat3) }
