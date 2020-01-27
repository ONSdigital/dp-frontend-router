package health

import (
	"context"
	"time"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
)

var hc healthcheck.HealthCheck

// InitializeHealthCheck initializes the HealthCheck object with startTime now
func InitializeHealthCheck(ctx context.Context, buildTime, gitCommit, version string) healthcheck.HealthCheck {

	versionInfo, err := healthcheck.CreateVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		log.Event(ctx, "failed to obtain versionInfo", log.Error(err))
	}

	criticalTimeout := time.Minute
	interval := 10 * time.Second

	// zc := client.New(config.ZebedeeURL)
	// zChecker := func(context.Context, *healthcheck.CheckState) error {
	// 	_, err := zc.Checker(ctx)
	// 	return err
	// }

	hc, err = healthcheck.Create(versionInfo, criticalTimeout, interval, nil)

	hc.Start(ctx)
	return hc
}
