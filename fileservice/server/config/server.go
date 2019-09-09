package config

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	Name                string `yaml:"name"`
	IP                  string `yaml:"ip"`
	Port                int    `yaml:"port"`
	ReadTimeout         int    `yaml:"read_timeout"`
	WriteTimeout        int    `yaml:"write_timeout"`
	IdleTimeout         int    `yaml:"idle_timeout"`
	MaxConnsPerHost     int    `yaml:"max_conns_per_host"`      // 每一个host对应的最大连接数
	MaxIdleConns        int    `yaml:"max_idle_conns"`          // 所有host对应的idle状态最大的连接总数
	MaxIdleConnsPerHost int    `yaml:"max_idle_conns_per_host"` // 每一个host对应idle状态的最大的连接数
}

func (p *Server) init() {
	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = p.MaxConnsPerHost
	http.DefaultTransport.(*http.Transport).MaxIdleConns = p.MaxIdleConns
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = p.MaxIdleConnsPerHost

	if p.Port == 0 {
		p.Port = 8080
	}
}

func (p *Server) Addr() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}

func (p *Server) NewServer(router *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:         p.Addr(),
		Handler:      router,
		ReadTimeout:  time.Second * time.Duration(p.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(p.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(p.IdleTimeout),
	}
}
