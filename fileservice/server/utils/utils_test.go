package utils

import (
	"os"
	"testing"
)

const filename1 = "1.wmv"
const filename2 = "2.wmv"

func TestMd5File(t *testing.T) {
	file1, err := os.OpenFile(filename1, os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("filename: %s, md5:%s", filename1, Md5File(file1))
}

func TestMd5File2(t *testing.T) {
	file1, err := os.OpenFile(filename1, os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}

	file2, err := os.OpenFile(filename2, os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("filename1: %s, md5:%s", filename1, Md5File(file1))
	t.Logf("filename2: %s, md5:%s", filename2, Md5File(file2))
}

func BenchmarkMd5File(b *testing.B) {
	for i := 0; i < b.N; i++ {
		file1, err := os.OpenFile(filename1, os.O_RDONLY, 0666)
		if err != nil {
			b.Fatal(err)
		}
		b.Log(Md5File(file1))
	}
}
