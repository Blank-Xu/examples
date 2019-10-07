package utils

import (
	"os"
)

func OpenFile(filename string) (*os.File, error) {
	return os.Open(filename)
}
