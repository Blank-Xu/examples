package ftp

import (
	"net"
)

var DefaultServer = NewServer(Config{})

type Server struct {
	serverName string
	conn       *net.TCPConn
	addr       string
}

func NewServer(cfg Config) *Server {
	cfg.init()

	s := &Server{
		serverName: cfg.ServerName,
		addr:       GetTcpAddr(cfg.HostName, cfg.Port),
	}

	return s
}

func ListenAndServe() error {
	return DefaultServer.ListenAndServe()
}

func (p *Server) ListenAndServe() error {
	addr, err := net.ResolveTCPAddr("tcp", p.addr)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	for {
		p.conn, err = listener.AcceptTCP()
		if err != nil {
			break
		}
		go p.handle()
	}

	return nil
}

func (p *Server) handle() {
	defer func() {
		if err := recover(); err != nil {

		}
		p.conn.Close()
	}()
	ctx := NewContext(p.conn)

	ctx.WriteMessage(220, p.serverName)

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
	}
}

func Close() error {
	return DefaultServer.Close()
}

func (p *Server) Close() error {
	if p.conn == nil {
		return nil
	}
	return p.conn.Close()
}
