package analytics

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/ONSdigital/log.go/log"
	jwt "github.com/dgrijalva/jwt-go"
)

const pageIndexParam = "pageIndex"
const pageSizeParam = "pageSize"
const linkIndexParam = "linkIndex"
const urlParam = "url"
const termParam = "term"
const searchTypeParam = "type"
const timestampKey = "timestamp"
const gaIDParam = "ga"
const gIDParam = "gid"

// Service - defines a Stats Service Interface
type Service interface {
	CaptureAnalyticsData(r *http.Request) (string, error)
}

// ServiceBackend is used to store data output by the analytics service
type ServiceBackend interface {
	Store(req *http.Request, url, term, listType, gaID string, gID string, pageIndex, linkIndex, pageSize float64)
}

// ServiceImpl - Implementation of the Analytics Service interface.
type ServiceImpl struct {
	backend        ServiceBackend
	redirectSecret string
}

// NewServiceImpl - Creates a new Analytics ServiceImpl.
func NewServiceImpl(backend ServiceBackend, redirectSecret string) *ServiceImpl {
	return &ServiceImpl{backend, redirectSecret}
}

// CaptureAnalyticsData - captures the analytics values
func (s *ServiceImpl) CaptureAnalyticsData(r *http.Request) (string, error) {
	data := r.URL.Query().Get(":data")
	fmt.Println("data", data)

	token, err := jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.redirectSecret), nil
	})

	if err != nil {
		return "", errors.Wrap(err, "Invalid redirect data")
	}

	log.Event(r.Context(), "jwt token", log.INFO, log.Data{"token": token})

	var url, term, listType, gaID, gID string
	var pageIndex, linkIndex, pageSize float64

	var claims jwt.MapClaims
	var ok bool
	if claims, ok = token.Claims.(jwt.MapClaims); !ok || !token.Valid {
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

	if c, err := r.Cookie("_ga"); err == nil && c != nil {
		// 2 year expiration cookie (_ga)
		gaID = c.Value
	}

	if c, err := r.Cookie("_gid"); err == nil && c != nil {
		// 24 hour expiration cookie (_gid)
		gID = c.Value
	}

	if len(url) == 0 {
		return "", errors.New("url is a mandatory parameter")
	}

	// FIXME do we want to log as well as store in backend?
	log.Event(r.Context(), "search analytics data", log.INFO, log.Data{
		urlParam:        url,
		termParam:       term,
		searchTypeParam: listType,
		pageIndexParam:  pageIndex,
		linkIndexParam:  linkIndex,
		pageSizeParam:   pageSize,
		gaIDParam:       gaID,
		gIDParam:        gID,
	})

	if s.backend != nil {
		s.backend.Store(r, url, term, listType, gaID, gID, pageIndex, linkIndex, pageSize)
	}

	return url, nil
}
