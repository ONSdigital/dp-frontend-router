package analytics

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const validURL = "http://localhost:20000/redir?url=/economy/inflationandpriceindices/bulletins/consumerpriceinflation/december2015&pageIndex=1&linkIndex=1&term=cpi"

// Mock
type mockAnalyticsService struct {
	args          []url.Values
	mockBehaviour func(r *http.Request, m *mockAnalyticsService) (string, error)
}

func (m *mockAnalyticsService) CaptureAnalyticsData(r *http.Request) (string, error) {
	m.args = append(m.args, r.URL.Query())
	return m.mockBehaviour(r, m)
}

func successBehavior(r *http.Request, m *mockAnalyticsService) (string, error) {
	return r.URL.Query().Get("url"), nil
}

func errorBehavior(r *http.Request, m *mockAnalyticsService) (string, error) {
	return "", errors.New("Error!")
}

type MockRedirArgs struct {
	response http.ResponseWriter
	request  *http.Request
	urlStr   string
	code     int
}

type MockHttpRedir struct {
	args []*MockRedirArgs
}

func (m *MockHttpRedir) mockRedirector(w http.ResponseWriter, r *http.Request, urlStr string, code int) {
	m.args = append(m.args, &MockRedirArgs{w, r, urlStr, code})
}

func TestHandleSearch(t *testing.T) {
	requestedURL, _ := url.Parse(validURL)

	Convey("Given valid input parameters", t, func() {

		serviceMock := &mockAnalyticsService{
			args:          make([]url.Values, 0),
			mockBehaviour: successBehavior,
		}

		mockRedir := &MockHttpRedir{args: make([]*MockRedirArgs, 0)}
		sh := &searchHandler{
			service:    serviceMock,
			redirector: mockRedir.mockRedirector,
		}

		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", requestedURL.RequestURI(), nil)

		Convey("When the search redirect handler is invoked", func() {
			sh.ServeHTTP(resp, req)

			Convey("Then service is called 1 time with the expected parameters", func() {
				So(len(serviceMock.args), ShouldEqual, 1)
				So(serviceMock.args[0], ShouldResemble, requestedURL.Query())
			})

			Convey("And Redirect is called 1 time with the expected parameters. ", func() {
				So(len(mockRedir.args), ShouldEqual, 1)
				args := mockRedir.args[0]
				So(args.response, ShouldResemble, resp)
				So(args.request, ShouldResemble, req)
				So(args.urlStr, ShouldResemble, requestedURL.Query().Get("url"))
				So(args.code, ShouldResemble, http.StatusTemporaryRedirect)
				So(len(serviceMock.args), ShouldEqual, 1)
			})
		})

		Convey("When the handler is invoked and the service returns an error", func() {
			serviceMock := &mockAnalyticsService{
				args:          make([]url.Values, 0),
				mockBehaviour: errorBehavior,
			}
			sh.service = serviceMock
			resp := httptest.NewRecorder()
			req := httptest.NewRequest("GET", requestedURL.RequestURI(), nil)

			sh.ServeHTTP(resp, req)

			Convey("Then a BAD REQUEST status is returned.", func() {
				So(resp.Code, ShouldEqual, 400)
			})

			Convey("And the service is called 1 time with the expected args", func() {
				So(len(serviceMock.args), ShouldEqual, 1)
				So(serviceMock.args[0], ShouldResemble, requestedURL.Query())
			})

			Convey("And the redirected is never invoked", func() {
				So(len(mockRedir.args), ShouldEqual, 0)
			})
		})
	})
}
