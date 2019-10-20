package ftp

import (
	"net"
)

type Server struct {
	listener net.Listener
}

func (p *Server) Listen() error {
	return nil
}
