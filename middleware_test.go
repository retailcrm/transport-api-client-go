package transport_api_client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("custom header middleware", func(t *testing.T) {
		t.Parallel()

		var called bool
		next := DoerFunc(func(req *http.Request) (*http.Response, error) {
			called = true
			require.Equal(t, "yes", req.Header.Get("X-Test"))
			return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
		})

		mw := func(next HttpRequestDoer) HttpRequestDoer {
			return DoerFunc(func(req *http.Request) (*http.Response, error) {
				req.Header.Set("X-Test", "yes")
				return next.Do(req)
			})
		}

		doer := mw(next)
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		resp, err := doer.Do(req)

		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.True(t, called)
	})

	t.Run("middleware order", func(t *testing.T) {
		t.Parallel()

		var order []string

		mw1 := func(next HttpRequestDoer) HttpRequestDoer {
			return DoerFunc(func(req *http.Request) (*http.Response, error) {
				order = append(order, "mw1-before")
				resp, err := next.Do(req)
				order = append(order, "mw1-after")
				return resp, err
			})
		}

		mw2 := func(next HttpRequestDoer) HttpRequestDoer {
			return DoerFunc(func(req *http.Request) (*http.Response, error) {
				order = append(order, "mw2-before")
				resp, err := next.Do(req)
				order = append(order, "mw2-after")
				return resp, err
			})
		}

		final := DoerFunc(func(req *http.Request) (*http.Response, error) {
			order = append(order, "final")
			return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
		})

		doer := mw1(mw2(final))
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		_, err := doer.Do(req)

		require.NoError(t, err)
		require.Equal(t,
			[]string{"mw1-before", "mw2-before", "final", "mw2-after", "mw1-after"},
			order,
		)
	})
}
