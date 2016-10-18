package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/ONSdigital/dp-frontend-router/handlers/homepage"
	"github.com/ONSdigital/go-ns/handlers/requestID"
	"github.com/ONSdigital/go-ns/handlers/timeout"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
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
	proxy := httputil.NewSingleHostReverseProxy(babbageURL)
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = babbageURL.Host
	}
	router.Handle("/{uri:.*}", proxy)

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
