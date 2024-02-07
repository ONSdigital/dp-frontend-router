package relcal

import (
	"net/http"
)

func Handler(useNewReleaseCalendar bool, newHandler, oldHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serve(w, r, newHandler, oldHandler, useNewReleaseCalendar)
	})
}

func serve(w http.ResponseWriter, req *http.Request, n, o http.Handler, useNewReleaseCalendar bool) {
	if useNewReleaseCalendar {
		n.ServeHTTP(w, req)
		return
	}

	o.ServeHTTP(w, req)
}
