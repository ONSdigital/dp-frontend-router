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
				So(cfg.CookiesControllerURL, ShouldEqual, "http://localhost:24100")
				So(cfg.DatasetControllerURL, ShouldEqual, "http://localhost:20200")
				So(cfg.FilterDatasetControllerURL, ShouldEqual, "http://localhost:20001")
				So(cfg.FeedbackControllerURL, ShouldEqual, "http://localhost:25200")
				So(cfg.FeedbackEnabled, ShouldBeFalse)
				So(cfg.HTTPMaxConnections, ShouldEqual, 0)
				So(cfg.SearchControllerURL, ShouldEqual, "http://localhost:25000")
				So(cfg.DataAggregationPagesEnabled, ShouldBeFalse)
				So(cfg.TopicAggregationPagesEnabled, ShouldBeFalse)
				So(cfg.RelatedDataRouteEnabled, ShouldBeFalse)
				So(cfg.SearchRoutesEnabled, ShouldBeTrue)
				So(cfg.ReleaseCalendarControllerURL, ShouldEqual, "http://localhost:27700")
				So(cfg.LegacySearchRedirectsEnabled, ShouldBeFalse)
				So(cfg.APIRouterURL, ShouldEqual, "http://localhost:23200/v1")
				So(cfg.DownloaderURL, ShouldEqual, "http://localhost:23400")
				So(cfg.SiteDomain, ShouldEqual, "ons.gov.uk")
				So(cfg.ContentTypeByteLimit, ShouldEqual, 5000000)
				So(cfg.HealthcheckInterval, ShouldEqual, 30*time.Second)
				So(cfg.HealthcheckCriticalTimeout, ShouldEqual, 90*time.Second)
				So(cfg.ZebedeeRequestMaximumTimeout, ShouldEqual, 5*time.Second)
				So(cfg.ZebedeeRequestMaximumRetries, ShouldEqual, 0)
				So(cfg.ProxyTimeout, ShouldEqual, 5*time.Second)
				So(cfg.LegacyCacheProxyEnabled, ShouldBeFalse)
				So(cfg.LegacyCacheProxyURL, ShouldEqual, "http://localhost:29200")
			})
		})
	})
}
