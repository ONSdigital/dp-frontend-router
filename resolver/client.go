package resolver

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/go-ns/log"
)

var ErrUnauthorised = errors.New("unauthorised")

var Client ResolverClient = &http.Client{
	Timeout: 5 * time.Second,
}

type responseBodyReaderFunc func(r io.Reader) ([]byte, error)

var responseBodyReader = ioutil.ReadAll

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

	jsonBytes, err = responseBodyReader(response.Body)
	if err != nil {
		log.ErrorC("Error reading body", err, nil)
		return jsonBytes, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		if response.StatusCode == 401 {
			return nil, ErrUnauthorised
		}

		err = errors.New("Response status code is not 200")
		log.ErrorR(request, err, nil)
		return jsonBytes, err
	}
	return jsonBytes, nil
}

func getRequest(uri string) (*http.Request, error) {
	request, err := http.NewRequest("GET", config.ResolverURL+uri, nil)
	if err != nil {
		log.Debug("Error creating new request", nil)
		log.ErrorR(request, err, nil)
		return nil, err
	}
	return request, nil
}
