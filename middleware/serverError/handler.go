package serverError

// This whole package and process needs a refactor. Re-added file - take very little responsibility for anything in here
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/lang"
	"github.com/ONSdigital/log.go/v2/log"
)

type responseInterceptor struct {
	http.ResponseWriter
	req            *http.Request
	intercepted    bool
	headersWritten bool
	headerCache    http.Header
}

func (rI *responseInterceptor) WriteHeader(status int) {
	if status >= 400 {
		log.Info(rI.req.Context(), "Intercepted error response", log.Data{"status": status})
		rI.intercepted = true
		if status == 404 {
			rI.renderErrorPage(404, "404 - The webpage you are requesting does not exist on the site", `<p> The page may have been moved, updated or deleted or you may have typed the web address incorrectly, please check the url and spelling. Alternatively, please try the <a href="#nav-search">search</a>, or return to the <a href="/" title="Our homepage" target="_self">homepage</a>.</p>`)
			return
		} else if status == 401 {
			rI.renderErrorPage(401, "401 - You do not have permission to view this web page", `<p>This page may exist, but you do not currently have permission to view it. If you believe this to be incorrect please contact a system administrator.</p>`)
			return
		}
	}
	rI.writeHeaders()
	rI.ResponseWriter.WriteHeader(status)
}

func (rI *responseInterceptor) renderErrorPage(code int, title, description string) {
	// Attempt to render an error page
	if err := rI.callRenderer(code, title, description); err != nil {
		// Calling the renderer failed, render the disaster page
		if err != nil {
			rI.writeHeaders()
			rI.ResponseWriter.WriteHeader(http.StatusInternalServerError)
			log.Error(rI.req.Context(), "error calling renderer", err)
			return
		}
	}
}

func (rI *responseInterceptor) callRenderer(code int, title, description string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	preferencesCookie := cookies.GetCookiePreferences(rI.req)
	data := map[string]interface{}{
		"error": map[string]interface{}{
			"title":       title,
			"description": description,
		},
		"cookies_preferences_set": preferencesCookie.IsPreferenceSet,
		"cookies_policy":          preferencesCookie.Policy,
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %s", err)
	}

	rendererReq, err := http.NewRequest("POST", cfg.RendererURL+"/error", bytes.NewReader(b))
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

	b, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	for hdr, v := range res.Header {
		for _, v2 := range v {
			rI.ResponseWriter.Header().Add(hdr, v2)
		}
	}

	log.Info(rI.req.Context(), "returning error page")
	rI.ResponseWriter.WriteHeader(code)
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
