package modules

type HealthCheckModule interface {
	// Register enables the healthcheck and repeatedly runs status checks. This needs to run right after creating a
	// healthcheck module, otherwise the ticker will not get started.
	Register()

	// IsHealthy returns a bool which represents the actual state of the healthcheck module (failing/passing).
	IsHealthy() bool

	// CheckStatus does the actual job of checking the internals of the module and decide on the health status. Returns
	// an error if needed, and it's primarily called by `IsHealthy`.
	CheckStatus() error

	// Identifier returns the name of this healthcheck module.
	Identifier() string

	// GetLastError returns the last error registered. Can be nil.
	GetLastError() error
}
