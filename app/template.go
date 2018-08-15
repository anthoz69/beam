package app

import (
	"html/template"
	"net/http"

	"github.com/acoshift/header"
	"github.com/acoshift/hime"
)

func initTemplates(app *hime.App) {
	app.Template().
		Dir("template").
		Root("root").
		Minify().
		Funcs(template.FuncMap{
			"static": func(s string) string {
				var assets string
				switch s {
				case "css":
					assets = "style.css"
				case "js":
					assets = "app.js"
				}
				return "/-/" + assets
			},
		}).
		Component(
			"main.tmpl",
		).
		Parse("app/notfound", "_layout/app.tmpl", "app/notfound.tmpl").
		Parse("app/index", "_layout/app.tmpl", "app/index.tmpl").
		Parse("app/profile", "_layout/app.tmpl", "app/profile.tmpl")

	app.BeforeRender(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(header.CacheControl, "no-cache, no-store, must-revalidate")
			h.ServeHTTP(w, r)
		})
	})

}
