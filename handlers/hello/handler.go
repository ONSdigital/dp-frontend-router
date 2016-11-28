package hello

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/lang"
	"github.com/ONSdigital/go-ns/log"
)

func Handler(w http.ResponseWriter, req *http.Request) {

	// Make call to controller
	model, _, err := doRequest(req, "GET", config.HelloWorldURL, nil)
	if err != nil {
		log.ErrorR(req, err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	rdr := bytes.NewReader(model)

	b, headers, err := doRequest(req, "POST", config.RendererURL+"/hello", rdr)
	if err != nil {
		log.ErrorR(req, err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	for hdr, v := range headers {
		for _, v2 := range v {
			w.Header().Add(hdr, v2)
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func doRequest(originalRequest *http.Request, method string, url string, body io.Reader) (responseBody []byte, headers http.Header, err error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	request.Header.Set("Accept-Language", string(lang.Get(originalRequest)))
	request.Header.Set("X-Request-Id", originalRequest.Header.Get("X-Request-Id"))

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("Handler.handler: unexpected status code: %d", res.StatusCode)
	}

	responseBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	headers = res.Header

	return
}
