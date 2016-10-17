package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/ONSdigital/dp-frontend-router/handlers/homepage"
	"github.com/gorilla/pat"
)

type config struct {
	BindAddr    string
	BabbageURL  string
	RendererURL string
}

func main() {
	cfg := config{
		BindAddr:    ":8080",
		BabbageURL:  "http://web.onsdigital.co.uk",
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

	babbageURL, err := url.Parse(cfg.BabbageURL)
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/", homepage.Handler(cfg.RendererURL))
	router.Handle("/{uri:.*}", newReverseProxy(httputil.NewSingleHostReverseProxy(babbageURL)))

	log.Printf("Starting server on %s\n", cfg.BindAddr)

	if err := http.ListenAndServe(cfg.BindAddr, router); err != nil {
		log.Fatal(err)
	}
}

func newReverseProxy(reverseProxy http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Server", "Go!!!!!")
		start := time.Now()
		reverseProxy.ServeHTTP(w, req)
		diff := time.Now().Sub(start)
		log.Printf("request completed: %v => %s", diff, req.URL.Path)
	}
}
