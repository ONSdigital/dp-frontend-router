package homepage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-frontend-renderer/handlers/homepage"
	"github.com/ONSdigital/go-ns/log"
)

func Handler(rendererURL string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// pretend we've called Zebedee and mapped it to the renderer model
		data := homepage.Homepage{
			Type: "home_page",
			URI:  "/",
			Metadata: homepage.Metadata{
				Title:       "Homepage",
				Description: "The homepage",
				Keywords:    []string{"homepage!"},
			},
		}

		b, err := json.Marshal(&data)
		if err != nil {
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		rdr := bytes.NewReader(b)

		rendererReq, err := http.NewRequest("POST", rendererURL+"/homepage", rdr)
		if err != nil {
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// FIXME there's other headers we want
		rendererReq.Header.Set("X-Request-Id", req.Header.Get("X-Request-Id"))

		res, err := http.DefaultClient.Do(rendererReq)
		if err != nil {
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			err = fmt.Errorf("unexpected status code: %d", res.StatusCode)
			log.ErrorR(req, err, nil)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// FIXME should stream this using a io.Reader etc
		b, err = ioutil.ReadAll(res.Body)
		if err != nil {
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
		w.WriteHeader(res.StatusCode)
		w.Write(b)
	}
}
