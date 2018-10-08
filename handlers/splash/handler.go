package splash

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/lang"
	"github.com/ONSdigital/go-ns/log"
)

func Handler(splashPage string, enabled bool) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if !enabled {
				if err := callRenderer(w, req, splashPage); err != nil {
					w.WriteHeader(500)
				}
				return
			}

			c, err := req.Cookie("splash")
			if err != nil && err != http.ErrNoCookie {
				log.ErrorR(req, err, nil)
			}
			if c == nil || c.Value != "y" {
				v := req.FormValue("confirm")
				if req.Method != "POST" || v != "y" {
					if err := callRenderer(w, req, splashPage); err != nil {
						w.WriteHeader(500)
					}
					return
				}
				http.SetCookie(w, &http.Cookie{Name: "splash", Value: "y"})
				http.Redirect(w, req, req.URL.String(), http.StatusFound)
				return
			}
			h.ServeHTTP(w, req)
		})
	}
}

func callRenderer(w http.ResponseWriter, req *http.Request, splashPage string) error {
	cfg, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		return err
	}

	rendererReq, err := http.NewRequest("POST", cfg.RendererURL+"/"+splashPage, bytes.NewReader([]byte(`{}`)))
	if err != nil {
		err = fmt.Errorf("error creating request: %s", err)
		return err
	}

	// FIXME there's other headers we want
	rendererReq.Header.Set("Accept-Language", string(lang.Get(req)))
	rendererReq.Header.Set("X-Request-Id", req.Header.Get("X-Request-Id"))

	res, err := http.DefaultClient.Do(rendererReq)
	if err != nil {
		return fmt.Errorf("error rendering page: %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	for hdr, v := range res.Header {
		for _, v2 := range v {
			w.Header().Add(hdr, v2)
		}
	}

	log.DebugR(req, "returning splash page", nil)
	w.WriteHeader(res.StatusCode)
	w.Write(b)

	return nil
}
