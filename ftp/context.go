package ftp

import (
	"bufio"
	"bytes"
	"errors"
	"net"
)

type Context struct {
	conn          net.Conn      // TCP connection
	writer        *bufio.Writer // Writer on the TCP connection
	reader        *bufio.Reader // Reader on the TCP connection
	user          []byte        // Authenticated user
	path          []byte        // Current path
	command       []byte        // Command received on the connection
	data          []byte        // Param of the FTP command
	connectedTime int64         // Date of connection
	timeoutSecond int64
	ctxRnfr       []byte  // Rename from
	ctxRest       []byte  // Restart point
	transferTLS   bool    // Use TLS for transfer connection
	log           *Logger // Client handler logging
}

func NewContext(conn net.Conn) (*Context, error) {
	p := &Context{
		conn:          conn,
		writer:        bufio.NewWriter(conn),
		reader:        bufio.NewReader(conn),
		user:          nil,
		path:          nil,
		command:       nil,
		data:          nil,
		connectedTime: 0,
		ctxRnfr:       nil,
		ctxRest:       nil,
		transferTLS:   false,
		log:           &Logger{},
	}
	return p, nil
}

func (p *Context) Read() ([]byte, error) {
	if p.reader == nil {
		return nil, errors.New("reader is nil")
	}

	data, err := p.reader.ReadBytes('\n')
	if err != nil {

		switch err.(type) {
		case net.Error:
			// if err.
		default:

		}

		return nil, err
	}
	p.data = data

	return data, nil
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

func (p *Context) Close() error {
	return p.conn.Close()
}
