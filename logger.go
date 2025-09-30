package transport_api_client

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Logger is an abstraction used by the Logging middleware.
// Implementations can log messages in any format or destination.
type Logger interface {
	Log(ctx context.Context, format string, args ...interface{})
}

// Logging is a middleware that logs outgoing HTTP requests and their results.
// It records the request method, URL, status code, and total duration
// (including waiting in other middlewares such as Limiter).
func Logging(l Logger) Middleware {
	return func(next HttpRequestDoer) HttpRequestDoer {
		return DoerFunc(func(req *http.Request) (*http.Response, error) {
			ctx := req.Context()
			start := time.Now()

			l.Log(WithLogLevel(ctx, LogLevelDebug), "HTTP %s %s - started", req.Method, req.URL.String())

			resp, err := next.Do(req)
			dur := time.Since(start)

			if err != nil {
				l.Log(
					WithLogLevel(ctx, LogLevelError), "HTTP %s %s - ERROR: %v (took %v)",
					req.Method, req.URL.String(), err, dur,
				)
				return nil, err
			}

			statusCode := 0
			statusText := ""
			if resp != nil {
				statusCode = resp.StatusCode
				statusText = resp.Status
			}

			l.Log(
				WithLogLevel(ctx, LogLevelDebug),
				"HTTP %s %s - %d %s (took %v)",
				req.Method, req.URL.String(), statusCode, statusText, dur,
			)

			return resp, nil
		})
	}
}

// defaultLogger is a Logger implementation based on the standard log.Logger.
// It automatically extracts log level from context (via LogLevelFromContext)
// and prefixes each message with the level string.
type defaultLogger struct{ *log.Logger }

// Log implements Logger by prefixing messages with a log level
// extracted from context, or falling back to INFO.
func (l *defaultLogger) Log(ctx context.Context, format string, args ...interface{}) {
	level := LogLevelFromContext(ctx)
	l.Logger.Printf("[%s] "+format, append([]interface{}{level}, args...)...)
}

// NewDefaultLogger wraps a standard log.Logger into a Logger compatible with the Logging middleware.
func NewDefaultLogger(l *log.Logger) Logger {
	return &defaultLogger{l}
}

// LogLevel defines severity levels for Logging middleware.
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type ctxKeyLogLevel struct{}

// WithLogLevel returns a new context with the given log level
func WithLogLevel(ctx context.Context, level LogLevel) context.Context {
	return context.WithValue(ctx, ctxKeyLogLevel{}, level)
}

// LogLevelFromContext extracts log level from context or returns LogLevelInfo
func LogLevelFromContext(ctx context.Context) LogLevel {
	if v := ctx.Value(ctxKeyLogLevel{}); v != nil {
		if lvl, ok := v.(LogLevel); ok {
			return lvl
		}
	}
	return LogLevelInfo
}
