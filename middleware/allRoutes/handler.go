package allRoutes

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/go-ns/log"
)

//Handler ...
func Handler(routesHandler map[string]http.Handler) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			//FIXME We should be doing a HEAD request but Restolino doesn't allow it - either wait for the
			// new Content API (https://github.com/ONSdigital/dp-content-api) to be in prod or update Restolino
			contentURL := config.ZebedeeURL + "/data?uri=" + req.URL.Path
			request, err := http.NewRequest("GET", contentURL, nil)
			if err != nil {
				log.ErrorR(req, err, nil)
				h.ServeHTTP(w, req)
				return
			}

			if c, err := req.Cookie(`access_token`); err == nil && len(c.Value) > 0 {
				request.Header.Set(`X-Florence-Token`, c.Value)
				log.Debug(c.Value, nil)
			}

			res, err := http.DefaultClient.Do(request)
			if err != nil {
				log.ErrorR(req, err, nil)
				h.ServeHTTP(w, req)
				return
			}

			statusCode := res.StatusCode
			if statusCode >= 400 {
				log.DebugR(req, "Unexpected status code", log.Data{"statusCode": statusCode, "url": contentURL})
				h.ServeHTTP(w, req)
				return
			}

			io.Copy(ioutil.Discard, res.Body)
			res.Body.Close()

			pageType := res.Header.Get("ONS-Page-Type")

			if h, ok := routesHandler[pageType]; ok {
				log.DebugR(req, "Using handler for page type", log.Data{"pageType": pageType, "url": contentURL})
				h.ServeHTTP(w, req)
				return
			}

			h.ServeHTTP(w, req)
		})
	}
}
