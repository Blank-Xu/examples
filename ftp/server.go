package ftp

import (
	"net"
	"time"
)

var DefaultServer, _ = NewServer(&Config{})

type HandlerFunc func(*Context)

type Server struct {
	config   *Config
	listener *net.TCPListener
}

func NewServer(cfg *Config) (*Server, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	return &Server{config: cfg}, cfg.init()
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
			return err
		}

		if err = conn.SetKeepAlive(true); err != nil {
			return err
		}

		var now = time.Now().Unix()

		if err = conn.SetDeadline(time.Unix(now+int64(p.config.DeadlineSeconds), 0)); err != nil {
			return err
		}
		if err = conn.SetReadDeadline(time.Unix(now+int64(p.config.ReadDeadlineSeconds), 0)); err != nil {
			return err
		}
		if err = conn.SetWriteDeadline(time.Unix(now+int64(p.config.WriteDeadlineSeconds), 0)); err != nil {
			return err
		}
		if err = conn.SetKeepAlivePeriod(time.Duration(p.config.KeepAlivePeriodSeconds)); err != nil {
			return err
		}

		go p.handle(conn)
	}
}

func (p *Server) handle(conn *net.TCPConn) {
	defer func() {
		if err := recover(); err != nil {

		}
		conn.Close()
	}()

	ctx := NewContext(p.config, conn)
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
		if !ok || fn == nil {
			ctx.WriteMessage(500, "Command not found")
			continue
		}

		fn.HandlerFunc(ctx)

		if len(ctx.errs) > 0 {

			ctx.errs = make([]error, 0, 10)
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
