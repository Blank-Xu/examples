package ftp

import (
	"strings"
)

type commandFunc struct {
	HandlerFunc
	NeedLogin bool
	NeedParam bool
}

var (
	routerMap = map[string]*commandFunc{
		"ALLO": &commandFunc{HandlerFunc: commandALLO},
		"CDUP": &commandFunc{HandlerFunc: commandCDUP, NeedLogin: true},
		"CWD":  &commandFunc{HandlerFunc: commandCWD, NeedLogin: true, NeedParam: true},
		"DELE": &commandFunc{HandlerFunc: commandDELE, NeedLogin: true, NeedParam: true},
		"EPRT": &commandFunc{HandlerFunc: commandEPRT, NeedLogin: true, NeedParam: true},
		"EPSV": &commandFunc{HandlerFunc: commandEPSV, NeedLogin: true},
		"FEAT": &commandFunc{HandlerFunc: commandFEAT},
		"LIST": &commandFunc{HandlerFunc: commandLIST, NeedLogin: true},
		"NLST": nil,
		"MDTM": nil,
		"MKD":  nil,
		"MODE": nil,
		"NOOP": nil,
		"OPTS": &commandFunc{HandlerFunc: commandOPTS},
		"PASS": &commandFunc{HandlerFunc: commandPASS, NeedLogin: false, NeedParam: true},
		"PASV": nil,
		"PORT": nil,
		"PWD":  nil,
		"QUIT": &commandFunc{HandlerFunc: commandQUIT},
		"RETR": nil,
		"RNFR": nil,
		"RNTO": nil,
		"RMD":  nil,
		"SIZE": nil,
		"STOR": nil,
		"STRU": nil,
		"SYST": nil,
		"TYPE": &commandFunc{HandlerFunc: commandTYPE},
		"USER": &commandFunc{HandlerFunc: commandUSER},
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
)

func init() {
	var m = make(map[string]*commandFunc, len(routerMap))
	for command, fn := range routerMap {
		if fn == nil {
			continue
		}
		if fn.NeedLogin {
			fn.HandlerFunc = checkLogin(fn.HandlerFunc)
		}
		if fn.NeedParam {
			fn.HandlerFunc = checkParam(fn.HandlerFunc)
		}

		m[command] = fn

		cmd := strings.ToLower(command)
		m[cmd] = fn
	}
	routerMap = m
}
