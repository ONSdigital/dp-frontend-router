package config

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	Convey("Given an environment with no environment variables set", t, func() {
		cfg, err := Get()

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("The values should be set to the expected defaults", func() {
				So(cfg.BindAddr, ShouldEqual, ":20000")
				So(cfg.BabbageURL, ShouldEqual, "http://localhost:8080")
				So(cfg.RendererURL, ShouldEqual, "http://localhost:20010")
				So(cfg.CookiesControllerURL, ShouldEqual, "http://localhost:24100")
				So(cfg.DatasetControllerURL, ShouldEqual, "http://localhost:20200")
				So(cfg.FilterDatasetControllerURL, ShouldEqual, "http://localhost:20001")
				So(cfg.GeographyControllerURL, ShouldEqual, "http://localhost:23700")
				So(cfg.FeedbackControllerURL, ShouldEqual, "http://localhost:25200")
				So(cfg.FeedbackEnabled, ShouldBeFalse)
				So(cfg.GeographyEnabled, ShouldBeFalse)
				So(cfg.SearchControllerURL, ShouldEqual, "http://localhost:25000")
				So(cfg.SearchRoutesEnabled, ShouldBeFalse)
				So(cfg.InteractivesControllerURL, ShouldEqual, "http://localhost:27300")
				So(cfg.InteractivesRoutesEnabled, ShouldBeFalse)
				So(cfg.APIRouterURL, ShouldEqual, "http://localhost:23200/v1")
				So(cfg.DownloaderURL, ShouldEqual, "http://localhost:23400")
				So(cfg.AreaProfilesControllerURL, ShouldEqual, "http://localhost:26600")
				So(cfg.AreaProfilesRoutesEnabled, ShouldBeFalse)
				So(cfg.FilterFlexDatasetServiceURL, ShouldEqual, "http://localhost:20100")
				So(cfg.FilterFlexRoutesEnabled, ShouldBeFalse)
				So(cfg.PatternLibraryAssetsPath, ShouldEqual, "https://cdn.ons.gov.uk/sixteens/f816ac8")
				So(cfg.SiteDomain, ShouldEqual, "ons.gov.uk")
				So(cfg.RedirectSecret, ShouldEqual, "secret")
				So(cfg.SQSAnalyticsURL, ShouldEqual, "")
				So(cfg.ContentTypeByteLimit, ShouldEqual, 5000000)
				So(cfg.HealthcheckInterval, ShouldEqual, 30*time.Second)
				So(cfg.HealthcheckCriticalTimeout, ShouldEqual, 90*time.Second)
				So(cfg.ZebedeeRequestMaximumTimeoutSeconds, ShouldEqual, 5*time.Second)
				So(cfg.ZebedeeRequestMaximumRetries, ShouldEqual, 0)
				So(cfg.ProxyTimeout, ShouldEqual, 5*time.Second)
			})
		})
	})
}

func TestIsEnabledABSearch(t *testing.T) {
	Convey("IsEnabledABSearch returns expected value", t, func() {
		Convey("false when EnableSearchABTest is false", func() {
			cfg := Config{EnableSearchABTest: false, SearchABTestPercentage: 10}
			result := IsEnableSearchABTest(cfg)
			So(result, ShouldBeFalse)
		})
		Convey("false when SearchABTestPercentage is below 0", func() {
			cfg := Config{EnableSearchABTest: true, SearchABTestPercentage: -10}
			result := IsEnableSearchABTest(cfg)
			So(result, ShouldBeFalse)
		})
		Convey("false when SearchABTestPercentage is over 100", func() {
			cfg := Config{EnableSearchABTest: true, SearchABTestPercentage: 110}
			result := IsEnableSearchABTest(cfg)
			So(result, ShouldBeFalse)
		})
		Convey("true when EnableSearchABTest is set and a sensible percentage int is used", func() {
			cfg := Config{EnableSearchABTest: true, SearchABTestPercentage: 10}
			result := IsEnableSearchABTest(cfg)
			So(result, ShouldBeTrue)
		})
	})
}
