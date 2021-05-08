package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
	"strings"
)

var knownBabbageEndpoints = []string{
	"/chartconfig",
	"/chartimage",
	"/generator",
	"/timeseriestool",
	"/search",
	"/file",
	"/resource",
	"/file",
	"/calendar",
	"/chart",
	"/embed",
	"/export",
	"/hash",
}

var knownBabbageEndpointPrefixes = []string{
	"/visualisations/",
	"/ons/",
}

// IsKnownBabbageEndpoint returns true if the given path matches a known babbage endpoint
func IsKnownBabbageEndpoint(path string) bool {
	for _, endpoint := range knownBabbageEndpoints {
		if path == endpoint {
			return true
		}
	}
	for _, endpoint := range knownBabbageEndpointPrefixes {
		if strings.HasPrefix(path, endpoint) {
			return true
		}
	}
	return false
}

// isKnownBabbageEndpointMatcher is a mux MatcherFunc implementation, allowing routes to be matched on known babbage endpoints
func isKnownBabbageEndpointMatcher(request *http.Request, match *mux.RouteMatch) bool {
	return IsKnownBabbageEndpoint(request.URL.Path)
}

// HasFileExt returns true if the given path has a file extension
func HasFileExt(path string) bool {
	return len(filepath.Ext(path)) > 0
}

// hasFileExtMatcher is a mux MatcherFunc implementation, allowing routes to be matched on having a file extension
func hasFileExtMatcher(request *http.Request, match *mux.RouteMatch) bool {
	return HasFileExt(request.URL.Path)
}
