package data

import "testing"

func BenchmarkNoraml(b *testing.B) {
	for i := 0; i < b.N; i++ {
		normal()
	}
}

func BenchmarkBlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		block()
	}
}
