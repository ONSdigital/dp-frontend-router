# dp-frontend-router

## Configuration

| Environment variable            | Default                                   | Description                                                                        |
|---------------------------------|-------------------------------------------|------------------------------------------------------------------------------------|
| BIND_ADDR                       | :20000                                    | The host and port to bind to.                                                      |
| BABBAGE_URL                     | <https://localhost:8080>                  | The URL of the babbage instance to use                                             |
| RENDERER_URL                    | <http://localhost:20010>                  | The URL of dp-frontend-renderer                                                    |
| COOKIES_CONTROLLER_URL          | <http://localhost:24100>                  | The URL of dp-frontend-cookie-controller                                           |
| HOMEPAGE_CONTROLLER_URL         | <http://localhost:24400>                  | The URL of dp-frontend-dataset-controller                                          |
| DATASET_CONTROLLER_URL          | <http://localhost:20200>                  | The URL of dp-frontend-dataset-controller                                          |
| FILTER_DATASET_CONTROLLER_URL   | <http://localhost:20001>                  | The URL of dp-frontend-filter-dataset-controller                                   |
| GEOGRAPHY_CONTROLLER_URL        | <http://localhost:23700>                  | The URL of dp-frontend-geography-controller                                        |
| GEOGRAPHY_ENABLED               | false                                     | Geography feature toggle                                                           |
| FEEDBACK_CONTROLLER_URL         | <http://localhost:25200>                  | The URL of dp-frontend-feedback-controller                                         |
| SEARCH_CONTROLLER_URL           | <http://localhost:25000>                  | The URL of dp-frontend-search-controller                                           |
| SEARCH_ROUTES_ENABLED           | false                                     | Search routes feature toggle                                                       |
| API_ROUTER_URL                  | <http://localhost:23200/v1>               | The API router URL                                                                 |
| DOWNLOADER_URL                  | <http://localhost:23400>                  | The URL of dp-file-downloader.                                                     |
| AREA_PROFILE_CONTROLLER_URL     | <http://localhost:26600>                  | The URL of dp-frontend-area-profiles.                                              |
| AREA_PROFILE_ROUTES_ENABLED     | false                                     | Area profiles routes enabled                                                       |
| INTERACTIVES_CONTROLLER_URL     | <http://localhost:27300>                  | The URL of dp-frontend-interactives-controller                                     |
| INTERACTIVES_ROUTES_ENABLED     | false                                     | Interactives routes enabled                                                        |
| PATTERN_LIBRARY_ASSETS_PATH     | <https://cdn.ons.gov.uk/sixteens/e42235b> | The URL to the sixteens build to use                                               |
| SITE_DOMAIN                     | ons.gov.uk                                | The domain hosting the site                                                        |
| REDIRECT_SECRET                 | secret                                    | Pre-shared key for signing/encrypting redirect data                                |
| ANALYTICS_SQS_URL               |                                           | SQS URL for search analytics; leave blank to disable                               |
| CONTENT_TYPE_BYTE_LIMIT         | 5000000 (5MB)                             | Response size at which we stop checking content-type to avoid oom errors           |
| HEALTHCHECK_INTERVAL            | 30s                                       | The period of time between health checks                                           |
| HEALTHCHECK_CRITICAL_TIMEOUT    | 90s                                       | The period of time after which failing checks will result in critical global check |
| ZEBEDEE_REQUEST_TIMEOUT_SECONDS | 5s                                        | The period of time to wait before timing out when communicating with Zebedee       |
| ZEBEDEE_REQUEST_MAXIMUM_RETRIES | 0                                         | The number of retry attempts to make to Zebedee                                    |
| ENABLE_SEARCH_AB_TEST           | false                                     | Enable AB search                                                                   |
| SEARCH_AB_TEST_PERCENTAGE       | 10                                        | AB search percentage                                                               |
| PROXY_TIMEOUT                   | 5s                                        | The write timeout for proxied requests |
| DATASET_ENABLED                 | false                                     | Flag to enable dataset page routing to dp-frontend-dataset-controller instead of babbage |

### Licence

Copyright © 2021, Office for National Statistics (<https://www.ons.gov.uk>)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
