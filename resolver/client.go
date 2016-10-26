package resolver

import (
	"errors"
	"github.com/ONSdigital/go-ns/log"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type config struct {
	renderUrl string
}

var cfg = config{renderUrl: "http://localhost:20020"}
var Client ResolverClient = &http.Client{
	Timeout: 5 * time.Second,
}

func init() {
	if renderUrl := os.Getenv("RESOLVER_URL"); len(renderUrl) > 0 {
		cfg.renderUrl = renderUrl
	}
}

type ResolverClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func Get(uri string) ([]byte, error) {
	var jsonBytes []byte

	request, err := getRequest(uri)
	if err != nil {
		return jsonBytes, err
	}

	response, err := Client.Do(request)
	if err != nil {
		log.ErrorR(request, err, nil)
		return jsonBytes, err
	}

	jsonBytes, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		log.ErrorC("Error reading body", err, nil)
		return jsonBytes, err
	}

	if response.StatusCode != 200 {
		err = errors.New("Response status code is not 200")
		log.ErrorR(request, err, nil)
		return jsonBytes, err
	}
	return jsonBytes, nil
}

func getRequest(uri string) (*http.Request, error) {
	request, err := http.NewRequest("GET", cfg.renderUrl+uri, nil)
	if err != nil {
		log.Debug("Error creating new request", nil)
		log.ErrorR(request, err, nil)
		return nil, err
	}
	return request, nil
}
