package ftp

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

type Context struct {
	config        *Config
	conn          *net.TCPConn
	connectedTime int64         // Date of connection
	timeoutSecond int64         //
	reader        *bufio.Reader // Reader on the TCP connection
	writer        *bufio.Writer // Writer on the TCP connection
	workDir       string        // work dir

	data    []byte // Request data bytes
	command string // Command received on the connection
	param   []byte // Param of the FTP command

	dataConn *net.TCPConn
	listener *net.TCPListener
	path     string // Current path
	user     string // Authenticated user
	pass     string //
	rnfr     string // Rename from command path
	rest     []byte // Restart point

	errs []error // record errors
	log  *Logger // Client handler logging
}

func NewContext(config *Config, conn *net.TCPConn) *Context {
	p := &Context{
		config:  config,
		conn:    conn,
		writer:  bufio.NewWriter(conn),
		reader:  bufio.NewReader(conn),
		workDir: config.Dir,
		path:    config.Dir,
		log:     &Logger{},
	}
	return p
}

func (p *Context) Read() error {
	if p.reader == nil {
		return errors.New("reader is nil")
	}

	var err error
	p.data, err = p.reader.ReadBytes('\n')
	if err != nil {

		switch err.(type) {
		case net.Error:
			// if err.
		default:

		}

		return err
	}
	return nil
}

func (p *Context) ParseParam() bool {
	p.data = bytes.Trim(p.data, "\r\n")
	params := bytes.SplitN(p.data, []byte{' '}, 2)
	switch len(params) {
	case 0:
		return false
	case 1:
		p.command = string(params[0])
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

func (p *Context) WriteLine(line []byte) (err error) {
	var buf bytes.Buffer
	buf.Grow(len(line) + 10)
	buf.Write(line)
	buf.WriteString("\r\n")
	if _, err = buf.WriteTo(p.writer); err != nil {
		p.Error(err)
	}
	return
}

func (p *Context) WriteBuffer(buf *bytes.Buffer) (err error) {
	buf.WriteString("\r\n")
	if _, err = buf.WriteTo(p.writer); err != nil {
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
	return
}

func (p *Context) ChangeDir(dir string) bool {
	dir = filepath.Join(p.workDir, dir)
	if _, err := os.Stat(dir); err != nil {
		p.Error(err)
		return false
	}
	p.path = dir
	return true
}

func (p *Context) GetAbsPath(path []byte) string {
	return GetAbsPath(p.workDir, string(path))
}

func (p *Context) Upload(write bool, append bool) error {

	return nil
}

func (p *Context) Close() {
	var err error
	if p.conn != nil {
		if err = p.conn.Close(); err != nil {
			p.Error(err)
		}
	}
	if p.dataConn != nil {
		if err = p.dataConn.Close(); err != nil {
			p.Error(err)
		}
	}
	if p.listener != nil {
		if err = p.listener.Close(); err != nil {
			p.Error(err)
		}
	}
}

func (p *Context) Error(err error) {
	p.errs = append(p.errs, err)
}
