package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

// Config represents service configuration for dp-sessions-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	ZebedeeURL                 string        `envconfig:"ZEBEDEE_URL"`
	ServiceAuthToken           string        `envconfig:"SERVICE_AUTH_TOKEN"				json:"-"`
	ElasticacheAddr            string        `envconfig:"ELASTICACHE_ADDR"`
	ElasticachePassword        string        `envconfig:"ELASTICACHE_PASSWORD"				json:"-"`
	ElasticacheDatabase        int           `envconfig:"ELASTICACHE_DATABASE"`
	ElasticacheTTL             time.Duration `envconfig:"ELASTICACHE_TTL"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg := &Config{
		BindAddr:                   ":24400",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		ZebedeeURL:                 "http://localhost:8082",
		ServiceAuthToken:           "",
		ElasticacheAddr:            "localhost:6379",
		ElasticachePassword:        "default",
		ElasticacheDatabase:        0,
		ElasticacheTTL:             30 * time.Minute,
	}

	return cfg, envconfig.Process("", cfg)
}
