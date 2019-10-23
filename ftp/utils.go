package ftp

import (
	"strconv"
	"strings"
)

func GetTcpAddr(host string, port uint32) (addr string) {
	if strings.IndexByte(host, ':') >= 0 {
		addr = "[" + host + "]"
		if port > 0 {
			addr += ":" + strconv.Itoa(int(port))
		}
		return
	}
	if port > 0 {
		addr = host + ":" + strconv.Itoa(int(port))
	}
	return host
}
