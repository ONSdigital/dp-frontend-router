package health

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
)

var hc *healthcheck.HealthCheck

// InitializeHealthCheck initializes the HealthCheck object with startTime now
func InitializeHealthCheck(BuildTime, GitCommit, Version string) {

	buildTime, err := strconv.Atoi(BuildTime)
	if err != nil {
		log.Event(nil, "failed to obtain build time", log.Error(err))
		buildTime = 0
	}
	log.Event(nil, "init Healthckeck", log.Data{"BuildTime": BuildTime, "GitCommit": GitCommit, "Version": Version})

	hc = &healthcheck.HealthCheck{
		Status: healthcheck.StatusOK,
		Version: healthcheck.CreateVersionInfo(
			time.Unix(int64(buildTime), 0),
			GitCommit,
			Version,
		),
		Uptime:    time.Duration(0),
		StartTime: time.Now().UTC(),
		Checks:    []*healthcheck.Check{},
	}
}

// Handler updates the HealthCheck current uptime, marshalls it, and writes it to the ResponseWriter.
func Handler(w http.ResponseWriter, req *http.Request) {

	if hc == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hc.Uptime = time.Since(hc.StartTime) / time.Millisecond

	marshaled, err := json.Marshal(hc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(marshaled)
}
