package tokenbucket

import (
	"errors"
	"sync"
	"time"
)

var (
	errOutOfTokens = errors.New("out of tokens")
)

type TokenBucket struct {
	capacity int
	rate     int
	tokens   map[string]int
	rw       sync.RWMutex
	stop     chan struct{}
}

// New create a Token Bucket rate limiter
// TokenBucket rate limiter is based on giving a number of (`capacity`)
// tokens to be used. Tokens are refilled at `rate` capacity in seconds.
func New(capacity, rate int) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		rate:     rate,
		tokens:   make(map[string]int),
		stop:     make(chan struct{}),
	}
}

// Add adds a token in the bucket for a given `key`
func (tb *TokenBucket) Add(key string) {
	tb.rw.Lock()
	defer tb.rw.Unlock()

	val, exists := tb.tokens[key]
	if !exists {
		tb.tokens[key] = tb.capacity
	} else if val < tb.capacity {
		tb.tokens[key] += 1
	}
}

// Remove removes a token in the bucket for a given `key`
// Here we decide to return an error in case the given `key`
// has no tokens to consume
func (tb *TokenBucket) Remove(key string) error {
	tb.rw.Lock()
	defer tb.rw.Unlock()

	val, exists := tb.tokens[key]
	if !exists {
		tb.tokens[key] = tb.capacity - 1
		return nil
	}

	if val == 0 {
		return errOutOfTokens
	}

	tb.tokens[key] -= 1

	return nil
}

// Remaining returns the total remaining tokens
// in the bucket for the given `key`
func (tb *TokenBucket) Remaining(key string) int {
	tb.rw.RLock()
	defer tb.rw.RUnlock()

	return tb.tokens[key]
}

// Refill start a routine to refill the tokens for all the available keys
func (tb *TokenBucket) Refill() {
	ticker := time.NewTicker(time.Duration(tb.rate) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tb.rw.Lock()
			for key, val := range tb.tokens {
				if val < tb.capacity {
					tb.tokens[key] = val + 1
				}
			}
			tb.rw.Unlock()
		case <-tb.stop:
			return
		}
	}
}

// Stop stops refilling tokens
func (tb *TokenBucket) Stop() {
	close(tb.stop)
}
