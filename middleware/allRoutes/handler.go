package allRoutes

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/go-ns/log"
)

//Handler ...
func Handler(routesHandler map[string]http.Handler) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			path := req.URL.Path

			// No point calling zebedee for these paths so skip middleware
			if ok, err := regexp.MatchString(`^\/(datasets|filter|feedback|healthcheck).*$`, path); ok && err == nil {
				log.Info("Skipping content specific handling as not relevant on this path.", log.Data{"url": path})
				h.ServeHTTP(w, req)
				return
			}

			// We can skip handling based on content type where the url points to a known/expected file extension
			if ok, err := regexp.MatchString(`^*(xls|zip|csv|xslx)$`, req.URL.String()); ok && err == nil {
				log.Info("Skipping content specific handling as it's a request to download a known file extension.", log.Data{"url": req.URL.String()})
				h.ServeHTTP(w, req)
				return
			}

			contentURL := config.ZebedeeURL + "/data"

			if c, err := req.Cookie(`collection`); err == nil && len(c.Value) > 0 {
				contentURL += "/" + c.Value + "?uri=" + path
			} else {
				contentURL += "?uri=" + path
			}
			log.Debug(contentURL, nil)

			//FIXME We should be doing a HEAD request but Restolino doesn't allow it - either wait for the
			// new Content API (https://github.com/ONSdigital/dp-content-api) to be in prod or update Restolino

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
				io.Copy(ioutil.Discard, res.Body)
				res.Body.Close()
				h.ServeHTTP(w, req)
				return
			}

			// Use a limited reader so we dont oom the router checking for conten-type
			limitReader := io.LimitReader(res.Body, int64(config.ContentTypeByteLimit+1))
			b, err := ioutil.ReadAll(limitReader)
			res.Body.Close()

			if len(b) == config.ContentTypeByteLimit+1 {
				log.Info("Response exceeds acceptable byte limit for assessing content-type. Falling through to default handling", nil)
				h.ServeHTTP(w, req)
				return
			}

			if err != nil {
				log.ErrorR(req, err, nil)
				h.ServeHTTP(w, req)
				return
			}

			zebResp := struct {
				Type      string `json:"type"`
				DatasetID string `json:"apiDatasetId"`
			}{}
			if err := json.Unmarshal(b, &zebResp); err != nil {
				log.ErrorR(req, err, nil)
				h.ServeHTTP(w, req)
				return
			}

			log.Debug(zebResp.Type, nil)
			log.Debug(zebResp.DatasetID, nil)

			pageType := res.Header.Get("ONS-Page-Type")

			if len(zebResp.DatasetID) > 0 && zebResp.Type == "api_dataset_landing_page" {
				http.Redirect(w, req, fmt.Sprintf("/datasets/%s", zebResp.DatasetID), 302)
				return
			}

			if h, ok := routesHandler[pageType]; ok {
				log.DebugR(req, "Using handler for page type", log.Data{"pageType": pageType, "url": contentURL})
				h.ServeHTTP(w, req)
				return
			}

			h.ServeHTTP(w, req)
		})
	}
}
