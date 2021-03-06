package utils

import (
	"os"
)

func MkdirAll(dir string) (err error) {
	if err = os.MkdirAll(dir, 0766); err != nil {
		return
	}

	return os.Chmod(dir, 0766)
}
