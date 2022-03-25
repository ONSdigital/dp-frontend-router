package abtest

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/log.go/log"
)

//go:generate moq -out searchtest/handler.go -pkg test . Handler
type Handler http.Handler

func SearchHandler(newSearch, oldSearch http.Handler, percentage int, domain string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logData := log.Data{"domain": domain, "percent_to_new_search": percentage} // Remove this line
		now := time.Now().UTC()
		log.Event(req.Context(), "did we get here 11111111", logData)

		// Check if user has choosen to exit new search
		exitNewSearch, ok := req.URL.Query()["exit-new-search"]
		if ok && len(exitNewSearch[0]) > 0 {
			log.Event(req.Context(), "we are exiting new search, nooooooooooooooooooooo!", logData) // Remove this line
			HandleSearchExit(w, req, oldSearch, now, domain)
			return
		}
		log.Event(req.Context(), "did we get here 222222222", logData) // Remove this line

		// Retrieve AB Test Cookie
		cookie, err := cookies.GetABTest(req)
		if err != nil && err != cookies.ErrABTestCookieNotFound {
			log.Event(req.Context(), "error getting a/b test cookie", log.WARN, log.Error(err), logData) // Remove logData
		}
		logData["cookie"] = cookie                                        // Remove this line
		logData["cookie_error"] = cookies.ErrABTestCookieNotFound.Error() // Remove this line
		log.Event(req.Context(), "did we get here 3333333333", logData)   // Remove this line

		// If AB Test cookie not set, create new AB Test cookie
		if cookie.NewSearch == nil && cookie.OldSearch == nil {
			log.Event(req.Context(), "ab tea=st cookie not set, about to call HandleCookieCreationAndServ", log.WARN, log.Error(err), logData) // Remove this line
			HandleCookieCreationAndServ(w, req, newSearch, oldSearch, percentage, domain, now)
			return
		}
		logData["cookie"] = cookie                                     // Remove this line
		log.Event(req.Context(), "did we get here 444444444", logData) // Remove this line

		// If AB Test cookie expired, set new AB Test cookie
		if cookie.NewSearch.Before(now) && cookie.OldSearch.Before(now) {
			log.Event(req.Context(), "ab test cookie expired, about to call HandleCookieCreationAndServ", log.WARN, log.Error(err), logData) // Remove this line
			HandleCookieCreationAndServ(w, req, newSearch, oldSearch, percentage, domain, now)
			return
		}
		logData["cookie"] = cookie                                      // Remove this line
		log.Event(req.Context(), "did we get here 5555555555", logData) // Remove this line

		servABTest(newSearch, oldSearch, w, req, cookie, now)
	})
}

func setTime24HoursAhead(now time.Time) time.Time {
	log.Event(context.Background(), "setting time 24 hours ahead", log.INFO, log.Data{"time_now": now}) // Remove this line
	return now.Add(24 * time.Duration(time.Hour))
}

// HandleCookieCreationAndServ calls randomise, creates a cookie and servs the request
func HandleCookieCreationAndServ(w http.ResponseWriter, req *http.Request, newSearch, oldSearch http.Handler, percentage int, domain string, now time.Time) {
	servs := RandomiseABTestCookie(percentage, now)
	err := cookies.SetABTest(w, servs, domain)
	if err != nil {
		log.Event(req.Context(), "error setting a/b test cookie. directing user to old search", log.ERROR, log.Error(err))
		oldSearch.ServeHTTP(w, req)
		return
	}

	log.Event(context.Background(), "AB test Cookie set", log.INFO, log.Data{"cookies": servs}) // Remove this line

	servABTest(newSearch, oldSearch, w, req, servs, now)

	log.Event(context.Background(), "New or old search served", log.INFO, log.Data{"cookies": servs}) // Remove this line
	return
}

// RandomiseABTestCookie randomly sets expiry times for new and old search services
func RandomiseABTestCookie(percentage int, now time.Time) cookies.ABServices {
	var newSearch time.Time
	var oldSearch time.Time
	rand.Seed(time.Now().UnixNano())

	log.Event(context.Background(), "did we get here", log.INFO, log.Data{"percentage_users_to_new_search": percentage}) // Remove this line
	if rand.Intn(100) < percentage {
		newSearch = setTime24HoursAhead(now)
		oldSearch = now
		log.Event(context.Background(), "Setting cookie for new search", log.INFO) // Remove this line
	} else {
		newSearch = now
		oldSearch = setTime24HoursAhead(now)
		log.Event(context.Background(), "Setting cookie for old search", log.INFO) // Remove this line
	}

	return cookies.ABServices{
		NewSearch: &newSearch,
		OldSearch: &oldSearch,
	}
}

func servABTest(newSearch, oldSearch http.Handler, w http.ResponseWriter, req *http.Request, cookie cookies.ABServices, now time.Time) {
	if cookie.NewSearch.After(now) {
		log.Event(context.Background(), "should be directing new search", log.INFO, log.Data{"cookies": cookie}) // Remove this line
		newSearch.ServeHTTP(w, req)
	}

	if cookie.OldSearch.After(now) {
		log.Event(context.Background(), "should be directing to old search", log.INFO, log.Data{"cookies": cookie}) // Remove this line
		oldSearch.ServeHTTP(w, req)
	}

	log.Event(context.Background(), "serving ab test cookie", log.INFO, log.Data{"cookies": cookie}) // Remove this line
}

func HandleSearchExit(w http.ResponseWriter, req *http.Request, oldSearch http.Handler, now time.Time, domain string) {
	tomorrow := setTime24HoursAhead(now)
	err := cookies.UpdateSearch(req, w, now, tomorrow, domain)
	if err != nil {
		log.Event(req.Context(), "error update new search value of a/b test cookie. directing user to old search", log.ERROR, log.Error(err))
	}
	oldSearch.ServeHTTP(w, req)
}
