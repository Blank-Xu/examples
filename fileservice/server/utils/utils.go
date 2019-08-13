package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math"
	"os"
)

const (
	fileChunk = 8 * 1024 //  8KB
)

func Md5(data string) string {
	return Md5Bytes([]byte(data))
}

func Md5Bytes(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func Md5File(file *os.File) string {
	// calculate the file size
	info, _ := file.Stat()

	filesize := info.Size()

	blocks := uint64(math.Ceil(float64(filesize) / float64(fileChunk)))

	h := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(fileChunk, float64(filesize-int64(i*fileChunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(h, string(buf)) // append into the hash
	}

	return hex.EncodeToString(h.Sum(nil))
}
