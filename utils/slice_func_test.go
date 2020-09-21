package utils

import (
	"testing"
)

var extraNumbers = 1

func BenchmarkContainsIntsSort(b *testing.B) {
	benchmarkContainsIntsHelper(ContainsIntsSort, b)

}
func BenchmarkContainsIntsRange(b *testing.B) {
	benchmarkContainsIntsHelper(ContainsIntsRange, b)
}

func BenchmarkContainsInts(b *testing.B) {
	benchmarkContainsIntsHelper(ContainsInts, b)
}

func benchmarkContainsIntsHelper(f func([]int, int) bool, b *testing.B) {
	max := b.N
	if max > 10000 {
		max = 10000
	}
	in := make([]int, b.N+extraNumbers)
	for i := 0; i < b.N+extraNumbers; i++ {
		in[i] = i
	}
	n := 100
	if len(in) < n {
		n = len(in)
	}
	b.ResetTimer()
	for _, i := range in[:n] {
		if got := f(in, i+len(in)); got {
			b.Errorf("want %v, got %v", false, got)
		}
	}
	for _, i := range in[:n] {
		if got := f(in, i); !got {
			b.Errorf("want %v, got %v", true, got)
		}
	}
}