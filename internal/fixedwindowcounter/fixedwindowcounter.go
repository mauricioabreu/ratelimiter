package fixedwindowcounter

import "time"

type FixedWindowCounter struct {
	size      int
	threshold int
}

func New(size, threshold int) *FixedWindowCounter {
	return &FixedWindowCounter{
		size:      size,
		threshold: threshold,
	}
}

// currWindow returns the current window index
func (f *FixedWindowCounter) currWindow() int64 {
	now := time.Now()
	return (now.Unix() / int64(f.size)) * int64(f.size)
}
