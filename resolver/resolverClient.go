package resolver

import (
	"fmt"
	"github.com/ONSdigital/go-ns/log"
	"io/ioutil"
	"net/http"
	"os"
)

type config struct {
	host string "http://localhost"
	port string ":8082"
	uri  string "/data"
}

var cfg = config{host: "http://localhost", port: ":8082", uri: "/data"}

func init() {
	host := os.Getenv("RESOLVER_HOST")
	port := os.Getenv("RESOLVER_HOST")
	uri := os.Getenv("RESOLVER_HOST")

	if len(host) > 0 && len(port) > 0 && len(uri) > 0 {
		log.Debug("Applying resolver.Client config from environment.", nil)
		cfg.host = host
		cfg.port = port
		cfg.uri = uri
	}
}

// ResolveContent ...
func ResolveContent(uri string) ([]byte, error) {
	var jsonBytes []byte

	request, err := cfg.getRequest(uri)
	if err != nil {
		log.Debug("Error creating new request", nil)
		log.ErrorR(request, err, nil)
		return jsonBytes, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Debug("Error performing request.", nil)
		log.ErrorR(request, err, nil)
		return jsonBytes, err
	}

	jsonBytes, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Debug("Response status code is not 200. "+string(jsonBytes), nil)
		log.ErrorR(request, err, nil)
		return jsonBytes, err
	}

	if err != nil {
		log.ErrorC("Error reading body", err, nil)
		return jsonBytes, err
	}
	return jsonBytes, nil
}

func (resolverConfig config) getRequest(uri string) (*http.Request, error) {
	urlStr := resolverConfig.host + resolverConfig.port + resolverConfig.uri
	request, err := http.NewRequest("GET", urlStr, nil)

	if err != nil {
		log.ErrorR(request, err, nil)
		return nil, err
	}
	q := request.URL.Query()
	q.Add("uri", uri)
	request.URL.RawQuery = q.Encode()
	return request, nil
}
