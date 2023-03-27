package abtest

import (
	"net/http"
	"time"

	"github.com/ONSdigital/dp-cookies/cookies"
)

const (
	RelcalTestCookieAspect = "dsa-release-calendar"
	RelcalNewExit          = "exit-new-relcal"
)

// Handler returns the relevant handler on the basis of the supplied parameters.
// It delegates to both abTestHandler and abTestPurgeHandler on the basis the abTest parameter, but it is really
// an encapsulation of the decision-making process as to what handler is used.
func Handler(abTest bool, newHandler, old http.Handler, percentage int, aspectID, domain, exitNew string) http.Handler {
	if abTest {
		return abTestHandler(newHandler, old, percentage, aspectID, domain, exitNew)
	}
	return abTestPurgeHandler(newHandler, aspectID, domain)
}

// abTestHandler routes requests to either the old or new handler, for a given aspectID, according to the given percentage
// i.e. for the given percentage of calls X, X% will be routed to the new handler, and the remainder to the old handler.
// Most of the functionality is provided by the dp-cookies library, which uses a single ab_test cookie to embed all aspects
// If the aspect does not exist or has expired, it is created/renewed according to a particular randomiser - in general
// the DefaultABTestRandomiser in the library is sufficient
// A well known string - the exitNew string -  can be used as a query parameter to the call, in order to definitively chose
// the old handler
func abTestHandler(newHandler, old http.Handler, percentage int, aspectID, domain, exitNew string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		now := time.Now().UTC()

		if _, ok := req.URL.Query()[exitNew]; ok {
			cookies.HandleABTestExit(w, req, old, aspectID, domain)
			return
		}

		aspect := cookies.GetABTestCookieAspect(req, aspectID)

		if (aspect.New.IsZero() && aspect.Old.IsZero()) || (aspect.New.Before(now) && aspect.Old.Before(now)) {
			cookies.HandleCookieAndServ(w, req, newHandler, old, aspectID, domain, cookies.DefaultABTestRandomiser(percentage))
			return
		}

		cookies.ServABTest(w, req, newHandler, old, aspect)
	})
}

// abTestPurgeHandler is used to remove a given AspectID from the single ab_test cookie handled by the dp-cookies library
// It is useful when AB Testing for a particular aspect has finished, but the aspect is still embedded in client's ab_test
// cookie - this handler will remove the aspect and can be left in use for several weeks after testing has finished to 'clean'
// the underlying ab_test cookie
func abTestPurgeHandler(newHandler http.Handler, aspectID, domain string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookies.RemoveABTestCookieAspect(w, req, aspectID, domain)
		newHandler.ServeHTTP(w, req)
	})
}
