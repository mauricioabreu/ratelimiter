package fixedwindowcounter_test

import (
	"testing"
	"time"

	"github.com/mauricioabreu/ratelimiter/internal/fixedwindowcounter"
	"github.com/stretchr/testify/assert"
)

func TestCurrWindow(t *testing.T) {
	mockTime := func() time.Time {
		// Sat Mar 02 2024 22:31:40 GMT+0000
		return time.Date(2024, 3, 2, 22, 33, 10, 0, time.UTC)
	}
	fwc := fixedwindowcounter.New(60, 100).WithMockedTime(mockTime)

	currWindow := fwc.CurrWindow()

	assert.Equal(t, int64(1709418780), currWindow)
}

func TestIncrement(t *testing.T) {
	fwc := fixedwindowcounter.New(60, 100).WithMockedTime(func() time.Time {
		return time.Date(2024, 3, 2, 22, 33, 10, 0, time.UTC)
	})

	assert.Equal(t, 0, fwc.Count("127.0.0.1"))

	fwc.Increment("127.0.0.1")
	fwc.Increment("127.0.0.1")
	fwc.Increment("127.0.0.1")

	assert.Equal(t, 3, fwc.Count("127.0.0.1"))

	// Increment time to next seconds
	// Check if the current window is the same
	fwc.WithMockedTime(func() time.Time {
		return time.Date(2024, 3, 2, 22, 33, 15, 0, time.UTC)
	})

	assert.Equal(t, 3, fwc.Count("127.0.0.1"))

	// Increment time to the next minute
	fwc.WithMockedTime(func() time.Time {
		return time.Date(2024, 3, 2, 22, 34, 10, 0, time.UTC)
	})

	assert.Equal(t, 0, fwc.Count("127.0.0.1"))

	fwc.Increment("127.0.0.1")
	fwc.Increment("127.0.0.1")

	assert.Equal(t, 2, fwc.Count("127.0.0.1"))
}

func TestIncrementSize(t *testing.T) {
	lowThreshold := 3
	fwc := fixedwindowcounter.New(60, lowThreshold).WithMockedTime(func() time.Time {
		return time.Date(2024, 3, 2, 22, 33, 10, 0, time.UTC)
	})

	assert.Equal(t, 0, fwc.Count("127.0.0.1"))

	assert.True(t, fwc.Increment("127.0.0.1"))
	assert.True(t, fwc.Increment("127.0.0.1"))
	assert.True(t, fwc.Increment("127.0.0.1"))
	// Exceedes threshold
	assert.False(t, fwc.Increment("127.0.0.1"))
}
