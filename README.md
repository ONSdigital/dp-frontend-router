# dp-frontend-router

## Configuration

| Environment variable            | Default                                   | Description                                                                              |
|---------------------------------|-------------------------------------------|------------------------------------------------------------------------------------------|
| BIND_ADDR                       | :20000                                    | The host and port to bind to.                                                            |
| HTTP_MAX_CONNECTIONS            | 0                                         | Limit the number of concurrent http connections (0 = unlimited)                          |
| BABBAGE_URL                     | <https://localhost:8080>                  | The URL of the babbage instance to use                                                   |
| COOKIES_CONTROLLER_URL          | <http://localhost:24100>                  | The URL of dp-frontend-cookie-controller                                                 |
| HOMEPAGE_CONTROLLER_URL         | <http://localhost:24400>                  | The URL of dp-frontend-dataset-controller                                                |
| DATASET_CONTROLLER_URL          | <http://localhost:20200>                  | The URL of dp-frontend-dataset-controller                                                |
| FILTER_DATASET_CONTROLLER_URL   | <http://localhost:20001>                  | The URL of dp-frontend-filter-dataset-controller                                         |
| FEEDBACK_CONTROLLER_URL         | <http://localhost:25200>                  | The URL of dp-frontend-feedback-controller                                               |
| SEARCH_CONTROLLER_URL           | <http://localhost:25000>                  | The URL of dp-frontend-search-controller                                                 |
| DATA_AGGREGATION_PAGES_ENABLED  | false                                     | Enables the new data aggregation pages                                                   |
| TOPIC_AGGREGATION_PAGES_ENABLED | false                                     | enables the topic data aggregation pages                                                 |
| RELATED_DATA_ROUTE_ENABLED      | false                                     | Flag to enable related data page routing to the dp-frontend-search-controller            |
| PREVIOUS_RELEASES_ROUTE_ENABLED | false                                     | Flag to enable previous releases page routing to the dp-frontend-search-controller       |
| SEARCH_ROUTES_ENABLED           | false                                     | Search routes feature toggle                                                             |
| API_ROUTER_URL                  | <http://localhost:23200/v1>               | The API router URL                                                                       |
| DOWNLOADER_URL                  | <http://localhost:23400>                  | The URL of dp-file-downloader.                                                           |
| SITE_DOMAIN                     | ons.gov.uk                                | The domain hosting the site                                                              |
| CONTENT_TYPE_BYTE_LIMIT         | 5000000 (5MB)                             | Response size at which we stop checking content-type to avoid oom errors                 |
| HEALTHCHECK_INTERVAL            | 30s                                       | The period of time between health checks                                                 |
| HEALTHCHECK_CRITICAL_TIMEOUT    | 90s                                       | The period of time after which failing checks will result in critical global check       |
| ZEBEDEE_REQUEST_TIMEOUT_SECONDS | 5s                                        | The period of time to wait before timing out when communicating with Zebedee             |
| ZEBEDEE_REQUEST_MAXIMUM_RETRIES | 0                                         | The number of retry attempts to make to Zebedee                                          |
| PROXY_TIMEOUT                   | 5s                                        | The write timeout for proxied requests                                                   |
| NEW_DATASET_ROUTING_ENABLED     | false                                     | Flag to enable dataset page routing to dp-frontend-dataset-controller instead of babbage |
| DATASET_FINDER_ENABLED          | false                                     | Flag to enabled routing to dataset finder page in search                                 |
| OTEL_EXPORTER_OTLP_ENDPOINT     | localhost:4317                            | Host and port for the OpenTelemetry endpoint                                             |
| OTEL_SERVICE_NAME               | dp-frontend-router                        | Service name to report to telemetry tools                                                |
| OTEL_BATCH_TIMEOUT              | 5s                                        | Interval between pushes to OT Collector                                                  |
| OTEL_ENABLED                    | false                                     | Feature flag to enable OpenTelemetry                                                     |
| LEGACY_CACHE_PROXY_ENABLED      | false                                     | Flag to enable requests to Babbage to go through the dp-legacy-cache-proxy instead.      |
| LEGACY_CACHE_PROXY_URL          | <http://localhost:29200>                  | The URL of dp-legacy-cache-proxy                                                         |
| RELEASE_CALENDAR_CONTROLLER_URL | <http://localhost:27700>                  | The URL of dp-frontend-release-calendar                                                  |

### Licence

Copyright © 2025, Office for National Statistics (<https://www.ons.gov.uk>)

Released under MIT license, see [LICENSE](LICENSE.md) for details
