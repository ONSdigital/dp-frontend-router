package allRoutes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	client "github.com/ONSdigital/dp-api-clients-go/zebedee"
	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/log.go/log"
)

// HeaderOnsPageType is the header name that defines the handler that will be used by the Middleware
const HeaderOnsPageType = "ONS-Page-Type"

//Handler implements the middleware for dp-frontend-router. It sets the locale code, obtains the necessary cookies for the request path and access_token,
// authenticates with Zebedee if required,  and obtains the "ONS-Page-Type" header to use the handler for the page type, if present.
func Handler(routesHandler map[string]http.Handler, zebedeeClient *client.Client, cfg *config.Config) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			path := req.URL.Path

			// Populate context here with language
			req = common.SetLocaleCode(req)

			// No point calling zebedee for these paths so skip middleware
			if ok, err := regexp.MatchString(`^\/(?:datasets|filter|feedback|healthcheck)`, path); ok && err == nil {
				log.Event(req.Context(), "Skipping content specific handling as not relevant on this path.", log.INFO, log.Data{"url": path})
				h.ServeHTTP(w, req)
				return
			}

			// We can skip handling based on content type where the url points to a known/expected file extension
			if ok, err := regexp.MatchString(`^*\.(?:xls|zip|csv|xslx)$`, req.URL.String()); ok && err == nil {
				log.Event(req.Context(), "Skipping content specific handling as it's a request to download a known file extension.",
					log.INFO, log.Data{"url": req.URL.String()})
				h.ServeHTTP(w, req)
				return
			}

			// Construct contentPath with any colletion if present in cookie
			contentPath := "/data"
			if c, err := req.Cookie(`collection`); err == nil && len(c.Value) > 0 {
				contentPath += "/" + c.Value + "?uri=" + path
			} else {
				contentPath += "?uri=" + path
			}
			log.Event(req.Context(), "generated from 'collection' cookie", log.INFO, log.Data{"contentPath": contentPath})

			//FIXME We should be doing a HEAD request but Restolino doesn't allow it - either wait for the
			// new Content API (https://github.com/ONSdigital/dp-content-api) to be in prod or update Restolino
			/// Update: Is this still needed when using the Zebedee client?

			// Obtain access_token from cookie
			userAccessToken := ""
			c, err := req.Cookie(`access_token`)
			if err != nil {
				log.Event(req.Context(), "Cookie error", log.WARN, log.Error(err))
			} else if len(c.Value) > 0 {
				userAccessToken = c.Value
				log.Event(req.Context(), "Obtained access_token Cookie", log.INFO, log.Data{"value": c.Value})
			}

			// Do the GET call using Zebedee Client and providing any access_token from cookie
			b, headers, err := zebedeeClient.GetWithHeaders(req.Context(), userAccessToken, contentPath)
			if err != nil {
				log.Event(req.Context(), "Zebedee GET error", log.ERROR, log.Error(err))
				h.ServeHTTP(w, req)
				return
			}

			if len(b) > cfg.ContentTypeByteLimit {
				log.Event(req.Context(), "Response exceeds acceptable byte limit for assessing content-type. Falling through to default handling", log.WARN)
				h.ServeHTTP(w, req)
				return
			}

			zebResp := struct {
				Type      string `json:"type"`
				DatasetID string `json:"apiDatasetId"`
			}{}
			if err := json.Unmarshal(b, &zebResp); err != nil {
				log.Event(req.Context(), "json unmarshal error", log.ERROR, log.Error(err))
				h.ServeHTTP(w, req)
				return
			}

			log.Event(req.Context(), "zebedee response", log.INFO, log.Data{"type": zebResp.Type, "datasetID": zebResp.DatasetID})

			pageType := headers.Get(HeaderOnsPageType)

			if len(zebResp.DatasetID) > 0 && zebResp.Type == "api_dataset_landing_page" {
				http.Redirect(w, req, fmt.Sprintf("/datasets/%s", zebResp.DatasetID), 302)
				return
			}

			if h, ok := routesHandler[pageType]; ok {
				log.Event(req.Context(), "Using handler for page type", log.INFO, log.Data{"pageType": pageType, "path": contentPath})
				h.ServeHTTP(w, req)
				return
			}

			h.ServeHTTP(w, req)
		})
	}
}
