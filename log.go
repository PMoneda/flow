package gonnie

func log(ctx *Context, s ...string) {
	if len(s) > 0 {
		ctx.log.Push(s[0])
	}
}
