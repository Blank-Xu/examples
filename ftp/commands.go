package ftp

import (
	"bytes"
	"strings"
)

// commandALLO responds 'ALLO' command
func commandALLO(ctx *Context) {
	ctx.WriteMessage(202, "Obsolete")
}

// commandCDUP responds 'CDUP' command
func commandCDUP(ctx *Context) {
	commandCWD(ctx)
}

// commandCWD responds 'CWD' command
func commandCWD(ctx *Context) {
	ctx.path = ctx.param
}

// commandDELE responds 'DELE' command
func commandDELE(ctx *Context) {

}

// commandEPRT responds 'EPRT' command
func commandEPRT(ctx *Context) {

}

// commandEPSV responds 'EPSV' command
func commandEPSV(ctx *Context) {

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
	buf.WriteTo(ctx.writer)
}

// commandLIST responds 'LIST' command
func commandLIST(ctx *Context) {
	ctx.WriteMessage(150, "Opening ASCII mode data connection for file list")

}

// commandUSER responds 'USER' command
func commandUSER(ctx *Context) {
	ctx.user = ctx.param
	ctx.WriteMessage(331, "OK")
}

// commandPASS responds 'PASS' command
func commandPASS(ctx *Context) {
	if ok := ctx.Authenticate(ctx.param); ok {
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
}
