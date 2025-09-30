package transport_api_client

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

type fakeLimiter struct {
	err error
}

func (f fakeLimiter) Wait(ctx context.Context) error {
	return f.err
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
