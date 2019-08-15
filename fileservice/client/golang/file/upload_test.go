package file

import (
	"testing"
)

func TestUpload(t *testing.T) {
	if err := Upload(*host, *filename); err != nil {
		t.Fatal(err)
	}

	t.Log("upload success")
}

func BenchmarkUpload(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := Upload(*host, *filename); err != nil {
			b.Fatal(err)
		}
	}
}
