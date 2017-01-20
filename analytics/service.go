package analytics

import (
	"fmt"
	"github.com/ONSdigital/go-ns/log"
	"net/url"
	"strconv"
	"time"
)

const pageIndexParam = "pageIndex"
const pageSizeParam = "pageSize"
const linkIndexParam = "linkIndex"
const urlParam = "url"
const termParam = "term"
const searchTypeParam = "type"
const timestampKey = "timestamp"

// Service - defines a Stats Service Interface
type Service interface {
	CaptureAnalyticsData(analytics *Model) string
}

// ServiceImpl - Implementation of the Analytics Service interface.
type ServiceImpl struct{}

// Model - Type to encapsulate Search Analytics data.
type Model struct {
	url        string
	term       string
	searchType string
	pageSize   int
	pageIndex  int
	linkIndex  int
}

// NewServiceImpl - Creates a new Analytics ServiceImpl.
func NewServiceImpl() *ServiceImpl {
	return &ServiceImpl{}
}

// GetURL - Get the URL of the search result the user clicked.
func (a *Model) GetURL() string {
	return a.url
}

// GetPageIndex - Get the Page index of the search result the use clicked.
func (a *Model) GetPageIndex() int {
	return a.pageIndex
}

// GetLinkIndex - Get the index of the link on the search result page the user clicked.
func (a *Model) GetLinkIndex() int {
	return a.linkIndex
}

// GetSearchTerm - Get the search term used.
func (a *Model) GetSearchTerm() string {
	return a.term
}

// GetSearchType - Get the type of search - search or list page etc.
func (a *Model) GetSearchType() string {
	return a.searchType
}

// GetPageSize - Get the Page size value used when searching.
func (a *Model) GetPageSize() int {
	return a.pageSize
}

// NewAnalyticsModel - Creates a new Analytics model to encapsulate the extracted values from the URL.
func NewAnalyticsModel(url *url.URL) *Model {
	return &Model{
		url:        url.Query().Get(urlParam),
		term:       url.Query().Get(termParam),
		searchType: url.Query().Get(searchTypeParam),
		pageIndex:  extractIntParam(url, pageIndexParam),
		linkIndex:  extractIntParam(url, linkIndexParam),
		pageSize:   extractIntParam(url, pageSizeParam),
	}
}

func extractIntParam(url *url.URL, name string) int {
	value := url.Query().Get(name)
	if len(value) == 0 {
		log.Debug(fmt.Sprintf("parameter '%s' was nil, default value will be used.", name), nil)
		return 0
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Debug(fmt.Sprintf("'%s' was could not be parsed to int. Default value will be used.", name), nil)
		return 0
	}
	return intValue
}

// CaptureAnalyticsData - captures the analytics values
func (s *ServiceImpl) CaptureAnalyticsData(analytics *Model) string {
	log.Debug("Capturing Search Results event.", log.Data{
		urlParam:        analytics.url,
		termParam:       analytics.term,
		pageIndexParam:  analytics.pageIndex,
		linkIndexParam:  analytics.linkIndex,
		searchTypeParam: analytics.searchType,
		pageSizeParam:   analytics.pageSize,
		timestampKey:    time.Now(),
	})
	return analytics.url
}
