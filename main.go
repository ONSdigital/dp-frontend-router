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
	"github.com/ONSdigital/dp-frontend-router/handlers/homepage"
	"github.com/ONSdigital/dp-frontend-router/handlers/splash"
	"github.com/ONSdigital/dp-frontend-router/middleware/allRoutes"
	"github.com/ONSdigital/dp-frontend-router/middleware/redirects"
	"github.com/ONSdigital/dp-frontend-router/middleware/serverError"
	"github.com/ONSdigital/go-ns/handlers/requestID"
	"github.com/ONSdigital/go-ns/handlers/reverseProxy"
	hc "github.com/ONSdigital/go-ns/healthcheck"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/render"
	"github.com/ONSdigital/go-ns/server"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
	unrolled "github.com/unrolled/render"
)

func main() {
	if v := os.Getenv("BIND_ADDR"); len(v) > 0 {
		config.BindAddr = v
	}
	if v := os.Getenv("BABBAGE_URL"); len(v) > 0 {
		config.BabbageURL = v
	}
	if v := os.Getenv("RESOLVER_URL"); len(v) > 0 {
		config.ResolverURL = v
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

	if v := os.Getenv("HOMEPAGE_AB_PERCENT"); len(v) > 0 {
		a, _ := strconv.Atoi(v)
		if a < 0 || a > 100 {
			log.Debug("HOMEPAGE_AB_PERCENT must be between 0 and 100", nil)
			os.Exit(1)
		}
		config.HomepageABPercent = int(a)
	}

	config.DebugMode = getEnvBool("DEBUG")

	config.GeoEnabled = getEnvBool("GEOGRAPHY_ENABLED")

	if v := os.Getenv("TAXONOMY_DOMAIN"); len(v) > 0 {
		config.TaxonomyDomain = v
	}

	log.Namespace = "dp-frontend-router"

	log.Debug("overriding default renderer with service assets", nil)
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
		log.Error(err, nil)
		os.Exit(1)
	}

	filterDatasetControllerURL, err := url.Parse(config.FilterDatasetControllerURL)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	geographyControllerURL, err := url.Parse(config.GeographyControllerURL)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	redirects.Init(assets.Asset)

	router := pat.New()

	router.Path("/healthcheck").HandlerFunc(hc.Do)

	middleware := []alice.Constructor{
		requestID.Handler(16),
		log.Handler,
		securityHandler,
		serverError.Handler,
		allRoutes.Handler(map[string]http.Handler{
			"dataset_landing_page": reverseProxy.Create(datasetControllerURL, nil),
		}),
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
		log.Error(err, nil)
		os.Exit(1)
	}

	downloaderURL, err := url.Parse(config.DownloaderURL)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	searchHandler, err := analytics.NewSearchHandler()
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	reverseProxy := createReverseProxy(babbageURL)
	router.Handle("/redir/{data:.*}", searchHandler)
	router.Handle("/download/{uri:.*}", createReverseProxy(downloaderURL))
	router.Handle("/", abHandler(http.HandlerFunc(homepage.Handler(reverseProxy)), reverseProxy, config.HomepageABPercent))
	router.Handle("/datasets/{uri:.*}", createReverseProxy(datasetControllerURL))
	router.Handle("/feedback{uri:.*}", createReverseProxy(datasetControllerURL))
	router.Handle("/filters/{uri:.*}", createReverseProxy(filterDatasetControllerURL))
	router.Handle("/filter-outputs/{uri:.*}", createReverseProxy(filterDatasetControllerURL))
	// remov geo from prod
	if config.GeoEnabled == true {
		router.Handle("/geography{uri:.*}", createReverseProxy(geographyControllerURL))
	}
	router.Handle("/{uri:.*}", reverseProxy)

	log.Debug("Starting server", log.Data{
		"bind_addr":                config.BindAddr,
		"babbage_url":              config.BabbageURL,
		"dataset_controller_url":   config.DatasetControllerURL,
		"geography_controller_url": config.GeographyControllerURL,
		"renderer_url":             config.RendererURL,
		"resolver_url":             config.ResolverURL,
		"homepage_ab_percent":      config.HomepageABPercent,
		"site_domain":              config.SiteDomain,
		"assets_path":              config.PatternLibraryAssetsPath,
		"splash_page":              config.SplashPage,
		"taxonomy_domain":          config.TaxonomyDomain,
		"analytics_sqs_url":        config.SQSAnalyticsURL,
	})

	s := server.New(config.BindAddr, alice)
	if err := s.ListenAndServe(); err != nil {
		log.Error(err, nil)
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
func abHandler(a, b http.Handler, percentA int) http.Handler {
	if percentA == 0 {
		log.Debug("percentA is 0, defaulting to handler B", nil)
		return b
	} else if percentA == 100 {
		log.Debug("percentA is 100, defaulting to handler A", nil)
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
			if rand.Intn(100) < percentA {
				cookieValue = "A"
			} else {
				cookieValue = "B"
			}

			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie = &http.Cookie{Name: "homepage-version", Value: cookieValue, Expires: expiration}
			http.SetCookie(w, cookie)
		}

		// Use cookie value to direct to a or b handler
		switch cookie.Value {
		case "A":
			a.ServeHTTP(w, req)
		case "B":
			b.ServeHTTP(w, req)
		default:
			log.Debug("invalid cookie value, reselecting", log.Data{"value": cookie.Value})
			cookie = nil
			goto RETRY
		}
	})
}

func createReverseProxy(proxyURL *url.URL) http.Handler {
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
		log.DebugR(req, "Proxying request", log.Data{
			"destination": proxyURL,
		})
		director(req)
	}
	return proxy
}

//error handling
func getEnvBool(key string) bool {
	res, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		log.Debug(key+" must be a bool", nil)
		return false
	}
	return res
}
