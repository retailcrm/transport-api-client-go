package transport_api_client

import (
	"net/http"
)

// Middleware wraps an HttpRequestDoer to add cross-cutting functionality
// such as logging, retries, tracing, or rate limiting.
type Middleware func(HttpRequestDoer) HttpRequestDoer

// WithMiddlewares applies a chain of middlewares to the client.
// Middlewares are applied in the order they are passed.
func WithMiddlewares(mws ...Middleware) ClientOption {
	return func(c *Client) error {
		if c.Client == nil {
			c.Client = &http.Client{}
		}

		for _, mw := range mws {
			c.Client = mw(c.Client)
		}
		return nil
	}
}

type DoerFunc func(*http.Request) (*http.Response, error)

func (f DoerFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }
