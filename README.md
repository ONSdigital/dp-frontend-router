dp-frontend-router
==================

### Configuration

| Environment variable          | Default                                 | Description
| ----------------------------- | --------------------------------------- | --------------------------------------
| BIND_ADDR                     | :20000                                  | The host and port to bind to.
| BABBAGE_URL                   | https://www.ons.gov.uk                  | The URL of the babbage instance to use
| RENDERER_URL                  | http://localhost:20010                  | The URL of dp-frontend-renderer
| DATASET_CONTROLLER_URL        | http://localhost:20200                  | The URL of dp-frontend-dataset-controller
| FILTER_DATASET_CONTROLLER_URL | http://localhost:20001                  | The URL of dp-frontend-filter-dataset-controller
| GEOGRAPHY_CONTROLLER_URL      | http://localhost:23700                  | The URL of dp-frontend-geography-controller
| ZEBEDEE_URL                   | http://localhost:8082                   | The URL of zebedee
| PATTERN_LIBRARY_ASSETS_PATH   | https://cdn.ons.gov.uk/sixteens/e42235b | The URL to the sixteens build to use
| SITE_DOMAIN                   | ons.gov.uk                              | The domain hosting the site
| HOMEPAGE_AB_PERCENT           | 0                                       | Percentage of users who get version B
| DOWNLOADER_URL                | http://localhost:23400                  | The URL of dp-file-downloader.
| ANALYTICS_SQS_URL             |                                         | SQS URL for search analytics; leave blank to disable
| AWS_ACCESS_KEY_ID             |                                         | Your AWS access key ID (required for SQS)
| AWS_SECRET_ACCESS_KEY         |                                         | Your AWS secret access key
| AWS_REGION                    |                                         | AWS region (normally eu-west-1)
| REDIRECT_SECRET               | secret                                  | Pre-shared key for signing/encrypting redirect data
| CONTENT_TYPE_BYTE_LIMIT       | 5000000 (5MB)                           | Response size at which we stop checking content-type to avoid oom errors
| HEALTHCHECK_INTERVAL          | 30s                                     | The period of time between health checks
| HEALTHCHECK_CRITICAL_TIMEOUT  | 90s                                     | The period of time after which failing checks will result in critical global check 
| ENABLE_PROFILER               | false                                   | Flag to enable go profiler
| PPROF_TOKEN                   | ""                                      | The profiling token to access service profiling

### Profiling

An optional `/debug` endpoint has been added, in order to profile this service via `pprof` go library.
In order to use this endpoint, you will need to enable profiler flag and set a PPROF_TOKEN:

```
export ENABLE_PROFILER=true
export PPROF_TOKEN={generated uuid}
```

Then you can us the profiler as follows:

1- Start service, load test or if on environment wait for a number of requests to be made.

2- Send authenticated request and store response in a file (this can be best done in command line like so: `curl <host>:<port>/debug/pprof/heap -H "Authorization: Bearer {generated uuid}" > heap.out` - see pprof documentation on other endpoints

3- View profile either using a web ui to navigate data (a) or using pprof on command line to navigate data (b) 
  a) `go tool pprof -http=:8080 heap.out`
  b) `go tool pprof heap.out`, -o flag to see various options

### Licence

Copyright ©‎ 2016 - 2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
