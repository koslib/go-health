package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/koslibpro/go-health"
	"github.com/koslibpro/go-health/modules/db"
)

var (
	MyDb            *sql.DB
	MyDbHealthCheck health.HealthCheckModule
)

// Handler holds information about your API handler. In our case it has router and healthchecker instances inside.
type Handler struct {
	router        *mux.Router
	healthchecker *health.HealthCheck
}

// NewHandler generates and returns a new Handler instance with a healthchecker
func NewHandler(healthchecker *health.HealthCheck) Handler {
	return Handler{
		router:        mux.NewRouter(),
		healthchecker: healthchecker,
	}
}

func main() {
	// Obviously you need to connect to your DB, which is not depicted in this example.

	// Create the healthchecks modules you need
	MyDbHealthCheck = db.New(
		MyDb,
		"MyDbHealthCheck",
		30*time.Second,
	)

	// Add them in the healthcheck
	healthchecker := health.New([]health.HealthCheckModule{MyDbHealthCheck})

	h := NewHandler(healthchecker)

	h.router.HandleFunc("/status", h.Status)
	http.Handle("/", h.router)
}

// Status is a mux handler func that calls the healthcheck status function and reflects the actual state of your app.
// The function has a simple purpose: if things are ok, then return status code 200, otherwise 400.
func (h Handler) Status(w http.ResponseWriter, r *http.Request) {
	responses := h.healthchecker.Status()
	for _, r := range responses {
		if r.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
