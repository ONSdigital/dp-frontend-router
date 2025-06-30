package helpers

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// ReturnSecondSegmentFromPath returns the second segment of a path and assumes the path is formed /firstSegment/secondSegment
func ReturnSecondSegmentFromPath(path string) (secondSegment string, err error) {
	subs := strings.Split(path, "/")
	if len(subs) < 3 {
		err = fmt.Errorf("unable to extract secondSegment from path: %s", path)
		return
	}
	return subs[2], nil
}

func ParseURL(ctx context.Context, cfgValue, configName string) (*url.URL, error) {
	parsedURL, err := url.Parse(cfgValue)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": configName, "value": cfgValue})
		return nil, err
	}
	return parsedURL, nil
}

func CreateReverseProxy(proxyName string, proxyURL *url.URL) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	director := proxy.Director
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       180 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	proxy.Director = func(req *http.Request) {
		log.Info(req.Context(), "proxying request", log.HTTP(req, 0, 0, nil, nil), log.Data{
			"destination": proxyURL,
			"proxy_name":  proxyName,
		})
		otel.GetTextMapPropagator().Inject(req.Context(), propagation.HeaderCarrier(req.Header))
		director(req)
	}
	return proxy
}
