package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/zebedee"
	"github.com/ONSdigital/dp-frontend-router/assets"
	"github.com/ONSdigital/dp-frontend-router/config"
	"github.com/ONSdigital/dp-frontend-router/handlers/analytics"
	"github.com/ONSdigital/dp-frontend-router/middleware/redirects"
	"github.com/ONSdigital/dp-frontend-router/router"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/log.go/v2/log"
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
		os.Exit(1)
	}

	log.Info(ctx, "got service configuration", log.Data{"config": cfg})

	cookiesControllerURL, err := url.Parse(cfg.CookiesControllerURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "CookiesControllerURL", "value": cfg.CookiesControllerURL})
		os.Exit(1)
	}

	datasetControllerURL, err := url.Parse(cfg.DatasetControllerURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "DatasetControllerURL", "value": cfg.DatasetControllerURL})
		os.Exit(1)
	}

	var prefixedDatasetURL = cfg.DatasetControllerURL + "/dataset"
	prefixDatasetControllerURL, err := url.Parse(prefixedDatasetURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "DatasetControllerURL", "value": cfg.DatasetControllerURL})
		os.Exit(1)
	}

	filterDatasetControllerURL, err := url.Parse(cfg.FilterDatasetControllerURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "FilterDatasetControllerURL", "value": cfg.FilterDatasetControllerURL})
		os.Exit(1)
	}

	geographyControllerURL, err := url.Parse(cfg.GeographyControllerURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "GeographyControllerURL", "value": cfg.GeographyControllerURL})
		os.Exit(1)
	}

	homepageControllerURL, err := url.Parse(cfg.HomepageControllerURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "HomepageControllerURL", "value": cfg.HomepageControllerURL})
		os.Exit(1)
	}

	searchControllerURL, err := url.Parse(cfg.SearchControllerURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "SearchControllerURL", "value": cfg.SearchControllerURL})
		os.Exit(1)
	}

	relcalControllerURL, err := url.Parse(cfg.ReleaseCalendarControllerURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "ReleaseCalendarControllerURL", "value": cfg.ReleaseCalendarControllerURL})
		os.Exit(1)
	}

	babbageURL, err := url.Parse(cfg.BabbageURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "BabbageURL", "value": cfg.BabbageURL})
		os.Exit(1)
	}

	downloaderURL, err := url.Parse(cfg.DownloaderURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "DownloaderURL", "value": cfg.DownloaderURL})
		os.Exit(1)
	}

	feedbackControllerURL, err := url.Parse(cfg.FeedbackControllerURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "FeedbackControllerURL", "value": cfg.FeedbackControllerURL})
		os.Exit(1)
	}

	areaProfileControllerURL, err := url.Parse(cfg.AreaProfilesControllerURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "AreaProfileControllerURL", "value": cfg.AreaProfilesControllerURL})
		os.Exit(1)
	}

	filterFlexDatasetServiceURL, err := url.Parse(cfg.FilterFlexDatasetServiceURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "FilterFlexDatasetServiceURL", "value": cfg.FilterFlexDatasetServiceURL})
		os.Exit(1)
	}

	interactivesControllerURL, err := url.Parse(cfg.InteractivesControllerURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": "InteractivesControllerURL", "value": cfg.InteractivesControllerURL})
		os.Exit(1)
	}

	censusAtlasURL := urlFromConfig(ctx, "CensusAtlas", cfg.CensusAtlasURL)

	enableSearchABTest := config.IsEnableSearchABTest(*cfg)
	enableRelCalABTest := config.IsEnabledRelCalABTest(*cfg)

	redirects.Init(assets.Asset)

	// create ZebedeeClient proxying calls through the API Router
	hcClienter := dphttp.NewClient()
	hcClienter.SetMaxRetries(cfg.ZebedeeRequestMaximumRetries)
	hcClienter.SetTimeout(cfg.ZebedeeRequestMaximumTimeoutSeconds)

	zebedeeClient := zebedee.NewClientWithClienter(cfg.APIRouterURL, hcClienter)

	hcClient := health.NewClient("api-router", cfg.APIRouterURL)
	filterClient := filter.NewWithHealthClient(hcClient)
	datasetClient := dataset.NewWithHealthClient(hcClient)

	// Healthcheck API
	versionInfo, err := healthcheck.NewVersionInfo(BuildTime, GitCommit, Version)
	if err != nil {
		log.Fatal(ctx, "Failed to obtain VersionInfo for healthcheck", err)
		os.Exit(1)
	}
	hc := healthcheck.New(versionInfo, cfg.HealthcheckCriticalTimeout, cfg.HealthcheckInterval)
	if err = hc.AddCheck("API router", zebedeeClient.Checker); err != nil {
		log.Fatal(ctx, "Failed to add api router checker to healthcheck", err)
		os.Exit(1)
	}

	analyticsHandler, err := analytics.NewSearchHandler(ctx, cfg.SQSAnalyticsURL, cfg.RedirectSecret)
	if err != nil {
		log.Fatal(ctx, "error creating search analytics handler", err)
		os.Exit(1)
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
	areaProfileHandler := createReverseProxy("areas", areaProfileControllerURL)
	filterFlexHandler := createReverseProxy("flex", filterFlexDatasetServiceURL)
	interactivesHandler := createReverseProxy("interactives", interactivesControllerURL)
	censusAtlasHandler := createReverseProxy("censusAtlas", censusAtlasURL)
	var geographyHandler http.Handler
	if cfg.AreaProfilesRoutesEnabled {
		geographyHandler = redirects.DynamicRedirectHandler("/geography", "/areas")
	} else {
		geographyHandler = createReverseProxy("geography", geographyControllerURL)
	}

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
		FilterFlexEnabled:            cfg.FilterFlexRoutesEnabled,
		FilterFlexHandler:            filterFlexHandler,
		GeographyEnabled:             cfg.GeographyEnabled,
		GeographyHandler:             geographyHandler,
		LegacySearchRedirectsEnabled: cfg.LegacySearchRedirectsEnabled,
		SearchRoutesEnabled:          cfg.SearchRoutesEnabled,
		SearchHandler:                searchHandler,
		RelCalHandler:                relcalHandler,
		RelCalEnabled:                cfg.ReleaseCalendarEnabled,
		RelCalRoutePrefix:            cfg.ReleaseCalendarRoutePrefix,
		RelCalEnableABTest:           enableRelCalABTest,
		RelCalABTestPercentage:       cfg.ReleaseCalendarABTestPercentage,
		InteractivesHandler:          interactivesHandler,
		InteractivesEnabled:          cfg.InteractivesRoutesEnabled,
		EnableSearchABTest:           enableSearchABTest,
		SearchABTestPercentage:       cfg.SearchABTestPercentage,
		SiteDomain:                   cfg.SiteDomain,
		HomepageHandler:              homepageHandler,
		BabbageHandler:               babbageHandler,
		ZebedeeClient:                zebedeeClient,
		ContentTypeByteLimit:         cfg.ContentTypeByteLimit,
		CensusAtlasHandler:           censusAtlasHandler,
		CensusAtlasEnabled:           cfg.CensusAtlasRoutesEnabled,
		DatasetFinderEnabled:         cfg.DatasetFinderEnabled,
	}

	httpHandler := router.New(routerConfig)

	log.Info(ctx, "Starting server", log.Data{"config": cfg})

	s := &http.Server{
		Addr:         cfg.BindAddr,
		Handler:      httpHandler,
		ReadTimeout:  cfg.ProxyTimeout,
		WriteTimeout: cfg.ProxyTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Start health check
	hc.Start(ctx)

	// Start server
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(ctx, "error starting server", err)
		hc.Stop()
		os.Exit(2)
	}
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
		director(req)
	}
	return proxy
}

func urlFromConfig(ctx context.Context, serviceName, serviceURL string) *url.URL {
	configuredServiceURL, err := url.Parse(serviceURL)
	if err != nil {
		log.Fatal(ctx, "configuration value is invalid", err, log.Data{"config_name": serviceName, "value": serviceURL})
		os.Exit(1)
	}
	return configuredServiceURL
}
