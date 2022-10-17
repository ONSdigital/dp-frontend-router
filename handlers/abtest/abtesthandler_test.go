package abtest

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type user struct {
	new    bool
	cookie http.Cookie
}

func TestABTestHandler(t *testing.T) {
	Convey("Given an old request handler and a new request handler", t, func() {
		oldHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { _, _ = w.Write([]byte("Old Served")) })
		newHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { _, _ = w.Write([]byte("New Served")) })

		var users []user
		Convey("And a given percentage in the AB split of a given number of requests to the new handler", func() {
			percentage := 40
			numberRequests := 500

			Convey("When the requests are made to the ABTestHandler", func() {
				handler := ABTestHandler(newHandler, oldHandler, percentage, "test-aspect", "my-domain", "exit-new-test")
				users = make([]user, numberRequests)
				nh := 0
				for i := 0; i < numberRequests; i++ {
					req := httptest.NewRequest("GET", "/", nil)
					resp := httptest.NewRecorder()
					handler.ServeHTTP(resp, req)
					So(resp.Result().Cookies(), ShouldHaveLength, 1)
					users[i] = user{cookie: *resp.Result().Cookies()[0]}

					b, _ := ioutil.ReadAll(resp.Result().Body)
					if string(b) == "New Served" {
						users[i].new = true
						nh++
					}
				}

				Convey("The ABTest cookies are analysed (and stored) to ensure the number of requests serviced by the new handler are within an acceptable deviation range of the requested percentage split))", func() {
					onePercent := numberRequests / 100
					expected := percentage * onePercent
					deviation := func(x int) int {
						if x < 0 {
							return -x
						}
						return x
					}(expected - nh)
					So(deviation, ShouldBeBetweenOrEqual, 0, 20*onePercent)
				})

				Convey("And when subsequent requests are made with the previously returned cookies, all are serviced by the correct handler", func() {
					for i := 0; i < numberRequests; i++ {
						req := httptest.NewRequest("GET", "/", nil)
						req.AddCookie(&users[i].cookie)
						resp := httptest.NewRecorder()
						handler.ServeHTTP(resp, req)

						b, _ := ioutil.ReadAll(resp.Result().Body)
						if users[i].new {
							So(string(b), ShouldEqual, "New Served")
						} else {
							So(string(b), ShouldEqual, "Old Served")
						}
					}
				})

				Convey("But when subsequent requests which would ordinarily be handled by the new handler are sent with the 'exit-new-test' query parameter,"+
					"the requests all are serviced by the old handler", func() {
					for i := 0; i < numberRequests; i++ {
						if users[i].new {
							req := httptest.NewRequest("GET", "/?exit-new-test", nil)
							req.AddCookie(&users[i].cookie)
							resp := httptest.NewRecorder()
							handler.ServeHTTP(resp, req)

							b, _ := ioutil.ReadAll(resp.Result().Body)
							So(string(b), ShouldEqual, "Old Served")
						}
					}
				})
			})
		})
	})
}

func TestABTestPurgeHandler(t *testing.T) {
	Convey("Given a new request handler", t, func() {
		newHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { _, _ = w.Write([]byte("New Served")) })

		Convey("And an ab_test cookie containing a now unused aspect for ab testing of the new/old handler", func() {
			c := &http.Cookie{
				Name:  "ab_test",
				Value: url.QueryEscape(`{"test-aspect":{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"},"second-aspect":{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`),
			}

			Convey("When a request with the cookie is made to the ABTestPurgeHandler", func() {
				req := httptest.NewRequest("GET", "/", nil)
				req.AddCookie(c)
				w := httptest.NewRecorder()
				ABTestPurgeHandler(newHandler, "test-aspect", "my-domain").ServeHTTP(w, req)

				Convey("The relevant aspect is removed from the ab_test cookie, and the request has been handled by the new handler", func() {
					So(w.Result().Cookies(), ShouldHaveLength, 1)
					So(w.Result().Cookies()[0].Name, ShouldEqual, "ab_test")
					So(w.Result().Cookies()[0].Value, ShouldEqual, url.QueryEscape(`{"second-aspect":{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`))

					b, _ := ioutil.ReadAll(w.Result().Body)
					So(string(b), ShouldEqual, "New Served")
				})
			})
		})
	})
}
