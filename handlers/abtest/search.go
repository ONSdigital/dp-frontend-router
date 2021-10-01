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
	rand.Seed(time.Now().UnixNano())

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Retrieve AB Test Cookie
		cookie, err := cookies.GetABTest(req)
		if err != nil && err != cookies.ErrABTestCookieNotFound {
			log.Event(req.Context(), "error getting a/b test cookie", log.WARN, log.Error(err))
		}

		// If AB Test cookie not set, create new AB Test cookie
		if cookie.NewSearch == nil && cookie.OldSearch == nil {
			servs := RandomiseABTestCookie(percentage, time.Now())
			cookieErr := cookies.SetABTest(w, servs, domain)
			if cookieErr != nil {
				log.Event(req.Context(), "error setting a/b test cookie. direct use to old search", log.ERROR, log.Error(err))
				oldSearch.ServeHTTP(w, req)
				return
			}
			servABTest(newSearch, oldSearch, w, req, servs)
			return
		}

		// If AB Test cookie expired, set new AB Test cookie
		if cookie.NewSearch.Before(time.Now().UTC()) && cookie.OldSearch.Before(time.Now().UTC()) {
			servs := RandomiseABTestCookie(percentage, time.Now())
			cookieErr := cookies.SetABTest(w, servs, domain)
			if cookieErr != nil {
				log.Event(req.Context(), "error setting a/b test cookie. direct use to old search", log.ERROR, log.Error(err))
				oldSearch.ServeHTTP(w, req)
				return
			}
			servABTest(newSearch, oldSearch, w, req, servs)
			return
		}

		servABTest(newSearch, oldSearch, w, req, cookie)
	})
}

func setTime24HoursAhead(now time.Time) time.Time {
	return now.UTC().Add(24 * time.Duration(time.Hour))
}

func RandomiseABTestCookie(percentage int, now time.Time) cookies.ABServices {
	var newSearch time.Time
	var oldSearch time.Time
	if rand.Intn(100) < percentage {
		newSearch = setTime24HoursAhead(now)
		oldSearch = now.UTC()
	} else {
		newSearch = now.UTC()
		oldSearch = setTime24HoursAhead(now)
	}

	return cookies.ABServices{
		NewSearch: &newSearch,
		OldSearch: &oldSearch,
	}
}

func servABTest(newSearch, oldSearch http.Handler, w http.ResponseWriter, req *http.Request, cookie cookies.ABServices) {
	if cookie.NewSearch.After(time.Now().UTC()) {
		newSearch.ServeHTTP(w, req)
	}
	if cookie.OldSearch.After(time.Now().UTC()) {
		oldSearch.ServeHTTP(w, req)
	}
}
