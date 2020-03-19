package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-frontend-router
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	BabbageURL                 string        `envconfig:"BABBAGE_URL"`
	RendererURL                string        `envconfig:"RENDERER_URL"`
	CookiesControllerURL       string        `envconfig:"COOKIES_CONTROLLER_URL"`
	CookiesRoutesEnabled       bool          `envconfig:"COOKIES_ROUTES_ENABLED"`
	DatasetRoutesEnabled       bool          `envconfig:"DATASET_ROUTES_ENABLED"`
	DatasetControllerURL       string        `envconfig:"DATASET_CONTROLLER_URL"`
	FilterDatasetControllerURL string        `envconfig:"FILTER_DATASET_CONTROLLER_URL"`
	GeographyControllerURL     string        `envconfig:"GEOGRAPHY_CONTROLLER_URL"`
	GeographyEnabled           bool          `envconfig:"GEOGRAPHY_ENABLED"`
	ZebedeeURL                 string        `envconfig:"ZEBEDEE_URL"`
	DownloaderURL              string        `envconfig:"DOWNLOADER_URL"`
	PatternLibraryAssetsPath   string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	SiteDomain                 string        `envconfig:"SITE_DOMAIN"`
	RedirectSecret             string        `envconfig:"REDIRECT_SECRET" json:"-"`
	SQSAnalyticsURL            string        `envconfig:"SQS_ANALYTICS_URL"`
	ContentTypeByteLimit       int           `envconfig:"CONTENT_TYPE_BYTE_LIMIT"`
	HealthckeckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	HealthckeckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
}

var cfg *Config

// Get returns the default config with any modifications made through environment variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg := &Config{
		BindAddr:                   ":20000",
		BabbageURL:                 "http://localhost:8080",
		RendererURL:                "http://localhost:20010",
		CookiesControllerURL:       "http://localhost:24100",
		CookiesRoutesEnabled:       false,
		DatasetRoutesEnabled:       false,
		DatasetControllerURL:       "http://localhost:20200",
		FilterDatasetControllerURL: "http://localhost:20001",
		GeographyControllerURL:     "http://localhost:23700",
		GeographyEnabled:           false,
		ZebedeeURL:                 "http://localhost:8082",
		DownloaderURL:              "http://localhost:23400",
		PatternLibraryAssetsPath:   "https://cdn.ons.gov.uk/sixteens/f816ac8",
		SiteDomain:                 "ons.gov.uk",
		RedirectSecret:             "secret",
		SQSAnalyticsURL:            "",
		ContentTypeByteLimit:       5000000,
		HealthckeckCriticalTimeout: time.Minute,
		HealthckeckInterval:        time.Second * 10,
	}

	return cfg, envconfig.Process("", cfg)
}
