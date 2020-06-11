package main

import (
	"context"
	"github.com/ONSdigital/dp-frontend-router/middleware/profiler"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ONSdigital/dp-frontend-router/middleware/serverError"

	client "github.com/ONSdigital/dp-api-clients-go/zebedee"
	"github.com/ONSdigital/dp-frontend-router/assets"
	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/handlers/analytics"
	"github.com/ONSdigital/dp-frontend-router/middleware/allRoutes"
	"github.com/ONSdigital/dp-frontend-router/middleware/redirects"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/go-ns/handlers/requestID"
	"github.com/ONSdigital/go-ns/handlers/reverseProxy"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"

	_ "net/http/pprof"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	log.Namespace = "dp-frontend-router"

	cfg, err := config.Get()
	if err != nil {
		log.Event(nil, "unable to retrieve service configuration", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	ctx := context.Background()

	log.Event(ctx, "got service configuration", log.INFO, log.Data{"config": cfg})

	cookiesControllerURL, err := url.Parse(cfg.CookiesControllerURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.FATAL, log.Data{"config_name": "CookiesControllerURL", "value": cfg.CookiesControllerURL}, log.Error(err))
		os.Exit(1)
	}

	datasetControllerURL, err := url.Parse(cfg.DatasetControllerURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.FATAL, log.Data{"config_name": "DatasetControllerURL", "value": cfg.DatasetControllerURL}, log.Error(err))
		os.Exit(1)
	}

	filterDatasetControllerURL, err := url.Parse(cfg.FilterDatasetControllerURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.FATAL, log.Data{"config_name": "FilterDatasetControllerURL", "value": cfg.FilterDatasetControllerURL}, log.Error(err))
		os.Exit(1)
	}

	geographyControllerURL, err := url.Parse(cfg.GeographyControllerURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.FATAL, log.Data{"config_name": "GeographyControllerURL", "value": cfg.GeographyControllerURL}, log.Error(err))
		os.Exit(1)
	}

	homepageControllerURL, err := url.Parse(cfg.HomepageControllerURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.FATAL, log.Data{"config_name": "HomepageControllerURL", "value": cfg.HomepageControllerURL}, log.Error(err))
		os.Exit(1)
	}

	babbageURL, err := url.Parse(cfg.BabbageURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.FATAL, log.Data{"config_name": "BabbageURL", "value": cfg.BabbageURL}, log.Error(err))
		os.Exit(1)
	}

	downloaderURL, err := url.Parse(cfg.DownloaderURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.FATAL, log.Data{"config_name": "DownloaderURL", "value": cfg.DownloaderURL}, log.Error(err))
		os.Exit(1)
	}

	redirects.Init(assets.Asset)

	router := pat.New()

	middleware := []alice.Constructor{
		requestID.Handler(16),
		log.Middleware,
		securityHandler,
		serverError.Handler,
		redirects.Handler,
	}

	zebedeeClient := client.New(cfg.ZebedeeURL)

	if cfg.DatasetRoutesEnabled {
		middleware = append(middleware, allRoutes.Handler(map[string]http.Handler{
			"dataset_landing_page": reverseProxy.Create(datasetControllerURL, nil),
		}, zebedeeClient, cfg))
	}

	aliceChain := alice.New(middleware...).Then(router)

	searchHandler, err := analytics.NewSearchHandler(cfg.SQSAnalyticsURL, cfg.RedirectSecret)
	if err != nil {
		log.Event(ctx, "error creating search analytics handler", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	// Healthcheck API
	versionInfo, err := healthcheck.NewVersionInfo(BuildTime, GitCommit, Version)
	if err != nil {
		log.Event(ctx, "Failed to obtain VersionInfo for healthcheck", log.FATAL, log.Error(err))
		os.Exit(1)
	}
	hc := healthcheck.New(versionInfo, cfg.HealthckeckCriticalTimeout, cfg.HealthckeckInterval)
	if err = hc.AddCheck("Zebedee", zebedeeClient.Checker); err != nil {
		log.Event(ctx, "Failed to add Zebedee checker to healthcheck", log.FATAL, log.Error(err))
		os.Exit(1)
	}
	router.HandleFunc("/health", hc.Handler)

	reverseProxy := createReverseProxy("babbage", babbageURL)
	router.Handle("/redir/{data:.*}", searchHandler)
	router.Handle("/download/{uri:.*}", createReverseProxy("download", downloaderURL))

	if cfg.CookiesRoutesEnabled {
		router.Handle("/cookies{uri:.*}", createReverseProxy("cookies", cookiesControllerURL))
	}

	if cfg.DatasetRoutesEnabled {
		router.Handle("/datasets/{uri:.*}", createReverseProxy("datasets", datasetControllerURL))
		router.Handle("/feedback{uri:.*}", createReverseProxy("feedback", datasetControllerURL))
		router.Handle("/filters/{uri:.*}", createReverseProxy("filters", filterDatasetControllerURL))
		router.Handle("/filter-outputs/{uri:.*}", createReverseProxy("filter-output", filterDatasetControllerURL))
	}
	// remove geo from prod
	if cfg.GeographyEnabled {
		router.Handle("/geography{uri:.*}", createReverseProxy("geography", geographyControllerURL))
	}

	if cfg.NewHomepageEnabled {
		router.Handle("/", createReverseProxy("homepage", homepageControllerURL))
	}

	if cfg.EnableProfiler {
		profilerChain := alice.New(profiler.Middleware(cfg.PprofToken)).Then(http.DefaultServeMux)
		router.PathPrefix("/debug").Handler(profilerChain)
	}

	router.Handle("/{uri:.*}", reverseProxy)

	log.Event(nil, "Starting server", log.INFO, log.Data{"config": cfg})

	s := &http.Server{
		Addr:         cfg.BindAddr,
		Handler:      aliceChain,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start healthcheck
	hc.Start(ctx)

	// Start server
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Event(ctx, "error starting server", log.FATAL, log.Error(err))
		hc.Stop()
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
		log.Event(nil, "abHandler decision", log.INFO, log.Data{"percentA": percentA, "destination": "b"})
		return b
	} else if percentA == 100 {
		log.Event(nil, "abHandler decision", log.INFO, log.Data{"percentA": percentA, "destination": "a"})
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

			log.Event(req.Context(), "abHandler decision", log.INFO, log.Data{"sel": sel, "handler": cookieValue})

			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie = &http.Cookie{Name: "homepage-version", Value: cookieValue, Expires: expiration}
			http.SetCookie(w, cookie)
		}

		// Use cookie value to direct to a or b handler
		switch cookie.Value {
		case "A":
			log.Event(req.Context(), "abHandler decision", log.INFO, log.Data{"cookie": "A", "destination": "a"})
			a.ServeHTTP(w, req)
		case "B":
			log.Event(req.Context(), "abHandler decision", log.INFO, log.Data{"cookie": "B", "destination": "b"})
			b.ServeHTTP(w, req)
		default:
			log.Event(req.Context(), "abHandler invalid cookie value, reselecting", log.INFO)
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
		log.Event(req.Context(), "proxying request", log.INFO, log.HTTP(req, 0, 0, nil, nil), log.Data{
			"destination": proxyURL,
			"proxy_name":  proxyName,
		})
		director(req)
	}
	return proxy
}
