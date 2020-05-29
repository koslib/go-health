package healthchecks

import "github.com/koslibpro/go-health/healthchecks/modules"

// HealthCheck is a collection of healthcheck modules.
type HealthCheck struct {
	modules []modules.HealthCheckModule
}

// HealthcheckResponse is the response that is produced by triggering each healthcheck module, and contains useful info.
type HealthcheckResponse struct {
	ModuleIdentifier string `json:"module_identifier"`
	Status           bool   `json:"status"`
	Error            error  `json:"error"`
}

// NewHealthCheck generates and returns a new HealthCheck instance with the given set of modules registered.
func NewHealthCheck(modules []modules.HealthCheckModule) (*HealthCheck, error) {
	healthcheck := &HealthCheck{
		modules: modules,
	}
	for _, module := range healthcheck.modules {
		if err := module.Register(); err != nil {
			return nil, err
		}
	}
	return healthcheck, nil
}

// Status returns the latest status of the healthcheck, which is a collection of HealthcheckResponses.
func (h *HealthCheck) Status() []HealthcheckResponse {
	result := make([]HealthcheckResponse, 0)

	// Iterate all modules registered for their latest health status
	for _, module := range h.modules {
		if healthy := module.IsHealthy(); !healthy {
			result = append(result, HealthcheckResponse{
				ModuleIdentifier: module.Identifier(),
				Status:           module.IsHealthy(),
				Error:            module.GetLastError(),
			})
		}
	}

	return result
}
