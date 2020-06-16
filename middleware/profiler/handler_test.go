package profiler

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {

	token := "123"
	middleware := Middleware(token)
	mockHandlerCalled := false
	handlerFunc := func(w http.ResponseWriter, r *http.Request) { mockHandlerCalled = true }
	mockHandler := http.HandlerFunc(handlerFunc)
	handler := middleware(mockHandler)

	Convey("Given a request without a pprof auth header", t, func() {

		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/", nil)

		Convey("When the profiler middleware is called", func() {

			handler.ServeHTTP(responseRecorder, request)

			Convey("Then the underlying handler is not called", func() {
				So(mockHandlerCalled, ShouldBeFalse)
			})
		})
	})

	Convey("Given a request with a pprof auth header", t, func() {

		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/", nil)
		request.Header.Set("Authorization", "Bearer 123")

		Convey("When the profiler middleware is called", func() {

			handler.ServeHTTP(responseRecorder, request)

			Convey("Then the underlying handler is not called", func() {
				So(mockHandlerCalled, ShouldBeTrue)
			})
		})
	})
}
