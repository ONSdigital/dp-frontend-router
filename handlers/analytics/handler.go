package analytics

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-router/analytics"
	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/go-ns/log"
)

type httpRedirector func(w http.ResponseWriter, r *http.Request, urlStr string, code int)

type searchHandler struct {
	service    analytics.Service
	redirector httpRedirector
}

// NewSearchHandler creates a new search handler
func NewSearchHandler() (http.Handler, error) {
	var b analytics.ServiceBackend
	var err error

	if len(config.SQSAnalyticsURL) > 0 {
		b, err = analytics.NewSQSBackend(config.SQSAnalyticsURL)
		if err != nil {
			return nil, err
		}
	}

	sh := &searchHandler{
		service:    analytics.NewServiceImpl(b),
		redirector: http.Redirect,
	}
	return sh, nil
}

// HandleSearch - http Handler func for dealing with Babbage Search requests. Captures search analytics data and redirects
// the user to the requested resource.
func (sh searchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("Handling search stats & redirect.", nil)
	redirectURL, err := sh.service.CaptureAnalyticsData(r)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	sh.redirector(w, r, redirectURL, http.StatusTemporaryRedirect)
}