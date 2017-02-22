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
	"time"

	"github.com/ONSdigital/dp-frontend-router/assets"
	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/handlers/homepage"
	"github.com/ONSdigital/dp-frontend-router/handlers/serverError"
	"github.com/ONSdigital/dp-frontend-router/handlers/splash"
	"github.com/ONSdigital/go-ns/handlers/requestID"
	"github.com/ONSdigital/go-ns/handlers/timeout"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/render"
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
	if v := os.Getenv("DATADISCOVERY_URL"); len(v) > 0 {
		config.DataDiscoveryURL = v
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
	if v := os.Getenv("DISABLED_PAGE"); len(v) > 0 {
		config.DisabledPage = v
	}

	if v := os.Getenv("HOMEPAGE_AB_PERCENT"); len(v) > 0 {
		a, _ := strconv.Atoi(v)
		if a < 0 || a > 100 {
			log.Debug("HOMEPAGE_AB_PERCENT must be between 0 and 100", nil)
			os.Exit(1)
		}
		config.HomepageABPercent = a
	}

	var err error
	config.DebugMode, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Error(err, nil)
	}

	log.Namespace = "dp-frontend-router"

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

	router := pat.New()
	middleware := []alice.Constructor{
		requestID.Handler(16),
		log.Handler,
		securityHandler,
		serverError.Handler,
		timeout.Handler(10 * time.Second),
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

	dataDiscoveryURL, err := url.Parse(config.DataDiscoveryURL)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	dataDiscoveryProxy := createReverseProxy(dataDiscoveryURL)

	babbageProxy := createReverseProxy(babbageURL)
	router.Handle("/robots.txt", http.HandlerFunc(robotsHandler))
	router.Handle("/", abHandler(http.HandlerFunc(homepage.Handler(babbageProxy)), babbageProxy, config.HomepageABPercent))
	router.Handle("/dd{uri:(|/.*)}", dataDiscoveryProxy)
	router.Handle("/{uri:.*}", babbageProxy)

	log.Debug("Starting server", log.Data{
		"bind_addr":           config.BindAddr,
		"babbage_url":         config.BabbageURL,
		"renderer_url":        config.RendererURL,
		"resolver_url":        config.ResolverURL,
		"datadiscovery_url":   config.DataDiscoveryURL,
		"homepage_ab_percent": config.HomepageABPercent,
		"site_domain":         config.SiteDomain,
		"assets_path":         config.PatternLibraryAssetsPath,
		"splash_page":         config.SplashPage,
	})

	server := &http.Server{
		Addr:         config.BindAddr,
		Handler:      alice,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error(err, nil)
		os.Exit(2)
	}
}

// securityHandler ...
func securityHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/embed" {
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
		panic("Percent 'a' but be between 0 and 100")
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

func createReverseProxy(babbageURL *url.URL) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(babbageURL)
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
			"destination": babbageURL,
		})
		director(req)
		req.Host = babbageURL.Host
	}
	return proxy
}

func robotsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`User-agent: *
Disallow: /
`))
}
