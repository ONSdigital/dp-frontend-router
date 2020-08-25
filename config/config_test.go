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
				So(cfg.DatasetRoutesEnabled, ShouldEqual, false)
				So(cfg.DatasetControllerURL, ShouldEqual, "http://localhost:20200")
				So(cfg.FilterDatasetControllerURL, ShouldEqual, "http://localhost:20001")
				So(cfg.GeographyControllerURL, ShouldEqual, "http://localhost:23700")
				So(cfg.GeographyEnabled, ShouldEqual, false)
				So(cfg.SearchControllerURL, ShouldEqual, "http://localhost:25000")
				So(cfg.SearchRoutesEnabled, ShouldEqual, false)
				So(cfg.ZebedeeURL, ShouldEqual, "http://localhost:8082")
				So(cfg.DownloaderURL, ShouldEqual, "http://localhost:23400")
				So(cfg.PatternLibraryAssetsPath, ShouldEqual, "https://cdn.ons.gov.uk/sixteens/f816ac8")
				So(cfg.SiteDomain, ShouldEqual, "ons.gov.uk")
				So(cfg.RedirectSecret, ShouldEqual, "secret")
				So(cfg.SQSAnalyticsURL, ShouldEqual, "")
				So(cfg.ContentTypeByteLimit, ShouldEqual, 5000000)
				So(cfg.HealthckeckInterval, ShouldEqual, 30*time.Second)
				So(cfg.HealthckeckCriticalTimeout, ShouldEqual, 90*time.Second)
			})
		})
	})
}
