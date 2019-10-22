package ftp

import (
	"strconv"
	"strings"
)

func GetTcpAddr(host string, port int) (addr string) {
	if strings.IndexByte(host, ':') >= 0 {
		addr = "[" + host + "]"
		if port > 0 {
			addr += ":" + strconv.Itoa(port)
		}
		return
	}
	if port > 0 {
		addr = host + ":" + strconv.Itoa(port)
	}
	return host
}
