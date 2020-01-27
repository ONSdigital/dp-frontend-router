package healthcheck

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ONSdigital/log.go/log"
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
		if !check.hasRun() {
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
		checkStatus := hc.getCheckStatus(check)
		if checkStatus == StatusCritical {
			return StatusCritical
		} else if checkStatus == StatusWarning {
			status = StatusWarning
		}
	}
	return status
}

// getCheckStatus returns a string for the status on if an individual check
func (hc HealthCheck) getCheckStatus(c *Check) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	switch c.state.Status {
	case StatusOK:
		return StatusOK
	case StatusWarning:
		return StatusWarning
	default:
		now := time.Now().UTC()
		status := StatusWarning
		criticalTimeThreshold := hc.TimeOfFirstCriticalError.Add(hc.CriticalErrorTimeout)
		if c.state.LastSuccess.Before(hc.TimeOfFirstCriticalError) && now.After(criticalTimeThreshold) {
			status = StatusCritical
		}
		// Set timestamp of first critical error to now
		if c.state.LastSuccess.After(hc.TimeOfFirstCriticalError) {
			hc.TimeOfFirstCriticalError = now
		}
		return status
	}
}
