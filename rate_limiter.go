package transport_api_client

import (
	"context"
	"golang.org/x/time/rate"
	"net/http"
)

// RateLimiter abstracts a token bucket rate limiter.
type RateLimiter interface {
	Wait(ctx context.Context) error
}

// Limiter is a middleware that applies a RateLimiter before forwarding the request.
// If the limiter denies the request (e.g. due to context cancellation), an error is returned.
func Limiter(l RateLimiter) Middleware {
	return func(next HttpRequestDoer) HttpRequestDoer {
		return DoerFunc(func(req *http.Request) (*http.Response, error) {
			if err := l.Wait(req.Context()); err != nil {
				return nil, err
			}

			return next.Do(req)
		})
	}
}

// defaultLimiter is an adapter around golang.org/x/time/rate.Limiter.
type defaultLimiter struct{ *rate.Limiter }

func (x defaultLimiter) Wait(ctx context.Context) error { return x.Limiter.Wait(ctx) }

// NewDefaultLimiter creates a new token-bucket limiter that allows events up to rate r,
// with a maximum burst size of burst.
func NewDefaultLimiter(r float64, burst int) RateLimiter {
	return defaultLimiter{Limiter: rate.NewLimiter(rate.Limit(r), burst)}
}
