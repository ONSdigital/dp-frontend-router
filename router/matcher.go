package router

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
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
	"/chart",
	"/embed",
	"/export",
	"/hash",
}

var knownBabbageEndpointSuffixes = []string{
	"/data",
	"/latest",
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
	for _, endpoint := range knownBabbageEndpointSuffixes {
		if strings.HasSuffix(path, endpoint) {
			return true
		}
	}
	return false
}

// isKnownBabbageEndpointMatcher is a mux MatcherFunc implementation, allowing routes to be matched on known babbage endpoints
func isKnownBabbageEndpointMatcher(request *http.Request, _ *mux.RouteMatch) bool {
	return IsKnownBabbageEndpoint(request.URL.Path)
}

// HasFileExt returns true if the given path has a file extension
func HasFileExt(path string) bool {
	return filepath.Ext(path) != ""
}

// hasFileExtMatcher is a mux MatcherFunc implementation, allowing routes to be matched on having a file extension
func hasFileExtMatcher(request *http.Request, _ *mux.RouteMatch) bool {
	return HasFileExt(request.URL.Path)
}
