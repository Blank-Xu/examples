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
	defer file1.Close()

	md5, err := Md5File(file1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("filename: %s, md5:%s", filename1, md5)
}

func TestMd5File2(t *testing.T) {
	file1, err := os.OpenFile(filename1, os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer file1.Close()

	file2, err := os.OpenFile(filename2, os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer file2.Close()

	md5, err := Md5File(file1)
	if err != nil {
		t.Fatal(err)
	}

	md52, err := Md5File(file2)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("filename1: %s, md5:%s", filename1, md5)
	t.Logf("filename2: %s, md5:%s", filename2, md52)
}

func TestMd5Filename(t *testing.T) {
	md5, err := Md5Filename(filename1)
	if err != nil {
		t.Error(err)
	}

	t.Logf("filename: %s, md5:%s", filename1, md5)
}

func BenchmarkMd5File(b *testing.B) {
	for i := 0; i < b.N; i++ {
		file1, err := os.OpenFile(filename1, os.O_RDONLY, 0666)
		if err != nil {
			b.Fatal(err)
		}

		md5, err := Md5File(file1)
		if err != nil {
			b.Fatal(err)
		}
		b.Log(md5)
	}
}

func BenchmarkMd5Filename(b *testing.B) {
	for i := 0; i < b.N; i++ {
		md5, err := Md5Filename(filename1)
		if err != nil {
			b.Fatal(err)
		}
		b.Log(md5)
	}
}
