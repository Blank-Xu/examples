package file

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var (
	host     = flag.String("h", "http://127.0.0.1:8080", "file server host")
	filename = flag.String("f", "1.wmv", "test filename")
)

func TestDownload(t *testing.T) {
	var (
		lfilename = filepath.Join(workDir, *filename)
		info      os.FileInfo

		err error
	)
	info, _ = os.Stat(lfilename)
	if len(info.Name()) > 0 {
		if err = os.Remove(lfilename); err != nil {
			t.Fatal(err)
		}
	}

	if err = Download(*host, *filename); err != nil {
		t.Fatal(err)
	}

	t.Log("download success")
}

func BenchmarkDownload(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var (
			lfilename = filepath.Join(workDir, *filename)
			info      os.FileInfo

			err error
		)
		info, _ = os.Stat(lfilename)
		if len(info.Name()) > 0 {
			if err = os.Remove(lfilename); err != nil {
				b.Fatal(err)
			}
		}

		b.StartTimer()
		if err = Download(*host, *filename); err != nil {
			b.Fatal(err)
		}
		b.StopTimer()

		b.Log("download success")
	}
}
