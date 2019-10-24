package ftp

import (
	"testing"
)

func TestGetAddress(t *testing.T) {

}

func TestGetRandomPort(t *testing.T) {
	var ports = []struct {
		minPort int
		maxPort int
	}{
		{0, 0},
		{0, 1},
		{1, 1},
		{1, 2},
		{2, 1},
	}

	for _, port := range ports {
		if GetRandomPort(port.minPort, port.maxPort) > port.maxPort {
			t.Errorf("port %v check failed", port)
		}
	}
}
