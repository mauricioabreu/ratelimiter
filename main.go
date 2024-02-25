package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mauricioabreu/ratelimiter/internal/tokenbucket"
)

const bucketCapacity = 10

func TokenBucketMiddleware(tb *tokenbucket.TokenBucket, keyExtractor func(echo.Context) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := keyExtractor(c)
			if err := tb.Remove(key); err != nil {
				return c.String(http.StatusTooManyRequests, "Rate limit exceeded!\n")
			}

			return next(c)
		}
	}
}

func tbLimitedHandler(c echo.Context) error {
	return c.String(http.StatusOK, "You still have requests to spend!\n")
}

func keyExtractor(c echo.Context) string {
	return c.RealIP()
}

func main() {
	tb := tokenbucket.New(bucketCapacity, 1, 60)
	go tb.Refill()

	e := echo.New()
	e.GET("/tb/limited", tbLimitedHandler, TokenBucketMiddleware(tb, keyExtractor))

	e.Logger.Fatal(e.Start(":8080"))
}
