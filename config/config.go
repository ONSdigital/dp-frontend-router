package config

import (
	"encoding/json"

	"github.com/kelseyhightower/envconfig"
)

// Configuration structure which hold information for configuring the import API
type Configuration struct {
	BindAddr                   string `envconfig:"BIND_ADDR"`
	BabbageURL                 string `envconfig:"BABBAGE_URL"`
	ResolverURL                string `envconfig:"RESOLVER_URL"`
	RendererURL                string `envconfig:"RENDERER_URL"`
	DatasetControllerURL       string `envconfig:"DATASET_CONTROLLER_URL"`
	FilterDatasetControllerURL string `envconfig:"FILTER_DATASET_CONTROLLER_URL"`
	GeographyControllerURL     string `envconfig:"GEOGRAPHY_CONTROLLER_URL"`
	GeoEnabled                 bool   `envconfig:"GEOGRAPHY_ENABLED"`
	ZebedeeURL                 string `envconfig:"ZEBEDEE_URL"`
	DownloaderURL              string `envconfig:"DOWNLOADER_URL"`
	HomepageABPercent          int    `envconfig:"HOMEPAGE_AB_PERCENT"`
	DebugMode                  bool   `envconfig:"DEBUG"`
	PatternLibraryAssetsPath   string `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	SiteDomain                 string `envconfig:"SITE_DOMAIN"`
	SplashPage                 string `envconfig:"SPLASH_PAGE"`
	RedirectSecret             string `envconfig:"REDIRECT_SECRET"`
	DisabledPage               string `envconfig:"DISABLED_PAGE"`
	TaxonomyDomain             string `envconfig:"TAXONOMY_DOMAIN"`
	SQSAnalyticsURL            string `envconfig:"ANALYTICS_SQS_URL"`
}

var cfg *Configuration

// Get the application and returns the configuration structure
func Get() (*Configuration, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Configuration{
		BindAddr:                   ":20000",
		BabbageURL:                 "http://localhost:8080",
		ResolverURL:                "http://localhost:20020",
		RendererURL:                "http://localhost:20010",
		DatasetControllerURL:       "http://localhost:20200",
		FilterDatasetControllerURL: "http://localhost:20001",
		GeographyControllerURL:     "http://localhost:23700",
		GeoEnabled:                 false,
		ZebedeeURL:                 "http://localhost:8082",
		DownloaderURL:              "http://localhost:23400",
		HomepageABPercent:          0,
		DebugMode:                  false,
		PatternLibraryAssetsPath:   "https://cdn.ons.gov.uk/sixteens/6cc1837",
		SiteDomain:                 "ons.gov.uk",
		SplashPage:                 "",
		RedirectSecret:             "secret",
		DisabledPage:               "",
		TaxonomyDomain:             "",
		SQSAnalyticsURL:            "",
	}

	return cfg, envconfig.Process("", cfg)
}

// String is implemented to prevent sensitive fields being logged.
// The config is returned as JSON with sensitive fields omitted.
func (config Configuration) String() string {
	json, _ := json.Marshal(config)
	return string(json)
}
