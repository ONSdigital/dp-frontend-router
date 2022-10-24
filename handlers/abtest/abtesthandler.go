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

func Handler(abTest bool, new, old http.Handler, percentage int, aspectID, domain, exitNew string) http.Handler {
	if abTest {
		return ABTestHandler(new, old, percentage, aspectID, domain, exitNew)
	} else {
		return ABTestPurgeHandler(new, aspectID, domain)
	}
}

// ABTestHandler routes requests to either the old or new handler, for a given aspectID, according to the given percentage
// i.e. for the given percentage of calls X, X% will be routed to the new handler, and the remainder to the old handler.
// Most of the functionality is provided by the dp-cookies library, which uses a single ab_test cookie to embed all aspects
// If the aspect does not exist or has expired, it is created/renewed according to a particular randomiser - in general
// the DefaultABTestRandomiser in the library is sufficient
// A well known string - the exitNew string -  can be used as a query parameter to the call, in order to definitively chose
// the old handler
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

// ABTestPurgeHandler is used to remove a given AspectID from the single ab_test cookie handled by the dp-cookies library
// It is useful when AB Testing for a particular aspect has finished, but the aspect is still embedded in client's ab_test
// cookie - this handler will remove the aspect and can be left in use for several weeks after testing has finished to 'clean'
// the underlying ab_test cookie
func ABTestPurgeHandler(new http.Handler, aspectID, domain string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookies.RemoveABTestCookieAspect(w, req, aspectID, domain)
		new.ServeHTTP(w, req)
	})
}
