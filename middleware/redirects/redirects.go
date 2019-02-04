package redirects

import (
	"bytes"
	"encoding/csv"
	"net/http"

	"github.com/ONSdigital/log.go/log"
)

var redirects = make(map[string]string)

// PanicOnInitError (when true) causes Init() to panic if redirects.csv
// contains invalid data
var PanicOnInitError = true

func Init(asset func(name string) ([]byte, error)) {
	b, err := asset("redirects/redirects.csv")
	if err != nil {
		log.Event(nil, "can't find redirects.csv", log.Error(err))
		if PanicOnInitError {
			panic("Can't find redirects.csv")
		}
		return
	}

	reader := csv.NewReader(bytes.NewReader(b))
	records, err := reader.ReadAll()
	if err != nil {
		log.Event(nil, "error reading redirects.csv", log.Error(err))
		if PanicOnInitError {
			panic("Unable to read CSV")
		}
		return
	}

	if len(records) == 0 {
		return
	}

	for line, record := range records {
		if len(record) > 0 {
			if len(record[0]) == 0 {
				log.Event(nil, "redirect 'from' URL empty", log.Data{"line": line})
				if PanicOnInitError {
					panic("redirect 'from' URL empty, check logs")
				}
				continue
			}
			if len(record) > 1 {
				if len(record[1]) == 0 {
					log.Event(nil, "redirect 'to' URL empty", log.Data{"line": line})
					if PanicOnInitError {
						panic("redirect 'to' URL empty, check logs")
					}
					continue
				}

				log.Event(nil, "adding redirct", log.Data{"from": record[0], "to": record[1]})
				redirects[record[0]] = record[1]
			} else {
				log.Event(nil, "redirect is missing 'to' value", log.Data{"line": line})
				if PanicOnInitError {
					panic("redirect 'to' URL empty, check logs")
				}
			}
		}
	}
}

//Handler ...
func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if redirect, ok := redirects[req.URL.Path]; ok {
			log.Event(req.Context(), "redirect found", log.Data{"location": redirect}, log.HTTP(req, 0, 0, nil, nil))
			http.Redirect(w, req, redirect, http.StatusTemporaryRedirect)
			return
		}

		h.ServeHTTP(w, req)
	})
}
