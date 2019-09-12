package httpserver

import (
	"fmt"
	"net/http"
	"time"
)

var Default *HttpServer

type HttpServer struct {
	BindAddr            string `json:"bind_addr" yaml:"bind_addr"`
	Port                int    `json:"port" yaml:"port"`
	Https               bool   `json:"https" yaml:"https"`
	CertFile            string `json:"cert_file" yaml:"cert_file"`
	KeyFile             string `json:"key_file" yaml:"key_file"`
	ReadTimeout         int    `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout        int    `json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout         int    `json:"idle_timeout" yaml:"idle_timeout"`
	MaxConnsPerHost     int    `json:"max_conns_per_host" yaml:"max_conns_per_host"`           // 每一个host对应的最大连接数
	MaxIdleConns        int    `json:"max_idle_conns" yaml:"max_idle_conns"`                   // 所有host对应的idle状态最大的连接总数
	MaxIdleConnsPerHost int    `json:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host"` // 每一个host对应idle状态的最大的连接数
}

func (p *HttpServer) Init() {
	if p.Port <= 0 {
		p.Port = 8080
	}

	var transport = http.DefaultTransport.(*http.Transport)
	transport.MaxConnsPerHost = p.MaxConnsPerHost
	transport.MaxIdleConns = p.MaxIdleConns
	transport.MaxIdleConnsPerHost = p.MaxIdleConnsPerHost
}

func (p *HttpServer) Addr() string {
	return fmt.Sprintf("%s:%d", p.BindAddr, p.Port)
}

func (p *HttpServer) NewHttpServer(router *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:         p.Addr(),
		Handler:      router,
		ReadTimeout:  time.Second * time.Duration(p.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(p.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(p.IdleTimeout),
	}
}
