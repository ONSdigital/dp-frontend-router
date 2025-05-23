package router

import (
	"context"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/middleware/allRoutes"
	"github.com/ONSdigital/dp-frontend-router/middleware/datasetType"
	"github.com/ONSdigital/dp-frontend-router/middleware/redirects"
	dprequest "github.com/ONSdigital/dp-net/v3/request"
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

	router.Handle("/", cfg.HomepageHandler)

	if cfg.CensusAtlasEnabled {
		router.Handle("/census/maps{uri:.*}", cfg.CensusAtlasHandler)
	}

	router.Handle("/census", cfg.HomepageHandler)

	if cfg.DatasetFinderEnabled {
		router.Handle("/census/find-a-dataset", cfg.SearchHandler)
	}

	router.Handle("/download/{uri:.*}", cfg.DownloadHandler)
	router.Handle("/cookies{uri:.*}", cfg.CookieHandler)
	router.Handle("/datasets/{uri:.*}", cfg.DatasetHandler)
	router.Handle("/filters/{uri:.*}", datasetType.Handler(cfg.FilterClient, cfg.DatasetClient)(cfg.FilterHandler, cfg.FilterFlexHandler))
	router.Handle("/filter-outputs/{uri:.*}", cfg.FilterHandler)
	router.Handle("/feedback{uri:.*}", cfg.FeedbackHandler)

	if cfg.LegacySearchRedirectsEnabled {
		searchDataHandler := redirects.DynamicRedirectHandler("/searchdata", "/search")
		searchPublicationHandler := redirects.DynamicRedirectHandler("/searchpublication", "/search")

		router.Handle("/searchdata", searchDataHandler)
		router.Handle("/searchpublication", searchPublicationHandler)
	}

	if cfg.SearchRoutesEnabled {
		// needs both the SearchRoutesEnabled and DataAggregationPagesEnabled since it relies on the SearchHandler
		if cfg.DataAggregationPagesEnabled {
			// These pages used to exist as subpages but we're redirecting to root now as they didn't
			// do anything.
			legacyAdhocPagesRedirectHandler := redirects.RouteRedirectHandler("/alladhocs")
			router.Handle("/{uri:.*}/alladhocs", legacyAdhocPagesRedirectHandler)

			legacyMethodologyRedirectHandler := redirects.RouteRedirectHandler("/allmethodologies")
			router.Handle("/{uri:.*}/allmethodologies", legacyMethodologyRedirectHandler)

			legacyPublishedRequestsRedirectHandler := redirects.RouteRedirectHandler("/publishedrequests")
			router.Handle("/{uri:.*}/publishedrequests", legacyPublishedRequestsRedirectHandler)

			router.Handle("/alladhocs", cfg.SearchHandler)
			router.Handle("/datalist", cfg.SearchHandler)
			router.Handle("/allmethodologies", cfg.SearchHandler)
			router.Handle("/publishedrequests", cfg.SearchHandler)
			router.Handle("/staticlist", cfg.SearchHandler)
			router.Handle("/publications", cfg.SearchHandler)
			router.Handle("/topicspecificmethodology", cfg.SearchHandler)
			router.Handle("/timeseriestool", cfg.SearchHandler)
		}

		if cfg.TopicAggregationPagesEnabled {
			// handle dynamic topic aggregation pages
			router.Handle("/{uri:.*}/datalist", cfg.SearchHandler)
			router.Handle("/{uri:.*}/publications", cfg.SearchHandler)
			router.Handle("/{uri:.*}/staticlist", cfg.SearchHandler)
			router.Handle("/{uri:.*}/topicspecificmethodology", cfg.SearchHandler)
		}
		router.Handle("/search", cfg.SearchHandler)
	}

	router.Handle("/calendar/releasecalendar", cfg.RelCalHandler)
	router.Handle("/releasecalendar", cfg.RelCalHandler)
	if cfg.LegacyCacheProxyEnabled {
		router.Handle("/releases/{uri:.*}", cfg.ProxyHandler)
	} else {
		router.Handle("/releases/{uri:.*}", cfg.RelCalHandler)
	}

	if cfg.PreviousReleasesRouteEnabled {
		if cfg.LegacyCacheProxyEnabled {
			router.Handle("/{uri:.*}/previousreleases", cfg.ProxyHandler)
		} else {
			router.Handle("/{uri:.*}/previousreleases", cfg.SearchHandler)
		}
	}

	if cfg.RelatedDataRouteEnabled {
		if cfg.LegacyCacheProxyEnabled {
			router.Handle("/{uri:.*}/relateddata", cfg.ProxyHandler)
			router.Handle("/{uri:.*}/relatedData", cfg.ProxyHandler)
		} else {
			router.Handle("/{uri:.*}/relateddata", cfg.SearchHandler)
			router.Handle("/{uri:.*}/relatedData", cfg.SearchHandler)
		}
	}

	var handler http.Handler
	if cfg.LegacyCacheProxyEnabled {
		handler = cfg.ProxyHandler
	} else {
		handler = cfg.BabbageHandler
	}

	// if the request is for a file go directly to babbage instead of using the allRoutesMiddleware
	router.MatcherFunc(hasFileExtMatcher).Handler(handler)
	// If it is a known babbage endpoint go directly to babbage instead of using the allRoutesMiddleware
	router.MatcherFunc(isKnownBabbageEndpointMatcher).Handler(handler)

	// all other requests go through the allRoutesMiddleware to check the page type first
	handlers := map[string]http.Handler{
		"dataset_landing_page": cfg.DatasetHandler,
	}
	if cfg.NewDatasetRoutingEnabled {
		handlers["dataset"] = cfg.PrefixDatasetHandler
	}
	allRoutesMiddleware := allRoutes.Handler(handlers, cfg.ZebedeeClient, cfg.ContentTypeByteLimit)

	babbageRouter := router.PathPrefix("/").Subrouter()
	babbageRouter.Use(allRoutesMiddleware)
	babbageRouter.PathPrefix("/").Handler(handler)

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
