package analytics

import (
	"fmt"
	"github.com/ONSdigital/go-ns/log"
	"net/url"
	"strconv"
	"net/http"
	"errors"
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
	CaptureAnalyticsData(r *http.Request) (string, error)
}

// ServiceImpl - Implementation of the Analytics Service interface.
type ServiceImpl struct{}

// NewServiceImpl - Creates a new Analytics ServiceImpl.
func NewServiceImpl() *ServiceImpl {
	return &ServiceImpl{}
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
func (s *ServiceImpl) CaptureAnalyticsData(r *http.Request) (string, error) {
	url := r.URL.Query().Get(urlParam)
	if len(url) == 0 {
		log.Error(errors.New("Failed to redirect to search results as parameter URL was missing."), nil)
		return "", errors.New("400: URL is a mandatory parameter.")
	}

	term := r.URL.Query().Get(termParam)
	searchType := r.URL.Query().Get(searchTypeParam)
	pageIndex := extractIntParam(r.URL, pageIndexParam)
	linkIndex := extractIntParam(r.URL, linkIndexParam)
	pageSize := extractIntParam(r.URL, pageSizeParam)

	// TODO implement.
	log.Debug("CaptureAnalyticsData", log.Data{
		urlParam:        url,
		termParam:       term,
		searchTypeParam: searchType,
		pageIndexParam:  pageIndex,
		linkIndexParam:  linkIndex,
		pageSizeParam:   pageSize,
	})
	return url, nil
}
