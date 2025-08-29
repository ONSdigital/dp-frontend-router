//nolint:revive,stylecheck // ignore, Package name "allRoutes" is kept for compatibility.
package allRoutes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	dprequest "github.com/ONSdigital/dp-net/v3/request"
	"github.com/ONSdigital/log.go/v2/log"
)

// HeaderOnsPageType is the header name that defines the handler that will be used by the Middleware
const HeaderOnsPageType = "ONS-Page-Type" // NOTE: when using the http method Add and Get, this returns the canonical format: "Ons-Page-Type"

//go:generate moq -out allroutestest/zebedeeclient.go -pkg allroutestest . ZebedeeClient
type ZebedeeClient interface {
	GetWithHeaders(ctx context.Context, userAccessToken, path string) ([]byte, http.Header, error)
}

// Handler implements the middleware for dp-frontend-router. It sets the locale code, obtains the necessary cookies for the request path and access_token,
// authenticates with Zebedee if required,  and obtains the "ONS-Page-Type" header to use the handler for the page type, if present.
func Handler(routesHandler map[string]http.Handler, zebedeeClient ZebedeeClient, contentTypeByteLimit int) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			path := req.URL.Path

			// Populate context here with language
			req = dprequest.SetLocaleCode(req)

			// Construct contentPath with any collection if present in cookie
			contentPath := constructContentPath(req, path)

			// FIXME We should be doing a HEAD request but Restolino doesn't allow it - either wait for the
			// new Content API (https://github.com/ONSdigital/dp-content-api) to be in prod or update Restolino
			/// Update: Is this still needed when using the Zebedee client?

			// Obtain access_token from cookie
			userAccessToken := ""
			if c, err := req.Cookie(`access_token`); err == nil && c.Value != "" {
				userAccessToken = c.Value
				log.Info(req.Context(), "obtained access_token Cookie")
			}

			// Do the GET call using Zebedee Client and providing any access_token from cookie
			b, headers, err := zebedeeClient.GetWithHeaders(req.Context(), userAccessToken, contentPath)
			if err != nil {
				// intentionally log as info with the error in log.data to prevent the full stack trace being logged as zebedee 404's are common
				log.Info(req.Context(), "zebedee GET failed", log.Data{"error": err.Error(), "path": path})
				h.ServeHTTP(w, req)
				return
			}

			if len(b) > contentTypeByteLimit {
				log.Warn(req.Context(), "response exceeds acceptable byte limit for assessing content-type. Falling through to default handling")
				h.ServeHTTP(w, req)
				return
			}

			var zebResp struct {
				Type      string `json:"type"`
				DatasetID string `json:"apiDatasetId"`
			}

			if err := json.Unmarshal(b, &zebResp); err != nil {
				log.Error(req.Context(), "json unmarshal error", err)
				h.ServeHTTP(w, req)
				return
			}

			log.Info(req.Context(), "zebedee response", log.Data{"type": zebResp.Type, "datasetID": zebResp.DatasetID})

			pageType := headers.Get(HeaderOnsPageType)

			if zebResp.DatasetID != "" && zebResp.Type == "api_dataset_landing_page" {
				http.Redirect(w, req, fmt.Sprintf("/datasets/%s", zebResp.DatasetID), http.StatusFound)
				return
			}

			if routesH, ok := routesHandler[pageType]; ok {
				log.Info(req.Context(), "using handler for page type", log.Data{"pageType": pageType, "path": contentPath})
				req.Header.Add(HeaderOnsPageType, pageType)
				routesH.ServeHTTP(w, req)
				return
			}

			h.ServeHTTP(w, req)
		})
	}
}

func constructContentPath(req *http.Request, path string) string {
	contentPath := "/data"
	if c, err := req.Cookie(`collection`); err == nil && c.Value != "" {
		contentPath += "/" + c.Value + "?uri=" + path
		log.Info(req.Context(), "generated from 'collection' cookie", log.Data{"contentPath": contentPath})
	} else {
		contentPath += "?uri=" + path
	}
	return contentPath
}
