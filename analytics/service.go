package analytics

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const pageIndexParam = "pageIndex"
const pageSizeParam = "pageSize"
const linkIndexParam = "linkIndex"
const urlParam = "url"
const termParam = "term"
const searchTypeParam = "type"
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
	vars := mux.Vars(r)
	data := vars["data"]

	token, err := jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.redirectSecret), nil
	})

	if err != nil {
		return "", errors.Wrap(err, "Invalid redirect data")
	}

	var claims jwt.MapClaims
	var ok bool
	if claims, ok = token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return "", errors.New("error validating token")
	}

	url := getStringClaim(claims, "uri")
	term := getStringClaim(claims, "term")
	listType := getStringClaim(claims, "listType")
	pageIndex := getFloat64Claim(claims, "page")
	linkIndex := getFloat64Claim(claims, "index")
	pageSize := getFloat64Claim(claims, "pageSize")

	gaID := getCookieValue(r, "_ga")
	gID := getCookieValue(r, "_gid")

	if url == "" {
		return "", errors.New("URL is a mandatory parameter")
	}

	logData := log.Data{
		urlParam:        url,
		termParam:       term,
		searchTypeParam: listType,
		pageIndexParam:  pageIndex,
		linkIndexParam:  linkIndex,
		pageSizeParam:   pageSize,
		gaIDParam:       gaID,
		gIDParam:        gID,
	}
	log.Info(r.Context(), "search analytics data", logData)

	if s.backend != nil {
		s.backend.Store(r, url, term, listType, gaID, gID, pageIndex, linkIndex, pageSize)
	}

	return url, nil
}

func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}

func getFloat64Claim(claims jwt.MapClaims, key string) float64 {
	if val, ok := claims[key].(float64); ok {
		return val
	}
	return 0
}

func getCookieValue(r *http.Request, cookieName string) string {
	if c, err := r.Cookie(cookieName); err == nil && c != nil {
		return c.Value
	}
	return ""
}
