package ftp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Context struct {
	config      *Config
	conn        *net.TCPConn
	connectTime time.Time     // Date of connection
	reader      *bufio.Reader // Reader on the TCP connection
	writer      *bufio.Writer // Writer on the TCP connection
	workDir     string        // work dir

	data    []byte // Request data bytes
	command string // Command received on the connection
	param   []byte // Param of the FTP command

	dataConn *net.TCPConn
	listener *net.TCPListener
	path     string // Current path
	user     string // Authenticated user
	pass     string //

	rnfr string // Rename from command path
	rest []byte // Restart point

	errs []error // record errors
	log  *Logger // Client handler logging
}

func NewContext(config *Config, conn *net.TCPConn) *Context {
	return &Context{
		config:      config,
		conn:        conn,
		connectTime: time.Now().UTC(),
		writer:      bufio.NewWriter(conn),
		reader:      bufio.NewReader(conn),
		workDir:     config.Dir,
		path:        "/",
		log:         &Logger{},
	}
}

func (p *Context) Read() (err error) {
	if p.reader == nil {
		return errors.New("reader is nil")
	}

	p.data, err = p.reader.ReadBytes('\n')

	fmt.Printf("\n[data: %s, err: %v]", p.data, err)

	if err != nil {
		switch terr := err.(type) {
		case net.Error:
			if terr.Timeout() {
				p.conn.SetDeadline(time.Now().Add(time.Minute))
				p.WriteMessage(421, "timeout")
			}
		default:
		}

		return
	}

	return
}

func (p *Context) ParseParam() bool {
	p.data = bytes.Trim(p.data, "\r\n")
	params := bytes.SplitN(p.data, []byte{' '}, 2)

	switch len(params) {
	case 0:
		return false
	case 1:
		p.command = string(params[0])
		p.param = make([]byte, 0, 20)
	case 2:
		p.command = string(params[0])
		p.param = make([]byte, len(params[1]))
		copy(p.param, params[1])
	}

	return true
}

func (p *Context) Authenticate(pass string) bool {
	account, ok := p.config.accountMap[p.user]
	if !ok {
		return false
	}
	if account.Password == pass {
		p.workDir = GetAbsPath(p.workDir, account.Dir)
		return true
	}

	return false
}

func (p *Context) WriteBuffer(buf *bytes.Buffer) (err error) {
	buf.WriteString("\r\n")
	if _, err = buf.WriteTo(p.writer); err != nil {
		p.Error(err)
		return
	}
	if err = p.writer.Flush(); err != nil {
		p.Error(err)
	}

	return
}

func (p *Context) WriteMessage(code int, msg string) (err error) {
	var buf bytes.Buffer

	buf.Grow(len(msg) + 10)
	buf.WriteString(strconv.Itoa(code))
	buf.WriteByte(' ')
	buf.WriteString(msg)
	buf.WriteString("\r\n")
	if _, err = buf.WriteTo(p.writer); err != nil {
		p.Error(err)
	}
	if err = p.writer.Flush(); err != nil {
		p.Error(err)
	}

	return
}

func (p *Context) GetAbsPath(path []byte) string {
	return GetAbsPath(p.workDir, string(path))
}

func (p *Context) TransferFile(path string, write, append bool) error {

	return nil
}

func (p *Context) SetDataConn(conn *net.TCPConn) {
	if p.dataConn != nil {
		p.dataConn.Close()
		p.dataConn = nil
	}
	p.dataConn = conn
}

func (p *Context) Abort() {
	if p.dataConn != nil {
		if err := p.dataConn.Close(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			p.Error(err)
		}
	}

	if p.listener != nil {
		if err := p.listener.Close(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			p.Error(err)
		}
	}
}

func (p *Context) Close() {
	p.Abort()

	if p.conn != nil {
		if err := p.conn.Close(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			p.Error(err)
		}
	}
}

func (p *Context) Error(err error) {
	if err != nil {
		p.errs = append(p.errs, err)
	}
}
