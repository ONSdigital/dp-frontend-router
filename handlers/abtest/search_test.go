package abtest

import (
	"net/http"
	"testing"
	"time"

	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

func TestABSearchHandler(t *testing.T) {

	domain := "www.ons.gov.uk"

	Convey("SearchHandler", t, func() {
		oldSearch, newSearch := http.Handler, http.Handler
		result := SearchHandler(newSearch, oldSearch, 1, domain)
		spew.Dump(result)

	})

	Convey("randomiseABTestCookie returns the correct result", t, func() {
		Convey("sets new search for twenty for hours when set to 100 of traffic", func() {
			now := time.Now().UTC()
			tomorrow := now.Add(24 * time.Duration(time.Hour))
			result := randomiseABTestCookie(100, now)
			spew.Dump(result)
			So(result, ShouldResemble, cookies.ABServices{NewSearch: &tomorrow, OldSearch: &now})
		})

		Convey("sets old search for twenty for hours when set to 0 of traffic", func() {
			now := time.Now().UTC()
			tomorrow := now.Add(24 * time.Duration(time.Hour))
			result := randomiseABTestCookie(0, now)
			spew.Dump(result)
			So(result, ShouldResemble, cookies.ABServices{NewSearch: &now, OldSearch: &tomorrow})
		})
	})

}
