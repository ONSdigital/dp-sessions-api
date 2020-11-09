dp-sessions-api
================
API for Sessions

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default   | Description
| ---------------------------- | --------- | -----------
| BIND_ADDR                    | :         | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s        | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 10s       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 1m        | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| ZEBEDEE_URL                  | http://localhost:8082 | URL for Zebedee
| SERVICE_AUTH_TOKEN           |           | Service Auth token for communicating with Zebedee
| ELASTICACHE_ADDR             | localhost:6379 | Address of Elasticache/Redis
| ELASTICACHE_PASSWORD         | default   | Password for Elasticache/Redis
| ELASTICACHE_DATABASE         | 0         | Database for Elasticache/Redis (`int` format)
| ELASTICACHE_TTL              | 30m       | Time before Elasticache/Redis key expires (`time.Duration` format)
| ENABLE_REDIS_TLS_CONFIG      | false     | Turn TLS configuration on or off (`bool` format)

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2020, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

