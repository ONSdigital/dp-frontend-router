package handlers

import (
	"github.com/ONSdigital/dp-frontend-router/search"
	"github.com/ONSdigital/go-ns/log"
	"net/http"
)

var SearchStatsService search.SearchStatsService

func CaptureSearchStats(w http.ResponseWriter, req *http.Request) {
	log.Debug("Handling search stats & redirect.", nil)
	searchStats := search.NewSearchStatistics(req.URL)
	SearchStatsService.CaptureAndRedirect(searchStats, w, req)
}
