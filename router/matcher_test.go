package router_test

import (
	"testing"

	"github.com/ONSdigital/dp-frontend-router/router"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHasFileExt(t *testing.T) {
	Convey("Given a URL that has a file extension", t, func() {
		Convey("When the HasFileExt function is called", func() {
			Convey("Then true is returned", func() {
				So(router.HasFileExt("/some/path.json"), ShouldBeTrue)
				So(router.HasFileExt("/main.css"), ShouldBeTrue)
				So(router.HasFileExt("/the/main/site.js"), ShouldBeTrue)
			})
		})
	})
	Convey("Given a URL that has no file extension", t, func() {
		Convey("When the HasFileExt function is called", func() {
			Convey("Then false is returned", func() {
				So(router.HasFileExt("/some/path"), ShouldBeFalse)
				So(router.HasFileExt("/"), ShouldBeFalse)
				So(router.HasFileExt("/test/module.test/data"), ShouldBeFalse)
			})
		})
	})
}

func TestIsKnownBabbageEndpoint(t *testing.T) {
	Convey("Given a URL that is for a known babbage endpoint", t, func() {
		Convey("When the IsKnownBabbageEndpoint function is called", func() {
			Convey("Then true is returned", func() {
				So(router.IsKnownBabbageEndpoint("/chartimage"), ShouldBeTrue)
				So(router.IsKnownBabbageEndpoint("/visualisations/the/name"), ShouldBeTrue)
				So(router.IsKnownBabbageEndpoint("/generator"), ShouldBeTrue)
				So(router.IsKnownBabbageEndpoint("/chartconfig"), ShouldBeTrue)
				So(router.IsKnownBabbageEndpoint("/search"), ShouldBeTrue)
				So(router.IsKnownBabbageEndpoint("/resource"), ShouldBeTrue)
				So(router.IsKnownBabbageEndpoint("/file"), ShouldBeTrue)
				So(router.IsKnownBabbageEndpoint("/ons/some/more/url"), ShouldBeTrue)
				So(router.IsKnownBabbageEndpoint("/timeseriestool"), ShouldBeTrue)
				So(router.IsKnownBabbageEndpoint("/economy/environmentalaccounts/bulletins/ukenvironmentalaccounts/latest"), ShouldBeTrue)
				So(router.IsKnownBabbageEndpoint("/economy/environmentalaccounts/bulletins/ukenvironmentalaccounts/2020/data"), ShouldBeTrue)
			})
		})
	})
	Convey("Given a URL that is not for a known babbage endpoint", t, func() {
		Convey("When the IsKnownBabbageEndpoint function is called", func() {
			Convey("Then false is returned", func() {
				So(router.IsKnownBabbageEndpoint("/some/path"), ShouldBeFalse)
				So(router.IsKnownBabbageEndpoint("/"), ShouldBeFalse)
				So(router.IsKnownBabbageEndpoint("/files/123"), ShouldBeFalse)
			})
		})
	})
}
