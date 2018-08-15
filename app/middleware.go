package app

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/acoshift/header"
)

func securityHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.XFrameOptions, "deny")
		w.Header().Set(header.XXSSProtection, "1; mode=block")
		w.Header().Set(header.XContentTypeOptions, "nosniff")
		h.ServeHTTP(w, r)
	})
}

// ErrorRecovery middleware
func ErrorRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err)
				debug.PrintStack()
			}
		}()
		h.ServeHTTP(w, r)
	})
}
