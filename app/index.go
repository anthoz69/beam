package app

import "github.com/acoshift/hime"

func indexGetHandler(ctx hime.Context) hime.Result {
	p := page()
	return ctx.View("app/index", p)
}
