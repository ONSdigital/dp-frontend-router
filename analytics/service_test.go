package analytics

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

const requestedURI = "/economy/inflationandpriceindices/bulletins/consumerpriceinflation/december2015"

const validURLBase = "http://localhost:20000/redir"

func Test_extractIntParam(t *testing.T) {
	Convey("Given a valid redirect URL", t, func() {
		requestedURL, _ := url.Parse(validURLBase + "?requestedURL=" + requestedURI + "&pageSize=10")

		Convey("When extractIntParam func is invoked with a parameter name", func() {
			intVal := extractIntParam(requestedURL, "pageSize")

			Convey("Then the requested value is returned as an int", func() {
				So(intVal, ShouldEqual, 10)
			})
		})
	})

	Convey("Given a valid redirect URL", t, func() {
		requestedURL, _ := url.Parse(validURLBase + "?requestedURL=" + requestedURI)

		Convey("When extractIntParam func is invoked for a parameter not in the query string", func() {
			intVal := extractIntParam(requestedURL, "pageIndex")

			Convey("Then the default value is returned.", func() {
				So(intVal, ShouldEqual, 0)
			})
		})
	})

	Convey("Given a valid redirect URL", t, func() {
		requestedURL, _ := url.Parse(validURLBase + "?requestedURL=" + requestedURI)

		Convey("When extractIntParam func is invoked for a parameter that is not an int", func() {
			intVal := extractIntParam(requestedURL, "requestedURL")

			Convey("Then the default value is returned.", func() {
				So(intVal, ShouldEqual, 0)
			})
		})
	})

}
