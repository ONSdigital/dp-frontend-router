package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/netutil"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-router/assets"
	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/helpers"
	"github.com/ONSdigital/dp-frontend-router/middleware/redirects"
	"github.com/ONSdigital/dp-frontend-router/router"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v3/http"
	dpotelgo "github.com/ONSdigital/dp-otel-go"
	"github.com/ONSdigital/log.go/v2/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	log.Namespace = "dp-frontend-router"

	ctx := context.Background()

	cfg, err := config.Get()
	if err != nil {
		log.Fatal(ctx, "unable to retrieve service configuration", err)
	}

	log.Info(ctx, "got service configuration", log.Data{"config": cfg})

	var otelShutdown func(context.Context) error
	if cfg.OtelEnabled {
		// Set up OpenTelemetry
		otelConfig := dpotelgo.Config{
			OtelServiceName:          cfg.OTServiceName,
			OtelExporterOtlpEndpoint: cfg.OTExporterOTLPEndpoint,
			OtelBatchTimeout:         cfg.OTBatchTimeout,
		}

		var oErr error
		otelShutdown, oErr = dpotelgo.SetupOTelSDK(ctx, otelConfig)
		if oErr != nil {
			log.Error(ctx, "error setting up OpenTelemetry - hint: ensure OTEL_EXPORTER_OTLP_ENDPOINT is set", oErr)
		}
		// Handle shutdown properly so nothing leaks.
		defer func() {
			err = errors.Join(err, otelShutdown(context.Background()))
		}()
	}

	cookiesControllerURL, _ := helpers.ParseURL(ctx, cfg.CookiesControllerURL, "CookiesControllerURL")
	datasetControllerURL, _ := helpers.ParseURL(ctx, cfg.DatasetControllerURL, "DatasetControllerURL")
	prefixedDatasetURL := cfg.DatasetControllerURL + "/dataset"
	prefixDatasetControllerURL, _ := helpers.ParseURL(ctx, prefixedDatasetURL, "DatasetControllerURL")
	filterDatasetControllerURL, _ := helpers.ParseURL(ctx, cfg.FilterDatasetControllerURL, "FilterDatasetControllerURL")
	homepageControllerURL, _ := helpers.ParseURL(ctx, cfg.HomepageControllerURL, "HomepageControllerURL")
	searchControllerURL, _ := helpers.ParseURL(ctx, cfg.SearchControllerURL, "SearchControllerURL")
	relcalControllerURL, _ := helpers.ParseURL(ctx, cfg.ReleaseCalendarControllerURL, "ReleaseCalendarControllerURL")
	legacyCacheProxyURL, _ := helpers.ParseURL(ctx, cfg.LegacyCacheProxyURL, "LegacyCacheProxyURL")
	babbageURL, _ := helpers.ParseURL(ctx, cfg.BabbageURL, "BabbageURL")
	downloaderURL, _ := helpers.ParseURL(ctx, cfg.DownloaderURL, "DownloaderURL")
	feedbackControllerURL, _ := helpers.ParseURL(ctx, cfg.FeedbackControllerURL, "FeedbackControllerURL")
	censusAtlasURL := urlFromConfig(ctx, "CensusAtlas", cfg.CensusAtlasURL)

	redirects.Init(assets.Asset)

	// create ZebedeeClient proxying calls through the API Router
	hcClienter := dphttp.NewClient()
	hcClienter.SetMaxRetries(cfg.ZebedeeRequestMaximumRetries)
	hcClienter.SetTimeout(cfg.ZebedeeRequestMaximumTimeout)

	zebedeeClient := zebedee.NewClientWithClienter(cfg.APIRouterURL, hcClienter)

	hcClient := health.NewClient("api-router", cfg.APIRouterURL)
	filterClient := filter.NewWithHealthClient(hcClient)
	datasetClient := dataset.NewWithHealthClient(hcClient)

	// Healthcheck API
	versionInfo, err := healthcheck.NewVersionInfo(BuildTime, GitCommit, Version)
	if err != nil {
		log.Fatal(ctx, "Failed to obtain VersionInfo for healthcheck", err)
	}
	hc := healthcheck.New(versionInfo, cfg.HealthcheckCriticalTimeout, cfg.HealthcheckInterval)
	if err = hc.AddCheck("API router", zebedeeClient.Checker); err != nil {
		log.Fatal(ctx, "Failed to add api router checker to healthcheck", err)
	}

	downloadHandler := helpers.CreateReverseProxy("download", downloaderURL)
	cookieHandler := helpers.CreateReverseProxy("cookies", cookiesControllerURL)
	datasetHandler := helpers.CreateReverseProxy("datasets", datasetControllerURL)
	prefixDatasetHandler := helpers.CreateReverseProxy("datasets", prefixDatasetControllerURL)
	feedbackHandler := helpers.CreateReverseProxy("feedback", feedbackControllerURL)
	filterHandler := helpers.CreateReverseProxy("filters", filterDatasetControllerURL)
	searchHandler := helpers.CreateReverseProxy("search", searchControllerURL)
	relcalHandler := helpers.CreateReverseProxy("relcal", relcalControllerURL)
	homepageHandler := helpers.CreateReverseProxy("homepage", homepageControllerURL)
	babbageHandler := helpers.CreateReverseProxy("babbage", babbageURL)
	proxyHandler := helpers.CreateReverseProxy("legacyCacheProxy", legacyCacheProxyURL)
	censusAtlasHandler := helpers.CreateReverseProxy("censusAtlas", censusAtlasURL)

	routerConfig := router.Config{
		DownloadHandler:              downloadHandler,
		CookieHandler:                cookieHandler,
		DatasetHandler:               datasetHandler,
		NewDatasetRoutingEnabled:     cfg.NewDatasetRoutingEnabled,
		PrefixDatasetHandler:         prefixDatasetHandler,
		DatasetClient:                datasetClient,
		HealthCheckHandler:           hc.Handler,
		FilterClient:                 filterClient,
		FeedbackHandler:              feedbackHandler,
		FilterHandler:                filterHandler,
		LegacySearchRedirectsEnabled: cfg.LegacySearchRedirectsEnabled,
		DataAggregationPagesEnabled:  cfg.DataAggregationPagesEnabled,
		TopicAggregationPagesEnabled: cfg.TopicAggregationPagesEnabled,
		SearchRoutesEnabled:          cfg.SearchRoutesEnabled,
		SearchHandler:                searchHandler,
		RelCalHandler:                relcalHandler,
		SiteDomain:                   cfg.SiteDomain,
		HomepageHandler:              homepageHandler,
		BabbageHandler:               babbageHandler,
		ProxyHandler:                 proxyHandler,
		ZebedeeClient:                zebedeeClient,
		ContentTypeByteLimit:         cfg.ContentTypeByteLimit,
		CensusAtlasHandler:           censusAtlasHandler,
		CensusAtlasEnabled:           cfg.CensusAtlasRoutesEnabled,
		DatasetFinderEnabled:         cfg.DatasetFinderEnabled,
		LegacyCacheProxyEnabled:      cfg.LegacyCacheProxyEnabled,
		PreviousReleasesRouteEnabled: cfg.PreviousReleasesRouteEnabled,
		RelatedDataRouteEnabled:      cfg.RelatedDataRouteEnabled,
	}

	httpHandler := router.New(routerConfig)

	if cfg.OtelEnabled {
		httpHandler = otelhttp.NewHandler(httpHandler, "/")
	}

	log.Info(ctx, "Starting server", log.Data{"config": cfg})

	s := &http.Server{
		Handler:      httpHandler,
		ReadTimeout:  cfg.ProxyTimeout,
		WriteTimeout: cfg.ProxyTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Start health check
	hc.Start(ctx)

	// Create a LimitListener to cap concurrent http connections
	l, err := net.Listen("tcp", cfg.BindAddr)
	if err != nil {
		log.Fatal(ctx, "error starting tcp listener", err)
	}

	if maxC := cfg.HTTPMaxConnections; maxC > 0 {
		l = netutil.LimitListener(l, maxC)
	}

	// Start server
	if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
		log.Fatal(ctx, "error starting server", err)
	}
	l.Close()

	if cfg.OtelEnabled {
		err = otelShutdown(ctx)
		if err != nil {
			log.Fatal(ctx, "error shutting down opentelemettry", err)
		}
	}
}

func urlFromConfig(ctx context.Context, serviceName, serviceURL string) *url.URL {
	configuredServiceURL, err := url.Parse(serviceURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": serviceName, "value": serviceURL})
	}
	return configuredServiceURL
}
