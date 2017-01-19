package handlers

import (
	"github.com/ONSdigital/dp-frontend-router/analytics"
	"github.com/ONSdigital/go-ns/log"
	"net/http"
)

var searchAnalyticsService analytics.AnalyticsService = analytics.NewAnalyticsServiceImpl()

func CaptureSearchStats(w http.ResponseWriter, req *http.Request) {
	log.Debug("Handling search stats & redirect.", nil)
	searchAnalyticsService.CaptureAndRedirect(analytics.NewSearchAnalytics(req.URL), w, req)
}
