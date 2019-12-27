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
		"ABOR": &commandFunc{HandlerFunc: commandABOR, NeedLogin: true},
		"ALLO": &commandFunc{HandlerFunc: commandALLO},
		"CDUP": &commandFunc{HandlerFunc: commandCDUP, NeedLogin: true},
		"CWD":  &commandFunc{HandlerFunc: commandCWD, NeedLogin: true, NeedParam: true},
		"DELE": &commandFunc{HandlerFunc: commandDELE, NeedLogin: true, NeedParam: true},
		"EPRT": &commandFunc{HandlerFunc: commandEPRT, NeedLogin: true, NeedParam: true},
		"EPSV": &commandFunc{HandlerFunc: commandEPSV, NeedLogin: true},
		"FEAT": &commandFunc{HandlerFunc: commandFEAT},
		"LIST": &commandFunc{HandlerFunc: commandLIST, NeedLogin: true},
		"NLST": &commandFunc{HandlerFunc: commandNLST, NeedLogin: true},
		"MDTM": &commandFunc{HandlerFunc: commandMDTM, NeedLogin: true, NeedParam: true},
		"MKD":  &commandFunc{HandlerFunc: commandMKD, NeedLogin: true, NeedParam: true},
		"MODE": &commandFunc{HandlerFunc: commandMODE, NeedLogin: true, NeedParam: true},
		"NOOP": &commandFunc{HandlerFunc: commandNOOP},
		"OPTS": &commandFunc{HandlerFunc: commandOPTS},
		"PASS": &commandFunc{HandlerFunc: commandPASS, NeedLogin: false, NeedParam: true},
		"PASV": &commandFunc{HandlerFunc: commandPASV, NeedLogin: true},
		"PORT": &commandFunc{HandlerFunc: commandPORT, NeedLogin: true, NeedParam: true},
		"PWD":  &commandFunc{HandlerFunc: commandPWD},
		"QUIT": &commandFunc{HandlerFunc: commandQUIT},
		"RETR": &commandFunc{HandlerFunc: commandRETR, NeedLogin: true, NeedParam: true},
		"RNFR": &commandFunc{HandlerFunc: commandRNFR, NeedLogin: true, NeedParam: true},
		"RNTO": &commandFunc{HandlerFunc: commandRNTO, NeedLogin: true, NeedParam: true},
		"RMD":  &commandFunc{HandlerFunc: commandRMD, NeedLogin: true, NeedParam: true},
		"SIZE": &commandFunc{HandlerFunc: commandSIZE, NeedLogin: true, NeedParam: true},
		"STOR": &commandFunc{HandlerFunc: commandSTOR, NeedLogin: true, NeedParam: true},
		"STRU": &commandFunc{HandlerFunc: commandSTRU, NeedLogin: true, NeedParam: true},
		"SYST": &commandFunc{HandlerFunc: commandSYST, NeedLogin: true},
		"TYPE": &commandFunc{HandlerFunc: commandTYPE, NeedLogin: true},
		"USER": &commandFunc{HandlerFunc: commandUSER, NeedLogin: false, NeedParam: true},
		"XCUP": &commandFunc{HandlerFunc: commandCDUP, NeedLogin: true},
		"XCWD": &commandFunc{HandlerFunc: commandCWD, NeedLogin: true, NeedParam: true},
		"XPWD": &commandFunc{HandlerFunc: commandPWD},
		"XRMD": &commandFunc{HandlerFunc: commandRMD, NeedLogin: true, NeedParam: true},
	}

	checkLogin = func(handler HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			if len(ctx.pass) == 0 {
				ctx.WriteMessage(530, "Action aborted, should login")
				return
			}

			handler(ctx)
		}
	}

	checkParam = func(handler HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			if len(ctx.param) == 0 {
				ctx.WriteMessage(553, "Action aborted, required param missing")
				return
			}

			handler(ctx)
		}
	}
)

func init() {
	m := make(map[string]*commandFunc, len(routerMap))
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
