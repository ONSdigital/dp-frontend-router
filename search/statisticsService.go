package search

import (
	"fmt"
	"github.com/ONSdigital/go-ns/log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const pageIndexParam = "pageIndex"
const linkIndexParam = "linkIndex"
const urlParam = "url"
const termParam = "term"

type HttpRedirect func(w http.ResponseWriter, r *http.Request, urlStr string, code int)

var redirector HttpRedirect = http.Redirect

type SearchStatsService interface {
	CaptureAndRedirect(searchStats *searchStatistics, w http.ResponseWriter, req *http.Request)
}

type SearchStatsServiceImpl struct {
}

type searchStatistics struct {
	url       string
	term      string
	pageIndex int
	linkIndex int
	timestamp time.Time
}

func (s *searchStatistics) GetURL() string {
	return s.url
}

func (s *searchStatistics) GetPageIndex() int {
	return s.pageIndex
}

func (s *searchStatistics) GetLinkIndex() int {
	return s.linkIndex
}

func (s *searchStatistics) GetSearchTerm() string {
	return s.term
}

func (s *searchStatistics) GetTimestamp() time.Time {
	return s.timestamp
}

func NewSearchStatistics(url *url.URL) *searchStatistics {
	return &searchStatistics{
		url:       url.Query().Get(urlParam),
		pageIndex: extractIntParam(url, pageIndexParam),
		linkIndex: extractIntParam(url, linkIndexParam),
		term:      url.Query().Get(termParam),
		timestamp: time.Now(),
	}
}

func extractIntParam(url *url.URL, name string) int {
	value := url.Query().Get(name)
	if len(value) == 0 {
		log.Debug(fmt.Sprintf("parameter '%s' was nil", name), nil)
		return 0
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Debug(fmt.Sprintf("'%s' was could not be parsed to int.", name), nil)
		return 0
	}
	return intValue
}

func (s *searchStatistics) IsValid() bool {
	return len(s.url) > 0 && len(s.term) > 0 && s.pageIndex > 0 && s.linkIndex > 0
}

func (s *SearchStatsServiceImpl) CaptureAndRedirect(searchStats *searchStatistics, w http.ResponseWriter, req *http.Request) {
	log.Debug("Capturing Search Results event.", log.Data{
		urlParam:       searchStats.url,
		termParam:      searchStats.term,
		pageIndexParam: searchStats.pageIndex,
		linkIndexParam: searchStats.linkIndex,
		"timestamp":    searchStats.timestamp,
	})
	redirector(w, req, searchStats.GetURL(), http.StatusTemporaryRedirect)
}
