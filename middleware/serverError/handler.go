package serverError

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		log.DebugR(rI.req, "Intercepted error response", log.Data{"url": status})
		rI.intercepted = true
		if status == 500 {
			rI.renderErrorPage(500, "Internal server error", "<p>We're currently experiencing some technical difficulties. You could try <a href='"+rI.req.Host+rI.req.URL.Path+"'>refreshing the page or trying again later.</a> </p>")
		} else {
			rI.renderErrorPage(503, "Service temporarily unavailable", `<p>The service is temporarily unavailable, please check our <a href="https://twitter.com/onsdigital">twitter</a> feed for updates.</p>`)
		}
		return
	}
	rI.writeHeaders()
	rI.ResponseWriter.WriteHeader(status)
}

func (rI *responseInterceptor) renderErrorPage(code int, title, description string) {
	// Attempt to render an error page
	if err := rI.callRenderer(code, title, description); err != nil {
		log.ErrorR(rI.req, err, nil)
		log.DebugR(rI.req, "rendering disaster page", nil)

		// Calling the renderer failed, render the disaster page
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
}

func (rI *responseInterceptor) callRenderer(code int, title, description string) error {
	data := map[string]interface{}{
		"error": map[string]interface{}{
			"title":       title,
			"description": description,
		},
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %s", err)
	}

	rendererReq, err := http.NewRequest("POST", config.RendererURL+"/error", bytes.NewReader(b))
	if err != nil {
		err = fmt.Errorf("error creating request: %s", err)
		return err
	}

	// FIXME there's other headers we want
	rendererReq.Header.Set("Accept-Language", string(lang.Get(rI.req)))
	rendererReq.Header.Set("X-Request-Id", rI.req.Header.Get("X-Request-Id"))

	res, err := http.DefaultClient.Do(rendererReq)
	if err != nil {
		return fmt.Errorf("error rendering page: %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	for hdr, v := range res.Header {
		for _, v2 := range v {
			rI.ResponseWriter.Header().Add(hdr, v2)
		}
	}

	log.DebugR(rI.req, "returning error page", nil)
	rI.ResponseWriter.WriteHeader(res.StatusCode)
	rI.ResponseWriter.Write(b)

	return nil
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
