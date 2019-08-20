package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

func Md5(data string) string {
	return Md5Bytes([]byte(data))
}

func Md5Bytes(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func Md5File(file *os.File) (string, error) {
	if file == nil {
		return "", errors.New("file is nil")
	}

	h := md5.New()

	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("md5 copy failed, err: %v", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
