package relcal

import (
	"net/http"
)

func Handler(useNewReleaseCalendar bool, newHandler, oldHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serve(w, r, newHandler, oldHandler, useNewReleaseCalendar)
	})
}

func serve(w http.ResponseWriter, req *http.Request, newHandler, oldHandler http.Handler, useNewReleaseCalendar bool) {
	if useNewReleaseCalendar {
		newHandler.ServeHTTP(w, req)
		return
	}

	oldHandler.ServeHTTP(w, req)
}
