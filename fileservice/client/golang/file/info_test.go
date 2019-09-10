package file

import (
	"testing"
)

func BenchmarkInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Info(*host, *filename, *username, *password, true); err != nil {
			b.Fatal(err)
		}
	}
}
