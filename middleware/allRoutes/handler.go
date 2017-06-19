package allRoutes

import (
	"net/http"

	"github.com/ONSdigital/go-ns/log"
)

//Handler ...
func Handler(routesHandler map[string]http.Handler) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			//FIXME This breaks Zebedee CMS because we need to send the access token header with the req
			//FIXME We should be doing a HEAD request but Restolino doesn't allow it - either wait for the
			// new Content API (https://github.com/ONSdigital/dp-content-api) to be in prod or update Restolino
			res, err := http.Get("http://localhost:8082/data?uri=" + req.URL.Path)
			if err != nil {
				log.Error(err, nil)
			}

			res.Body.Close()

			pageType := res.Header.Get("ONS-Page-Type")

			if h, ok := routesHandler[pageType]; ok {
				log.Debug("Using handler for page type "+pageType, nil)
				h.ServeHTTP(w, req)
				return
			}

			h.ServeHTTP(w, req)
		})
	}
}
