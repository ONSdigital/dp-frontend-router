package main

import (
	"context"
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
	"github.com/ONSdigital/dp-frontend-router/middleware/allRoutes"
	"github.com/ONSdigital/dp-frontend-router/middleware/redirects"
	"github.com/ONSdigital/go-ns/handlers/requestID"
	"github.com/ONSdigital/go-ns/handlers/reverseProxy"
	hc "github.com/ONSdigital/go-ns/healthcheck"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
)

func main() {
	log.Namespace = "dp-frontend-router"

	cfg, err := config.Get()
	if err != nil {
		log.Event(nil, "unable to retrieve service configuration", log.Error(err))
		os.Exit(1)
	}

	ctx := context.Background()

	log.Event(ctx, "got service configuration", log.Data{"config": cfg})

	cookiesControllerURL, err := url.Parse(cfg.CookiesControllerURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.Data{"config_name": "CookiesControllerURL", "value": cfg.CookiesControllerURL}, log.Error(err))
		os.Exit(1)
	}

	datasetControllerURL, err := url.Parse(cfg.DatasetControllerURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.Data{"config_name": "DatasetControllerURL", "value": cfg.DatasetControllerURL}, log.Error(err))
		os.Exit(1)
	}

	filterDatasetControllerURL, err := url.Parse(cfg.FilterDatasetControllerURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.Data{"config_name": "FilterDatasetControllerURL", "value": cfg.FilterDatasetControllerURL}, log.Error(err))
		os.Exit(1)
	}

	geographyControllerURL, err := url.Parse(cfg.GeographyControllerURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.Data{"config_name": "GeographyControllerURL", "value": cfg.GeographyControllerURL}, log.Error(err))
		os.Exit(1)
	}

	babbageURL, err := url.Parse(cfg.BabbageURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.Data{"config_name": "BabbageURL", "value": cfg.BabbageURL}, log.Error(err))
		os.Exit(1)
	}

	downloaderURL, err := url.Parse(cfg.DownloaderURL)
	if err != nil {
		log.Event(nil, "configuration value is invalid", log.Data{"config_name": "DownloaderURL", "value": cfg.DownloaderURL}, log.Error(err))
		os.Exit(1)
	}

	redirects.Init(assets.Asset)

	router := pat.New()

	router.Path("/healthcheck").HandlerFunc(hc.Do)

	middleware := []alice.Constructor{
		requestID.Handler(16),
		log.Middleware,
		securityHandler,
		redirects.Handler,
	}

	if cfg.DatasetRoutesEnabled {
		middleware = append(middleware, allRoutes.Handler(map[string]http.Handler{
			"dataset_landing_page": reverseProxy.Create(datasetControllerURL, nil),
		}, cfg))
	}

	alice := alice.New(middleware...).Then(router)

	searchHandler, err := analytics.NewSearchHandler(cfg.SQSAnalyticsURL, cfg.RedirectSecret)
	if err != nil {
		log.Event(nil, "error creating search analytics handler", log.Error(err))
		os.Exit(1)
	}

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
	router.Handle("/{uri:.*}", reverseProxy)

	log.Event(nil, "Starting server", log.Data{"config": cfg})

	s := &http.Server{
		Addr:         cfg.BindAddr,
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
