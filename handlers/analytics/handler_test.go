package analytics

import (
	"github.com/ONSdigital/dp-frontend-router/analytics"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const validURL = "http://localhost:20000/redir?url=/economy/inflationandpriceindices/bulletins/consumerpriceinflation/december2015&pageIndex=1&linkIndex=1&term=cpi"

// Mock
type mockAnalyticsService struct {
	args []*analytics.Model
}

func (m *mockAnalyticsService) CaptureAnalyticsData(analytics *analytics.Model) string {
	m.args = append(m.args, analytics)
	return analytics.GetURL()
}

type MockArgs struct {
	response http.ResponseWriter
	request  *http.Request
	urlStr   string
	code     int
}

type MockHttpRedir struct {
	args []*MockArgs
}

func (m *MockHttpRedir) mockRedirector(w http.ResponseWriter, r *http.Request, urlStr string, code int) {
	m.args = append(m.args, &MockArgs{w, r, urlStr, code})
}

func TestCaptureSearchStats(t *testing.T) {
	requestedURL, _ := url.Parse(validURL)

	Convey("Given valid input parameters", t, func() {

		analyticsServiceMock := &mockAnalyticsService{make([]*analytics.Model, 0)}
		analyticsService = analyticsServiceMock

		mockRedir := &MockHttpRedir{args: make([]*MockArgs, 0)}
		redirector = mockRedir.mockRedirector

		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", requestedURL.RequestURI(), nil)

		Convey("When the search redirect handler is invoked", func() {
			HandleSearch(resp, req)

			Convey("Then analyticsService is called 1 time with the expected parameters", func() {
				So(len(analyticsServiceMock.args), ShouldEqual, 1)
				So(analyticsServiceMock.args[0], ShouldResemble, analytics.NewAnalyticsModel(requestedURL))
			})

			Convey("And Redirect is called 1 time with the expected parameters. ", func() {
				So(len(mockRedir.args), ShouldEqual, 1)
				args := mockRedir.args[0]
				So(args.response, ShouldResemble, resp)
				So(args.request, ShouldResemble, req)
				So(args.urlStr, ShouldResemble, requestedURL.Query().Get("url"))
				So(args.code, ShouldResemble, http.StatusTemporaryRedirect)
				So(len(analyticsServiceMock.args), ShouldEqual, 1)
			})
		})
	})
}
