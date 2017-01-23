package serverError

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/lang"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/render"
)

type responseInterceptor struct {
	http.ResponseWriter
	req            *http.Request
	intercepted    bool
	headersWritten bool
	headerCache    http.Header
}

func (rI *responseInterceptor) WriteHeader(status int) {
	if status >= 500 {
		log.DebugR(rI.req, "Intercepted error response", nil)
		rI.intercepted = true
		rI.renderErrorPage(500, "Internal server error", "<p>We're currently experiencing some technical difficulties.</p>")
		return
	}
	rI.writeHeaders()
	rI.ResponseWriter.WriteHeader(status)
}

func (rI *responseInterceptor) renderErrorPage(code int, title, description string) {
	// TODO ask dp-frontend-renderer to render the error page
	// failing that, render this "disaster" page

	render.HTML(rI.ResponseWriter, code, "error", map[string]interface{}{
		"URI":                      rI.req.URL.Path,
		"Language":                 lang.Get(rI.req),
		"PatternLibraryAssetsPath": config.PatternLibraryAssetsPath,
		"SiteDomain":               config.SiteDomain,
		"Error": map[string]interface{}{
			"Title":       title,
			"Description": description,
		},
	})
}

func (rI *responseInterceptor) Write(b []byte) (int, error) {
	if rI.intercepted {
		return len(b), nil
	}
	rI.writeHeaders()
	return rI.ResponseWriter.Write(b)
}

func (rI *responseInterceptor) writeHeaders() {
	if rI.headersWritten {
		return
	}

	// Overwrite the server header
	rI.headerCache.Set("Server", "dp-frontend-router")

	for k, v := range rI.headerCache {
		for _, v2 := range v {
			rI.ResponseWriter.Header().Add(k, v2)
		}
	}

	rI.headersWritten = true
}

func (rI *responseInterceptor) Header() http.Header {
	return rI.headerCache
}

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(&responseInterceptor{w, req, false, false, make(http.Header)}, req)
	})
}
