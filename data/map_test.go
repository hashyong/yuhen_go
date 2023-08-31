package data

import "testing"

func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test(0)
	}
}

func BenchmarkTestCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test(max)
	}
}

// go test -bench . -benchmem
/*
goos: linux
goarch: amd64
pkg: yuhen/data
cpu: Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz
BenchmarkTest-12            1761            687613 ns/op          687321 B/op        278 allocs/op
BenchmarkTestCap-12         3168            335979 ns/op          322282 B/op         12 allocs/op
*/

func BenchmarkTest1128(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test1([128]byte{1, 2, 3})
	}
}

func BenchmarkTest1129(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test1([129]byte{1, 2, 3})
	}
}

/*
BenchmarkTest1128-12              182816              6316 ns/op           21266 B/op          8 allocs/op
BenchmarkTest1129-12              137486              9134 ns/op           17287 B/op        103 allocs/op
*/
