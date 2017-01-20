package analytics

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

const requestedURI = "/economy/inflationandpriceindices/bulletins/consumerpriceinflation/december2015"
const pageSize = 10
const searchType = "search"
const pageIndex = 1
const linkIndex = 1
const term = "cpi"

const validURLBase = "http://localhost:20000/redir"

func TestNewSearchAnalytics(t *testing.T) {
	url, _ := url.Parse(validURLBase + "?url=" + requestedURI + "&term=" + term + "&type=" + searchType + "&pageIndex=1" + "&linkIndex=1" + "&pageSize=10")
	Convey("When NewSearchAnalytics is invoked with a valid redirect URL", t, func() {
		result := NewAnalyticsModel(url)

		expected := &Model{
			url:        requestedURI,
			term:       term,
			searchType: searchType,
			pageIndex:  pageIndex,
			linkIndex:  linkIndex,
			pageSize:   pageSize,
		}

		Convey("Then then analytics are as expected", func() {
			So(result, ShouldResemble, expected)
		})
	})

	Convey("When NewAnalyticsModel is invoked with a redirect URL missing a string parameter (url, term, type).", t, func() {
		url, _ := url.Parse(validURLBase + "?term=" + term + "&type=" + searchType + "&pageIndex=1" + "&linkIndex=1" + "&pageSize=10")
		result := NewAnalyticsModel(url)
		expected := &Model{
			url:        "",
			term:       term,
			searchType: searchType,
			pageIndex:  pageIndex,
			linkIndex:  linkIndex,
			pageSize:   pageSize,
		}

		Convey("Then the missing field is the default string value.", func() {
			So(result, ShouldResemble, expected)
		})
	})

	Convey("When NewAnalyticsModel is invoked with a redirect URL missing an int parameter(pageIndex, linkIndex, pageSize).", t, func() {
		url, _ := url.Parse(validURLBase + "?url=" + requestedURI + "&term=" + term + "&type=" + searchType + "&linkIndex=1" + "&pageSize=10")
		result := NewAnalyticsModel(url)
		expected := &Model{
			url:        requestedURI,
			term:       term,
			searchType: searchType,
			pageIndex:  0,
			linkIndex:  linkIndex,
			pageSize:   pageSize,
		}

		Convey("Then the missing field is the default int value.", func() {
			So(result, ShouldResemble, expected)
		})
	})

	Convey("When NewAnalyticsModel is invoked with a redirect URL with non int value for an int parameter (pageIndex, linkIndex, pageSize).", t, func() {
		url, _ := url.Parse(validURLBase + "?url=" + requestedURI + "&term=" + term + "&type=" + searchType + "&pageIndex=abcd" + "&linkIndex=1" + "&pageSize=10")
		result := NewAnalyticsModel(url)
		expected := &Model{
			url:        requestedURI,
			term:       term,
			searchType: searchType,
			pageIndex:  0,
			linkIndex:  linkIndex,
			pageSize:   pageSize,
		}

		Convey("Then the missing field is the default int value.", func() {
			So(result, ShouldResemble, expected)
		})
	})
}
