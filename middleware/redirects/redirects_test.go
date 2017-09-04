package redirects

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-frontend-router/handlers/serverError"
	"github.com/ONSdigital/go-ns/handlers/requestID"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRedirect(t *testing.T) {

	router := pat.New()
	middleware := []alice.Constructor{
		Handler,
	}
	alice := alice.New(middleware...).Then(router)
	var handled bool
	router.HandleFunc("/{uri:.*}", func(w http.ResponseWriter, req *http.Request) {
		handled = true
	})

	redirects["/redirect"] = "/redirected"

	Convey("Test that a non redirect request reaches the handler", t, func() {
		handled = false
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		alice.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, 200)
		So(handled, ShouldBeTrue)
	})

	Convey("Test that a redirect request returns a redirect", t, func() {
		handled = false
		req, _ := http.NewRequest("GET", "/redirect", nil)
		w := httptest.NewRecorder()
		alice.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, 307)
		So(handled, ShouldBeFalse)
		So(w.HeaderMap, ShouldContainKey, "Location")
	})

}

func TestInit(t *testing.T) {
	var shouldError bool
	var returnBytes []byte
	var called bool
	var asset = func(name string) ([]byte, error) {
		called = true
		if shouldError {
			return nil, errors.New("Error")
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
		So(func() { Init(asset) }, ShouldPanicWith, "Redirects must have two fields")
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
		So(func() { Init(asset) }, ShouldPanicWith, "Redirect from URL must not be empty")
	})

	Convey("Init should panic if redirect has no to url", t, func() {
		shouldError = false
		returnBytes = []byte(`a,`)
		So(func() { Init(asset) }, ShouldPanicWith, "Redirect to URL must not be empty")
	})

}

func BenchmarkWithoutRedirectMiddleware(b *testing.B) {

	router := pat.New()
	middleware := []alice.Constructor{
		requestID.Handler(16),
		log.Handler,
		//securityHandler,
		serverError.Handler,
	}
	alice := alice.New(middleware...).Then(router)
	router.HandleFunc("/{uri:.*}", func(w http.ResponseWriter, req *http.Request) {})
	req, _ := http.NewRequest("GET", "/", nil)

	for n := 0; n < b.N; n++ {
		alice.ServeHTTP(nil, req)
	}
}

func BenchmarkWithoutRedirects(b *testing.B) {

	router := pat.New()
	middleware := []alice.Constructor{
		requestID.Handler(16),
		log.Handler,
		//securityHandler,
		serverError.Handler,
		Handler,
	}
	alice := alice.New(middleware...).Then(router)
	router.HandleFunc("/{uri:.*}", func(w http.ResponseWriter, req *http.Request) {})
	req, _ := http.NewRequest("GET", "/", nil)

	for n := 0; n < b.N; n++ {
		alice.ServeHTTP(nil, req)
	}
}

func BenchmarkWith100Redirects(b *testing.B) {
	redirects = make(map[string]string)

	for i := 0; i < 100; i++ {
		redirects[fmt.Sprintf("/test/%d", i)] = "/"
	}

	router := pat.New()
	middleware := []alice.Constructor{
		requestID.Handler(16),
		log.Handler,
		//securityHandler,
		serverError.Handler,
		Handler,
	}
	alice := alice.New(middleware...).Then(router)
	router.HandleFunc("/{uri:.*}", func(w http.ResponseWriter, req *http.Request) {})
	req, _ := http.NewRequest("GET", "/", nil)

	for n := 0; n < b.N; n++ {
		alice.ServeHTTP(nil, req)
	}
}

func BenchmarkWith10000Redirects(b *testing.B) {
	redirects = make(map[string]string)

	for i := 0; i < 10000; i++ {
		redirects[fmt.Sprintf("/test/%d", i)] = "/"
	}

	router := pat.New()
	middleware := []alice.Constructor{
		requestID.Handler(16),
		log.Handler,
		//securityHandler,
		serverError.Handler,
		Handler,
	}
	alice := alice.New(middleware...).Then(router)
	router.HandleFunc("/{uri:.*}", func(w http.ResponseWriter, req *http.Request) {})
	req, _ := http.NewRequest("GET", "/", nil)

	for n := 0; n < b.N; n++ {
		alice.ServeHTTP(nil, req)
	}
}

func BenchmarkWith1000000Redirects(b *testing.B) {
	redirects = make(map[string]string)

	for i := 0; i < 1000000; i++ {
		redirects[fmt.Sprintf("/test/%d", i)] = "/"
	}

	router := pat.New()
	middleware := []alice.Constructor{
		requestID.Handler(16),
		log.Handler,
		//securityHandler,
		serverError.Handler,
		Handler,
	}
	alice := alice.New(middleware...).Then(router)
	router.HandleFunc("/{uri:.*}", func(w http.ResponseWriter, req *http.Request) {})
	req, _ := http.NewRequest("GET", "/", nil)

	for n := 0; n < b.N; n++ {
		alice.ServeHTTP(nil, req)
	}
}
