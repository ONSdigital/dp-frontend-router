package analytics

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
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

type MockArgs struct {
	response http.ResponseWriter
	request  *http.Request
	urlStr   string
	code     int
}

type MockHttpRedir struct {
	args []*MockArgs
}

func (m *MockHttpRedir) mockHttpRedirector(w http.ResponseWriter, r *http.Request, urlStr string, code int) {
	m.args = append(m.args, &MockArgs{w, r, urlStr, code})
}

func TestStatsServiceImpl_CaptureAndRedirect(t *testing.T) {
	mock := &MockHttpRedir{make([]*MockArgs, 0)}
	service := &AnalyticsServiceImpl{Redirect: mock.mockHttpRedirector}

	Convey("Given a valid parameters", t, func() {
		url, _ := url.Parse(validURLBase + "?url=" + requestedURI + "&term=" + term + "&type=" + searchType + "&pageIndex=1" + "&linkIndex=1" + "&pageSize=10")
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("/GET", requestedURI, nil)
		analytics := NewSearchAnalytics(url)

		Convey("When CaptureAndRedirect is called.", func() {
			service.CaptureAndRedirect(analytics, resp, req)

			Convey("Then Redirect is called 1 time. ", func() {
				So(len(mock.args), ShouldEqual, 1)
			})

			Convey("And redirect is called with the expected parameters. ", func() {
				args := mock.args[0]
				So(args.response, ShouldResemble, resp)
				So(args.request, ShouldResemble, req)
				So(args.urlStr, ShouldResemble, analytics.url)
				So(args.code, ShouldResemble, http.StatusTemporaryRedirect)
				So(len(mock.args), ShouldEqual, 2)
			})
		})
	})
}

func TestNewSearchAnalytics(t *testing.T) {
	url, _ := url.Parse(validURLBase + "?url=" + requestedURI + "&term=" + term + "&type=" + searchType + "&pageIndex=1" + "&linkIndex=1" + "&pageSize=10")
	Convey("When NewSearchAnalytics is invoked with a valid redirect URL", t, func() {
		result := NewSearchAnalytics(url)

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

	Convey("When NewSearchAnalytics is invoked with a redirect URL missing a string parameter (url, term, type).", t, func() {
		url, _ := url.Parse(validURLBase + "?term=" + term + "&type=" + searchType + "&pageIndex=1" + "&linkIndex=1" + "&pageSize=10")
		result := NewSearchAnalytics(url)
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

	Convey("When NewSearchAnalytics is invoked with a redirect URL missing an int parameter(pageIndex, linkIndex, pageSize).", t, func() {
		url, _ := url.Parse(validURLBase + "?url=" + requestedURI + "&term=" + term + "&type=" + searchType + "&linkIndex=1" + "&pageSize=10")
		result := NewSearchAnalytics(url)
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

	Convey("When NewSearchAnalytics is invoked with a redirect URL with non int value for an int parameter (pageIndex, linkIndex, pageSize).", t, func() {
		url, _ := url.Parse(validURLBase + "?url=" + requestedURI + "&term=" + term + "&type=" + searchType + "&pageIndex=abcd" + "&linkIndex=1" + "&pageSize=10")
		result := NewSearchAnalytics(url)
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
