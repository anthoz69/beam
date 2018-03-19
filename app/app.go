package app

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"github.com/acoshift/csrf"
	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/acoshift/methodmux"
	"github.com/acoshift/middleware"
	"github.com/acoshift/session"
	redisstore "github.com/acoshift/session/store/redis"
	"github.com/garyburd/redigo/redis"
	"github.com/polycurve/servestatic"
)

// Config is app's config
type Config struct {
	DB            *sql.DB
	SessionKey    []byte
	SessionSecret []byte
	RedisPool     *redis.Pool
	RedisPrefix   string
}

var (
	static      = make(map[string]string)
	db          *sql.DB
	redisPool   *redis.Pool
	redisPrefix string
)

// New creates new app's factory
func New(c *Config) hime.HandlerFactory {
	db = c.DB
	redisPool = c.RedisPool
	redisPrefix = c.RedisPrefix

	return func(app hime.App) http.Handler {
		app.
			TemplateFuncs(template.FuncMap{
				"static": func(s string) string {
					p := static[s]
					if p == "" {
						return "/-/" + s
					}
					return p
				},
			}).
			TemplateDir("template").
			TemplateRoot("root").
			// Component()
			Template("app/index", "_layout/app.tmpl", "app/index.tmpl").
			BeforeRender(func(h http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// getSession(r.Context()).Flash().Clear()
					w.Header().Set(header.CacheControl, "no-cache, no-store, must-revalidate")
					h.ServeHTTP(w, r)
				})
			}).
			Routes(hime.Routes{
				"index": "/",
			})

		mux := http.NewServeMux()
		mux.Handle("/healthz", http.HandlerFunc(healthzHandler))

		// only use on development
		mux.Handle("/-/", http.StripPrefix("/-", servestatic.New("assets")))

		m := http.NewServeMux()
		m.Handle("/", methodmux.Get(
			hime.H(indexGetHandler),
		))

		mux.Handle("/", middleware.Chain(
			session.Middleware(session.Config{
				Store: redisstore.New(redisstore.Config{
					Pool:   redisPool,
					Prefix: redisPrefix,
				}),
				HTTPOnly: true,
				Secure:   session.PreferSecure,
				Proxy:    true,
				MaxAge:   60 * 24 * time.Hour,
				Path:     "/",
				Rolling:  true,
				Keys:     [][]byte{c.SessionKey},
				Secret:   c.SessionSecret,
				SameSite: session.SameSiteLax,
			}),
		)(m))

		return middleware.Chain(
			errorRecovery,
			noCORS,
			securityHeaders,
			methodFilter,
			csrf.New(csrf.Config{Origins: []string{"http://localhost:8000"}}),
		)(mux)

	}
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if err := db.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c := redisPool.Get()
	defer c.Close()
	if _, err := c.Do("PING"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func notFoundHandler(ctx hime.Context) hime.Result {
	return ctx.Redirect("/")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
