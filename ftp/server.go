package ftp

import (
	"log"
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

	return &Server{config: cfg}, cfg.checkAll()
}

// func SetConfig(cfg Config)  {
// 	DefaultServer.SetConfig(cfg)
// }
//
// func (p *Server) SetConfig(cfg Config) {
// 	cfg.checkAll()
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

	if p.listener, err = net.ListenTCP("tcp", addr); err != nil {
		return err
	}

	log.Println("server listening address: " + p.config.addr)

	for {
		conn, err := p.listener.AcceptTCP()
		if err != nil {
			return err
		}

		go p.handle(conn)
	}
}

func (p *Server) handle(conn *net.TCPConn) {
	var ctx *Context

	defer func() {
		if ctx != nil {
			ctx.Close()
			if len(ctx.errs) > 0 {
				log.Println(ctx.errs)
			}
			ctx = nil
		}

		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	var err error
	if err = conn.SetKeepAlive(true); err != nil {
		return
	}
	var now = time.Now().Unix()
	if err = conn.SetDeadline(time.Unix(now+int64(p.config.DeadlineSeconds), 0)); err != nil {
		return
	}
	if err = conn.SetReadDeadline(time.Unix(now+int64(p.config.ReadDeadlineSeconds), 0)); err != nil {
		return
	}
	if err = conn.SetWriteDeadline(time.Unix(now+int64(p.config.WriteDeadlineSeconds), 0)); err != nil {
		return
	}
	if err = conn.SetKeepAlivePeriod(time.Duration(p.config.KeepAlivePeriodSeconds)); err != nil {
		return
	}

	ctx = NewContext(p.config, conn)
	if err = ctx.WriteMessage(220, p.config.ServerName); err != nil {
		ctx.WriteMessage(550, "refuse")
		return
	}

	for {
		if err := ctx.Read(); err != nil {
			log.Println(err)
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
			log.Println(ctx.errs)
		}

		if ctx.command == "QUIT" || ctx.command == "quit" {
			break
		}

		ctx.errs = make([]error, 0, 10)
	}
}

func Close() error {
	return DefaultServer.Close()
}

func (p *Server) Close() error {
	if p.listener != nil {
		return p.listener.Close()
	}
	return nil
}
