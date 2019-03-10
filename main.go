package main

import (
	"html/template"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ONSdigital/dp-frontend-router/assets"
	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/handlers/analytics"
	"github.com/ONSdigital/dp-frontend-router/handlers/serverError"
	"github.com/ONSdigital/dp-frontend-router/handlers/splash"
	"github.com/ONSdigital/dp-frontend-router/middleware/redirects"
	"github.com/ONSdigital/go-ns/handlers/requestID"
	"github.com/ONSdigital/go-ns/render"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
	unrolled "github.com/unrolled/render"
)

func main() {
	log.Namespace = "dp-frontend-router"

	if v := os.Getenv("BIND_ADDR"); len(v) > 0 {
		config.BindAddr = v
	}
	if v := os.Getenv("BABBAGE_URL"); len(v) > 0 {
		config.BabbageURL = v
	}
	if v := os.Getenv("RENDERER_URL"); len(v) > 0 {
		config.RendererURL = v
	}
	if v := os.Getenv("DOWNLOADER_URL"); len(v) > 0 {
		config.DownloaderURL = v
	}
	if v := os.Getenv("PATTERN_LIBRARY_ASSETS_PATH"); len(v) > 0 {
		config.PatternLibraryAssetsPath = v
	}
	if v := os.Getenv("SITE_DOMAIN"); len(v) > 0 {
		config.SiteDomain = v
	}
	if v := os.Getenv("SPLASH_PAGE"); len(v) > 0 {
		config.SplashPage = v
	}

	if v := os.Getenv("REDIRECT_SECRET"); len(v) > 0 {
		config.RedirectSecret = v
	}

	if v := os.Getenv("ANALYTICS_SQS_URL"); len(v) > 0 {
		config.SQSAnalyticsURL = v
	}

	if v := os.Getenv("DISABLED_PAGE"); len(v) > 0 {
		config.DisabledPage = v
	}

	var err error
	config.DebugMode, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Event(nil, "DEBUG is not a boolean", log.Data{"value": os.Getenv("DEBUG")}, log.Error(err))
	}

	render.Renderer = unrolled.New(unrolled.Options{
		Asset:         assets.Asset,
		AssetNames:    assets.AssetNames,
		IsDevelopment: config.DebugMode,
		Layout:        "main",
		Funcs: []template.FuncMap{{
			"safeHTML": func(s string) template.HTML {
				return template.HTML(s)
			},
		}},
	})

	redirects.Init(assets.Asset)

	router := pat.New()
	middleware := []alice.Constructor{
		requestID.Handler(16),
		log.Middleware,
		securityHandler,
		serverError.Handler,
		redirects.Handler,
	}
	if len(config.DisabledPage) > 0 {
		middleware = append(middleware, splash.Handler(config.DisabledPage, false))
	} else if len(config.SplashPage) > 0 {
		middleware = append(middleware, splash.Handler(config.SplashPage, true))
	}
	alice := alice.New(middleware...).Then(router)

	babbageURL, err := url.Parse(config.BabbageURL)
	if err != nil {
		log.Event(nil, "error parsing babbage URL", log.Error(err), log.Data{"url": config.BabbageURL})
		os.Exit(1)
	}

	downloaderURL, err := url.Parse(config.DownloaderURL)
	if err != nil {
		log.Event(nil, "error parsing download URL", log.Error(err), log.Data{"url": config.DownloaderURL})
		os.Exit(1)
	}

	searchHandler, err := analytics.NewSearchHandler()
	if err != nil {
		log.Event(nil, "error creating search analytics handler", log.Error(err))
		os.Exit(1)
	}

	reverseProxy := createReverseProxy("babbage", babbageURL)
	router.Handle("/redir/{data:.*}", searchHandler)
	router.Handle("/download/{uri:.*}", createReverseProxy("download", downloaderURL))
	router.Handle("/{uri:.*}", reverseProxy)

	log.Event(nil, "starting server", log.Data{
		"bind_addr":         config.BindAddr,
		"babbage_url":       config.BabbageURL,
		"renderer_url":      config.RendererURL,
		"downloader_url":    config.DownloaderURL,
		"site_domain":       config.SiteDomain,
		"assets_path":       config.PatternLibraryAssetsPath,
		"splash_page":       config.SplashPage,
		"analytics_sqs_url": config.SQSAnalyticsURL,
	})

	server := &http.Server{
		Addr:         config.BindAddr,
		Handler:      alice,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Event(nil, "error starting server", log.Error(err))
			os.Exit(2)
		}
	}
}

// securityHandler ...
func securityHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/embed" && !strings.HasPrefix(req.URL.Path, "/visualisations/") {
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		}
		h.ServeHTTP(w, req)
	})
}

//abHandler ... percentA is the percentage of request that handler 'a' is used
//
// FIXME this isn't used anymore, it could be removed, but seems like it might be useful?
func abHandler(a, b http.Handler, percentA int) http.Handler {
	if percentA == 0 {
		log.Event(nil, "abHandler decision", log.Data{"percentA": percentA, "destination": "b"})
		return b
	} else if percentA == 100 {
		log.Event(nil, "abHandler decision", log.Data{"percentA": percentA, "destination": "a"})
		return a
	}

	if percentA < 0 || percentA > 100 {
		panic("Percent 'a' must be between 0 and 100")
	}
	rand.Seed(time.Now().UnixNano())

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Detect cookie
		cookie, _ := req.Cookie("homepage-version")

	RETRY:
		if cookie == nil {
			var cookieValue string
			sel := rand.Intn(100)
			if sel < percentA {
				cookieValue = "A"
			} else {
				cookieValue = "B"
			}

			log.Event(nil, "abHandler decision", log.Data{"sel": sel, "handler": cookieValue})

			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie = &http.Cookie{Name: "homepage-version", Value: cookieValue, Expires: expiration}
			http.SetCookie(w, cookie)
		}

		// Use cookie value to direct to a or b handler
		switch cookie.Value {
		case "A":
			log.Event(nil, "abHandler decision", log.Data{"cookie": "A", "destination": "a"})
			a.ServeHTTP(w, req)
		case "B":
			log.Event(nil, "abHandler decision", log.Data{"cookie": "B", "destination": "b"})
			b.ServeHTTP(w, req)
		default:
			log.Event(nil, "abHandler invalid cookie value, reselecting")
			cookie = nil
			goto RETRY
		}
	})
}

func createReverseProxy(proxyName string, proxyURL *url.URL) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	director := proxy.Director
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	proxy.Director = func(req *http.Request) {
		log.Event(req.Context(), "proxying request", log.HTTP(req, 0, 0, nil, nil), log.Data{
			"destination": proxyURL,
			"proxy_name":  proxyName,
		})
		director(req)
	}
	return proxy
}
