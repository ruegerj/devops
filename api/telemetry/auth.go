package telemetry

import "github.com/prometheus/client_golang/prometheus"

var SuccessfulAuthCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "auth_successful_count",
		Help: "Number of times the authentication was successfult for a request",
	},
)

var FailedAuthCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "auth_failed_count",
		Help: "Number of times the authentication failed for a request.",
	},
)
