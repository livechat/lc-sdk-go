package internal

import (
	"time"
)

// internal stats struct
type stats struct {
	sampleChan chan *APICallStats
}

func (s *stats) sample(sample *APICallStats) {
	select {
	case s.sampleChan <- sample:
	default:
	}
}

func newStats() *stats {
	return &stats{
		sampleChan: make(chan *APICallStats, 100),
	}
}

type APICallStats struct {
	Method        string
	ExecutionTime time.Duration
}
