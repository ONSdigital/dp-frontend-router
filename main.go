package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"golang.org/x/net/netutil"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-router/assets"
	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/handlers/analytics"
	"github.com/ONSdigital/dp-frontend-router/middleware/redirects"
	"github.com/ONSdigital/dp-frontend-router/router"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	dpotelgo "github.com/ONSdigital/dp-otel-go"
	"github.com/ONSdigital/log.go/v2/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

	cookiesControllerURL, _ := parseURL(ctx, cfg.CookiesControllerURL, "CookiesControllerURL")
	datasetControllerURL, _ := parseURL(ctx, cfg.DatasetControllerURL, "DatasetControllerURL")
	prefixedDatasetURL := cfg.DatasetControllerURL + "/dataset"
	prefixDatasetControllerURL, _ := parseURL(ctx, prefixedDatasetURL, "DatasetControllerURL")
	filterDatasetControllerURL, _ := parseURL(ctx, cfg.FilterDatasetControllerURL, "FilterDatasetControllerURL")
	homepageControllerURL, _ := parseURL(ctx, cfg.HomepageControllerURL, "HomepageControllerURL")
	searchControllerURL, _ := parseURL(ctx, cfg.SearchControllerURL, "SearchControllerURL")
	relcalControllerURL, _ := parseURL(ctx, cfg.ReleaseCalendarControllerURL, "ReleaseCalendarControllerURL")
	legacyCacheProxyURL, _ := parseURL(ctx, cfg.LegacyCacheProxyURL, "LegacyCacheProxyURL")
	babbageURL, _ := parseURL(ctx, cfg.BabbageURL, "BabbageURL")
	downloaderURL, _ := parseURL(ctx, cfg.DownloaderURL, "DownloaderURL")
	feedbackControllerURL, _ := parseURL(ctx, cfg.FeedbackControllerURL, "FeedbackControllerURL")
	areaProfileControllerURL, _ := parseURL(ctx, cfg.AreaProfilesControllerURL, "AreaProfileControllerURL")
	filterFlexDatasetServiceURL, _ := parseURL(ctx, cfg.FilterFlexDatasetServiceURL, "FilterFlexDatasetServiceURL")
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

	analyticsHandler, err := analytics.NewSearchHandler(ctx, cfg.SQSAnalyticsURL, cfg.RedirectSecret)
	if err != nil {
		log.Fatal(ctx, "error creating search analytics handler", err)
	}

	downloadHandler := createReverseProxy("download", downloaderURL)
	cookieHandler := createReverseProxy("cookies", cookiesControllerURL)
	datasetHandler := createReverseProxy("datasets", datasetControllerURL)
	prefixDatasetHandler := createReverseProxy("datasets", prefixDatasetControllerURL)
	filterHandler := createReverseProxy("filters", filterDatasetControllerURL)
	feedbackHandler := createReverseProxy("feedback", feedbackControllerURL)
	searchHandler := createReverseProxy("search", searchControllerURL)
	relcalHandler := createReverseProxy("relcal", relcalControllerURL)
	homepageHandler := createReverseProxy("homepage", homepageControllerURL)
	babbageHandler := createReverseProxy("babbage", babbageURL)
	proxyHandler := createReverseProxy("legacyCacheProxy", legacyCacheProxyURL)
	areaProfileHandler := createReverseProxy("areas", areaProfileControllerURL)
	filterFlexHandler := createReverseProxy("flex", filterFlexDatasetServiceURL)
	censusAtlasHandler := createReverseProxy("censusAtlas", censusAtlasURL)

	routerConfig := router.Config{
		AnalyticsHandler:             analyticsHandler,
		AreaProfileEnabled:           cfg.AreaProfilesRoutesEnabled,
		AreaProfileHandler:           areaProfileHandler,
		DownloadHandler:              downloadHandler,
		CookieHandler:                cookieHandler,
		DatasetHandler:               datasetHandler,
		NewDatasetRoutingEnabled:     cfg.NewDatasetRoutingEnabled,
		PrefixDatasetHandler:         prefixDatasetHandler,
		DatasetClient:                datasetClient,
		HealthCheckHandler:           hc.Handler,
		FilterHandler:                filterHandler,
		FilterClient:                 filterClient,
		FeedbackHandler:              feedbackHandler,
		FilterFlexHandler:            filterFlexHandler,
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

func parseURL(ctx context.Context, cfgValue, configName string) (*url.URL, error) {
	parsedURL, err := url.Parse(cfgValue)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": configName, "value": cfgValue})
		return nil, err
	}
	return parsedURL, nil
}

func createReverseProxy(proxyName string, proxyURL *url.URL) http.Handler {
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

func urlFromConfig(ctx context.Context, serviceName, serviceURL string) *url.URL {
	configuredServiceURL, err := url.Parse(serviceURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": serviceName, "value": serviceURL})
	}
	return configuredServiceURL
}
