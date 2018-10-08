package homepage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/lang"
	"github.com/ONSdigital/dp-frontend-router/resolver"
	"github.com/ONSdigital/go-ns/log"
)

const xRequestIDHeaderParam = "X-Request-Id"

func Handler(babbageProxy http.Handler) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		b, err := resolver.Get("/", req.Header.Get(xRequestIDHeaderParam))
		if err == resolver.ErrUnauthorised {
			err = fmt.Errorf("unauthorised user: %s", err)
			log.ErrorR(req, err, nil)
			babbageProxy.ServeHTTP(w, req)
			return
		} else if err != nil {
			err = fmt.Errorf("failed to resolve request: %s", err)
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		rdr := bytes.NewReader(b)

		cfg, err := config.Get()
		if err != nil {
			log.Error(err, nil)
			return
		}

		rendererReq, err := http.NewRequest("POST", cfg.RendererURL+"/homepage", rdr)
		if err != nil {
			err = fmt.Errorf("error creating request: %s", err)
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// FIXME there's other headers we want
		rendererReq.Header.Set("Accept-Language", string(lang.Get(req)))
		rendererReq.Header.Set("X-Request-Id", req.Header.Get("X-Request-Id"))

		res, err := http.DefaultClient.Do(rendererReq)
		if err != nil {
			err = fmt.Errorf("error rendering page: %s", err)
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			err = fmt.Errorf("Handler.handler: unexpected status code: %d", res.StatusCode)
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// FIXME should stream this using a io.Reader etc
		b, err = ioutil.ReadAll(res.Body)
		if err != nil {
			err = fmt.Errorf("error reading response body: %s", err)
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		for hdr, v := range res.Header {
			for _, v2 := range v {
				w.Header().Add(hdr, v2)
			}
		}

		log.DebugR(req, "returning homepage", nil)
		w.WriteHeader(res.StatusCode)
		w.Write(b)
	}
}
