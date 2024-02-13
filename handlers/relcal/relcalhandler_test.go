package relcal

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-frontend-router/config"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	newServed = "New Served"
	oldServed = "Old Served"
)

func TestRelcalHandler(t *testing.T) {
	Convey("Given an old request handler and a new request handler", t, func() {
		oldHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte(oldServed))
		})
		newHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte(newServed))
		})

		totalRequests := 10
		cfg, _ := config.Get()

		Convey("And the config  to use the new release calendar is set to false", func() {
			cfg.UseNewReleaseCalendar = false

			Convey("When the requests are made to the relcalHandler", func() {
				handler := Handler(cfg.UseNewReleaseCalendar, newHandler, oldHandler)

				nh := 0
				for i := 0; i < totalRequests; i++ {
					req := httptest.NewRequest("GET", "/", http.NoBody)
					resp := httptest.NewRecorder()
					handler.ServeHTTP(resp, req)
					b, _ := io.ReadAll(resp.Result().Body)
					if string(b) == oldServed {
						nh++
					}
				}

				Convey("All the requests should be to the old handler", func() {
					So(nh, ShouldEqual, totalRequests)
				})
			})
		})

		Convey("And the config  to use the new release calendar is set to true", func() {
			cfg.UseNewReleaseCalendar = true

			Convey("When the requests are made to the relcalHandler", func() {
				handler := Handler(cfg.UseNewReleaseCalendar, newHandler, oldHandler)

				nh := 0
				for i := 0; i < totalRequests; i++ {
					req := httptest.NewRequest("GET", "/", http.NoBody)
					resp := httptest.NewRecorder()
					handler.ServeHTTP(resp, req)
					b, _ := io.ReadAll(resp.Result().Body)
					if string(b) == newServed {
						nh++
					}
				}

				Convey("All the requests should be to the new handler", func() {
					So(nh, ShouldEqual, totalRequests)
				})
			})
		})
	})
}
