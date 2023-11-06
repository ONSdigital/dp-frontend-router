package config

import (
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-frontend-router
type Config struct {
	AWS                             AWS
	APIRouterURL                    string        `envconfig:"API_ROUTER_URL"`
	AreaProfilesControllerURL       string        `envconfig:"AREA_PROFILE_CONTROLLER_URL"`
	AreaProfilesRoutesEnabled       bool          `envconfig:"AREA_PROFILE_ROUTES_ENABLED"`
	BabbageURL                      string        `envconfig:"BABBAGE_URL"`
	BindAddr                        string        `envconfig:"BIND_ADDR"`
	CensusAtlasRoutesEnabled        bool          `envconfig:"CENSUS_ATLAS_ROUTES_ENABLED"`
	CensusAtlasURL                  string        `envconfig:"CENSUS_ATLAS_URL"`
	ContentTypeByteLimit            int           `envconfig:"CONTENT_TYPE_BYTE_LIMIT"`
	CookiesControllerURL            string        `envconfig:"COOKIES_CONTROLLER_URL"`
	DatasetControllerURL            string        `envconfig:"DATASET_CONTROLLER_URL"`
	DatasetFinderEnabled            bool          `envconfig:"DATASET_FINDER_ENABLED"`
	DownloaderURL                   string        `envconfig:"DOWNLOADER_URL"`
	EnableReleaseCalendarABTest     bool          `envconfig:"ENABLE_RELEASE_CALENDAR_AB_TEST"`
	FeedbackControllerURL           string        `envconfig:"FEEDBACK_CONTROLLER_URL"`
	FeedbackEnabled                 bool          `envconfig:"FEEDBACK_ENABLED"`
	FilterDatasetControllerURL      string        `envconfig:"FILTER_DATASET_CONTROLLER_URL"`
	FilterFlexDatasetServiceURL     string        `envconfig:"FILTER_FLEX_DATASET_SERVICE_URL"`
	GeographyControllerURL          string        `envconfig:"GEOGRAPHY_CONTROLLER_URL"`
	GeographyEnabled                bool          `envconfig:"GEOGRAPHY_ENABLED"`
	HealthcheckCriticalTimeout      time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	HealthcheckInterval             time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HomepageControllerURL           string        `envconfig:"HOMEPAGE_CONTROLLER_URL"`
	HTTPMaxConnections              int           `envconfig:"HTTP_MAX_CONNECTIONS"`
	LegacySearchRedirectsEnabled    bool          `envconfig:"LEGACY_SEARCH_REDIRECTS_ENABLED"`
	NewDatasetRoutingEnabled        bool          `envconfig:"NEW_DATASET_ROUTING_ENABLED"`
	PatternLibraryAssetsPath        string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	ProxyTimeout                    time.Duration `envconfig:"PROXY_TIMEOUT"`
	RedirectSecret                  string        `envconfig:"REDIRECT_SECRET" json:"-"`
	ReleaseCalendarControllerURL    string        `envconfig:"RELEASE_CALENDAR_CONTROLLER_URL"`
	ReleaseCalendarEnabled          bool          `envconfig:"RELEASE_CALENDAR_ENABLED"`
	ReleaseCalendarRoutePrefix      string        `envconfig:"RELEASE_CALENDAR_ROUTE_PREFIX"`
	ReleaseCalendarABTestPercentage int           `envconfig:"RELEASE_CALENDAR_AB_TEST_PERCENTAGE"`
	SearchControllerURL             string        `envconfig:"SEARCH_CONTROLLER_URL"`
	SearchRoutesEnabled             bool          `envconfig:"SEARCH_ROUTES_ENABLED"`
	SiteDomain                      string        `envconfig:"SITE_DOMAIN"`
	SQSAnalyticsURL                 string        `envconfig:"SQS_ANALYTICS_URL"`
	ZebedeeRequestMaximumRetries    int           `envconfig:"ZEBEDEE_REQUEST_MAXIMUM_RETRIES"`
	ZebedeeRequestMaximumTimeout    time.Duration `envconfig:"ZEBEDEE_REQUEST_TIMEOUT_SECONDS"`
}

type AWS struct {
	AccessKeyID     string `envconfig:"AWS_ACCESS_KEY_ID"      json:"-"`
	Region          string `envconfig:"AWS_REGION"`
	SecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY"  json:"-"`
}

var cfg *Config

// Get returns the default config with any modifications made through environment variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg := &Config{
		APIRouterURL:                 "http://localhost:23200/v1",
		AreaProfilesControllerURL:    "http://localhost:26600",
		AreaProfilesRoutesEnabled:    false,
		BabbageURL:                   "http://localhost:8080",
		BindAddr:                     ":20000",
		CensusAtlasRoutesEnabled:     false,
		CensusAtlasURL:               "http://localhost:28100",
		ContentTypeByteLimit:         5000000,
		CookiesControllerURL:         "http://localhost:24100",
		DatasetControllerURL:         "http://localhost:20200",
		DatasetFinderEnabled:         false,
		DownloaderURL:                "http://localhost:23400",
		FeedbackControllerURL:        "http://localhost:25200",
		FeedbackEnabled:              false,
		FilterDatasetControllerURL:   "http://localhost:20001",
		FilterFlexDatasetServiceURL:  "http://localhost:20100",
		GeographyControllerURL:       "http://localhost:23700",
		GeographyEnabled:             false,
		HealthcheckCriticalTimeout:   90 * time.Second,
		HealthcheckInterval:          30 * time.Second,
		HomepageControllerURL:        "http://localhost:24400",
		HTTPMaxConnections:           0,
		LegacySearchRedirectsEnabled: false,
		NewDatasetRoutingEnabled:     false,
		PatternLibraryAssetsPath:     "https://cdn.ons.gov.uk/sixteens/f816ac8",
		ProxyTimeout:                 5 * time.Second,
		RedirectSecret:               "secret",
		ReleaseCalendarControllerURL: "http://localhost:27700",
		ReleaseCalendarEnabled:       false,
		SearchControllerURL:          "http://localhost:25000",
		SearchRoutesEnabled:          true,
		SiteDomain:                   "ons.gov.uk",
		SQSAnalyticsURL:              "",
		ZebedeeRequestMaximumRetries: 0,
		ZebedeeRequestMaximumTimeout: 5 * time.Second,
	}

	cfg.AWS = AWS{
		AccessKeyID:     "",
		Region:          "eu-west-2",
		SecretAccessKey: "",
	}

	if err := envconfig.Process("", cfg); err != nil {
		return cfg, err
	}

	cfg.ReleaseCalendarRoutePrefix = validatePrivatePrefix(cfg.ReleaseCalendarRoutePrefix)

	return cfg, nil
}

// IsEnabledRelCalABTest checks whether ab test is enabled and that percentage is a sensible int value
func IsEnabledRelCalABTest(cfg Config) bool {
	percentage := cfg.ReleaseCalendarABTestPercentage
	if cfg.EnableReleaseCalendarABTest && percentage > 0 && percentage < 100 {
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
