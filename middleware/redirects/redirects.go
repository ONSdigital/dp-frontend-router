package redirects

import (
	"bytes"
	"encoding/csv"
	"net/http"

	"github.com/ONSdigital/go-ns/log"
)

var redirects = make(map[string]string)

func Init(asset func(name string) ([]byte, error)) {
	b, err := asset("redirects/redirects.csv")
	if err != nil {
		log.Error(err, nil)
		panic("Can't find redirects.csv")
	}
	reader := csv.NewReader(bytes.NewReader(b))
	records, err := reader.ReadAll()
	if err != nil {
		log.Error(err, nil)
		panic("Unable to read CSV")
	}

	if len(records) == 0 {
		return
	}

	if len(records[0]) != 2 {
		panic("Redirects must have two fields")
	}

	for _, record := range records {
		if len(record[0]) == 0 {
			panic("Redirect from URL must not be empty")
		}
		if len(record[1]) == 0 {
			panic("Redirect to URL must not be empty")
		}
		redirects[record[0]] = record[1]
	}
}

//Handler ...
func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		log.DebugCtx(req.Context(), "Request is hitting redirect handler", nil)

		if redirect, ok := redirects[req.URL.Path]; ok {
			log.DebugCtx(req.Context(), "Redirected "+req.URL.Path+" to "+redirect, nil)
			http.Redirect(w, req, redirect, http.StatusTemporaryRedirect)
			return
		}

		log.DebugCtx(req.Context(), "Request is serving from redirect handler", nil)

		h.ServeHTTP(w, req)
	})
}
