package db

import (
	"database/sql"
	"time"
)

type HealthCheckModule struct {
	db            *sql.DB
	identifier    string
	checkInterval time.Duration
	ticker        *time.Ticker

	LastStatus bool
	LastError  error
}

// New generates a new HealthCheckModule instance with the given params
func New(db *sql.DB, identifier string, interval time.Duration) *HealthCheckModule {
	return &HealthCheckModule{
		db:            db,
		identifier:    identifier,
		checkInterval: interval,
	}
}

// Register enables the healthcheck and repeatedly runs status checks.
func (h *HealthCheckModule) Register() {
	h.ticker = time.NewTicker(h.checkInterval)

	go func() {
		for {
			select {
			case <-h.ticker.C:
				go func() {
					h.CheckStatus()
				}()
			}
		}
	}()
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
	err := h.db.Ping()
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
