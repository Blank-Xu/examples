package ftp

type HandlerFunc func(*Context)

var RouterMap = map[string]HandlerFunc{
	"ALLO": commandALLO,
	"CDUP": nil,
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

var (
	funcAuthenticate = func(handler HandlerFunc) HandlerFunc {

	}
)

func init() {

}
