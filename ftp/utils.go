package ftp

import (
	"bytes"
	"errors"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetAddress(host string, port int) (addr string) {
	var buf bytes.Buffer
	if strings.IndexByte(host, ':') > -1 {
		buf.WriteByte('[')
		buf.WriteString(host)
		buf.WriteByte(']')
	} else {
		buf.WriteString(host)
	}
	if port > 0 {
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(port))
	}
	return buf.String()
}

func GetRandomPort(minPort, maxPort int) (port int) {
	var n = maxPort - minPort
	if n == 0 {
		if minPort > 0 {
			port = minPort
		}
	} else if n > 0 {
		port = minPort + random.Intn(n+1)
	}
	return
}

func NewActiveTCPConn(host string, port int) (*net.TCPConn, error) {
	addr := GetAddress(host, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	return net.DialTCP("tcp", nil, tcpAddr)
}

func NewPassiveTCPListener(host string, minPort, maxPort int) (*net.TCPListener, error) {
	for i := 0; i < 100; i++ {
		listener, err := NewTcpListener(host, GetRandomPort(minPort, maxPort))
		if err == nil {
			return listener, nil
		}
	}
	return nil, errors.New("Unable to find an available port")
}

func NewTcpListener(host string, port int) (*net.TCPListener, error) {
	addr := GetAddress(host, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", tcpAddr)
}
