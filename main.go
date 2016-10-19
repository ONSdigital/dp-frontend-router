package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/pat"
	"github.com/justinas/alice"
	"github.com/onsdigital/dp-frontend-router/handlers/homepage"
	"github.com/onsdigital/go-ns/handlers/requestID"
	"github.com/onsdigital/go-ns/handlers/timeout"
	"github.com/onsdigital/go-ns/log"
)

type config struct {
	BindAddr    string
	BabbageURL  string
	RendererURL string
}

func main() {
	cfg := config{
		BindAddr:    ":8080",
		BabbageURL:  "https://www.ons.gov.uk",
		RendererURL: "http://localhost:8081",
	}

	if v := os.Getenv("BIND_ADDR"); len(v) > 0 {
		cfg.BindAddr = v
	}
	if v := os.Getenv("BABBAGE_URL"); len(v) > 0 {
		cfg.BabbageURL = v
	}
	if v := os.Getenv("RENDERER_URL"); len(v) > 0 {
		cfg.RendererURL = v
	}

	log.Namespace = "dp-frontend-router"

	router := pat.New()
	alice := alice.New(
		timeout.Handler(10*time.Second),
		log.Handler,
		requestID.Handler(16),
	).Then(router)

	babbageURL, err := url.Parse(cfg.BabbageURL)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	router.HandleFunc("/", homepage.Handler(cfg.RendererURL))
	router.Handle("/{uri:.*}", createReverseProxy(babbageURL))

	log.Debug("Starting server", log.Data{"bind_addr": cfg.BindAddr})

	server := &http.Server{
		Addr:         cfg.BindAddr,
		Handler:      alice,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error(err, nil)
		os.Exit(2)
	}
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
		director(req)
		req.Host = babbageURL.Host
	}
	return proxy
}
