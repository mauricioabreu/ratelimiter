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
	fwc := fixedwindowcounter.New(100, 60).WithMockedTime(mockTime)

	currWindow := fwc.CurrWindow()

	assert.Equal(t, int64(1709418700), currWindow)
}

func TestIncrement(t *testing.T) {
	mockTime := func() time.Time {
		return time.Date(2024, 3, 2, 22, 33, 10, 0, time.UTC)
	}
	fwc := fixedwindowcounter.New(100, 60).WithMockedTime(mockTime)

	assert.Equal(t, 0, fwc.Count("127.0.0.1"))

	fwc.Increment("127.0.0.1")
	fwc.Increment("127.0.0.1")
	fwc.Increment("127.0.0.1")

	assert.Equal(t, 3, fwc.Count("127.0.0.1"))

	// Increment time to the next minute
	nextMockTime := func() time.Time {
		return time.Date(2024, 3, 2, 22, 34, 10, 0, time.UTC)
	}
	fwc.WithMockedTime(nextMockTime)

	assert.Equal(t, 0, fwc.Count("127.0.0.1"))

	fwc.Increment("127.0.0.1")
	fwc.Increment("127.0.0.1")

	assert.Equal(t, 2, fwc.Count("127.0.0.1"))
}
