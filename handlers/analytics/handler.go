package analytics

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-router/analytics"
	"github.com/ONSdigital/log.go/log"
)

type httpRedirector func(w http.ResponseWriter, r *http.Request, urlStr string, code int)

type searchHandler struct {
	service    analytics.Service
	redirector httpRedirector
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(SQSanalyticsURL, RedirectSecret string) (http.Handler, error) {
	var b analytics.ServiceBackend
	var err error

	if len(SQSanalyticsURL) > 0 {
		b, err = analytics.NewSQSBackend(SQSanalyticsURL)
		if err != nil {
			return nil, err
		}
	}

	sh := &searchHandler{
		service:    analytics.NewServiceImpl(b, RedirectSecret),
		redirector: http.Redirect,
	}
	return sh, nil
}

// HandleSearch - http Handler func for dealing with Babbage Search requests. Captures search analytics data and redirects
// the user to the requested resource.
func (sh searchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Event(r.Context(), "capturing search analytics data", log.INFO)
	redirectURL, err := sh.service.CaptureAnalyticsData(r)

	if err != nil {
		log.Event(r.Context(), "error capturing analytics data", log.ERROR, log.Error(err))
		w.WriteHeader(400)
		// FIXME probably want to display a better error page than this
		w.Write([]byte(err.Error()))
		return
	}

	sh.redirector(w, r, redirectURL, http.StatusTemporaryRedirect)
}
