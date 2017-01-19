package handlers

import (
	"github.com/ONSdigital/dp-frontend-router/statistics"
	"github.com/ONSdigital/go-ns/log"
	"net/http"
)

var searchStatsService search.AnalyticsService

func CaptureSearchStats(w http.ResponseWriter, req *http.Request) {
	log.Debug("Handling search stats & redirect.", nil)
	searchStats := search.NewSearchAnalytics(req.URL)
	searchStatsService.CaptureAndRedirect(searchStats, w, req)
}

// SetStatsService - Set the StatsService implementation to use.
func SetAnalyticsService(service search.AnalyticsService) {
	searchStatsService = service
}
