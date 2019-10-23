package ftp

import (
	"strings"
)

var (
	routerMap = map[string]HandlerFunc{
		"ALLO": commandALLO,
		"CDUP": checkLogin(commandCDUP),
		"CWD":  checkLoginAndParam(commandCWD),
		"DELE": checkLoginAndParam(commandDELE),
		"EPRT": checkLoginAndParam(commandEPRT),
		"EPSV": checkLogin(commandEPSV),
		"FEAT": commandFEAT,
		"LIST": checkLogin(commandLIST),
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

	checkLogin = func(handler HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			if len(ctx.pass) == 0 {
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

	checkLoginAndParam = func(handler HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			if len(ctx.pass) == 0 {
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

func init() {
	var m = make(map[string]HandlerFunc, len(routerMap))
	for command, fn := range routerMap {
		cmd := strings.ToLower(command)
		m[command] = fn
		m[cmd] = fn
	}
	routerMap = m
}
