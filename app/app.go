package app

import (
	"net/http"

	"github.com/acoshift/hime"
	"github.com/acoshift/httprouter"
	"github.com/acoshift/middleware"
	"github.com/acoshift/webstatic"
)

// New is init all config
func New(app *hime.App) http.Handler {

	initRoutes(app)
	initTemplates(app)

	mux := http.NewServeMux()
	mux.Handle("/-/", http.StripPrefix("/-", webstatic.New(webstatic.Config{
		Dir:          "assets",
		CacheControl: "public, max-age=3600",
	})))

	m := httprouter.New()
	m.HandleMethodNotAllowed = false
	m.NotFound = hime.Handler(notFoundHandler)

	m.Get(app.Route("index"), hime.Handler(indexHandler))
	m.Get("/profile/:id", hime.Handler(profileHandler))

	mux.Handle("/", m)

	return middleware.Chain(
		ErrorRecovery,
		securityHeaders,
	)(mux)
}

func notFoundHandler(ctx *hime.Context) error {
	return ctx.View("app/notfound", page(ctx))
}

func indexHandler(ctx *hime.Context) error {
	p := page(ctx)
	return ctx.View("app/index", p)
}

func profileHandler(ctx *hime.Context) error {
	p := page(ctx)
	p["ID"] = httprouter.GetParams(ctx).ByName("id")
	return ctx.View("app/profile", p)
}
