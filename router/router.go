package router

import (
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-frontend-router/handlers/abtest"
	"github.com/ONSdigital/dp-frontend-router/middleware/allRoutes"
	"github.com/ONSdigital/dp-frontend-router/middleware/datasetType"
	"github.com/ONSdigital/dp-frontend-router/middleware/redirects"
	"github.com/ONSdigital/dp-frontend-router/middleware/serverError"
	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

const (
	HttpHeaderKeyXFrameOptions = "X-Frame-Options"
)

//go:generate moq -out routertest/handler.go -pkg routertest . Handler
type Handler http.Handler

type Config struct {
	HealthCheckHandler       func(w http.ResponseWriter, req *http.Request)
	AnalyticsHandler         http.Handler
	AreaProfileEnabled       bool
	AreaProfileHandler       http.Handler
	DownloadHandler          http.Handler
	DatasetHandler           http.Handler
	NewDatasetRoutingEnabled bool
	PrefixDatasetHandler     http.Handler
	DatasetClient            datasetType.DatasetClient
	CookieHandler            http.Handler
	FilterHandler            http.Handler
	FilterFlexHandler        http.Handler
	FilterFlexEnabled        bool
	FilterClient             datasetType.FilterClient
	FeedbackHandler          http.Handler
	ContentTypeByteLimit     int
	ZebedeeClient            allRoutes.ZebedeeClient
	GeographyEnabled         bool
	GeographyHandler         http.Handler
	InteractivesEnabled      bool
	InteractivesHandler      http.Handler
	SearchRoutesEnabled      bool
	EnableSearchABTest       bool
	SearchABTestPercentage   int
	SiteDomain               string
	SearchHandler            http.Handler
	RelCalHandler            http.Handler
	RelCalEnabled            bool
	RelCalPrivatePrefix      string
	HomepageHandler          http.Handler
	BabbageHandler           http.Handler
	CensusAtlasHandler       http.Handler
	CensusAtlasEnabled       bool
}

func New(cfg Config) http.Handler {

	router := mux.NewRouter()

	middleware := []alice.Constructor{
		dprequest.HandlerRequestID(16),
		log.Middleware,
		SecurityHandler,
		healthcheckHandler(cfg.HealthCheckHandler),
		serverError.Handler,
		redirects.Handler,
	}

	alice := alice.New(middleware...).Then(router)

	router.Handle("/", cfg.HomepageHandler)
	router.Handle("/census", cfg.HomepageHandler)

	router.Handle("/redir/{data:.*}", cfg.AnalyticsHandler)
	router.Handle("/download/{uri:.*}", cfg.DownloadHandler)
	router.Handle("/cookies{uri:.*}", cfg.CookieHandler)
	router.Handle("/datasets/{uri:.*}", cfg.DatasetHandler)
	if cfg.FilterFlexEnabled {
		router.Handle("/filters/{uri:.*}", datasetType.Handler(cfg.FilterClient, cfg.DatasetClient)(cfg.FilterHandler, cfg.FilterFlexHandler))
	} else {
		router.Handle("/filters/{uri:.*}", cfg.FilterHandler)
	}
	router.Handle("/filter-outputs/{uri:.*}", cfg.FilterHandler)
	router.Handle("/feedback{uri:.*}", cfg.FeedbackHandler)

	if cfg.AreaProfileEnabled {
		router.Handle("/areas{uri:.*}", cfg.AreaProfileHandler)
		router.Handle("/geography{uri:.*}", cfg.GeographyHandler)
	} else if cfg.GeographyEnabled {
		router.Handle("/geography{uri:.*}", cfg.GeographyHandler)
	}

	if cfg.SearchRoutesEnabled {
		if cfg.EnableSearchABTest {
			router.Handle("/search", abtest.SearchHandler(cfg.SearchHandler, cfg.BabbageHandler, cfg.SearchABTestPercentage, cfg.SiteDomain))
		} else {
			router.Handle("/search", cfg.SearchHandler)
		}
	}

	if cfg.RelCalEnabled {
		router.Handle(cfg.RelCalRoutePrefix+"/releasecalendar", cfg.RelCalHandler)
		router.Handle(cfg.RelCalRoutePrefix+"/calendar/releasecalendar", cfg.RelCalHandler)
		router.Handle(cfg.RelCalRoutePrefix+"/releases/{uri:.*}", cfg.RelCalHandler)
	}

	if cfg.InteractivesEnabled {
		router.Handle("/interactives/{uri:.*}", cfg.InteractivesHandler)
	}

	if cfg.CensusAtlasEnabled {
		router.Handle("/census-atlas{uri:.*}", cfg.CensusAtlasHandler)
	}

	// if the request is for a file go directly to babbage instead of using the allRoutesMiddleware
	router.MatcherFunc(hasFileExtMatcher).Handler(cfg.BabbageHandler)

	// If it is a known babbage endpoint go directly to babbage instead of using the allRoutesMiddleware
	router.MatcherFunc(isKnownBabbageEndpointMatcher).Handler(cfg.BabbageHandler)

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
	babbageRouter.PathPrefix("/").Handler(cfg.BabbageHandler)

	return alice
}

// SecurityHandler ...
func SecurityHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/embed" &&
			!strings.HasPrefix(req.URL.Path, "/visualisations/") &&
			!strings.HasPrefix(req.URL.Path, "/interactives/") {
			w.Header().Set(HttpHeaderKeyXFrameOptions, "SAMEORIGIN")
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
