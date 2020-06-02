# go-health: healthchecks for Go services

This is an attempt to gather in a small lib healthcheck modules and healthchecks for services written in Golang. 
The idea is that instead of following default tutorials' advice on just returning "healthy" as a string, have something 
more meaningful which would reflect the actual health of your system.

# Usage instructions

## Install

Install with
```bash
go get github.com/koslibpro/go-health
```

## Usage example
There is a full code example [here](example/main.go). The example is based on `mux` router, however you can use this lib
 with any http server or router of your choice.

Basically you create healthcheck module objects for the healthchecks you want to enable. Eg. for a database: 

```go
import "github.com/koslibpro/go-health/healthchecks/modules/db"

myDbHealthCheck := db.NewHealthCheckModule(
		myDb,
		"MyDbHealthCheck",
		30*time.Second,
	)
```

and then create a healthchecker instance:

```go
healthchecker, err := healthchecks.NewHealthCheck([]modules.HealthCheckModule{myDbHealthCheck})
```

You can then freely use this `healthchecker` instance in your `/health` API endpoint by calling the `Status()` func.

```go
healthchecker.Status()
```

An example function you can use as-is with your http handler is: 

```go
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
```

As you will notice, this simple function returns status code `200` if everything is ok in your app and no errors are 
returned by the `Status()` call, or `400` if something was wrong. 

Feel free to edit the actual response body of it, depending on what you need to display.

# Extending modules

I have added a basic set of modules initially, which cover database and redis connections. However, it's very easy to 
create your own modules. You can do that by implementing the `HealthCheck` interface (found [here](healthchecks/healthcheck.go)).

# Contributing

Feel free to send in PRs for modules you'd like to use, or improvements in general.  

# TODOs

1. Add tests (quick and dirty is fun, but let's tidy up).
2. Add modules for: mongodb, rabbitmq, (find more).