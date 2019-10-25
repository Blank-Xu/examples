package ftp

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var listRegexp = regexp.MustCompile("^-[alt]+$")

// commandALLO responds 'ALLO' command
func commandALLO(ctx *Context) {
	ctx.WriteMessage(202, "Obsolete")
}

// commandCDUP responds 'CDUP' command
func commandCDUP(ctx *Context) {
	ctx.param = []byte("..")
	commandCWD(ctx)
}

// commandCWD responds 'CWD' command
func commandCWD(ctx *Context) {
	path := ctx.GetAbsPath(ctx.param)
	if ctx.ChangeDir(path) {
		ctx.WriteMessage(250, "Directory changed to "+path)
		return
	}
	ctx.WriteMessage(550, "Action not taken")
}

// commandDELE responds 'DELE' command
func commandDELE(ctx *Context) {
	path := ctx.GetAbsPath(ctx.param)
	if err := os.RemoveAll(path); err != nil {
		ctx.Error(err)
		ctx.WriteMessage(550, "Action not taken")
		return
	}
	ctx.WriteMessage(250, "File deleted")
}

// commandEPRT responds 'EPRT' command
func commandEPRT(ctx *Context) {
	parts := bytes.Split(ctx.param, []byte{ctx.param[0]})
	if len(parts) < 3 {
		ctx.WriteMessage(553, "action aborted, required param missing")
		return
	}
	addressFamily := string(parts[1])
	if addressFamily != "1" && addressFamily != "2" {
		ctx.WriteMessage(522, "Network protocol not supported, use (1,2)")
		return
	}

	port, err := strconv.Atoi(string(parts[3]))
	if err != nil {
		ctx.WriteMessage(553, "action aborted, required param missing")
		return
	}

	host := string(parts[2])
	conn, err := NewActiveTCPConn(host, port)
	if err != nil {
		ctx.WriteMessage(425, "Data connection failed")
		return
	}

	if ctx.dataConn != nil {
		ctx.dataConn.Close()
		ctx.dataConn = nil
	}
	ctx.dataConn = conn

	ctx.WriteMessage(200, fmt.Sprintf("Connection established (%d)", port))
}

// commandEPSV responds 'EPSV' command
func commandEPSV(ctx *Context) {
	listener, err := NewPassiveTCPListener(ctx.config.Host, int(ctx.config.PasvMinPort), int(ctx.config.PasvMaxPort))
	if err != nil {
		ctx.WriteMessage(425, "Data connection failed")
		return
	}
	if ctx.listener != nil {
		ctx.listener.Close()
		ctx.listener = nil
	}
	ctx.listener = listener

	addr := listener.Addr().(*net.TCPAddr)
	ctx.WriteMessage(229, fmt.Sprintf("Entering Extended Passive Mode (|||%d|)", addr.Port))
}

const _msgFEAT = "211-Features supported:\r\n" +
	" EPRT\r\n" +
	" EPSV\r\n" +
	" MDTM\r\n" +
	" SIZE\r\n" +
	" UTF8\r\n" +
	"211 End FEAT.\r\n"

// commandFEAT responds 'FEAT' command
func commandFEAT(ctx *Context) {
	var buf = bytes.NewBufferString(_msgFEAT)
	_, err := buf.WriteTo(ctx.writer)
	ctx.Error(err)
}

// commandLIST responds 'LIST' command
func commandLIST(ctx *Context) {
	ctx.WriteMessage(150, "Opening ASCII mode data connection for file list")

	if listRegexp.Match(ctx.param) {
		ctx.param = nil
	}

	absPath := ctx.GetAbsPath(ctx.param)
	files, err := GetFileList(absPath, string(ctx.param))
	if err != nil {
		ctx.WriteMessage(450, "Could not STAT: "+err.Error())
		return
	}

	var buf bytes.Buffer
	buf.Grow(1024)
	for _, file := range files {
		buf.WriteString(file.Mode().String())
		buf.WriteString(" 1 owner group ")
		buf.WriteString(strconv.FormatInt(file.Size(), 10))
		buf.WriteByte(' ')
		buf.WriteString(strconv.FormatInt(file.ModTime().UTC().Unix(), 10))
		buf.WriteByte(' ')
		buf.WriteString(file.Name())
		buf.WriteString("\r\n")
	}
	ctx.WriteBuffer(&buf)
}

// commandNLST responds 'NLST' command
func commandNLST(ctx *Context) {
	ctx.WriteMessage(150, "Opening ASCII mode data connection for file list")

	if listRegexp.Match(ctx.param) {
		ctx.param = nil
	}

	absPath := ctx.GetAbsPath(ctx.param)
	files, err := filepath.Glob(absPath)
	if err != nil {
		ctx.WriteMessage(450, "Could not STAT: "+err.Error())
		return
	}

	var buf bytes.Buffer
	buf.Grow(1024)
	for _, file := range files {
		buf.WriteString(file)
		buf.WriteString("\r\n")
	}
	ctx.WriteBuffer(&buf)
}

// commandMDTM responds 'MDTM' command
func commandMDTM(ctx *Context) {
	absPath := ctx.GetAbsPath(ctx.param)
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		ctx.WriteMessage(450, "File not available")
		return
	}
	ctx.WriteMessage(213, fileInfo.ModTime().Format("%Y%m%d%H%M%S"))
}

// commandMKD responds 'MKD' command
func commandMKD(ctx *Context) {
	absPath := ctx.GetAbsPath(ctx.param)
	if err := os.Mkdir(absPath, 0666); err != nil {
		ctx.WriteMessage(550, "Action not taken")
		return
	}
	ctx.WriteMessage(257, "Directory created")
}

// commandMODE responds 'MODE' command
func commandMODE(ctx *Context) {
	mode := string(ctx.param)
	switch mode {
	case "s", "S":
		ctx.WriteMessage(200, "OK")
	default:
		ctx.WriteMessage(504, "MODE is an obsolete command")
	}
}

// commandNOOP responds 'NOOP' command
func commandNOOP(ctx *Context) {
	ctx.WriteMessage(200, "OK")
}

// commandOPTS responds 'OPTS' command
func commandOPTS(ctx *Context) {
	opts := string(ctx.param)
	if opts == "UTF8" || opts == "UTF8 ON" {
		ctx.WriteMessage(200, "OK")
		return
	}
	ctx.WriteMessage(500, "Command not found")
}

// commandUSER responds 'USER' command
func commandUSER(ctx *Context) {
	ctx.user = string(ctx.param)
	ctx.WriteMessage(331, "OK")
}

// commandPASS responds 'PASS' command
func commandPASS(ctx *Context) {
	pass := string(ctx.param)
	if ok := ctx.Authenticate(pass); ok {
		ctx.pass = pass
		ctx.WriteMessage(230, "Password ok, continue")
		return
	}

	ctx.WriteMessage(530, "Incorrect password, not logged in")
	commandQUIT(ctx)
}

// commandPASV responds 'PASV' command
func commandPASV(ctx *Context) {
	listener, err := NewPassiveTCPListener(ctx.config.Host, int(ctx.config.PasvMinPort), int(ctx.config.PasvMaxPort))
	if err != nil {
		ctx.WriteMessage(425, "Data connection failed")
		return
	}

	if ctx.listener != nil {
		ctx.listener.Close()
		ctx.listener = nil
	}
	ctx.listener = listener

	addr := listener.Addr().(*net.TCPAddr)
	p1 := addr.Port / 256
	p2 := addr.Port - (p1 * 256)

	var buf bytes.Buffer
	buf.Grow(256)
	buf.WriteString("227 Entering Passive Mode (")
	buf.WriteString(strings.ReplaceAll(ctx.config.Host, ",", "."))
	buf.WriteByte(',')
	buf.WriteString(strconv.Itoa(p1))
	buf.WriteByte(',')
	buf.WriteString(strconv.Itoa(p2))
	buf.WriteByte(')')

	ctx.WriteBuffer(&buf)
}

// commandPORT responds 'PORT' command
func commandPORT(ctx *Context) {
	param := string(ctx.param)
	params := strings.Split(param, ",")
	if len(params) < 5 {
		ctx.WriteMessage(500, "")
		return
	}

	var buf bytes.Buffer
	buf.Grow(len(param))
	buf.WriteString(params[0])
	buf.WriteByte('.')
	buf.WriteString(params[1])
	buf.WriteByte('.')
	buf.WriteString(params[2])
	buf.WriteByte('.')
	buf.WriteString(params[3])

	p1, _ := strconv.Atoi(params[4])
	p2, _ := strconv.Atoi(params[5])
	port := p1*256 + p2

	dataConn, err := NewActiveTCPConn(buf.String(), port)
	if err != nil {
		ctx.WriteMessage(425, "Data connection failed")
		return
	}
	defer dataConn.Close()

	if ctx.dataConn != nil {
		ctx.dataConn.Close()
		ctx.dataConn = nil
	}
	ctx.dataConn = dataConn

	ctx.WriteMessage(200, fmt.Sprintf("Connection established (%d)", port))
}

// commandPWD responds 'PWD' command
func commandPWD(ctx *Context) {
	ctx.WriteMessage(257, "\""+ctx.path+"\" is the current directory")
}

// commandQUIT for 'QUIT' command
func commandQUIT(ctx *Context) {
	ctx.WriteMessage(221, "Goodbye.")
	ctx.Close()
	ctx = nil
}

// commandRETR responds 'RETR' command
func commandRETR(ctx *Context) {
	if ctx.dataConn == nil {
		ctx.WriteMessage(551, "Data connection invalid.")
		return
	}

	absPath := ctx.GetAbsPath(ctx.param)
	file, err := os.Open(absPath)
	if err != nil {
		ctx.WriteMessage(551, "File not available")
		return
	}
	defer file.Close()

	ctx.WriteMessage(150, "Data connection open. Transfer starting.")

	defer ctx.dataConn.Close()

	if _, err = io.Copy(ctx.dataConn, file); err != nil {
		ctx.Error(err)
		ctx.WriteMessage(550, "Action not taken")
		return
	}

	ctx.WriteMessage(226, "Transfer complete.")

	time.Sleep(time.Millisecond * 10)
}

// commandRNFR responds 'RNFR' command
func commandRNFR(ctx *Context) {
	absPath := ctx.GetAbsPath(ctx.param)
	if _, err := os.Stat(absPath); err != nil {
		if os.IsNotExist(err) {

		}
	}

	ctx.fnfr = absPath
	ctx.WriteMessage(200, "")
}

// commandTYPE responds 'TYPE' command
func commandTYPE(ctx *Context) {

}
