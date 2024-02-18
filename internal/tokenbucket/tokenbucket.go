package tokenbucket

import (
	"errors"
	"sync"
)

var (
	errOutOfTokens = errors.New("out of tokens")
)

type TokenBucket struct {
	capacity int
	rate     int
	tokens   map[string]int
	rw       sync.RWMutex
}

func New(capacity, rate int) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		rate:     rate,
		tokens:   make(map[string]int),
	}
}

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

func (tb *TokenBucket) Remaining(key string) int {
	tb.rw.RLock()
	defer tb.rw.RUnlock()

	return tb.tokens[key]
}
