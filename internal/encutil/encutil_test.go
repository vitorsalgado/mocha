package encutil

import "testing"

func BenchmarkJoin(b *testing.B) {
	b.ReportAllocs()

	cases := [][]string{
		nil,
		{"one"},
		{"a", "b", "c"},
		{"dev", "test", "ops"},
		{"hello", "world", "http", "mock", "library"},
		{"hello", "world", "http", "mock", "library", "hello", "world", "http", "mock", "library"},
	}

	for i := 0; i < b.N; i++ {
		for _, c := range cases {
			Join(",", c...)
		}
	}
}
