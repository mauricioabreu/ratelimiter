package tokenbucket_test

import (
	"testing"
	"time"

	"github.com/mauricioabreu/ratelimiter/internal/tokenbucket"
	"github.com/stretchr/testify/assert"
)

func TestRemainingToken(t *testing.T) {
	tb := tokenbucket.New(10, 1)
	ip := "127.0.0.1"

	assert.Equal(t, tb.Remaining(ip), 0)
}

func TestAddToken(t *testing.T) {
	tb := tokenbucket.New(10, 1)
	ip := "127.0.0.1"

	assert.Equal(t, tb.Remaining(ip), 0)

	tb.Add(ip)
	tb.Add(ip)
	tb.Add(ip)

	assert.Equal(t, tb.Remaining(ip), 10)
}

func TestRemoveToken(t *testing.T) {
	tb := tokenbucket.New(10, 1)
	ip := "127.0.0.1"

	assert.Equal(t, tb.Remaining(ip), 0)

	err := tb.Remove(ip)

	assert.NoError(t, err)
	assert.Equal(t, tb.Remaining(ip), 9)
}

func TestRefill(t *testing.T) {
	tb := tokenbucket.New(10, 1)
	ip := "127.0.0.1"

	assert.Equal(t, tb.Remaining(ip), 0)

	tb.Remove(ip)

	assert.Equal(t, tb.Remaining(ip), 9)

	go tb.Refill()

	time.Sleep(2 * time.Second)

	tb.Stop()

	assert.True(t, tb.Remaining(ip) > 9)
}
