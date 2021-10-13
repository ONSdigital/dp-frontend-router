package abtest

import (
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
		now := time.Now().UTC()

		// Check if user has choosen to exit new search
		exitNewSearch, ok := req.URL.Query()["exit-new-search"]
		if ok && len(exitNewSearch[0]) > 0 {
			HandleSearchExit(w, req, oldSearch, now, domain)
			return
		}

		// Retrieve AB Test Cookie
		cookie, err := cookies.GetABTest(req)
		if err != nil && err != cookies.ErrABTestCookieNotFound {
			log.Event(req.Context(), "error getting a/b test cookie", log.WARN, log.Error(err))
		}

		// If AB Test cookie not set, create new AB Test cookie
		if cookie.NewSearch == nil && cookie.OldSearch == nil {
			HandleCookieCreationAndServ(w, req, newSearch, oldSearch, percentage, domain, now)
			return
		}

		// If AB Test cookie expired, set new AB Test cookie
		if cookie.NewSearch.Before(now) && cookie.OldSearch.Before(now) {
			HandleCookieCreationAndServ(w, req, newSearch, oldSearch, percentage, domain, now)
			return
		}

		servABTest(newSearch, oldSearch, w, req, cookie, now)
	})
}

func setTime24HoursAhead(now time.Time) time.Time {
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
	servABTest(newSearch, oldSearch, w, req, servs, now)
	return
}

// RandomiseABTestCookie randomly sets expiry times for new and old search services
func RandomiseABTestCookie(percentage int, now time.Time) cookies.ABServices {
	var newSearch time.Time
	var oldSearch time.Time
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) < percentage {
		newSearch = setTime24HoursAhead(now)
		oldSearch = now
	} else {
		newSearch = now
		oldSearch = setTime24HoursAhead(now)
	}

	return cookies.ABServices{
		NewSearch: &newSearch,
		OldSearch: &oldSearch,
	}
}

func servABTest(newSearch, oldSearch http.Handler, w http.ResponseWriter, req *http.Request, cookie cookies.ABServices, now time.Time) {
	if cookie.NewSearch.After(now) {
		newSearch.ServeHTTP(w, req)
	}
	if cookie.OldSearch.After(now) {
		oldSearch.ServeHTTP(w, req)
	}
}

func HandleSearchExit(w http.ResponseWriter, req *http.Request, oldSearch http.Handler, now time.Time, domain string) {
	tomorrow := setTime24HoursAhead(now)
	err := cookies.UpdateNewSearch(req, w, now, domain)
	if err != nil {
		log.Event(req.Context(), "error update new search value of a/b test cookie. directing user to old search", log.ERROR, log.Error(err))
	}
	err = cookies.UpdateOldSearch(req, w, tomorrow, domain)
	if err != nil {
		log.Event(req.Context(), "error update old search value of a/b test cookie. directing user to old search", log.ERROR, log.Error(err))
	}
	oldSearch.ServeHTTP(w, req)
}
