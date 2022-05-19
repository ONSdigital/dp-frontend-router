package config

import (
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-frontend-router
type Config struct {
	BindAddr                            string        `envconfig:"BIND_ADDR"`
	BabbageURL                          string        `envconfig:"BABBAGE_URL"`
	RendererURL                         string        `envconfig:"RENDERER_URL"`
	CookiesControllerURL                string        `envconfig:"COOKIES_CONTROLLER_URL"`
	HomepageControllerURL               string        `envconfig:"HOMEPAGE_CONTROLLER_URL"`
	DatasetControllerURL                string        `envconfig:"DATASET_CONTROLLER_URL"`
	FilterDatasetControllerURL          string        `envconfig:"FILTER_DATASET_CONTROLLER_URL"`
	GeographyControllerURL              string        `envconfig:"GEOGRAPHY_CONTROLLER_URL"`
	GeographyEnabled                    bool          `envconfig:"GEOGRAPHY_ENABLED"`
	FeedbackControllerURL               string        `envconfig:"FEEDBACK_CONTROLLER_URL"`
	FeedbackEnabled                     bool          `envconfig:"FEEDBACK_ENABLED"`
	SearchControllerURL                 string        `envconfig:"SEARCH_CONTROLLER_URL"`
	SearchRoutesEnabled                 bool          `envconfig:"SEARCH_ROUTES_ENABLED"`
	ReleaseCalendarControllerURL        string        `envconfig:"RELEASE_CALENDAR_CONTROLLER_URL"`
	ReleaseCalendarEnabled              bool          `envconfig:"RELEASE_CALENDAR_ENABLED"`
	ReleaseCalendarRoutePrefix          string        `envconfig:"RELEASE_CALENDAR_ROUTE_PREFIX"`
	APIRouterURL                        string        `envconfig:"API_ROUTER_URL"`
	DownloaderURL                       string        `envconfig:"DOWNLOADER_URL"`
	AreaProfilesControllerURL           string        `envconfig:"AREA_PROFILE_CONTROLLER_URL"`
	AreaProfilesRoutesEnabled           bool          `envconfig:"AREA_PROFILE_ROUTES_ENABLED"`
	InteractivesControllerURL           string        `envconfig:"INTERACTIVES_CONTROLLER_URL"`
	InteractivesRoutesEnabled           bool          `envconfig:"INTERACTIVES_ROUTES_ENABLED"`
	FilterFlexDatasetServiceURL         string        `envconfig:"FILTER_FLEX_DATASET_SERVICE_URL"`
	FilterFlexRoutesEnabled             bool          `envconfig:"FILTER_FLEX_ROUTES_ENABLED"`
	CensusAtlasURL                      string        `envconfig:"CENSUS_ATLAS_URL"`
	CensusAtlasRoutesEnabled            bool          `envconfig:"CENSUS_ATLAS_ROUTES_ENABLED"`
	PatternLibraryAssetsPath            string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	SiteDomain                          string        `envconfig:"SITE_DOMAIN"`
	RedirectSecret                      string        `envconfig:"REDIRECT_SECRET" json:"-"`
	SQSAnalyticsURL                     string        `envconfig:"SQS_ANALYTICS_URL"`
	ContentTypeByteLimit                int           `envconfig:"CONTENT_TYPE_BYTE_LIMIT"`
	HealthcheckCriticalTimeout          time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	HealthcheckInterval                 time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	ZebedeeRequestMaximumTimeoutSeconds time.Duration `envconfig:"ZEBEDEE_REQUEST_TIMEOUT_SECONDS"`
	ZebedeeRequestMaximumRetries        int           `envconfig:"ZEBEDEE_REQUEST_MAXIMUM_RETRIES"`
	EnableSearchABTest                  bool          `envconfig:"ENABLE_SEARCH_AB_TEST"`
	SearchABTestPercentage              int           `envconfig:"SEARCH_AB_TEST_PERCENTAGE"`
	ProxyTimeout                        time.Duration `envconfig:"PROXY_TIMEOUT"`
}

var cfg *Config

// Get returns the default config with any modifications made through environment variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg := &Config{
		BindAddr:                            ":20000",
		BabbageURL:                          "http://localhost:8080",
		RendererURL:                         "http://localhost:20010",
		CookiesControllerURL:                "http://localhost:24100",
		HomepageControllerURL:               "http://localhost:24400",
		DatasetControllerURL:                "http://localhost:20200",
		FilterDatasetControllerURL:          "http://localhost:20001",
		GeographyControllerURL:              "http://localhost:23700",
		GeographyEnabled:                    false,
		FeedbackControllerURL:               "http://localhost:25200",
		SearchControllerURL:                 "http://localhost:25000",
		SearchRoutesEnabled:                 false,
		ReleaseCalendarControllerURL:        "http://localhost:27700",
		ReleaseCalendarEnabled:              false,
		APIRouterURL:                        "http://localhost:23200/v1",
		DownloaderURL:                       "http://localhost:23400",
		AreaProfilesControllerURL:           "http://localhost:26600",
		AreaProfilesRoutesEnabled:           false,
		InteractivesControllerURL:           "http://localhost:27300",
		InteractivesRoutesEnabled:           false,
		FilterFlexDatasetServiceURL:         "http://localhost:20100",
		FilterFlexRoutesEnabled:             false,
		CensusAtlasURL:                      "http://localhost:28100",
		CensusAtlasRoutesEnabled:            false,
		PatternLibraryAssetsPath:            "https://cdn.ons.gov.uk/sixteens/f816ac8",
		SiteDomain:                          "ons.gov.uk",
		RedirectSecret:                      "secret",
		SQSAnalyticsURL:                     "",
		ContentTypeByteLimit:                5000000,
		HealthcheckCriticalTimeout:          90 * time.Second,
		HealthcheckInterval:                 30 * time.Second,
		ZebedeeRequestMaximumTimeoutSeconds: 5 * time.Second,
		ZebedeeRequestMaximumRetries:        0,
		EnableSearchABTest:                  false,
		SearchABTestPercentage:              10,
		ProxyTimeout:                        5 * time.Second,
	}

	if err := envconfig.Process("", cfg); err != nil {
		return cfg, err
	}

	cfg.ReleaseCalendarRoutePrefix = validatePrivatePrefix(cfg.ReleaseCalendarRoutePrefix)

	return cfg, nil
}

// IsEnableSearchABTest checks whether ab test is enabled and that percentage is a sensible int value
func IsEnableSearchABTest(cfg Config) bool {
	percentage := cfg.SearchABTestPercentage
	if cfg.EnableSearchABTest && percentage > 0 && percentage < 100 {
		return true
	}
	return false
}

// validatePrivatePrefix ensures that a non-empty private path prefix starts with a '/'
func validatePrivatePrefix(prefix string) string {
	if prefix != "" && !strings.HasPrefix(prefix, "/") {
		return "/" + prefix
	}

	return prefix
}
