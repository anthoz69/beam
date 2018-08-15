package app

import (
	"github.com/acoshift/hime"
)

func initRoutes(app *hime.App) {
	app.Routes(hime.Routes{
		"index":     "/",
		"dashboard": "/dashboard",
	})
}
