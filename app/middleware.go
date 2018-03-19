package app

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/acoshift/header"
)

func errorRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err)
				debug.PrintStack()
				// go service.SendSlackMessage(slackURL, fmt.Sprintf("biri: error; %v", err))
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func noCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func securityHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.XFrameOptions, "deny")
		w.Header().Set(header.XXSSProtection, "1; mode=block")
		w.Header().Set(header.XContentTypeOptions, "nosniff")
		h.ServeHTTP(w, r)
	})
}

func methodFilter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodPost:
			h.ServeHTTP(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}
