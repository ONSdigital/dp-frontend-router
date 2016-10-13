package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/pat"
)

func main() {
	bindAddr := os.Getenv("BIND_ADDR")
	if len(bindAddr) == 0 {
		bindAddr = ":8080"
	}

	p := pat.New()

	babbageURL, err := url.Parse("http://web.onsdigital.co.uk")
	if err != nil {
		log.Fatal(err)
	}

	p.Handle("/{uri:.*}", newReverseProxy(httputil.NewSingleHostReverseProxy(babbageURL)))

	if err := http.ListenAndServe(bindAddr, p); err != nil {
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
