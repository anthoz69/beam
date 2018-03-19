package app

import "github.com/acoshift/hime"

func indexGetHandler(ctx hime.Context) hime.Result {
	return ctx.View("app/index", "data")
}
