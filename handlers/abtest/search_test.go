package abtest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-router/handlers/abtest"
	"github.com/ONSdigital/dp-frontend-router/handlers/abtest/searchtest"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/pat"
	. "github.com/smartystreets/goconvey/convey"
)

func NewHandlerMock() *searchtest.HandlerMock {
	return &searchtest.HandlerMock{
		ServeHTTPFunc: func(in1 http.ResponseWriter, in2 *http.Request) {},
	}
}

func TestABSearchHandler(t *testing.T) {

	domain := "www.ons.gov.uk"
	now := time.Now().UTC()
	tomorrow := now.Add(24 * time.Duration(time.Hour))

	Convey("SearchHandler", t, func() {
		newSearch, oldSearch := NewHandlerMock(), NewHandlerMock()
		// _ = abtest.SearchHandler(newSearch, oldSearch, 1, domain)

		// TODO call abtesthandler

		Convey("request to old or new search if no cookie is set", func() {
			req := httptest.NewRequest("GET", "/search", nil)
			res := httptest.NewRecorder()

			router := mockRouter(newSearch, oldSearch, 10, domain)
			router.ServeHTTP(res, req)
			So(len(oldSearch.ServeHTTPCalls())+len(newSearch.ServeHTTPCalls()), ShouldBeGreaterThan, 0)
		})

		Convey("request to old search if cookie is set to old search", func() {
			req := httptest.NewRequest("GET", "/search", nil)
			testServices := cookies.ABServices{NewSearch: &now, OldSearch: &tomorrow}
			testMarshalledValue, err := json.Marshal(testServices)
			if err != nil {
				t.Fatal("failed to marshal test data", err)
			}
			testEscapedValue := url.QueryEscape(string(testMarshalledValue))
			cookie := http.Cookie{Name: "ab_test", Value: testEscapedValue}
			req.AddCookie(&cookie)
			res := httptest.NewRecorder()

			router := mockRouter(newSearch, oldSearch, 10, domain)
			router.ServeHTTP(res, req)

			So(len(oldSearch.ServeHTTPCalls()), ShouldEqual, 1)
			So(oldSearch.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, "/search")

		})

		Convey("request to new search if cookie is set to new search", func() {
			req := httptest.NewRequest("GET", "/search", nil)
			testServices := cookies.ABServices{NewSearch: &tomorrow, OldSearch: &now}
			testMarshalledValue, err := json.Marshal(testServices)
			if err != nil {
				t.Fatal("failed to marshal test data", err)
			}
			testEscapedValue := url.QueryEscape(string(testMarshalledValue))
			cookie := http.Cookie{Name: "ab_test", Value: testEscapedValue}
			req.AddCookie(&cookie)
			res := httptest.NewRecorder()

			router := mockRouter(newSearch, oldSearch, 10, domain)
			router.ServeHTTP(res, req)

			So(len(newSearch.ServeHTTPCalls()), ShouldEqual, 1)
			So(newSearch.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, "/search")
		})
	})

	Convey("randomiseABTestCookie returns the correct result", t, func() {
		Convey("sets new search for twenty for hours when set to 100 of traffic", func() {
			result := abtest.RandomiseABTestCookie(100, now)
			spew.Dump(result)
			So(result, ShouldResemble, cookies.ABServices{NewSearch: &tomorrow, OldSearch: &now})
		})

		Convey("sets old search for twenty for hours when set to 0 of traffic", func() {
			result := abtest.RandomiseABTestCookie(0, now)
			spew.Dump(result)
			So(result, ShouldResemble, cookies.ABServices{NewSearch: &now, OldSearch: &tomorrow})
		})
	})
}

func mockRouter(searchHandler, babbageHandler http.Handler, newSearchABTestPercentage int, siteDomain string) http.Handler {
	router := pat.New()

	router.Handle("/search", abtest.SearchHandler(searchHandler, babbageHandler, newSearchABTestPercentage, siteDomain))

	return router
}
