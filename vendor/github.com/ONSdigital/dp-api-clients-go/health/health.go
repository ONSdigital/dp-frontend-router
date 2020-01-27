package health

import (
	"context"
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
	Check  *health.Check
	Name   string
	URL    string
}

// NewClient creates a new instance of Client with a given app url
func NewClient(name, url string) *Client {
	c := &Client{
		Client: rchttp.NewClient(),
		Name:   name,
		URL:    url,
		Check: &health.Check{
			Name: name,
		},
	}

	// healthcheck client should not retry when calling a healthcheck endpoint,
	// append to current paths as to not change the client setup by service
	paths := c.Client.GetPathsWithNoRetries()
	paths = append(paths, "/health", "/healthcheck")
	c.Client.SetPathsWithNoRetries(paths)

	return c
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
func (c *Client) Checker(ctx context.Context) (*health.Check, error) {
	logData := log.Data{
		"service": c.Name,
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
	c.Check.StatusCode = code
	c.Check.LastChecked = &currentTime

	switch code {
	case 0: // When there is a problem with the client return error in message
		c.Check.Message = err.Error()
		c.Check.Status = health.StatusCritical
		c.Check.LastFailure = &currentTime
	case 200:
		c.Check.Message = StatusMessage[health.StatusOK]
		c.Check.Status = health.StatusOK
		c.Check.LastSuccess = &currentTime
	case 429:
		c.Check.Message = StatusMessage[health.StatusWarning]
		c.Check.Status = health.StatusWarning
		c.Check.LastFailure = &currentTime
	default:
		c.Check.Message = StatusMessage[health.StatusCritical]
		c.Check.Status = health.StatusCritical
		c.Check.LastFailure = &currentTime
	}

	return c.Check, nil
}

func (c *Client) get(ctx context.Context, path string) (int, error) {
	clientlog.Do(ctx, "retrieving dataset", c.Name, c.URL)

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
