package transport_api_client

import (
	"context"
	"net/http"
)

func WithTransportToken(token string) ClientOption {
	return WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("X-Transport-Token", token)

		return nil
	})
}
