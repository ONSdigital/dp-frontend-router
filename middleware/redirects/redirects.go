package redirects

import (
	"bytes"
	"context"
	"encoding/csv"
	"net/http"
	"strings"

	"github.com/ONSdigital/log.go/v2/log"
)

var redirects = make(map[string]string)

// PanicOnInitError (when true) causes Init() to panic if redirects.csv
// contains invalid data
var PanicOnInitError = true

func Init(asset func(name string) ([]byte, error)) {
	b, err := asset("redirects/redirects.csv")
	if err != nil {
		log.Error(context.Background(), "can't find redirects.csv", err)
		if PanicOnInitError {
			panic("Can't find redirects.csv")
		}
		return
	}

	reader := csv.NewReader(bytes.NewReader(b))
	records, err := reader.ReadAll()
	if err != nil {
		log.Error(context.Background(), "error reading redirects.csv", err)
		if PanicOnInitError {
			panic("Unable to read CSV")
		}
		return
	}

	if len(records) == 0 {
		return
	}

	handleRecords(records)
}

func handleRecords(records [][]string) {
	for line, record := range records {
		if len(record) > 0 {
			if record[0] == "" {
				log.Warn(context.Background(), "redirect 'from' URL empty", log.Data{"line": line})
				if PanicOnInitError {
					panic("redirect 'from' URL empty, check logs")
				}
				continue
			}
			if len(record) > 1 {
				if record[1] == "" {
					log.Warn(context.Background(), "redirect 'to' URL empty", log.Data{"line": line})
					if PanicOnInitError {
						panic("redirect 'to' URL empty, check logs")
					}
					continue
				}

				log.Info(context.Background(), "adding redirect", log.Data{"from": record[0], "to": record[1]})
				redirects[record[0]] = record[1]
			} else {
				log.Warn(context.Background(), "redirect is missing 'to' value", log.Data{"line": line})
				if PanicOnInitError {
					panic("redirect 'to' URL empty, check logs")
				}
			}
		}
	}
}

// Handler with temporary redirect
func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if redirect, ok := redirects[req.URL.Path]; ok {
			log.Info(req.Context(), "redirect found", log.Data{"location": redirect}, log.HTTP(req, 0, 0, nil, nil))
			http.Redirect(w, req, redirect, http.StatusTemporaryRedirect)
			return
		}
		h.ServeHTTP(w, req)
	})
}

// DynamicRedirectHandler redirects requests to the provide 'to' base path whilst keeping all the other information of the request url
func DynamicRedirectHandler(redirectFrom, redirectTo string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		redirect := *req.URL
		redirect.Path = strings.Replace(req.URL.Path, redirectFrom, redirectTo, 1)
		redirectURL := redirect.String()

		log.Info(req.Context(), "redirect found", log.Data{"location": redirectURL}, log.HTTP(req, 0, 0, nil, nil))
		http.Redirect(w, req, redirectURL, http.StatusMovedPermanently)
	})
}

// RouteRedirectHandler will redirect to the new path provided, this can allow for wildcard redirecting. It also preserves query strings.
func RouteRedirectHandler(redirectTo string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		redirectURL := redirectTo

		if req.URL.RawQuery != "" {
			redirectURL += "?" + req.URL.RawQuery
		}

		log.Info(req.Context(), "redirect found", log.Data{"location": redirectURL}, log.HTTP(req, 0, 0, nil, nil))
		http.Redirect(w, req, redirectURL, http.StatusMovedPermanently)
	})
}
