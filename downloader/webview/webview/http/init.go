package http

import (
	"log"
	"net"
	"net/http"
)

func Init() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		log.Fatal(http.Serve(ln, router()))
	}()
	return "http://" + ln.Addr().String()
}
