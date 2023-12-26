package hashutil

import "testing"

func BenchmarkEncode(b *testing.B) {
	baseURL := "http://example.com"
	for i := 0; i < b.N; i++ {
		_ = Encode([]byte(baseURL))
	}
}
