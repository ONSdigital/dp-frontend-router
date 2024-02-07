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
				So(cfg.SearchRoutesEnabled, ShouldBeTrue)
				So(cfg.ReleaseCalendarControllerURL, ShouldEqual, "http://localhost:27700")
				So(cfg.ReleaseCalendarEnabled, ShouldBeFalse)
				So(cfg.ReleaseCalendarRoutePrefix, ShouldEqual, "")
				So(cfg.LegacySearchRedirectsEnabled, ShouldBeFalse)
				So(cfg.APIRouterURL, ShouldEqual, "http://localhost:23200/v1")
				So(cfg.DownloaderURL, ShouldEqual, "http://localhost:23400")
				So(cfg.AreaProfilesControllerURL, ShouldEqual, "http://localhost:26600")
				So(cfg.AreaProfilesRoutesEnabled, ShouldBeFalse)
				So(cfg.FilterFlexDatasetServiceURL, ShouldEqual, "http://localhost:20100")
				So(cfg.PatternLibraryAssetsPath, ShouldEqual, "https://cdn.ons.gov.uk/sixteens/f816ac8")
				So(cfg.SiteDomain, ShouldEqual, "ons.gov.uk")
				So(cfg.RedirectSecret, ShouldEqual, "secret")
				So(cfg.SQSAnalyticsURL, ShouldEqual, "")
				So(cfg.ContentTypeByteLimit, ShouldEqual, 5000000)
				So(cfg.HealthcheckInterval, ShouldEqual, 30*time.Second)
				So(cfg.HealthcheckCriticalTimeout, ShouldEqual, 90*time.Second)
				So(cfg.ZebedeeRequestMaximumTimeout, ShouldEqual, 5*time.Second)
				So(cfg.ZebedeeRequestMaximumRetries, ShouldEqual, 0)
				So(cfg.ProxyTimeout, ShouldEqual, 5*time.Second)
			})
		})
	})
}

func TestValidatePrivatePrefix(t *testing.T) {
	Convey("given an empty private path prefix", t, func() {
		prefix := ""
		Convey("validatePrivatePrefix returns the empty private path prefix", func() {
			So(validatePrivatePrefix(prefix), ShouldEqual, "")
		})
	})

	Convey("given a private path prefix is set without an initial '/'", t, func() {
		prefix := "a-prefix"
		Convey("validatePrivatePrefix return the given private path prefix with an initial '/'", func() {
			So(validatePrivatePrefix(prefix), ShouldEqual, "/a-prefix")
		})
	})

	Convey("given a private path prefix is set with an initial '/'", t, func() {
		prefix := "/a-prefix"
		Convey("validatePrivatePrefix returns the given private path prefix as valid", func() {
			So(validatePrivatePrefix(prefix), ShouldEqual, "/a-prefix")
		})
	})
}
