package transport_api_client

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const rateLimitPacingFactor = 0.99

// RateLimiter waits until a request can be sent.
type RateLimiter interface {
	Wait(ctx context.Context) error
}

// ResponseRateLimiter updates its state after each HTTP response.
type ResponseRateLimiter interface {
	RateLimiter
	Update(response *http.Response)
}

// Limiter applies a RateLimiter before forwarding a request and updates a ResponseRateLimiter after receiving a response.
func Limiter(l RateLimiter) Middleware {
	return func(next HttpRequestDoer) HttpRequestDoer {
		return DoerFunc(func(req *http.Request) (*http.Response, error) {
			if err := l.Wait(req.Context()); err != nil {
				return nil, err
			}

			response, err := next.Do(req)
			if limiter, ok := l.(ResponseRateLimiter); ok {
				limiter.Update(response)
			}
			return response, err
		})
	}
}

// defaultLimiter adapts golang.org/x/time/rate.Limiter to RateLimiter.
type defaultLimiter struct{ *rate.Limiter }

func (x defaultLimiter) Wait(ctx context.Context) error { return x.Limiter.Wait(ctx) }

// NewDefaultLimiter creates a token-bucket limiter with the specified rate and maximum burst size.
func NewDefaultLimiter(r float64, burst int) RateLimiter {
	return defaultLimiter{Limiter: rate.NewLimiter(rate.Limit(r), burst)}
}

// headerRateLimiter limits requests according to the current server rate-limit window.
type headerRateLimiter struct {
	mu          sync.Mutex
	pacer       *rate.Limiter
	limit       int64
	remaining   int64
	windowReset time.Time
}

// NewHeaderRateLimiter creates a limiter that follows server rate-limit response headers.
func NewHeaderRateLimiter(initialLimit int) RateLimiter {
	limiter := &headerRateLimiter{}
	if initialLimit > 0 {
		limiter.pacer = newPacer(int64(initialLimit))
	}
	return limiter
}

func (l *headerRateLimiter) Wait(ctx context.Context) error {
	for {
		l.mu.Lock()
		if l.windowReset.IsZero() {
			pacer := l.pacer
			l.mu.Unlock()
			if pacer == nil {
				return nil
			}
			return pacer.Wait(ctx)
		}
		if !time.Now().Before(l.windowReset) && l.remaining == 0 {
			l.remaining = l.limit
		}
		if l.remaining == 0 {
			wait := time.Until(l.windowReset)
			l.mu.Unlock()
			if err := waitForReset(ctx, wait); err != nil {
				return err
			}
			continue
		}
		if l.pacer == nil {
			l.pacer = newPacer(l.limit)
		}
		pacer := l.pacer
		l.mu.Unlock()
		if err := pacer.Wait(ctx); err != nil {
			return err
		}
		l.mu.Lock()
		if l.pacer == pacer && l.remaining > 0 {
			l.remaining--
			l.mu.Unlock()
			return nil
		}
		l.mu.Unlock()
	}
}

func waitForReset(ctx context.Context, wait time.Duration) error {
	timer := time.NewTimer(wait)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (l *headerRateLimiter) Update(response *http.Response) {
	limit, remaining, resetAt := rateLimitValuesFromResponse(response)
	l.mu.Lock()
	defer l.mu.Unlock()
	if resetAt.IsZero() {
		return
	}
	if l.windowReset.IsZero() {
		l.limit = limit
		l.remaining = remaining
		l.windowReset = resetAt
		l.setPacerLimit(limit)
		return
	}
	if resetAt.After(l.windowReset) {
		l.limit = limit
		l.remaining = remaining
		l.windowReset = resetAt
		l.setPacerLimit(limit)
		return
	}
	if resetAt.Equal(l.windowReset) {
		l.limit = min(l.limit, limit)
		l.remaining = min(l.remaining, remaining)
		l.setPacerLimit(l.limit)
	}
}

func (l *headerRateLimiter) setPacerLimit(limit int64) {
	if l.pacer == nil {
		l.pacer = newPacer(limit)
		return
	}
	l.pacer.SetLimit(rate.Limit(limit) * rateLimitPacingFactor)
}

func newPacer(limit int64) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(limit)*rateLimitPacingFactor, 1)
}

func rateLimitValuesFromResponse(response *http.Response) (int64, int64, time.Time) {
	if response == nil {
		return 0, 0, time.Time{}
	}

	limit, err := strconv.ParseInt(response.Header.Get("X-RateLimit-Limit"), 10, 64)
	if err != nil || limit <= 0 {
		return 0, 0, time.Time{}
	}
	remaining, err := strconv.ParseInt(response.Header.Get("X-RateLimit-Remaining"), 10, 64)
	if err != nil || remaining < 0 {
		return 0, 0, time.Time{}
	}
	reset, err := strconv.ParseInt(response.Header.Get("X-RateLimit-Reset"), 10, 64)
	if err != nil {
		return 0, 0, time.Time{}
	}
	return limit, remaining, time.Unix(reset, 0).Add(time.Second)
}
