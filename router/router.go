package router

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/middleware/allRoutes"
	"github.com/ONSdigital/dp-frontend-router/middleware/datasetType"
	"github.com/ONSdigital/dp-frontend-router/middleware/redirects"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	HTTPHeaderKeyXFrameOptions = "X-Frame-Options"
)

//go:generate moq -out routertest/handler.go -pkg routertest . Handler
type Handler http.Handler

type Config struct {
	HealthCheckHandler           func(w http.ResponseWriter, req *http.Request)
	AnalyticsHandler             http.Handler
	DownloadHandler              http.Handler
	DatasetHandler               http.Handler
	DatasetClient                datasetType.DatasetClient
	NewDatasetRoutingEnabled     bool
	PrefixDatasetHandler         http.Handler
	CookieHandler                http.Handler
	FilterHandler                http.Handler
	FilterFlexHandler            http.Handler
	FilterClient                 datasetType.FilterClient
	FeedbackHandler              http.Handler
	ContentTypeByteLimit         int
	ZebedeeClient                allRoutes.ZebedeeClient
	LegacySearchRedirectsEnabled bool
	DataAggregationPagesEnabled  bool
	TopicAggregationPagesEnabled bool
	SearchRoutesEnabled          bool
	SiteDomain                   string
	SearchHandler                http.Handler
	RelCalHandler                http.Handler
	HomepageHandler              http.Handler
	BabbageHandler               http.Handler
	ProxyHandler                 http.Handler
	CensusAtlasHandler           http.Handler
	CensusAtlasEnabled           bool
	DatasetFinderEnabled         bool
	LegacyCacheProxyEnabled      bool
	PerformanceTestHandler       http.Handler
	PreviousReleasesRouteEnabled bool
	RelatedDataRouteEnabled      bool
}

//nolint:gocyclo // This will be reduced once we complete decomm of legacy search and flags can be removed.
func New(cfg Config) http.Handler {
	router := mux.NewRouter()
	middleware := []alice.Constructor{
		dprequest.HandlerRequestID(16),
		log.Middleware,
		SecurityHandler,
		healthcheckHandler(cfg.HealthCheckHandler),
		redirects.Handler,
	}

	appConfig, err := config.Get()
	if err != nil {
		log.Error(context.Background(), "error getting config", err)
	}

	if appConfig.OtelEnabled {
		middleware = append(middleware, otelhttp.NewMiddleware("dp-frontend-router"))
	}

	newAlice := alice.New(middleware...).Then(router)

	router.Handle("/", cfg.PerformanceTestHandler)
	router.Handle("/beans", cfg.PerformanceTestHandler)
	router.Handle("/means", cfg.PerformanceTestHandler)
	router.Handle("/heinz", cfg.PerformanceTestHandler)

	for i := 0; i < 1000; i++ {
		router.Handle("/beans"+strconv.Itoa(i), cfg.PerformanceTestHandler)
	}

	router.Handle("/{uri:.*}/beans", cfg.PerformanceTestHandler)
	router.Handle("/{uri:.*}/means", cfg.PerformanceTestHandler)
	router.Handle("/{uri:.*}/heinz", cfg.PerformanceTestHandler)

	for i := 0; i < 1000; i++ {
		router.Handle("/{uri:.*}/beans"+strconv.Itoa(i), cfg.PerformanceTestHandler)
	}

	router.Handle("/beans/{uri:.*}", cfg.PerformanceTestHandler)
	router.Handle("/means/{uri:.*}", cfg.PerformanceTestHandler)
	router.Handle("/heinz/{uri:.*}", cfg.PerformanceTestHandler)

	for i := 0; i < 1000; i++ {
		router.Handle("/beans"+strconv.Itoa(i)+"/{uri:.*}", cfg.PerformanceTestHandler)
	}

	perfTestsRedirectHandler := redirects.RouteRedirectHandler("/beans")
	router.Handle("/cola", perfTestsRedirectHandler)
	router.Handle("/capri-sun", perfTestsRedirectHandler)
	router.Handle("/choc pudding", perfTestsRedirectHandler)

	for i := 0; i < 1000; i++ {
		router.Handle("/cola"+strconv.Itoa(i), perfTestsRedirectHandler)
	}

	perfTestsRedirectSuffixHandler := redirects.RouteRedirectHandler("/means")
	router.Handle("/{uri:.*}/cola", perfTestsRedirectSuffixHandler)
	router.Handle("/{uri:.*}/capri-sun", perfTestsRedirectSuffixHandler)
	router.Handle("/{uri:.*}/choc pudding", perfTestsRedirectSuffixHandler)

	for i := 0; i < 1000; i++ {
		router.Handle("/{uri:.*}/cola"+strconv.Itoa(i), perfTestsRedirectHandler)
	}

	perfTestsRedirectPrefixHandler := redirects.RouteRedirectHandler("/heinz")
	router.Handle("/cola/{uri:.*}", perfTestsRedirectPrefixHandler)
	router.Handle("/capri-sun/{uri:.*}", perfTestsRedirectPrefixHandler)
	router.Handle("/choc pudding/{uri:.*}", perfTestsRedirectPrefixHandler)

	for i := 0; i < 1000; i++ {
		router.Handle("/cola"+strconv.Itoa(i)+"/{uri:.*}", perfTestsRedirectHandler)
	}

	return newAlice
}

// SecurityHandler is the custom handler for for setting frame options
func SecurityHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/embed" &&
			!strings.HasPrefix(req.URL.Path, "/visualisations/") &&
			!strings.HasPrefix(req.URL.Path, "/census/maps/") {
			w.Header().Set(HTTPHeaderKeyXFrameOptions, "SAMEORIGIN")
		}
		h.ServeHTTP(w, req)
	})
}

// healthcheckHandler uses the provided handler for /health endpoint, and serves any other traffic to the next handler in chain
func healthcheckHandler(hc func(w http.ResponseWriter, req *http.Request)) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/health" {
				hc(w, req)
				return
			}
			h.ServeHTTP(w, req)
		})
	}
}
