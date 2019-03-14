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
	"github.com/pkg/errors"
	"math/rand"
)

//Handler ...
func Handler(routesHandler map[string]http.Handler) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			path := req.URL.Path

			// ----------------------------------

			// Lifted from the requestID handler
			// will remove once I know its not required
			var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			requestID := req.Header.Get("X-Request-Id")

			if len(requestID) == 0 {
				b := make([]rune, 16)
				for i := range b {
					b[i] = letters[rand.Intn(len(letters))]
				}
				req.Header.Set("X-Request-Id", string(b))
				log.Info("allRoutes handler: RequestID needed to be set.",nil)
			}

			// ---------------------------------

			ctx := req.Context()
			logData := log.Data{
				"path":path,
				"url":req.URL.String(),
			}

			// No point calling zebedee for these paths so skip middleware
			if ok, err := regexp.MatchString(`^\/(?:datasets|filter|feedback|healthcheck)`, path); ok && err == nil {
				log.InfoCtx(ctx, "Skipping content specific handling as not relevant on this path.", logData)
				h.ServeHTTP(w, req)
				return
			}

			// We can skip handling based on content type where the url points to a known/expected file extension
			if ok, err := regexp.MatchString(`^*\.(?:xls|zip|csv|xlsx)$`, req.URL.String()); ok && err == nil {
				log.InfoCtx(ctx,"Skipping content specific handling as it's a request to download a known file extension.", logData)
				h.ServeHTTP(w, req)
				return
			}

			contentURL := config.ZebedeeURL + "/data"

			if c, err := req.Cookie(`collection`); err == nil && len(c.Value) > 0 {
				contentURL += "/" + c.Value + "?uri=" + path
			} else {
				contentURL += "?uri=" + path
			}
			logData["content_url"] = contentURL
			log.DebugCtx(ctx, "Allroutes handler: created content URL.", logData)

			//FIXME We should be doing a HEAD request but Restolino doesn't allow it - either wait for the
			// new Content API (https://github.com/ONSdigital/dp-content-api) to be in prod or update Restolino

			request, err := http.NewRequest("GET", contentURL, nil)
			if err != nil {
				log.ErrorCtx(ctx, errors.WithMessage(err, "Allroutes handler. Failed to get content"), logData)
				h.ServeHTTP(w, req)
				return
			}

			if c, err := req.Cookie(`access_token`); err == nil && len(c.Value) > 0 {
				request.Header.Set(`X-Florence-Token`, c.Value)
				logData["X-Florence-Token"] = c.Value
				log.InfoCtx(ctx, "Allroutes handler. Created new token value.", logData)
			}

			res, err := http.DefaultClient.Do(request)
			if err != nil {
				log.ErrorCtx(ctx, errors.WithMessage(err, "allRoutes handler: Error while attepting to get content"), nil)
				h.ServeHTTP(w, req)
				return
			}

			logData["status_code"] = res.StatusCode
			log.InfoCtx(ctx, "allRoutes: content request complete", logData)

			statusCode := res.StatusCode
			if statusCode >= 400 {
				log.DebugCtx(ctx, "Unexpected status code", log.Data{"statusCode": statusCode, "url": contentURL})
				io.Copy(ioutil.Discard, res.Body)
				res.Body.Close()
				h.ServeHTTP(w, req)
				return
			}

			// Use a limited reader so we dont oom the router checking for content-type
			limitReader := io.LimitReader(res.Body, int64(config.ContentTypeByteLimit+1))
			defer io.Copy(ioutil.Discard, res.Body)
			b, err := ioutil.ReadAll(limitReader)
			res.Body.Close()

			if len(b) == config.ContentTypeByteLimit+1 {
				log.InfoCtx(ctx,"Response exceeds acceptable byte limit for assessing content-type. Falling through to default handling", logData)
				h.ServeHTTP(w, req)
				return
			}

			if err != nil {
				log.ErrorCtx(ctx, errors.WithMessage(err, "allRoutes handler: error while attmepting limited read."), logData)
				h.ServeHTTP(w, req)
				return
			}

			zebResp := struct {
				Type      string `json:"type"`
				DatasetID string `json:"apiDatasetId"`
			}{}
			if err := json.Unmarshal(b, &zebResp); err != nil {
				log.ErrorCtx(ctx, errors.WithMessage(err, "allRoutes handler: error while attempting to unmarshall json."), logData)
				h.ServeHTTP(w, req)
				return
			}

			logData["zebedee response type"] = zebResp.Type
			logData["zebedee response dataset id"] = zebResp.DatasetID

			log.Info("Zebedee content response recieved.", logData)
			log.Debug(zebResp.DatasetID, nil)


			pageType := res.Header.Get("ONS-Page-Type")

			if len(zebResp.DatasetID) > 0 && zebResp.Type == "api_dataset_landing_page" {
				http.Redirect(w, req, fmt.Sprintf("/datasets/%s", zebResp.DatasetID), 302)
				return
			}

			if h, ok := routesHandler[pageType]; ok {
				log.DebugCtx(ctx, "Using handler for page type", log.Data{"pageType": pageType, "url": contentURL})
				h.ServeHTTP(w, req)
				return
			}

			h.ServeHTTP(w, req)
		})
	}
}
