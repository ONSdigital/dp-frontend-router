package resolver

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/go-ns/log"
)

// ErrUnauthorised error for http status code 401.
var ErrUnauthorised = errors.New("unauthorised")

const xRequestIDHeaderParam = "X-Request-Id"

// Client client for sending requests to the content resolver.
var Client ResolverClient = &http.Client{
	Timeout: 5 * time.Second,
}

type responseBodyReaderFunc func(r io.Reader) ([]byte, error)

var responseBodyReader = ioutil.ReadAll

// ResolverClient definition of a Content Resolver client
type ResolverClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Get resolve the requested content.
func Get(uri string, xRequestID string) ([]byte, error) {
	var jsonBytes []byte

	request, err := getRequest(uri, xRequestID)
	if err != nil {
		return jsonBytes, err
	}

	log.DebugC(xRequestID, "resolver.client", log.Data{
		"method": "GET",
		"uri":    request.URL.Path,
		"query":  request.URL.RawQuery,
	})

	response, err := Client.Do(request)
	if err != nil {
		log.ErrorC(xRequestID, err, nil)
		return jsonBytes, err
	}

	jsonBytes, err = responseBodyReader(response.Body)
	if err != nil {
		log.ErrorC(xRequestID, err, nil)
		return jsonBytes, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		if response.StatusCode == 401 {
			return nil, ErrUnauthorised
		}

		err = fmt.Errorf("response status code is %d", response.StatusCode)
		log.ErrorC(xRequestID, err, nil)
		return jsonBytes, err
	}
	return jsonBytes, nil
}

func getRequest(uri string, xRequestID string) (*http.Request, error) {
	request, err := http.NewRequest("GET", config.ResolverURL+uri, nil)
	if err != nil {
		err = fmt.Errorf("error creating new request: %s", err)
		log.ErrorC(xRequestID, err, nil)
		return nil, err
	}
	request.Header.Add(xRequestIDHeaderParam, xRequestID)
	return request, nil
}
