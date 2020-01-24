package healthcheck

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// Checker represents the interface all checker functions abide to
type Checker func(context.Context) (*CheckState, error)

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
	State   *CheckState
	Checker *Checker
	mutex   *sync.Mutex
}

// newCheck returns a pointer to a new instantiated Check with
// the provided checker function
func newCheck(checker *Checker) (*Check, error) {
	if checker == nil {
		return nil, errors.New("expected checker but none provided")
	}

	return &Check{
		State:   nil,
		Checker: checker,
		mutex:   &sync.Mutex{},
	}, nil
}

func (c *Check) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(c.State)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Check) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &c.State); err != nil {
		return err
	}
	return nil
}
