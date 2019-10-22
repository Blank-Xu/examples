package ftp

type HandlerFunc func(*Context)

var (
	routerMap = map[string]HandlerFunc{
		"ALLO": commandALLO,
		"CDUP": checkUser(commandCDUP),
		"CWD":  nil,
		"DELE": nil,
		"EPRT": nil,
		"EPSV": nil,
		"FEAT": nil,
		"LIST": nil,
		"NLST": nil,
		"MDTM": nil,
		"MKD":  nil,
		"MODE": nil,
		"NOOP": nil,
		"OPTS": commandOPTS,
		"PASS": commandPASS,
		"PASV": nil,
		"PORT": nil,
		"PWD":  nil,
		"QUIT": commandQUIT,
		"RETR": nil,
		"RNFR": nil,
		"RNTO": nil,
		"RMD":  nil,
		"SIZE": nil,
		"STOR": nil,
		"STRU": nil,
		"SYST": nil,
		"TYPE": commandTYPE,
		"USER": commandUSER,
		"XCUP": nil,
		"XCWD": nil,
		"XPWD": nil,
		"XRMD": nil,
	}

	checkUser = func(handler HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			if len(ctx.user) == 0 {
				ctx.WriteMessage(530, "not logged in")
				return
			}
			handler(ctx)
		}
	}

	checkParam = func(handler HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			if len(ctx.param) == 0 {
				ctx.WriteMessage(553, "action aborted, required param missing")
				return
			}
			handler(ctx)
		}
	}

	checkUserAndParam = func(handler HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			if len(ctx.user) == 0 {
				ctx.WriteMessage(530, "not logged in")
				return
			}

			if len(ctx.param) == 0 {
				ctx.WriteMessage(553, "action aborted, required param missing")
				return
			}
			handler(ctx)
		}
	}
)
