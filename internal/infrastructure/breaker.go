package infrastructure

import (
	"errors"
	"sync"
	"time"
)

var ErrCircuitOpen = errors.New("circuit breaker is OPEN")

const (
	StateClosed   = "CLOSED"
	StateOpen     = "OPEN"
	StateHalfOpen = "HALF-OPEN"
)

type CircuitBreaker struct {
	mu            sync.Mutex
	state         string
	failures      int
	threshold     int
	timeout       time.Duration
	lastErrorTime time.Time
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:     StateClosed,
		threshold: threshold,
		timeout:   timeout,
	}
}

func (cb *CircuitBreaker) Execute(req func() error) error {

	cb.mu.Lock()
	if cb.state == StateClosed {
	}

	if cb.state == StateHalfOpen {
		cb.mu.Unlock()
		return ErrCircuitOpen
	}
	if cb.state == StateOpen {
		if time.Since(cb.lastErrorTime) > cb.timeout {
			cb.state = StateHalfOpen
		} else {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
	}
	cb.mu.Unlock()

	err := req()
	cb.mu.Lock()
	if err == nil {
		cb.failures = 0
		cb.state = StateClosed
	} else {
		cb.failures++
		if cb.state == StateHalfOpen || cb.failures >= cb.threshold {
			cb.state = StateOpen
			cb.lastErrorTime = time.Now()
		}
	}
	cb.mu.Unlock()
	return err
}
