package analytics

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/go-ns/log"
	"github.com/dgrijalva/jwt-go"
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

// CaptureAnalyticsData - captures the analytics values
func (s *ServiceImpl) CaptureAnalyticsData(r *http.Request) (string, error) {
	data := r.URL.Query().Get(":data")

	token, err := jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.RedirectSecret), nil
	})

	if err != nil {
		log.ErrorR(r, err, nil)
		return "", errors.New("Invalid redirect data")
	}

	log.DebugR(r, "token", log.Data{"token": token})

	var url, term, listType string
	var pageIndex, linkIndex, pageSize float64

	var claims jwt.MapClaims
	var ok bool
	if claims, ok = token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		log.ErrorR(r, errors.New("error validating token"), nil)
		return "", errors.New("error validating token")
	}

	if s, ok := claims["uri"].(string); ok {
		url = s
	}
	if s, ok := claims["term"].(string); ok {
		term = s
	}
	if s, ok := claims["listType"].(string); ok {
		listType = s
	}
	if s, ok := claims["page"].(float64); ok {
		pageIndex = s
	}
	if s, ok := claims["index"].(float64); ok {
		linkIndex = s
	}
	if s, ok := claims["pageSize"].(float64); ok {
		pageSize = s
	}

	if len(url) == 0 {
		log.ErrorR(r, errors.New("Failed to redirect to search results as parameter URL was missing."), nil)
		return "", errors.New("400: URL is a mandatory parameter.")
	}

	// TODO implement.
	log.DebugR(r, "CaptureAnalyticsData", log.Data{
		urlParam:        url,
		termParam:       term,
		searchTypeParam: listType,
		pageIndexParam:  pageIndex,
		linkIndexParam:  linkIndex,
		pageSizeParam:   pageSize,
	})
	return url, nil
}
