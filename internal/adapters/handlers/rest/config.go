package rest

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the configuration for the REST server.
type Config struct {
	// Port is the port on which the server listens.
	Port int
	// ReadTimeout is the maximum duration for reading the entire request.
	ReadTimeout time.Duration
	// WriteTimeout is the maximum duration before timing out writes of the response.
	WriteTimeout time.Duration
	// IdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled.
	IdleTimeout time.Duration
}

const (
	// Default values for the server configuration.
	// These can be overridden by environment variables.

	// defaultPort is the default port for the server.
	defaultPort = 8080
	// defaultReadTimeout is the default read timeout for the server.
	defaultReadTimeout = 5 * time.Second
	// defaultWriteTimeout is the default write timeout for the server.
	defaultWriteTimeout = 10 * time.Second
	// defaultIdleTimeout is the default idle timeout for the server.
	defaultIdleTimeout = 120 * time.Second
)

// LoadConfig loads the server configuration from environment variables.
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	port, exists := os.LookupEnv("SERVER_PORT")
	if exists {
		var err error
		cfg.Port, err = strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
		}
	} else {
		cfg.Port = defaultPort
	}

	readTimeout, exists := os.LookupEnv("SERVER_READ_TIMEOUT")
	if exists {
		var err error
		cfg.ReadTimeout, err = time.ParseDuration(readTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid SERVER_READ_TIMEOUT: %w", err)
		}
	} else {
		cfg.ReadTimeout = defaultReadTimeout
	}

	writeTimeout, exists := os.LookupEnv("SERVER_WRITE_TIMEOUT")
	if exists {
		var err error
		cfg.WriteTimeout, err = time.ParseDuration(writeTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid SERVER_WRITE_TIMEOUT: %w", err)
		}
	} else {
		cfg.WriteTimeout = defaultWriteTimeout
	}

	idleTimeout, exists := os.LookupEnv("SERVER_IDLE_TIMEOUT")
	if exists {
		var err error
		cfg.IdleTimeout, err = time.ParseDuration(idleTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid SERVER_IDLE_TIMEOUT: %w", err)
		}
	} else {
		cfg.IdleTimeout = defaultIdleTimeout
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid port number")
	}

	if cfg.ReadTimeout <= 0 || cfg.WriteTimeout <= 0 || cfg.IdleTimeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive")
	}

	return cfg, nil
}
