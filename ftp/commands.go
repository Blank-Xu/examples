package ftp

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

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
	conn, err := NewPassiveTCPListener(ctx.config.Host, int(ctx.config.PasvMinPort), int(ctx.config.PasvMaxPort))
	if err != nil {
		ctx.WriteMessage(425, "Data connection failed")
		return
	}
	if ctx.dataConn != nil {
		ctx.dataConn.Close()
		ctx.dataConn = nil
	}
	// ctx.dataConn = conn
	// ctx.WriteMessage(229,fmt.Sprintf("Entering Extended Passive Mode (|||%d|)", conn.Addr()))
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

}

// commandUSER responds 'USER' command
func commandUSER(ctx *Context) {
	ctx.user = string(ctx.param)
	ctx.WriteMessage(331, "OK")
}

// commandPASS responds 'PASS' command
func commandPASS(ctx *Context) {
	var pass = string(ctx.param)
	if ok := ctx.Authenticate(pass); ok {
		ctx.pass = pass
		ctx.WriteMessage(230, "Password ok, continue")
		return
	}

	ctx.WriteMessage(530, "Incorrect password, not logged in")
	commandQUIT(ctx)
}

func commandOPTS(ctx *Context) {
	var data = string(ctx.param)
	if data == "UTF8" || data == "UTF8 ON" {
		ctx.WriteMessage(200, "OK")
		return
	}
	ctx.WriteMessage(500, "Command not found")
}

// commandTYPE responds 'TYPE' command
func commandTYPE(ctx *Context) {

}

// commandQUIT for 'QUIT' command
func commandQUIT(ctx *Context) {
	ctx.WriteMessage(221, "Goodbye.")
	ctx.Close()
	ctx = nil
}
