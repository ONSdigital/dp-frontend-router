package health

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/clientlog"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	rchttp "github.com/ONSdigital/dp-rchttp"
	"github.com/ONSdigital/log.go/log"
)

var (
	// StatusMessage contains a map of messages to service response statuses
	StatusMessage = map[string]string{
		health.StatusOK:       "Everything is ok",
		health.StatusWarning:  "Things are degraded, but at least partially functioning",
		health.StatusCritical: "The checked functionality is unavailable or non-functioning",
	}
)

// ErrInvalidAppResponse is returned when an app does not respond
// with a valid status
type ErrInvalidAppResponse struct {
	ExpectedCode int
	ActualCode   int
	URI          string
}

// Client represents an app client
type Client struct {
	Client rchttp.Clienter
	URL    string
	name   string
}

// NewClient creates a new instance of Client with a given app url
func NewClient(url string) *Client {
	c := &Client{
		Client: rchttp.NewClient(),
		URL:    url,
	}

	// healthcheck client should not retry when calling a healthcheck endpoint,
	// append to current paths as to not change the client setup by service
	paths := c.Client.GetPathsWithNoRetries()
	paths = append(paths, "/health", "/healthcheck")
	c.Client.SetPathsWithNoRetries(paths)

	return c
}

// CreateCheck creates a new check state object
func CreateCheck(service string) (check health.CheckState) {
	check.Name = service
	return check
}

// Error should be called by the user to print out the stringified version of the error
func (e ErrInvalidAppResponse) Error() string {
	return fmt.Sprintf("invalid response from downstream service - should be: %d, got: %d, path: %s",
		e.ExpectedCode,
		e.ActualCode,
		e.URI,
	)
}

// Checker calls an app health endpoint and returns a check object to the caller
func (c *Client) Checker(ctx context.Context, state *health.CheckState) error {
	if state.Name == "" {
		return errors.New("missing service name in state")
	}
	c.name = state.Name
	logData := log.Data{
		"service": state.Name,
	}

	code, err := c.get(ctx, "/health")
	// Apps may still have /healthcheck endpoint
	// instead of a /health one
	if code == http.StatusNotFound {
		code, err = c.get(ctx, "/healthcheck")
	}
	if err != nil {
		log.Event(ctx, "failed to request service health", log.Error(err), logData)
	}

	currentTime := time.Now().UTC()
	state.StatusCode = code
	state.LastChecked = &currentTime

	switch code {
	case 0: // When there is a problem with the client return error in message
		state.Message = err.Error()
		state.Status = health.StatusCritical
		state.LastFailure = &currentTime
	case 200:
		state.Message = StatusMessage[health.StatusOK]
		state.Status = health.StatusOK
		state.LastSuccess = &currentTime
	case 429:
		state.Message = StatusMessage[health.StatusWarning]
		state.Status = health.StatusWarning
		state.LastFailure = &currentTime
	default:
		state.Message = StatusMessage[health.StatusCritical]
		state.Status = health.StatusCritical
		state.LastFailure = &currentTime
	}

	return nil
}

func (c *Client) get(ctx context.Context, path string) (int, error) {
	clientlog.Do(ctx, "retrieving dataset", c.name, c.URL)

	req, err := http.NewRequest("GET", c.URL+path, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.Client.Do(ctx, req)
	if err != nil {
		return 0, err
	}
	defer closeResponseBody(ctx, resp)

	if resp.StatusCode < 200 || (resp.StatusCode > 399 && resp.StatusCode != 429) {
		return resp.StatusCode, ErrInvalidAppResponse{http.StatusOK, resp.StatusCode, req.URL.Path}
	}

	return resp.StatusCode, nil
}

func closeResponseBody(ctx context.Context, resp *http.Response) {
	if resp.Body == nil {
		return
	}

	if err := resp.Body.Close(); err != nil {
		log.Event(ctx, "error closing http response body", log.Error(err))
	}
}
