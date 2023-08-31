package data

import (
	"testing"
)

func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !normalConv() {
			b.Fatal()
		}
	}
}

func BenchmarkUnsafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !unsafeConv() {
			b.Fatal()
		}
	}
}

/*
go test -bench . -benchmem -v
goos: linux
goarch: amd64
pkg: yuhen
cpu: Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz
BenchmarkNormal
BenchmarkNormal-12      10761168                94.05 ns/op          224 B/op          2 allocs/op
BenchmarkUnsafe
BenchmarkUnsafe-12      708710486                1.736 ns/op           0 B/op          0 allocs/op
PASS
ok      yuhen   2.538s
*/

func BenchmarkConcat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !concat() {
			b.Fatal()
		}
	}
}

func BenchmarkJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !join() {
			b.Fatal()
		}
	}
}

func BenchmarkBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !buffer() {
			b.Fatal()
		}
	}
}

/*
BenchmarkConcat
BenchmarkConcat-12          8161            145586 ns/op          530279 B/op        999 allocs/op
BenchmarkJoin
BenchmarkJoin-12          139759              7512 ns/op            1024 B/op          1 allocs/op
BenchmarkBuffer
BenchmarkBuffer-12        221882              4978 ns/op            2048 B/op          2 allocs/op
*/
