package analytics

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_extractIntParam(t *testing.T) {
	s := &ServiceImpl{
		redirectSecret: "secret",
	}
	Convey("Given a valid redirect URL", t, func() {
		r, err := http.NewRequest("GET", "/redir/{:data}", nil)
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)

		q := r.URL.Query()
		q.Set(":data", "eyJhbGciOiJIUzI1NiJ9.eyJpbmRleCI6MSwicGFnZVNpemUiOjEwLCJ0ZXJtIjoiSW50ZWdyYXRlZCIsInBhZ2UiOjEsInVyaSI6Ii9wZW9wbGVwb3B1bGF0aW9uYW5kY29tbXVuaXR5L2hvdXNpbmcvYnVsbGV0aW5zL2ludGVncmF0ZWRob3VzZWhvbGRzdXJ2ZXlleHBlcmltZW50YWxzdGF0aXN0aWNzLzIwMTQtMTAtMDciLCJsaXN0VHlwZSI6InNlYXJjaCJ9.MQnW73Zca_7DZbYXjQC9FMIbCiJjNe--AKcCpLU2azw")
		r.URL.RawQuery = q.Encode()

		url, err := s.CaptureAnalyticsData(r)
		So(err, ShouldBeNil)
		So(url, ShouldEqual, "/peoplepopulationandcommunity/housing/bulletins/integratedhouseholdsurveyexperimentalstatistics/2014-10-07")
	})
}
