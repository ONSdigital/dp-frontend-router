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



### Licence

Copyright ©‎ 2016 - 2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
