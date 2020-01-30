package health

import (
	"context"
	"time"

	client "github.com/ONSdigital/dp-api-clients-go/zebedee"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
)

const criticalTimeout = time.Minute
const interval = 10 * time.Second

var hc healthcheck.HealthCheck

// StartHealthCheck instantiates and starts HealthCheck
func StartHealthCheck(ctx context.Context, buildTime, gitCommit, version string, zc *client.Client) healthcheck.HealthCheck {

	// Create VersionInfo
	versionInfo, err := healthcheck.CreateVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		log.Event(ctx, "failed to obtain versionInfo", log.Error(err))
	}

	// Instantiate and start Healthcheck
	hc, err = healthcheck.Create(versionInfo, criticalTimeout, interval, zc.Checker)
	hc.Start(ctx)
	return hc
}
