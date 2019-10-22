package ftp

import (
	"sync"
)

type connections struct {
	conns sync.Map
}
