package analytics

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-frontend-router/analytics"
	"github.com/ONSdigital/log.go/v2/log"
)

type httpRedirector func(w http.ResponseWriter, r *http.Request, urlStr string, code int)

type searchHandler struct {
	service    analytics.Service
	redirector httpRedirector
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(ctx context.Context, sqSanalyticsURL, redirectSecret string) (http.Handler, error) {
	var b analytics.ServiceBackend
	var err error

	if sqSanalyticsURL != "" {
		b, err = analytics.NewSQSBackend(ctx, sqSanalyticsURL)
		if err != nil {
			return nil, err
		}
	}

	sh := &searchHandler{
		service:    analytics.NewServiceImpl(b, redirectSecret),
		redirector: http.Redirect,
	}
	return sh, nil
}

// HandleSearch - http Handler func for dealing with Babbage Search requests. Captures search analytics data and redirects
// the user to the requested resource.
func (sh searchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Context(), "capturing search analytics data")
	redirectURL, err := sh.service.CaptureAnalyticsData(r)

	if err != nil {
		log.Error(r.Context(), "error capturing analytics data", err)
		w.WriteHeader(400)
		// FIXME probably want to display a better error page than this
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Error(r.Context(), "error writing response", err)
		}
		return
	}

	sh.redirector(w, r, redirectURL, http.StatusTemporaryRedirect)
}
