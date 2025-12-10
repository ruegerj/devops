package telemetry

import "github.com/prometheus/client_golang/prometheus"

var UnlockCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "vault_unlock_count",
		Help: "Number of times the secret has been accessed sucessfully.",
	},
)
