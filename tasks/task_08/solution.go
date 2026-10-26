package main

import (
	"sync"
	"time"
)

type Clock interface{ Now() time.Time }
type Limiter struct {
	mu     sync.Mutex
	clock  Clock
	rate   float64
	burst  int
	tokens float64
	last   time.Time
}

func NewLimiter(clock Clock, ratePerSec float64, burst int) *Limiter {
	var tokens float64 = float64(burst)
	limiter := &Limiter{clock: clock, rate: ratePerSec, burst: burst, tokens: tokens}
	if clock != nil {
		limiter.last = clock.Now()
		return limiter
	}
	return limiter
}

func (l *Limiter) AllowN(n int) bool {
	if n <= 0 || n > l.burst || l.burst <= 0 || l.clock == nil {
		return false
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	currtime := l.clock.Now()
	if currtime.After(l.last) && l.rate > 0 {
		difftime := currtime.Sub(l.last).Seconds()
		l.tokens += difftime * l.rate
		if l.tokens > float64(l.burst) {
			l.tokens = float64(l.burst)
		}
		l.last = currtime
	}

	if l.tokens < float64(n) {
		return false
	}

	l.tokens -= float64(n)
	return true
}
