package fixedwindowcounter

import (
	"sync"
	"time"
)

type WindowEntry struct {
	count      int
	lastUpdate time.Time
}

type FixedWindowCounter struct {
	size      int
	threshold int
	counter   map[string]map[int64]*WindowEntry
	rw        sync.RWMutex
	getTime   func() time.Time
}

func New(size, threshold int) *FixedWindowCounter {
	return &FixedWindowCounter{
		size:      size,
		threshold: threshold,
		counter:   make(map[string]map[int64]*WindowEntry),
		getTime:   time.Now,
	}
}

func (f *FixedWindowCounter) WithMockedTime(mt func() time.Time) *FixedWindowCounter {
	f.getTime = mt
	return f
}

// currWindow returns the current window index
func (f *FixedWindowCounter) CurrWindow() int64 {
	now := f.getTime()
	return (now.Unix() / int64(f.size)) * int64(f.size)
}

// Increment increments counter for current window
// When counter cannot be incremeted because it has exceeded
// the size, return false
func (f *FixedWindowCounter) Increment(key string) bool {
	f.rw.Lock()
	defer f.rw.Unlock()

	currWindow := f.CurrWindow()

	if _, exists := f.counter[key]; !exists {
		f.counter[key] = make(map[int64]*WindowEntry)
	}

	if _, exists := f.counter[key][currWindow]; !exists {
		f.counter[key][currWindow] = &WindowEntry{count: 1, lastUpdate: f.getTime()}
		return true
	}

	entry := f.counter[key][currWindow]
	if entry.count >= f.threshold {
		return false
	}

	entry.count++
	entry.lastUpdate = f.getTime()

	return true
}

func (f *FixedWindowCounter) Count(key string) int {
	f.rw.RLock()
	defer f.rw.RUnlock()

	window := f.CurrWindow()
	if entry, exists := f.counter[key][window]; exists {
		return entry.count
	}

	return 0
}

func (f *FixedWindowCounter) Reset() {

}
