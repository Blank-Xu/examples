package utils

import (
	"os"
	"path/filepath"
)

func GetPath()(string,error)  {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}
