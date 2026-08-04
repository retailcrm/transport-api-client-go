package transport_api_client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeLimiter struct {
	err error
}

func (f fakeLimiter) Wait(ctx context.Context) error {
	return f.err
}

type fakeResponseLimiter struct {
	response *http.Response
}

func (f *fakeResponseLimiter) Wait(context.Context) error { return nil }

func (f *fakeResponseLimiter) Update(response *http.Response) {
	f.response = response
}

func TestLimiterMiddleware(t *testing.T) {
	t.Run("request passes when limiter allows", func(t *testing.T) {
		lim := fakeLimiter{err: nil}

		var called bool
		next := DoerFunc(func(req *http.Request) (*http.Response, error) {
			called = true
			return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
		})

		doer := Limiter(lim)(next)
		req, _ := http.NewRequest("GET", "http://example.com", nil)

		resp, err := doer.Do(req)

		require.NoError(t, err)
		require.True(t, called)
		require.Equal(t, 200, resp.StatusCode)
	})

	t.Run("response limiter receives response", func(t *testing.T) {
		limiter := &fakeResponseLimiter{}
		response := &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}
		doer := Limiter(limiter)(DoerFunc(func(*http.Request) (*http.Response, error) {
			return response, nil
		}))

		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com", nil)
		require.NoError(t, err)
		_, err = doer.Do(request)

		require.NoError(t, err)
		require.Same(t, response, limiter.response)
	})

	t.Run("request fails when limiter blocks", func(t *testing.T) {
		lim := fakeLimiter{err: errors.New("rate limit exceeded")}

		var called bool
		next := DoerFunc(func(req *http.Request) (*http.Response, error) {
			called = true
			return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
		})

		doer := Limiter(lim)(next)
		req, _ := http.NewRequest("GET", "http://example.com", nil)

		resp, err := doer.Do(req)

		require.Error(t, err)
		require.Nil(t, resp)
		require.False(t, called, "next.Do() should not be called if limiter blocks")
	})

	t.Run("NewDefaultLimiter allows only burst requests", func(t *testing.T) {
		lim := NewDefaultLimiter(2, 2)

		next := DoerFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
		})

		doer := Limiter(lim)(next)

		req1, _ := http.NewRequest("GET", "http://example.com/1", nil)
		start1 := time.Now()
		resp1, err1 := doer.Do(req1)
		require.NoError(t, err1)
		require.Equal(t, 200, resp1.StatusCode)
		require.Less(t, time.Since(start1), 50*time.Millisecond)

		req2, _ := http.NewRequest("GET", "http://example.com/2", nil)
		start2 := time.Now()
		resp2, err2 := doer.Do(req2)
		require.NoError(t, err2)
		require.Equal(t, 200, resp2.StatusCode)
		require.Less(t, time.Since(start2), 50*time.Millisecond)

		req3, _ := http.NewRequest("GET", "http://example.com/3", nil)
		start3 := time.Now()
		resp3, err3 := doer.Do(req3)
		require.NoError(t, err3)
		require.Equal(t, 200, resp3.StatusCode)

		delay := time.Since(start3)
		require.GreaterOrEqual(t, delay, 400*time.Millisecond)
	})
}

func TestHeaderRateLimiter(t *testing.T) {
	t.Parallel()

	t.Run("paces initial requests", func(t *testing.T) {
		t.Parallel()

		limiter := NewHeaderRateLimiter(10)
		require.NoError(t, limiter.Wait(t.Context()))

		startedAt := time.Now()
		require.NoError(t, limiter.Wait(t.Context()))
		assert.GreaterOrEqual(t, time.Since(startedAt), 90*time.Millisecond)
	})

	t.Run("paces requests", func(t *testing.T) {
		t.Parallel()

		limiter := NewHeaderRateLimiter(0).(*headerRateLimiter)
		limiter.Update(rateLimitResponse("10", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10)))

		require.NoError(t, limiter.Wait(t.Context()))
		startedAt := time.Now()
		require.NoError(t, limiter.Wait(t.Context()))

		assert.GreaterOrEqual(t, time.Since(startedAt), 90*time.Millisecond)
	})

	t.Run("ignores malformed headers", func(t *testing.T) {
		t.Parallel()

		limiter := NewHeaderRateLimiter(0).(*headerRateLimiter)
		limiter.Update(rateLimitResponse("invalid", "invalid"))

		assert.True(t, limiter.windowReset.IsZero())
	})

	t.Run("waits when server quota is exhausted", func(t *testing.T) {
		limiter := &headerRateLimiter{
			pacer:       newPacer(10),
			limit:       10,
			windowReset: time.Now().Add(time.Second),
		}
		ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond)
		defer cancel()

		assert.Error(t, limiter.Wait(ctx))
	})

	t.Run("prevents parallel requests after response", func(t *testing.T) {
		limiter := NewHeaderRateLimiter(0).(*headerRateLimiter)
		reset := strconv.FormatInt(time.Now().Add(2*time.Second).Unix(), 10)
		var calls atomic.Int64
		doer := Limiter(limiter)(DoerFunc(func(*http.Request) (*http.Response, error) {
			remaining := "0"
			if calls.Add(1) == 1 {
				remaining = "1"
			}
			return rateLimitResponse(remaining, reset), nil
		}))

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com", nil)
		require.NoError(t, err)
		_, err = doer.Do(req)
		require.NoError(t, err)

		var wg sync.WaitGroup
		results := make(chan error, 10)
		for range 10 {
			wg.Go(func() {
				ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
				defer cancel()
				request, requestErr := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
				if requestErr != nil {
					results <- requestErr
					return
				}
				_, requestErr = doer.Do(request)
				results <- requestErr
			})
		}
		wg.Wait()
		close(results)

		var successful int
		for result := range results {
			if result == nil {
				successful++
			}
		}
		assert.Equal(t, 1, successful)
		assert.EqualValues(t, 2, calls.Load())
	})

	t.Run("updates pacing after server limit changes", func(t *testing.T) {
		t.Parallel()

		var requestTimes []time.Time
		var requestTimesMu sync.Mutex
		server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
			requestTimesMu.Lock()
			requestTimes = append(requestTimes, time.Now())
			requestNumber := len(requestTimes)
			requestTimesMu.Unlock()

			if requestNumber == 1 {
				response.Header().Set("X-RateLimit-Limit", "100")
				response.Header().Set("X-RateLimit-Remaining", "99")
				response.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))
			} else {
				response.Header().Set("X-RateLimit-Limit", "5")
				response.Header().Set("X-RateLimit-Remaining", "4")
				response.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix()+1, 10))
			}
		}))
		defer server.Close()

		doer := Limiter(NewHeaderRateLimiter(0))(&http.Client{})
		for range 3 {
			request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL, nil)
			require.NoError(t, err)
			response, err := doer.Do(request)
			require.NoError(t, err)
			require.NoError(t, response.Body.Close())
		}

		requestTimesMu.Lock()
		defer requestTimesMu.Unlock()
		require.Len(t, requestTimes, 3)
		assert.GreaterOrEqual(t, requestTimes[2].Sub(requestTimes[1]), 150*time.Millisecond)
	})

	t.Run("keeps throughput with concurrent requests", func(t *testing.T) {
		dataProvider := []int{1, 5, 20, 50}
		for _, concurrency := range dataProvider {
			t.Run(strconv.Itoa(concurrency), func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
					response.Header().Set("X-RateLimit-Limit", "100")
					response.Header().Set("X-RateLimit-Remaining", "99")
					response.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10))
				}))
				defer server.Close()

				doer := Limiter(NewHeaderRateLimiter(0))(&http.Client{})
				warmupRequest, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL, nil)
				require.NoError(t, err)
				warmupResponse, err := doer.Do(warmupRequest)
				require.NoError(t, err)
				require.NoError(t, warmupResponse.Body.Close())

				const requests = 50
				var next atomic.Int64
				requestErrors := make(chan error, requests)
				startedAt := time.Now()
				var wg sync.WaitGroup
				for range concurrency {
					wg.Go(func() {
						for {
							if next.Add(1) > requests {
								return
							}
							request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL, nil)
							if err != nil {
								requestErrors <- err
								continue
							}
							response, err := doer.Do(request)
							if err == nil {
								err = response.Body.Close()
							}
							if err != nil {
								requestErrors <- err
							}
						}
					})
				}
				wg.Wait()
				close(requestErrors)
				for err := range requestErrors {
					require.NoError(t, err)
				}

				requestsPerSecond := float64(requests) / time.Since(startedAt).Seconds()
				assert.GreaterOrEqual(t, requestsPerSecond, 90.0)
			})
		}
	})

	t.Run("uses reset timestamp", func(t *testing.T) {
		t.Parallel()

		limiter := NewHeaderRateLimiter(0).(*headerRateLimiter)
		limiter.Update(rateLimitResponse("0", "100"))

		assert.Equal(t, time.Unix(100, 0).Add(time.Second), limiter.windowReset)
	})

	t.Run("accepts response from a new window after local wake", func(t *testing.T) {
		t.Parallel()

		windowReset := time.Now().Add(time.Second).Unix()
		limiter := &headerRateLimiter{
			limit:       10,
			windowReset: time.Now().Add(-time.Second),
		}
		require.NoError(t, limiter.Wait(t.Context()))
		limiter.Update(rateLimitResponse("9", strconv.FormatInt(windowReset+1, 10)))

		assert.Equal(t, time.Unix(windowReset+1, 0).Add(time.Second), limiter.windowReset)
		assert.Equal(t, int64(9), limiter.remaining)
	})
}

func TestRateLimitErrorRetryAfter(t *testing.T) {
	t.Run("uses reset", func(t *testing.T) {
		err := NewRateLimitError(rateLimitResponse("0", "100"), nil, nil)

		retryAfter, ok := err.RetryAfter(time.Unix(100, 0))
		require.True(t, ok)
		assert.Equal(t, time.Second, retryAfter)
	})
}

func rateLimitResponse(remaining, reset string) *http.Response {
	header := make(http.Header)
	header.Set("X-RateLimit-Limit", "10")
	header.Set("X-RateLimit-Remaining", remaining)
	header.Set("X-RateLimit-Reset", reset)
	return &http.Response{Header: header}
}
