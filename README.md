dp-frontend-router
==================

### Configuration

| Environment variable  | Default                 | Description
| --------------------- | ----------------------- | --------------------------------------
| BIND_ADDR             | :20000                  | The host and port to bind to.
| BABBAGE_URL           | https://www.ons.gov.uk  | The URL of the babbage instance to use.
| RENDERER_URL          | http://localhost:20010  | The URL of dp-frontend-renderer.
| RESOLVER_URL          | http://localhost:20020  | The URL of dp-content-resolver.
| DOWNLOADER_URL        | http://localhost:23400  | The URL of dp-file-downloader.
| ANALYTICS_SQS_URL     |                         | SQS URL for search analytics; leave blank to disable
| AWS_ACCESS_KEY_ID     |                         | Your AWS access key ID (required for SQS)
| AWS_SECRET_ACCESS_KEY |                         | Your AWS secret access key
| AWS_REGION            |                         | AWS region (normally eu-west-1)
| REDIRECT_SECRET       | secret                  | Pre-shared key for signing/encrypting redirect data

### Licence

Copyright ©‎ 2016, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
