package metrics

import "time"

type APICallStats struct {
	Method        string
	ExecutionTime time.Duration
	Success       bool
}
