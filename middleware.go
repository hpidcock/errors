package errors

import (
	"sync"
)

// Middleware are called to determine how to track an error.
type Middleware func(error) Synthesizer

// Synthesizer is returned by middleware to capture the original error and create a container.
type Synthesizer func(TrackedError) CauserWithCallers

var (
	mtx        sync.RWMutex
	middleware []Middleware
)

// RegisterMiddleware adds a middleware handler to the stack.
func RegisterMiddleware(m ...Middleware) {
	mtx.Lock()
	defer mtx.Unlock()
	middleware = append(middleware, m...)
}

func defaultSynthesizer(c TrackedError) CauserWithCallers {
	return &c
}

func synthesizer(err error) Synthesizer {
	mtx.RLock()
	defer mtx.RUnlock()
	for _, m := range middleware {
		s := m(err)
		if s != nil {
			return s
		}
	}
	return defaultSynthesizer
}
