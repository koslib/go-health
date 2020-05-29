package redis

import (
	"time"

	"github.com/go-redis/redis"
)

type HealthCheckModule struct {
	redisClient   *redis.Client
	identifier    string
	checkInterval time.Duration
	ticker        *time.Ticker

	LastStatus bool
	LastError  error
}

// NewHealthCheckModule generates a new HealthCheckModule instance with the given params
func NewHealthCheckModule(client *redis.Client, identifier string, interval time.Duration) *HealthCheckModule {
	return &HealthCheckModule{
		redisClient:   client,
		checkInterval: interval,
		identifier:    identifier,
	}
}

// Register enables the healthcheck and repeatedly runs status checks.
func (h *HealthCheckModule) Register() error {
	h.ticker = time.NewTicker(h.checkInterval)

	for {
		select {
		case <-h.ticker.C:
			go func() {
				h.IsHealthy()
			}()
		}
	}
}

// IsHealthy returns a bool which represents the actual state of the healthcheck module (failing/passing).
func (h *HealthCheckModule) IsHealthy() bool {
	if err := h.CheckStatus(); err != nil {
		return false
	}
	return true
}

// CheckStatus does the actual job of checking the internals of the module and decide on the health status.
func (h *HealthCheckModule) CheckStatus() error {
	cmd := h.redisClient.Ping()
	err := cmd.Err()
	if err != nil {
		h.LastError = err
		return err
	}
	h.LastError = nil
	return nil
}

// Identifier returns the name of this healthcheck module.
func (h *HealthCheckModule) Identifier() string {
	return h.identifier
}

// GetLastError returns the last error registered.
func (h *HealthCheckModule) GetLastError() error {
	return h.LastError
}
