package transport_api_client

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func fakeDoer(status int, err error) HttpRequestDoer {
	return DoerFunc(func(req *http.Request) (*http.Response, error) {
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: status,
			Status:     http.StatusText(status),
			Body:       http.NoBody,
		}, nil
	})
}

func TestLoggingMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("success request is logged", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		logger := NewDefaultLogger(log.New(&buf, "", 0))

		doer := Logging(logger)(fakeDoer(200, nil))
		req, _ := http.NewRequest("GET", "http://example.com/test", nil)

		resp, err := doer.Do(req)

		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)

		logs := buf.String()
		require.Contains(t, logs, "HTTP GET http://example.com/test - started")
		require.Contains(t, logs, "HTTP GET http://example.com/test - 200 OK")
	})

	t.Run("error is logged", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		logger := NewDefaultLogger(log.New(&buf, "", 0))

		expectedErr := context.DeadlineExceeded
		doer := Logging(logger)(fakeDoer(0, expectedErr))
		req, _ := http.NewRequest("POST", "http://example.com/fail", nil)

		_, err := doer.Do(req)

		require.ErrorIs(t, err, expectedErr)

		logs := buf.String()
		require.Contains(t, logs, "HTTP POST http://example.com/fail - ERROR")
	})

	t.Run("log level from context", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		logger := NewDefaultLogger(log.New(&buf, "", 0))

		doer := Logging(logger)(fakeDoer(201, nil))
		ctx := WithLogLevel(context.Background(), LogLevelDebug)
		req, _ := http.NewRequestWithContext(ctx, "PUT", "http://example.com/item", nil)

		resp, err := doer.Do(req)

		require.NoError(t, err)
		require.Equal(t, 201, resp.StatusCode)

		logs := buf.String()
		require.Contains(t, logs, "[DEBUG] HTTP PUT http://example.com/item - started")
		require.Contains(t, logs, "[DEBUG] HTTP PUT http://example.com/item - 201 Created")
	})

	t.Run("log level string values", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, "DEBUG", LogLevelDebug.String())
		require.Equal(t, "INFO", LogLevelInfo.String())
		require.Equal(t, "WARN", LogLevelWarn.String())
		require.Equal(t, "ERROR", LogLevelError.String())
		require.Equal(t, "FATAL", LogLevelFatal.String())
		require.Equal(t, "UNKNOWN", LogLevel(999).String())
	})

	t.Run("log level from empty context", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		require.Equal(t, LogLevelInfo, LogLevelFromContext(ctx))
	})
}
