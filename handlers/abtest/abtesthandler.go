package abtest

import (
	"net/http"
	"time"

	"github.com/ONSdigital/dp-cookies/cookies"
)

const (
	SearchTestCookieAspect = "dsa-search"
	SearchNewExit          = "exit-new-search"

	RelcalTestCookieAspect = "dsa-release-calendar"
	RelcalNewExit          = "exit-new-relcal"
)

func ABTestHandler(new, old http.Handler, percentage int, aspectID, domain, exitNew string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		now := time.Now().UTC()

		if _, ok := req.URL.Query()[exitNew]; ok {
			cookies.HandleABTestExit(w, req, old, aspectID, domain)
			return
		}

		aspect := cookies.GetABTestCookieAspect(req, aspectID)

		if (aspect.New.IsZero() && aspect.Old.IsZero()) || (aspect.New.Before(now) && aspect.Old.Before(now)) {
			cookies.HandleCookieAndServ(w, req, new, old, aspectID, domain, cookies.DefaultABTestRandomiser(percentage))
			return
		}

		cookies.ServABTest(w, req, new, old, aspect)
	})
}

func ABTestPurgeHandler(new http.Handler, aspectID, domain string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookies.RemoveABTestCookieAspect(w, req, aspectID, domain)
		new.ServeHTTP(w, req)
	})
}
