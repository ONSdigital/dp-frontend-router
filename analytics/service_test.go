package analytics

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"testing"
)

const requestedURI = "/economy/inflationandpriceindices/bulletins/consumerpriceinflation/december2015"
const pageSize = 10
const searchType = "search"
const pageIndex = 1
const linkIndex = 1
const term = "cpi"

const validURLBase = "http://localhost:20000/redir"

func slice(value string) []string {
	return []string{value}
}

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

func TestNewSearchAnalytics(t *testing.T) {
	requestedURL, _ := url.Parse(validURLBase + "?requestedURL=" + requestedURI + "&term=" + term + "&type=" + searchType + "&pageIndex=1" + "&linkIndex=1" + "&pageSize=10")
	Convey("When NewSearchAnalytics is invoked with a valid redirect URL", t, func() {

		query := requestedURL.Query()
		expected := url.Values{
			"requestedURL": slice(requestedURI),
			"term":         slice(term),
			"type":         slice(searchType),
			"pageIndex":    slice(strconv.Itoa(pageIndex)),
			"linkIndex":    slice(strconv.Itoa(linkIndex)),
			"pageSize":     slice(strconv.Itoa(pageSize)),
		}

		Convey("Then then analytics are as expected", func() {
			So(query, ShouldResemble, expected)
		})
	})
}
