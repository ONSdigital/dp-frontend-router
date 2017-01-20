package analytics

import (
	"github.com/ONSdigital/dp-frontend-router/analytics"
	"github.com/ONSdigital/go-ns/log"
	"net/http"
)

type HttpRedirector func(w http.ResponseWriter, r *http.Request, urlStr string, code int)

var analyticsService analytics.Service = analytics.NewServiceImpl()
var redirector HttpRedirector = http.Redirect

// HandleSearch - http Handler func for dealing with Babbage Search requests. Captures search analytics data and redirects
// the user to the requested resource.
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	log.Debug("Handling search stats & redirect.", nil)
	redirectURL := analyticsService.CaptureAnalyticsData(analytics.NewAnalyticsModel(r.URL))
	redirector(w, r, redirectURL, http.StatusTemporaryRedirect)
}
