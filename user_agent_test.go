package transport_api_client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAgent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		options  []ClientOption
		expected string
	}{
		{name: "default", expected: DefaultUserAgent()},
		{
			name: "request editor",
			options: []ClientOption{WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
				req.Header.Set("User-Agent", "custom-client/1.0.0")
				return nil
			})},
			expected: "custom-client/1.0.0",
		},
		{name: "user agent option", options: []ClientOption{WithUserAgent("custom-client/2.0.0")}, expected: "custom-client/2.0.0"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, err := NewClient("https://example.com", append(tc.options, WithHTTPClient(DoerFunc(func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, tc.expected, req.Header.Get("User-Agent"))
				assert.Equal(t, []string{tc.expected}, req.Header.Values("User-Agent"))
				return &http.Response{StatusCode: http.StatusNoContent, Body: io.NopCloser(strings.NewReader(""))}, nil
			})))...)
			require.NoError(t, err)
			_, err = client.ListChannels(t.Context(), nil)
			require.NoError(t, err)
		})
	}
}
