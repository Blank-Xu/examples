package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

func Md5(data string) string {
	return Md5Bytes([]byte(data))
}

func Md5Bytes(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

func Md5File(file *os.File) (string, error) {
	if file == nil {
		return "", errors.New("file is nil")
	}

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil && err != io.EOF {
		return "", fmt.Errorf("io copy failed, err: %v", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func Md5Filename(filename string) (string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return Md5File(file)
}
