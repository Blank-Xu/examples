package ftp

type Auth struct {
	User string
	Pass string
}

func middlewareAuth(auths []*Auth, ctx *Context) CommandFunc {

	return func(ctx *Context) {

	}
}
