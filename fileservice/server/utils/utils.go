package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
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
	if file == nil {
		return ""
	}

	h := md5.New()

	if _, err := io.Copy(h, file); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}
