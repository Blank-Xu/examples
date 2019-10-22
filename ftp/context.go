package ftp

import (
	"bufio"
	"bytes"
	"errors"
	"net"
)

type Context struct {
	conn          net.Conn      // TCP connection
	addr          string        //
	writer        *bufio.Writer // Writer on the TCP connection
	reader        *bufio.Reader // Reader on the TCP connection
	user          []byte        // Authenticated user
	path          []byte        // Current path
	data          []byte        // request data
	command       string        // Command received on the connection
	param         []byte        // Param of the FTP command
	connectedTime int64         // Date of connection
	timeoutSecond int64         //
	ctxRnfr       []byte        // Rename from
	ctxRest       []byte        // Restart point
	transferTLS   bool          // Use TLS for transfer connection
	log           *Logger       // Client handler logging
}

func NewContext(conn net.Conn) *Context {
	p := &Context{
		conn:          conn,
		addr:          conn.RemoteAddr().String(),
		writer:        bufio.NewWriter(conn),
		reader:        bufio.NewReader(conn),
		user:          nil,
		path:          nil,
		param:         nil,
		connectedTime: 0,
		ctxRnfr:       nil,
		ctxRest:       nil,
		transferTLS:   false,
		log:           &Logger{},
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
	params := bytes.SplitN(p.data, []byte(" "), 2)
	var l = len(params)
	switch l {
	case 0:
		return false
	case 1:
		p.command = string(params[0])
	case 2:
		p.command = string(params[0])
		p.param = params[1]
	}
	return true
}

func (p *Context) Authenticate(pass []byte) bool {

	return false
}

func (p *Context) WriteLine(line []byte) (err error) {
	var buf bytes.Buffer
	buf.Write(line)
	buf.WriteString("\r\n")
	_, err = buf.WriteTo(p.writer)
	return
}

func (p *Context) WriteBuffer(buf *bytes.Buffer) (err error) {
	buf.WriteString("\r\n")
	_, err = buf.WriteTo(p.writer)
	return
}

func (p *Context) WriteMessage(code int32, msg string) (err error) {
	var buf bytes.Buffer
	buf.WriteRune(rune(code))
	buf.WriteByte(' ')
	buf.WriteString(msg)
	buf.WriteString("\r\n")
	_, err = buf.WriteTo(p.writer)
	return
}

// func (p *Context) Close() error {
// 	return p.conn.Close()
// }
