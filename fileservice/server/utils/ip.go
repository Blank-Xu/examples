package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetIp(r *http.Request) (string, error) {
	var remoteAddr string

	if proxyIp := strings.Split(r.Header.Get("X-Forwarded-For"), ","); len(proxyIp) > 0 && len(proxyIp[0]) > 0 {
		remoteAddr = proxyIp[0]
	} else {
		remoteAddr = r.RemoteAddr
	}

	ip, _, err := net.SplitHostPort(remoteAddr)

	return ip, err
}
