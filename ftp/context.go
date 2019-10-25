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
	config *Config
	conn   *net.TCPConn

	dataConn      *net.TCPConn
	listener      *net.TCPListener
	writer        *bufio.Writer // Writer on the TCP connection
	reader        *bufio.Reader // Reader on the TCP connection
	user          string        // Authenticated user
	pass          string        //
	workDir       string        // work dir
	path          string        // Current path
	data          []byte        // request data
	command       string        // Command received on the connection
	param         []byte        // Param of the FTP command
	connectedTime int64         // Date of connection
	timeoutSecond int64         //
	fnfr          string        // Rename from
	rest          []byte        // Restart point
	errs          []error       // record errors
	log           *Logger       // Client handler logging
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
		p.workDir = filepath.Join(p.workDir, account.Dir)
		return true
	}
	return false
}

func (p *Context) WriteLine(line []byte) (err error) {
	var buf bytes.Buffer
	buf.Grow(len(line) + 10)
	buf.Write(line)
	buf.WriteString("\r\n")
	_, err = buf.WriteTo(p.writer)
	return p.Error(err)
}

func (p *Context) WriteBuffer(buf *bytes.Buffer) (err error) {
	buf.WriteString("\r\n")
	_, err = buf.WriteTo(p.writer)
	return p.Error(err)
}

func (p *Context) WriteMessage(code int, msg string) (err error) {
	var buf bytes.Buffer
	buf.Grow(len(msg) + 10)
	buf.WriteString(strconv.Itoa(code))
	buf.WriteByte(' ')
	buf.WriteString(msg)
	buf.WriteString("\r\n")
	_, err = buf.WriteTo(p.writer)
	return p.Error(err)
}

func (p *Context) ChangeDir(dir string) bool {
	dir = filepath.Join(p.workDir, dir)
	_, err := os.Stat(dir)
	if err != nil {
		p.Error(err)
		p.path = dir
		return true
	}
	return false
}

func (p *Context) GetAbsPath(path []byte) string {
	var newPath string
	if len(path) == 0 {
		return p.path
	} else if path[0] == '/' {
		return p.path
	} else {
		newPath = filepath.Join(p.path, string(path))
	}

	if len(newPath) <= len(p.workDir) {
		return p.workDir
	}

	return newPath
}

func (p *Context) Close() {
	if p.conn != nil {
		p.Error(p.conn.Close())
	}
	if p.dataConn != nil {
		p.Error(p.dataConn.Close())
	}
	if p.listener != nil {
		p.Error(p.listener.Close())
	}
}

func (p *Context) Error(err error) error {
	if err != nil {
		p.errs = append(p.errs, err)
	}
	return err
}
