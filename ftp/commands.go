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

// commandABOR responds 'ABOR' command
func commandABOR(ctx *Context) {
	ctx.WriteMessage(200, "OK")
	ctx.Abort()
}

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
	absPath := ctx.GetAbsPath(ctx.param)
	if _, err := os.Stat(absPath); err != nil {
		ctx.Error(err)
		ctx.WriteMessage(550, "Action not taken")
		return
	}
	ctx.path = string(ctx.param)
	ctx.WriteMessage(250, "Directory changed to "+absPath)
}

// commandDELE responds 'DELE' command
func commandDELE(ctx *Context) {
	absPath := ctx.GetAbsPath(ctx.param)
	if err := os.RemoveAll(absPath); err != nil {
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

	ctx.SetDataConn(conn)

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
	ctx.WriteBuffer(buf)
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

	// TODO: client can't get info
	for _, file := range files {
		s := fmt.Sprintf("%s 1 ftp ftp %12d %s %s\r\n",
			file.Mode(),
			file.Size(),
			GetFileModTime(time.Now(), file.ModTime()),
			file.Name())
		fmt.Println(s)

		ctx.writer.WriteString(s)
		ctx.writer.Flush()

		// _, err = fmt.Fprintf(ctx.writer,
		// 	"%s 1 ftp ftp %12d %s %s\r\n",
		// 	file.Mode(),
		// 	file.Size(),
		// 	GetFileModTime(time.Now(), file.ModTime()),
		// 	file.Name())
		// if err != nil {
		// 	ctx.Error(err)
		// 	ctx.WriteMessage(550, "Transfer failed.")
		// 	return
		// }
	}
	// ctx.WriteMessage(226, "Transfer complete.")

	return

	var (
		now = time.Now()
		buf bytes.Buffer
	)
	buf.Grow(1024)
	for _, file := range files {
		buf.WriteString(file.Mode().String())
		buf.WriteString(" 1 owner group ")
		buf.WriteString(fmt.Sprintf("%12d", file.Size()))
		buf.WriteByte(' ')
		buf.WriteString(GetFileModTime(now, file.ModTime()))
		buf.WriteByte(' ')
		buf.WriteString(file.Name())
		buf.WriteString("\r\n")
	}
	if err = ctx.WriteBuffer(&buf); err != nil {
		ctx.WriteMessage(550, "Transfer failed.")
		return
	}
	ctx.WriteMessage(226, "Transfer complete.")

	time.Sleep(10 * time.Millisecond)
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
	if err := os.MkdirAll(absPath, 0766); err != nil {
		ctx.WriteMessage(550, "Action not taken")
		return
	}
	ctx.WriteMessage(257, "Directory created")
}

// commandMODE responds 'MODE' command
func commandMODE(ctx *Context) {
	param := string(ctx.param)
	if param == "S" || param == "s" {
		ctx.WriteMessage(200, "OK")
		return
	}
	ctx.WriteMessage(504, "MODE is an obsolete command")
}

// commandNOOP responds 'NOOP' command
func commandNOOP(ctx *Context) {
	ctx.WriteMessage(200, "OK")
}

// commandOPTS responds 'OPTS' command
func commandOPTS(ctx *Context) {
	param := string(ctx.param)
	if param == "UTF8" || param == "UTF8 ON" {
		ctx.WriteMessage(200, "OK")
		return
	}
	ctx.WriteMessage(500, "Command not found")
}

// commandPASS responds 'PASS' command
func commandPASS(ctx *Context) {
	if len(ctx.user) == 0 {
		ctx.WriteMessage(503, "User required")
		return
	}

	param := string(ctx.param)
	if ok := ctx.Authenticate(param); ok {
		ctx.pass = param
		ctx.WriteMessage(230, "Password verified, continue")
		return
	}

	ctx.WriteMessage(530, "Incorrect password")
	ctx.command = "QUIT"
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

	// if err = listener.SetDeadline(time.Now().Add(time.Duration(ctx.config.DeadlineSeconds))); err != nil {
	// 	ctx.Error(err)
	// }

	addr := listener.Addr().(*net.TCPAddr)
	p1 := addr.Port / 256
	p2 := addr.Port - (p1 * 256)

	var buf bytes.Buffer
	buf.Grow(256)
	buf.WriteString("227 Entering Passive Mode (")
	buf.WriteString(ctx.config.externalIP)
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
		ctx.WriteMessage(500, "params invalid")
		return
	}
	p1, err := strconv.Atoi(params[4])
	if err != nil {
		ctx.WriteMessage(500, "params invalid")
		return
	}
	p2, err := strconv.Atoi(params[5])
	if err != nil {
		ctx.WriteMessage(500, "params invalid")
		return
	}
	port := p1*256 + p2

	var buf bytes.Buffer
	buf.Grow(len(param))
	buf.WriteString(params[0])
	buf.WriteByte('.')
	buf.WriteString(params[1])
	buf.WriteByte('.')
	buf.WriteString(params[2])
	buf.WriteByte('.')
	buf.WriteString(params[3])

	conn, err := NewActiveTCPConn(buf.String(), port)
	if err != nil {
		ctx.WriteMessage(425, "Data connection failed")
		return
	}
	ctx.SetDataConn(conn)

	ctx.WriteMessage(200, fmt.Sprintf("Connection established (%d)", port))
}

// commandPWD responds 'PWD' command
func commandPWD(ctx *Context) {
	ctx.WriteMessage(257, "\""+ctx.path+"\" is the current directory")
}

// commandQUIT for 'QUIT' command
// http://cr.yp.to/ftp/quit.html
func commandQUIT(ctx *Context) {
	ctx.WriteMessage(221, "Bye.")
	ctx.Close()
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
		ctx.WriteMessage(550, fmt.Sprintf("Couldn't access %s: %v", absPath, err))
		return
	}
	ctx.rnfr = absPath
	ctx.WriteMessage(200, "Sure, give me a target")
}

// commandRNTO responds 'RNTO' command
func commandRNTO(ctx *Context) {
	if len(ctx.rnfr) == 0 {
		ctx.WriteMessage(503, "Bad sequence of commands: use RNFR first.")
		return
	}

	absPath := ctx.GetAbsPath(ctx.param)
	if err := os.Rename(ctx.rnfr, absPath); err != nil {
		ctx.WriteMessage(550, "Action not taken")
		return
	}

	ctx.rnfr = ""
	ctx.WriteMessage(250, "File renamed")
}

// commandRMD responds 'RMD' command
func commandRMD(ctx *Context) {
	absPath := ctx.GetAbsPath(ctx.param)
	if err := os.RemoveAll(absPath); err != nil {
		ctx.WriteMessage(550, "Action not taken")
		return
	}
	ctx.WriteMessage(250, "Directory deleted")
}

// commandSIZE responds 'SIZE' command
func commandSIZE(ctx *Context) {
	absPath := ctx.GetAbsPath(ctx.param)
	file, err := os.Stat(absPath)
	if err != nil {
		ctx.WriteMessage(450, "file not available")
		return
	}
	ctx.WriteMessage(213, strconv.FormatInt(file.Size(), 10))
}

// commandSTOR responds 'STOR' command
func commandSTOR(ctx *Context) {
	if ctx.dataConn == nil {
		ctx.WriteMessage(450, "have no connection to transfer")
		return
	}

	ctx.WriteMessage(150, "Data transfer starting")
	absPath := ctx.GetAbsPath(ctx.param)
	if err := ctx.TransferFile(absPath, true, false); err != nil {
		ctx.WriteMessage(550, "Transfer failed, err: "+err.Error())
		return
	}
	ctx.WriteMessage(226, "Transfer complete.")
}

// commandSTRU responds 'STRU' command
func commandSTRU(ctx *Context) {
	param := string(ctx.param)
	if param == "F" || param == "f" {
		ctx.WriteMessage(200, "OK")
		return
	}
	ctx.WriteMessage(504, "STRU is an obsolete command")
}

// commandSYST responds 'SYST' command
func commandSYST(ctx *Context) {
	ctx.WriteMessage(215, "UNIX Type: L8")
}

// commandTYPE responds 'TYPE' command
func commandTYPE(ctx *Context) {
	param := string(ctx.param)
	switch param {
	case "A", "a":
		ctx.WriteMessage(200, "Type set to ASCII")
	case "I", "i":
		ctx.WriteMessage(200, "Type set to binary")
	default:
		ctx.WriteMessage(500, "Invalid type")
	}
}

// commandUSER responds 'USER' command
// http://cr.yp.to/ftp/user.html
func commandUSER(ctx *Context) {
	ctx.user = string(ctx.param)
	ctx.WriteMessage(331, "OK")
}
