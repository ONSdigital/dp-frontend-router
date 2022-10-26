package redirects

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRedirect(t *testing.T) {
	router := mux.NewRouter()
	middleware := []alice.Constructor{
		Handler,
	}
	testAlice := alice.New(middleware...).Then(router)
	var handled bool
	router.HandleFunc("/{uri:.*}", func(w http.ResponseWriter, req *http.Request) {
		handled = true
	})

	redirects["/redirect"] = "/redirected"

	Convey("Test that a non redirect request reaches the handler", t, func() {
		handled = false
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		testAlice.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, 200)
		So(handled, ShouldBeTrue)
	})

	Convey("Test that a redirect request returns a redirect", t, func() {
		handled = false
		req, _ := http.NewRequest("GET", "/redirect", nil)
		w := httptest.NewRecorder()
		testAlice.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, 307)
		So(handled, ShouldBeFalse)
		So(w.Header(), ShouldContainKey, "Location")
	})
}

func TestDynamicRedirect(t *testing.T) {
	router := mux.NewRouter()
	middleware := []alice.Constructor{
		Handler,
	}
	testAlice := alice.New(middleware...).Then(router)
	router.Handle("/original{uri:.*}", DynamicRedirectHandler("/original", "/redirected"))
	router.HandleFunc("/redirected{uri:.*}", func(w http.ResponseWriter, req *http.Request) {
	})

	Convey("Test that a redirect request is redirected to the new url", t, func() {
		req, _ := http.NewRequest("GET", "/original", nil)
		w := httptest.NewRecorder()
		testAlice.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, 301)
		So(w.Header(), ShouldContainKey, "Location")
		So(w.Header()["Location"], ShouldContain, "/redirected")
	})

	Convey("Test that a redirect request with a path extension is redirected to the new url with the same path extension", t, func() {
		req, _ := http.NewRequest("GET", "/original/extension", nil)
		w := httptest.NewRecorder()
		testAlice.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, 301)
		So(w.Header(), ShouldContainKey, "Location")
		So(w.Header()["Location"], ShouldContain, "/redirected/extension")
	})

	Convey("Test that a redirect request with parameters is redirected to the new url with the same parameters", t, func() {
		req, _ := http.NewRequest("GET", "/original?q=test&page=2", nil)
		w := httptest.NewRecorder()
		testAlice.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, 301)
		So(w.Header(), ShouldContainKey, "Location")
		So(w.Header()["Location"], ShouldContain, "/redirected?q=test&page=2")
	})
}

func TestInit(t *testing.T) {
	var shouldError bool
	var returnBytes []byte
	var called bool
	var asset = func(name string) ([]byte, error) {
		called = true
		if shouldError {
			return nil, errors.New("error")
		}
		return returnBytes, nil
	}

	Convey("Init calls asset function and panics on error", t, func() {
		shouldError = true
		So(func() { Init(asset) }, ShouldPanicWith, "Can't find redirects.csv")
		So(called, ShouldBeTrue)
	})

	Convey("Init function doesn't panic if asset doesn't return error", t, func() {
		shouldError = false
		So(func() { Init(asset) }, ShouldNotPanic)
	})

	Convey("Init should panic on invalid CSV", t, func() {
		shouldError = false
		returnBytes = []byte(`a,b
a,b,c`)
		So(func() { Init(asset) }, ShouldPanicWith, "Unable to read CSV")
	})

	Convey("Init should panic if redirects doesn't have two fields", t, func() {
		shouldError = false
		returnBytes = []byte(`a
a`)
		So(func() { Init(asset) }, ShouldPanicWith, "redirect 'to' URL empty, check logs")
	})

	Convey("Init should add entries to redirects", t, func() {
		shouldError = false
		returnBytes = []byte(`a,b
c,d`)
		So(func() { Init(asset) }, ShouldNotPanic)
		So(redirects, ShouldContainKey, "a")
		So(redirects["a"], ShouldEqual, "b")
		So(redirects, ShouldContainKey, "c")
		So(redirects["c"], ShouldEqual, "d")
	})

	Convey("Init should panic if redirect has no from url", t, func() {
		shouldError = false
		returnBytes = []byte(`,b`)
		So(func() { Init(asset) }, ShouldPanicWith, "redirect 'from' URL empty, check logs")
	})

	Convey("Init should panic if redirect has no to url", t, func() {
		shouldError = false
		returnBytes = []byte(`a,`)
		So(func() { Init(asset) }, ShouldPanicWith, "redirect 'to' URL empty, check logs")
	})
}

func BenchmarkWithoutRedirectMiddleware(b *testing.B) {
	router := mux.NewRouter()
	middleware := []alice.Constructor{
		dprequest.HandlerRequestID(16),
		log.Middleware,
		// securityHandler,
	}
	testAlice := alice.New(middleware...).Then(router)
	router.HandleFunc("/{uri:.*}", func(w http.ResponseWriter, req *http.Request) {})
	req, _ := http.NewRequest("GET", "/", nil)

	for n := 0; n < b.N; n++ {
		testAlice.ServeHTTP(nil, req)
	}
}

func BenchmarkWithoutRedirects(b *testing.B) {
	router := mux.NewRouter()
	middleware := []alice.Constructor{
		dprequest.HandlerRequestID(16),
		log.Middleware,
		// securityHandler,
		Handler,
	}
	testAlice := alice.New(middleware...).Then(router)
	router.HandleFunc("/{uri:.*}", func(w http.ResponseWriter, req *http.Request) {})
	req, _ := http.NewRequest("GET", "/", nil)

	for n := 0; n < b.N; n++ {
		testAlice.ServeHTTP(nil, req)
	}
}

func BenchmarkWith100Redirects(b *testing.B) {
	redirects = make(map[string]string)

	for i := 0; i < 100; i++ {
		redirects[fmt.Sprintf("/test/%d", i)] = "/"
	}

	router := mux.NewRouter()
	middleware := []alice.Constructor{
		dprequest.HandlerRequestID(16),
		log.Middleware,
		// securityHandler,
		Handler,
	}
	testAlice := alice.New(middleware...).Then(router)
	router.HandleFunc("/{uri:.*}", func(w http.ResponseWriter, req *http.Request) {})
	req, _ := http.NewRequest("GET", "/", nil)

	for n := 0; n < b.N; n++ {
		testAlice.ServeHTTP(nil, req)
	}
}

func BenchmarkWith10000Redirects(b *testing.B) {
	redirects = make(map[string]string)

	for i := 0; i < 10000; i++ {
		redirects[fmt.Sprintf("/test/%d", i)] = "/"
	}

	router := mux.NewRouter()
	middleware := []alice.Constructor{
		dprequest.HandlerRequestID(16),
		log.Middleware,
		// securityHandler,
		Handler,
	}
	testAlice := alice.New(middleware...).Then(router)
	router.HandleFunc("/{uri:.*}", func(w http.ResponseWriter, req *http.Request) {})
	req, _ := http.NewRequest("GET", "/", nil)

	for n := 0; n < b.N; n++ {
		testAlice.ServeHTTP(nil, req)
	}
}

func BenchmarkWith1000000Redirects(b *testing.B) {
	redirects = make(map[string]string)

	for i := 0; i < 1000000; i++ {
		redirects[fmt.Sprintf("/test/%d", i)] = "/"
	}

	router := mux.NewRouter()
	middleware := []alice.Constructor{
		dprequest.HandlerRequestID(16),
		log.Middleware,
		// securityHandler,
		Handler,
	}
	testAlice := alice.New(middleware...).Then(router)
	router.HandleFunc("/{uri:.*}", func(w http.ResponseWriter, req *http.Request) {})
	req, _ := http.NewRequest("GET", "/", nil)

	for n := 0; n < b.N; n++ {
		testAlice.ServeHTTP(nil, req)
	}
}
