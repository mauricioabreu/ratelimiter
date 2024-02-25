package tokenbucket

import (
	"errors"
	"sync"
	"time"
)

var (
	errOutOfTokens = errors.New("out of tokens")
)

type Entry struct {
	size       int
	lastUpdate time.Time
}

type TokenBucket struct {
	capacity int
	rate     int
	bucket   map[string]*Entry
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
		bucket:   make(map[string]*Entry),
		stop:     make(chan struct{}),
	}
}

// Add adds a token in the bucket for a given `key`
func (tb *TokenBucket) Add(key string) {
	tb.rw.Lock()
	defer tb.rw.Unlock()

	entry, exists := tb.bucket[key]
	if !exists {
		tb.bucket[key] = &Entry{size: tb.capacity, lastUpdate: time.Now()}
	} else if entry.size < tb.capacity {
		entry.size += 1
	}
}

// Remove removes a token in the bucket for a given `key`
// Here we decide to return an error in case the given `key`
// has no tokens to consume
func (tb *TokenBucket) Remove(key string) error {
	tb.rw.Lock()
	defer tb.rw.Unlock()

	entry, exists := tb.bucket[key]
	if !exists {
		tb.bucket[key] = &Entry{size: tb.capacity - 1, lastUpdate: time.Now()}
		return nil
	}

	if entry.size == 0 {
		return errOutOfTokens
	}

	entry.size -= 1

	return nil
}

// Remaining returns the total remaining tokens
// in the bucket for the given `key`
func (tb *TokenBucket) Remaining(key string) int {
	tb.rw.RLock()
	defer tb.rw.RUnlock()

	entry, exists := tb.bucket[key]
	if !exists {
		return 0
	}

	return entry.size
}

// Refill start a routine to refill the tokens for all the available keys
func (tb *TokenBucket) Refill() {
	ticker := time.NewTicker(time.Duration(tb.rate) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tb.rw.Lock()
			for _, entry := range tb.bucket {
				if entry.size < tb.capacity {
					entry.size += 1
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
