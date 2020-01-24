package healthcheck

import (
	"context"
	"errors"
	"runtime"
	"time"
)

const language = "go"

// HealthCheck represents the app's health check, including its component checks
type HealthCheck struct {
	Status                   string        `json:"status"`
	Version                  VersionInfo   `json:"version"`
	Uptime                   time.Duration `json:"uptime"`
	StartTime                time.Time     `json:"start_time"`
	Checks                   []*Check      `json:"checks"`
	Started                  bool          `json:"-"`
	Interval                 time.Duration `json:"-"`
	CriticalErrorTimeout     time.Duration `json:"-"`
	TimeOfFirstCriticalError time.Time     `json:"-"`
	Tickers                  []*ticker     `json:"-"`
}

// VersionInfo represents the version information of an app
type VersionInfo struct {
	BuildTime       time.Time `json:"build_time"`
	GitCommit       string    `json:"git_commit"`
	Language        string    `json:"language"`
	LanguageVersion string    `json:"language_version"`
	Version         string    `json:"version"`
}

// Create returns a new instantiated HealthCheck object. Caller to provide:
// version information of the app,
// criticalTimeout for how long to wait until an unhealthy dependent propagates its state to make this app unhealthy
// interval in which to check health of dependencies
// checkers which implement the checker interface and can run a checkup to determine the health of the app and/or its dependencies
func Create(version VersionInfo, criticalTimeout, interval time.Duration, checkers ...*Checker) HealthCheck {

	var checks []*Check

	hc := HealthCheck{
		Started:              false,
		Checks:               checks,
		Version:              version,
		CriticalErrorTimeout: criticalTimeout,
		Interval:             interval,
	}

	for _, checker := range checkers {
		hc.AddCheck(checker)
	}

	return hc
}

// CreateVersionInfo returns a health check version info object. Caller to provide:
// buildTime for when the app was built
// gitCommit the SHA-1 commit hash of the built app
// version the semantic version of the built app
func CreateVersionInfo(buildTime time.Time, gitCommit, version string) VersionInfo {
	return VersionInfo{
		BuildTime:       buildTime,
		GitCommit:       gitCommit,
		Language:        language,
		LanguageVersion: runtime.Version(),
		Version:         version,
	}
}

// AddCheck adds a provided checker to the health check
func (hc *HealthCheck) AddCheck(checker *Checker) (err error) {
	if hc.Started {
		err := errors.New("unable to add new check, health check has already started")
		return err
	}

	check, err := newCheck(checker)
	if err != nil {
		return err
	}

	hc.Checks = append(hc.Checks, check)

	return nil
}

// newTickers returns an array of tickers based on the number of checks in the checks parameter.
// Each check is executed at the given interval also passed into the function
func newTickers(interval time.Duration, checks []*Check) []*ticker {
	var tickers []*ticker
	for _, check := range checks {
		tickers = append(tickers, createTicker(interval, check))
	}
	return tickers
}

// Start begins each ticker, this is used to run the health checks on dependent apps
// takes argument context and should utilise contextWithCancel
func (hc *HealthCheck) Start(ctx context.Context) {
	hc.Started = true
	hc.Tickers = newTickers(hc.Interval, hc.Checks)
	hc.StartTime = time.Now().UTC()
	for _, ticker := range hc.Tickers {
		ticker.start(ctx)
	}
}

// Stop will cancel all tickers and thus stop all health checks
func (hc *HealthCheck) Stop() {
	for _, ticker := range hc.Tickers {
		ticker.stop()
	}
}
