package handlers

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"github.com/ONSdigital/dp-frontend-router/analytics"
)

const validURL = "http://localhost:20000/redir?url=/economy/inflationandpriceindices/bulletins/consumerpriceinflation/december2015&pageIndex=1&linkIndex=1&term=cpi"

type MockSearchStatsService struct {
	args []*analytics.Model
}

func (m *MockSearchStatsService) CaptureAndRedirect(searchStats *analytics.Model, w http.ResponseWriter, req *http.Request) {
	m.args = append(m.args, searchStats)
}

func TestCaptureSearchStats(t *testing.T) {
	requestedURL, _ := url.Parse(validURL)

	Convey("When the search redirect handler is invoked", t, func() {

		mock := &MockSearchStatsService{make([]*analytics.Model, 0)}
		searchAnalyticsService = mock

		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", requestedURL.RequestURI(), nil)
		CaptureSearchStats(resp, req)

		Convey("Then searchStatsService is called 1 time", func() {
			So(len(mock.args), ShouldEqual, 1)
		})

		Convey("And searchStatsService is called with the expected parameter", func() {
			So(mock.args[0], ShouldResemble, analytics.NewSearchAnalytics(requestedURL))
		})
	})
}
