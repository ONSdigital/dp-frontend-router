package healthcheck

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// A list of possible check statuses
const (
	StatusOK       = "OK"
	StatusWarning  = "WARNING"
	StatusCritical = "CRITICAL"
)

// Checker represents the interface all checker functions abide to
type Checker func(context.Context, *CheckState) error

// CheckState represents the health status returned by a checker
type CheckState struct {
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	StatusCode  int        `json:"status_code,omitempty"`
	Message     string     `json:"message"`
	LastChecked *time.Time `json:"last_checked"`
	LastSuccess *time.Time `json:"last_success"`
	LastFailure *time.Time `json:"last_failure"`
}

// Check represents a check performed by the health check
type Check struct {
	state   *CheckState
	checker Checker
	mutex   *sync.Mutex
}

// hasRun returns true if the check has been run and has state
func (c *Check) hasRun() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.state.LastChecked == nil {
		return false
	}
	return true
}

// newCheck returns a pointer to a new instantiated Check with
// the provided checker function
func newCheck(checker Checker) (*Check, error) {
	if checker == nil {
		return nil, errors.New("expected checker but none provided")
	}

	return &Check{
		state:   &CheckState{},
		checker: checker,
		mutex:   &sync.Mutex{},
	}, nil
}

func (c *Check) MarshalJSON() ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	b, err := json.Marshal(c.state)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Check) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &c.state); err != nil {
		return err
	}
	return nil
}
