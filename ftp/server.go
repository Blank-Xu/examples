package ftp

import (
	"net"
)

var DefaultServer = NewServer(&Config{})

type HandlerFunc func(*Context)

type Server struct {
	config   *Config
	listener *net.TCPListener
}

func NewServer(cfg *Config) *Server {
	if cfg == nil {
		cfg = &Config{}
	}

	cfg.init()

	return &Server{config: cfg}
}

// func SetConfig(cfg Config)  {
// 	DefaultServer.SetConfig(cfg)
// }
//
// func (p *Server) SetConfig(cfg Config) {
// 	cfg.init()
// 	p.config = cfg
// }

func ListenAndServe() error {
	return DefaultServer.ListenAndServe()
}

func (p *Server) ListenAndServe() error {
	addr, err := net.ResolveTCPAddr("tcp", p.config.addr)
	if err != nil {
		return err
	}

	p.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := p.listener.AcceptTCP()
		if err != nil {
			break
		}
		go p.handle(conn)
	}

	return nil
}

func (p *Server) handle(conn *net.TCPConn) {
	defer func() {
		if err := recover(); err != nil {

		}
		conn.Close()
	}()
	ctx := NewContext(p, conn)

	ctx.WriteMessage(220, p.config.ServerName)

	for {
		if err := ctx.Read(); err != nil {
			break
		}

		if !ctx.ParseParam() {
			ctx.WriteMessage(500, "Command not found")
			continue
		}

		fn, ok := routerMap[ctx.command]
		if !ok {
			ctx.WriteMessage(500, "Command not found")
			continue
		}

		fn(ctx)

		if len(ctx.errs) > 0 {

		}
	}
}

func Close() error {
	return DefaultServer.Close()
}

func (p *Server) Close() error {
	if p.listener == nil {
		return nil
	}
	return p.listener.Close()
}
