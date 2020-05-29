# go-health: healthchecks for Go services

This is an attempt to gather in a small lib healthcheck modules and healthchecks for services written in Golang. 
The idea is that instead of following default tutorials' advice on just returning "healthy" as a string, have something 
more meaningful which would reflect the actual health of your system.

# Usage instructions
There is a full code example [here](example/main.go). The example is based on `mux` router, however you can use this lib
 with any http server or router of your choice.

Basically you create healthcheck module objects for the healthchecks you want to enable. Eg. for a database: 

```go
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

# Extending modules

I have added a basic set of modules initially, which cover database and redis connections. However, it's very easy to 
create your own modules. You can do that by implementing the `HealthCheck` interface (found [here](healthchecks/healthcheck.go)).

# Contributing

Feel free to send in PRs for modules you'd like to use, or improvements in general.  

# TODOs

1. Add tests (quick and dirty is fun, but let's tidy up).
2. Add modules for: mongodb, rabbitmq, (find more).