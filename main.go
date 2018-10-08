package main

import (
	"html/template"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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

	log.Namespace = "dp-frontend-router"

	cfg, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	log.Info("config on startup", log.Data{"config": cfg})

	log.Debug("overriding default renderer with service assets", nil)
	render.Renderer = unrolled.New(unrolled.Options{
		Asset:         assets.Asset,
		AssetNames:    assets.AssetNames,
		IsDevelopment: cfg.DebugMode,
		Layout:        "main",
		Funcs: []template.FuncMap{{
			"safeHTML": func(s string) template.HTML {
				return template.HTML(s)
			},
		}},
	})

	datasetControllerURL, err := url.Parse(cfg.DatasetControllerURL)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	filterDatasetControllerURL, err := url.Parse(cfg.FilterDatasetControllerURL)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	geographyControllerURL, err := url.Parse(cfg.GeographyControllerURL)
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
	if len(cfg.DisabledPage) > 0 {
		middleware = append(middleware, splash.Handler(cfg.DisabledPage, false))
	} else if len(cfg.SplashPage) > 0 {
		middleware = append(middleware, splash.Handler(cfg.SplashPage, true))
	}
	alice := alice.New(middleware...).Then(router)

	babbageURL, err := url.Parse(cfg.BabbageURL)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	downloaderURL, err := url.Parse(cfg.DownloaderURL)
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
	router.Handle("/", abHandler(http.HandlerFunc(homepage.Handler(reverseProxy)), reverseProxy, cfg.HomepageABPercent))
	router.Handle("/datasets/{uri:.*}", createReverseProxy(datasetControllerURL))
	router.Handle("/geography{uri:.*}", createReverseProxy(geographyControllerURL))
	router.Handle("/feedback{uri:.*}", createReverseProxy(datasetControllerURL))
	router.Handle("/filters/{uri:.*}", createReverseProxy(filterDatasetControllerURL))
	router.Handle("/filter-outputs/{uri:.*}", createReverseProxy(filterDatasetControllerURL))
	router.Handle("/{uri:.*}", reverseProxy)

	log.Debug("Starting server", log.Data{
		"bind_addr":                cfg.BindAddr,
		"babbage_url":              cfg.BabbageURL,
		"dataset_controller_url":   cfg.DatasetControllerURL,
		"geography_controller_url": cfg.GeographyControllerURL,
		"renderer_url":             cfg.RendererURL,
		"resolver_url":             cfg.ResolverURL,
		"homepage_ab_percent":      cfg.HomepageABPercent,
		"site_domain":              cfg.SiteDomain,
		"assets_path":              cfg.PatternLibraryAssetsPath,
		"splash_page":              cfg.SplashPage,
		"taxonomy_domain":          cfg.TaxonomyDomain,
		"analytics_sqs_url":        cfg.SQSAnalyticsURL,
	})

	s := server.New(cfg.BindAddr, alice)
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
