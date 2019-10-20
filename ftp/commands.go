package ftp

import (
	"bytes"
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

}

// commandUSER responds 'USER' command
func commandUSER(ctx *Context) {
	ctx.user = ctx.data
	ctx.WriteMessage(331, "OK")
}

// commandPASS responds 'PASS' command
func commandPASS(ctx *Context) {
	if ok := ctx.Authenticate(ctx.data); ok {
		ctx.WriteMessage(230, "Password ok, continue")
		return
	}

	// sort.SearchStrings()

	ctx.WriteMessage(530, "Incorrect password, not logged in")
	commandQUIT(ctx)
}

func commandOPTS(ctx *Context) {
	if bytes.EqualFold(ctx.data, []byte("UTF8")) ||
		bytes.EqualFold(ctx.data, []byte("UTF8 ON")) {
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
