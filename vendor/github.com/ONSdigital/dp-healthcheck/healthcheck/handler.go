package healthcheck

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ONSdigital/log.go/log"
)

// A list of possible health statuses
const (
	StatusOK       = "OK"
	StatusCritical = "CRITICAL"
	StatusWarning  = "WARNING"
)

// Handler responds to an http request for the current health status
func (hc HealthCheck) Handler(w http.ResponseWriter, req *http.Request) {
	now := time.Now().UTC()
	ctx := req.Context()

	hc.Status = hc.getStatus(ctx)
	hc.Uptime = now.Sub(hc.StartTime)

	b, err := json.Marshal(hc)
	if err != nil {
		log.Event(ctx, "failed to marshal json", log.Error(err), log.Data{"health_check_response": hc})
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Event(ctx, "failed to write bytes for http response", log.Error(err))
		return
	}
}

// isAppStartingUp returns false when all clients have completed at least one check
func (hc HealthCheck) isAppStartingUp() bool {
	for _, check := range hc.Checks {
		if check.State == nil {
			return true
		}
	}
	return false
}

// getStatus returns a status as string as to the overall current apps health based on its dependent apps health
func (hc HealthCheck) getStatus(ctx context.Context) string {
	if hc.isAppStartingUp() {
		log.Event(ctx, "a dependency is still starting up")
		return StatusWarning
	}
	return hc.isAppHealthy()
}

// isAppHealthy checks every check for their health then produces and returns a status for this apps health
func (hc HealthCheck) isAppHealthy() string {
	status := StatusOK
	for _, check := range hc.Checks {
		appHealth := hc.readAppHealth(check)
		if appHealth == StatusCritical {
			return StatusCritical
		} else if appHealth == StatusWarning {
			status = StatusWarning
		}
	}
	return status
}

// readAppHealth locks mutex then reads a check finally it unlocks the mutex.
func (hc HealthCheck) readAppHealth(check *Check) string {
	check.mutex.Lock()
	defer check.mutex.Unlock()
	return hc.isCheckHealthy(check.State)
}

// isCheckHealthy returns a string for the status on if an individual dependent apps health
func (hc HealthCheck) isCheckHealthy(c *CheckState) string {
	now := time.Now().UTC()
	switch c.Status {
	case StatusWarning:
		return StatusWarning
	case StatusOK:
		return StatusOK
	default:
		status := StatusWarning
		criticalTimeThreshold := hc.TimeOfFirstCriticalError.Add(hc.CriticalErrorTimeout)
		if c.LastSuccess.Before(hc.TimeOfFirstCriticalError) && now.After(criticalTimeThreshold) {
			status = StatusCritical
		}
		// Set timestamp of first critical error to now
		if c.LastSuccess.After(hc.TimeOfFirstCriticalError) {
			hc.TimeOfFirstCriticalError = now
		}
		return status
	}
}
