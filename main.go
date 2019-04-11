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
	"github.com/ONSdigital/dp-frontend-router/handlers/splash"
	"github.com/ONSdigital/dp-frontend-router/middleware/allRoutes"
	"github.com/ONSdigital/dp-frontend-router/middleware/redirects"
	"github.com/ONSdigital/dp-frontend-router/middleware/serverError"
	"github.com/ONSdigital/go-ns/handlers/requestID"
	"github.com/ONSdigital/go-ns/handlers/reverseProxy"
	hc "github.com/ONSdigital/go-ns/healthcheck"
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
	if v := os.Getenv("DATASET_CONTROLLER_URL"); len(v) > 0 {
		config.DatasetControllerURL = v
	}
	if v := os.Getenv("FILTER_DATASET_CONTROLLER_URL"); len(v) > 0 {
		config.FilterDatasetControllerURL = v
	}
	if v := os.Getenv("GEOGRAPHY_CONTROLLER_URL"); len(v) > 0 {
		config.GeographyControllerURL = v
	}
	if v := os.Getenv("ZEBEDEE_URL"); len(v) > 0 {
		config.ZebedeeURL = v
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

	if v := os.Getenv("CONTENT_TYPE_BYTE_LIMIT"); len(v) > 0 {
		a, err := strconv.Atoi(v)
		if err == nil {
			config.ContentTypeByteLimit = int(a)
		}
	}

	var err error
	config.DebugMode, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Event(nil, "DEBUG is not a boolean", log.Data{"value": os.Getenv("DEBUG")}, log.Error(err))
	}

	config.GeographyEnabled, err = strconv.ParseBool(os.Getenv("GEOGRAPHY_ENABLED"))
	if err != nil {
		log.Event(nil, "error parsing GEOGRAPHY_ENABLED value", log.Error(err), log.Data{"value": os.Getenv("GEOGRAPHY_ENABLED")})
	}

	config.DatasetRoutesEnabled, err = strconv.ParseBool(os.Getenv("DATASET_ROUTES_ENABLED"))
	if err != nil {
		log.Event(nil, "error parsing DATASET_ROUTES_ENABLED value", log.Error(err), log.Data{"value": os.Getenv("DATASET_ROUTES_ENABLED")})
	}

	if v := os.Getenv("TAXONOMY_DOMAIN"); len(v) > 0 {
		config.TaxonomyDomain = v
	}

	log.Namespace = "dp-frontend-router"

	log.Event(nil, "overriding default renderer with service assets")

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

	datasetControllerURL, err := url.Parse(config.DatasetControllerURL)
	if err != nil {
		log.Event(nil, "error parsing dataset controller url", log.Error(err), log.Data{"url": config.DatasetControllerURL})
		os.Exit(1)
	}

	filterDatasetControllerURL, err := url.Parse(config.FilterDatasetControllerURL)
	if err != nil {
		log.Event(nil, "error parsing filter dataset controller url", log.Error(err), log.Data{"url": config.FilterDatasetControllerURL})
		os.Exit(1)
	}

	geographyControllerURL, err := url.Parse(config.GeographyControllerURL)
	if err != nil {
		log.Event(nil, "error parsing geography controller url", log.Error(err), log.Data{"url": config.GeographyControllerURL})
		os.Exit(1)
	}

	redirects.Init(assets.Asset)

	router := pat.New()

	router.Path("/healthcheck").HandlerFunc(hc.Do)

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

	if config.DatasetRoutesEnabled == true {
		middleware = append(middleware, allRoutes.Handler(map[string]http.Handler{
			"dataset_landing_page": reverseProxy.Create(datasetControllerURL, nil),
		}))
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

	if config.DatasetRoutesEnabled == true {
		router.Handle("/datasets/{uri:.*}", createReverseProxy("datasets", datasetControllerURL))
		router.Handle("/feedback{uri:.*}", createReverseProxy("feedback", datasetControllerURL))
		router.Handle("/filters/{uri:.*}", createReverseProxy("filters", filterDatasetControllerURL))
		router.Handle("/filter-outputs/{uri:.*}", createReverseProxy("filter-output", filterDatasetControllerURL))
	}
	// remove geo from prod
	if config.GeographyEnabled == true {
		router.Handle("/geography{uri:.*}", createReverseProxy("geography", geographyControllerURL))
	}
	router.Handle("/{uri:.*}", reverseProxy)

	log.Event(nil, "Starting server", log.Data{
		"bind_addr":                config.BindAddr,
		"babbage_url":              config.BabbageURL,
		"dataset_controller_url":   config.DatasetControllerURL,
		"geography_controller_url": config.GeographyControllerURL,
		"downloader_url":           config.DownloaderURL,
		"renderer_url":             config.RendererURL,
		"site_domain":              config.SiteDomain,
		"assets_path":              config.PatternLibraryAssetsPath,
		"splash_page":              config.SplashPage,
		"taxonomy_domain":          config.TaxonomyDomain,
		"analytics_sqs_url":        config.SQSAnalyticsURL,
	})

	s := &http.Server{
		Addr:         config.BindAddr,
		Handler:      alice,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Event(nil, "error starting server", log.Error(err))
		os.Exit(2)
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
