package profiler

import (
	"errors"
	"github.com/ONSdigital/log.go/log"
	"net/http"
)

// profileMiddleware to validate auth token before accessing endpoint
func Middleware(token string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			pprofToken := req.Header.Get("Authorization")
			if pprofToken == "Bearer " || pprofToken != "Bearer "+token {
				log.Event(ctx, "invalid pprof auth token", log.ERROR, log.Error(errors.New("invalid pprof auth token")))
				w.WriteHeader(403)
				return
			}

			log.Event(ctx, "accessing profiling endpoint", log.INFO)
			h.ServeHTTP(w, req)
		})
	}
}
