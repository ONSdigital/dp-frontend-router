package analytics

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.
var hmacSampleSecret []byte = []byte("secret")
var hmacBadSampleSecret []byte = []byte("bad")

func Test_extractIntParam(t *testing.T) {
	s := &ServiceImpl{
		redirectSecret: "secret",
	}
	Convey("Given a valid redirect URL", t, func() {
		data := "eyJhbGciOiJIUzI1NiJ9.eyJpbmRleCI6MSwicGFnZVNpemUiOjEwLCJ0ZXJtIjoiSW50ZWdyYXRlZCIsInBhZ2UiOjEsInVyaSI6Ii9wZW9wbGVwb3B1bGF0aW9uYW5kY29tbXVuaXR5L2hvdXNpbmcvYnVsbGV0aW5zL2ludGVncmF0ZWRob3VzZWhvbGRzdXJ2ZXlleHBlcmltZW50YWxzdGF0aXN0aWNzLzIwMTQtMTAtMDciLCJsaXN0VHlwZSI6InNlYXJjaCJ9.MQnW73Zca_7DZbYXjQC9FMIbCiJjNe--AKcCpLU2azw"
		r, err := http.NewRequest("GET", "/redir/"+data, nil)
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/redir/{data:.*}", func(writer http.ResponseWriter, request *http.Request) {
			url, err := s.CaptureAnalyticsData(request)
			So(err, ShouldBeNil)
			So(url, ShouldEqual, "/peoplepopulationandcommunity/housing/bulletins/integratedhouseholdsurveyexperimentalstatistics/2014-10-07")
		})
		router.ServeHTTP(rr, r)
	})

	Convey("Given a valid redirect URL, 2nd version", t, func() {
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"index":    1,
			"pageSize": 10,
			"term":     "Integrated",
			"page":     1,
			"uri":      "/peoplepopulationandcommunity/housing/bulletins/integratedhouseholdsurveyexperimentalstatistics/2014-10-07",
			"listType": "search",
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		So(err, ShouldBeNil)

		r, err := http.NewRequest("GET", "/redir/"+tokenString, nil)
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/redir/{data:.*}", func(writer http.ResponseWriter, request *http.Request) {
			url, err := s.CaptureAnalyticsData(request)
			So(err, ShouldBeNil)
			So(url, ShouldEqual, "/peoplepopulationandcommunity/housing/bulletins/integratedhouseholdsurveyexperimentalstatistics/2014-10-07")
		})
		router.ServeHTTP(rr, r)
	})

	Convey("Given a valid redirect URL, but payload has no uri", t, func() {
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"index":    1,
			"pageSize": 10,
			"term":     "Integrated",
			"page":     1,
			"listType": "search",
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		So(err, ShouldBeNil)

		r, err := http.NewRequest("GET", "/redir/"+tokenString, nil)
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/redir/{data:.*}", func(writer http.ResponseWriter, request *http.Request) {
			_, err = s.CaptureAnalyticsData(request)
			So(err, ShouldNotBeNil)
			es := fmt.Sprintf("%s", err)
			So(es, ShouldEqual, "url is a mandatory parameter")
		})
		router.ServeHTTP(rr, r)
	})

	Convey("Given a valid redirect URL, with bad secret", t, func() {
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"index": 1,
		})

		// Sign and get the complete encoded token as a string using the BAD secret
		tokenString, err := token.SignedString(hmacBadSampleSecret)
		So(err, ShouldBeNil)

		r, err := http.NewRequest("GET", "/redir/"+tokenString, nil)
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/redir/{data:.*}", func(writer http.ResponseWriter, request *http.Request) {
			_, err = s.CaptureAnalyticsData(request)
			So(err, ShouldNotBeNil)
			es := fmt.Sprintf("%s", err)
			// The following contains two errors, ONS's + lib error
			// We can only look for ONS's error, as the library error might change
			So(es, ShouldContainSubstring, "Invalid redirect data")
		})
		router.ServeHTTP(rr, r)
	})

	Convey("Given a valid redirect URL, and nil aud", t, func() {
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"index": 1,
			"aud":   nil,
			"uri":   "/peoplepopulationandcommunity",
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		So(err, ShouldBeNil)

		r, err := http.NewRequest("GET", "/redir/"+tokenString, nil)
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/redir/{data:.*}", func(writer http.ResponseWriter, request *http.Request) {
			url, err := s.CaptureAnalyticsData(request)
			So(err, ShouldBeNil)
			So(url, ShouldEqual, "/peoplepopulationandcommunity")
		})
		router.ServeHTTP(rr, r)
	})

	Convey("Given a valid redirect URL, and Empty aud", t, func() {
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"index": 1,
			"aud":   "",
			"uri":   "/peoplepopulationandcommunity",
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		So(err, ShouldBeNil)

		r, err := http.NewRequest("GET", "/redir/"+tokenString, nil)
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/redir/{data:.*}", func(writer http.ResponseWriter, request *http.Request) {
			url, err := s.CaptureAnalyticsData(request)
			So(err, ShouldBeNil)
			So(url, ShouldEqual, "/peoplepopulationandcommunity")
		})
		router.ServeHTTP(rr, r)
	})

	Convey("Given a valid redirect URL, and aud has a string", t, func() {
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"index": 1,
			"aud":   "hello",
			"uri":   "/peoplepopulationandcommunity",
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		So(err, ShouldBeNil)

		r, err := http.NewRequest("GET", "/redir/"+tokenString, nil)
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/redir/{data:.*}", func(writer http.ResponseWriter, request *http.Request) {
			url, err := s.CaptureAnalyticsData(request)
			So(err, ShouldBeNil)
			So(url, ShouldEqual, "/peoplepopulationandcommunity")
		})
		router.ServeHTTP(rr, r)
	})

	Convey("Given a valid redirect URL, and aud has no elements in []string", t, func() {
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"index": 1,
			"aud":   []string{},
			"uri":   "/peoplepopulationandcommunity",
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		So(err, ShouldBeNil)

		r, err := http.NewRequest("GET", "/redir/"+tokenString, nil)
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/redir/{data:.*}", func(writer http.ResponseWriter, request *http.Request) {
			url, err := s.CaptureAnalyticsData(request)
			So(err, ShouldBeNil)
			So(url, ShouldEqual, "/peoplepopulationandcommunity")
		})
		router.ServeHTTP(rr, r)
	})

	Convey("Given a valid redirect URL, and aud has one element in []string", t, func() {
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"index": 1,
			"aud":   []string{"hello2"},
			"uri":   "/peoplepopulationandcommunity",
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		So(err, ShouldBeNil)

		r, err := http.NewRequest("GET", "/redir/"+tokenString, nil)
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/redir/{data:.*}", func(writer http.ResponseWriter, request *http.Request) {
			url, err := s.CaptureAnalyticsData(request)
			So(err, ShouldBeNil)
			So(url, ShouldEqual, "/peoplepopulationandcommunity")
		})
		router.ServeHTTP(rr, r)
	})

	Convey("Given a valid redirect URL, and aud has two elements in []string", t, func() {
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"index": 1,
			"aud":   []string{"hello3", "hello4"},
			"uri":   "/peoplepopulationandcommunity",
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		So(err, ShouldBeNil)

		r, err := http.NewRequest("GET", "/redir/"+tokenString, nil)
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/redir/{data:.*}", func(writer http.ResponseWriter, request *http.Request) {
			url, err := s.CaptureAnalyticsData(request)
			So(err, ShouldBeNil)
			So(url, ShouldEqual, "/peoplepopulationandcommunity")
		})
		router.ServeHTTP(rr, r)
	})
}

// The following code is borrowed from v3.2.2 of library:
// github.com/form3tech-oss/jwt-go
// and tweaked a little to demonstrate issues with the previously used library:
// github.com/dgrijalva/jwt-go
// BUT the code won't work with v4 of original library

func Test_mapClaims_list_aud(t *testing.T) {
	mapClaims := jwt.MapClaims{
		"aud": []string{"foo"},
	}
	want := true
	got := mapClaims.VerifyAudience("foo", true)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}
func Test_mapClaims_string_aud(t *testing.T) {
	mapClaims := jwt.MapClaims{
		"aud": "foo",
	}
	want := true
	got := mapClaims.VerifyAudience("foo", true)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}

func Test_mapClaims_list_aud_no_match(t *testing.T) {
	mapClaims := jwt.MapClaims{
		"aud": []string{"bar"},
	}
	want := false
	got := mapClaims.VerifyAudience("foo", true)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}
func Test_mapClaims_string_aud_fail(t *testing.T) {
	mapClaims := jwt.MapClaims{
		"aud": "bar",
	}
	want := false
	got := mapClaims.VerifyAudience("foo", true)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}

func Test_mapClaims_string_aud_no_claim(t *testing.T) {
	mapClaims := jwt.MapClaims{}
	want := false
	got := mapClaims.VerifyAudience("foo", true)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}

func Test_mapClaims_string_aud_no_claim_not_required(t *testing.T) {
	mapClaims := jwt.MapClaims{}
	want := false
	got := mapClaims.VerifyAudience("foo", false)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}
