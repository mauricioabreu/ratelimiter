package fixedwindowcounter

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
