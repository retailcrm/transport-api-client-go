package transport_api_client

import (
	"context"
	"net/http"
)

const (
	transportTokenHeader = "X-Transport-Token"
)

// WithTransportToken sets a transport token in the HTTP request header for authentication purposes.
func WithTransportToken(token string) ClientOption {
	return WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set(transportTokenHeader, token)

		return nil
	})
}
