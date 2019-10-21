package ftp

type Auth struct {
	User string
	Pass string
}

func middlewareAuth(auths []*Auth, ctx *Context) HandlerFunc {

	return func(ctx *Context) {

	}
}
