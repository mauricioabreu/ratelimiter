package fixedwindowcounter_test

import (
	"testing"
	"time"

	"github.com/mauricioabreu/ratelimiter/internal/fixedwindowcounter"
	"github.com/stretchr/testify/assert"
)

func TestCurrWindow(t *testing.T) {
	fwc := fixedwindowcounter.New(60, 100).WithMockedTime(func() time.Time {
		return time.Date(2024, 3, 2, 22, 33, 10, 0, time.UTC)
	})

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

func TestIncrementMulitpleKeys(t *testing.T) {
	fwc := fixedwindowcounter.New(60, 100).WithMockedTime(func() time.Time {
		return time.Date(2024, 3, 2, 22, 33, 10, 0, time.UTC)
	})

	assert.Equal(t, 0, fwc.Count("127.0.0.1"))
	assert.Equal(t, 0, fwc.Count("127.0.0.2"))

	fwc.Increment("127.0.0.1")
	fwc.Increment("127.0.0.2")
	fwc.Increment("127.0.0.2")

	assert.Equal(t, 1, fwc.Count("127.0.0.1"))
	assert.Equal(t, 2, fwc.Count("127.0.0.2"))
}

func TestIncrementThreshold(t *testing.T) {
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

func TestReset(t *testing.T) {
	fwc := fixedwindowcounter.New(60, 100).WithMockedTime(func() time.Time {
		return time.Date(2024, 3, 2, 22, 33, 10, 0, time.UTC)
	})

	fwc.Increment("127.0.0.1")
	fwc.Increment("127.0.0.1")
	fwc.Increment("127.0.0.1")

	assert.Equal(t, 3, fwc.Count("127.0.0.1"))

	fwc.ExpirePastWindows()

	// Expiring the current window should not change anything
	assert.Equal(t, 3, fwc.Count("127.0.0.1"))

	fwc.WithMockedTime(func() time.Time {
		return time.Date(2024, 3, 2, 22, 34, 10, 0, time.UTC)
	})

	fwc.Increment("127.0.0.1")

	assert.Equal(t, 1, fwc.Count("127.0.0.1"))

	fwc.ExpirePastWindows()

	snapshot := fwc.EntriesByKey("127.0.0.1")

	assert.Nil(t, snapshot[1709418780])
}
