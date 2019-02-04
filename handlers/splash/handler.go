package splash

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/lang"
	"github.com/ONSdigital/log.go/log"
)

func Handler(splashPage string, enabled bool) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if !enabled {
				log.Event(req.Context(), "rendering disabled page")
				if err := callRenderer(w, req, splashPage); err != nil {
					log.Event(req.Context(), "error rendering disabled page", log.Error(err))
					w.WriteHeader(500)
				}
				return
			}

			c, err := req.Cookie("splash")
			if err != nil && err != http.ErrNoCookie {
				log.Event(req.Context(), "error fetching splash cookie", log.Error(err))
			}

			if c == nil || c.Value != "y" {
				v := req.FormValue("confirm")
				if req.Method != "POST" || v != "y" {
					log.Event(req.Context(), "rendering splash page")
					if err := callRenderer(w, req, splashPage); err != nil {
						log.Event(req.Context(), "error rendering splash page", log.Error(err))
						w.WriteHeader(500)
					}
					return
				}
				log.Event(req.Context(), "splash confirmed, redirecting", log.Data{"location": req.URL.String()})
				http.SetCookie(w, &http.Cookie{Name: "splash", Value: "y"})
				http.Redirect(w, req, req.URL.String(), http.StatusFound)
				return
			}

			log.Event(req.Context(), "splash not required")
			h.ServeHTTP(w, req)
		})
	}
}

// ErrCreatingRequest is returned when http.NewRequest() fails
type ErrCreatingRequest struct{ err error }

func (e ErrCreatingRequest) Error() string { return "error creating request" }

// ErrRenderingPage is returned when calling the renderer fails
type ErrRenderingPage struct{ err error }

func (e ErrRenderingPage) Error() string { return "error rendering page" }

// ErrUnexpectedStatus is returned when the renderer returns an unexpected status code
type ErrUnexpectedStatus struct{ statusCode int }

func (e ErrUnexpectedStatus) Error() string { return "unexpected status code" }

// ErrReadingBody is returned when reading the renderer response body fails
type ErrReadingBody struct{ err error }

func (e ErrReadingBody) Error() string { return "error reading body" }

func callRenderer(w http.ResponseWriter, req *http.Request, splashPage string) error {
	rendererReq, err := http.NewRequest("POST", config.RendererURL+"/"+splashPage, bytes.NewReader([]byte(`{}`)))
	if err != nil {
		return ErrCreatingRequest{err}
	}

	// FIXME there's other headers we want
	rendererReq.Header.Set("Accept-Language", string(lang.Get(req)))
	rendererReq.Header.Set("X-Request-Id", req.Header.Get("X-Request-Id"))

	res, err := http.DefaultClient.Do(rendererReq)
	if err != nil {
		return ErrRenderingPage{err}
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return ErrUnexpectedStatus{res.StatusCode}
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ErrReadingBody{err}
	}

	for hdr, v := range res.Header {
		for _, v2 := range v {
			w.Header().Add(hdr, v2)
		}
	}

	log.Event(req.Context(), "returning splash page")
	w.WriteHeader(res.StatusCode)
	w.Write(b)

	return nil
}
