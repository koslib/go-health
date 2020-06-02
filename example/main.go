package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/koslibpro/go-health/healthchecks"
	"github.com/koslibpro/go-health/healthchecks/modules"
	"github.com/koslibpro/go-health/healthchecks/modules/db"
)

var (
	MyDb            *sql.DB
	MyDbHealthCheck modules.HealthCheckModule
)

// Handler holds information about your API handler. In our case it has router and healthchecker instances inside.
type Handler struct {
	router        *mux.Router
	healthchecker *healthchecks.HealthCheck
}

// NewHandler generates and returns a new Handler instance with a healthchecker
func NewHandler(healthchecker *healthchecks.HealthCheck) Handler {
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
	healthchecker := healthchecks.New([]modules.HealthCheckModule{MyDbHealthCheck})

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
