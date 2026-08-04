package transport_api_client

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type RateLimitError struct {
	Response  *http.Response
	Body      []byte
	Limit     *int64
	Remaining *int64
	Reset     *time.Time
	Cause     error
}

func NewRateLimitError(response *http.Response, body []byte, cause error) *RateLimitError {
	err := &RateLimitError{Response: response, Body: body, Cause: cause}
	if response == nil {
		return err
	}
	err.Limit = parseRateLimitInt(response.Header.Get("X-RateLimit-Limit"))
	err.Remaining = parseRateLimitInt(response.Header.Get("X-RateLimit-Remaining"))
	if reset := parseRateLimitInt(response.Header.Get("X-RateLimit-Reset")); reset != nil {
		value := time.Unix(*reset, 0).Add(time.Second)
		err.Reset = &value
	}
	return err
}

func (e *RateLimitError) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return fmt.Sprintf("rate limit exceeded: %s", http.StatusText(http.StatusTooManyRequests))
}

func (e *RateLimitError) Unwrap() error { return e.Cause }

func (e *RateLimitError) RetryAfter(now time.Time) (time.Duration, bool) {
	if e.Reset == nil {
		return 0, false
	}
	return max(e.Reset.Sub(now), time.Duration(0)), true
}

func responseDetails(value any) (*http.Response, []byte, bool) {
	carrier, ok := value.(HTTPResponseCarrier)
	if !ok {
		return nil, nil, false
	}
	return carrier.Response(), carrier.ResponseBody(), true
}

func parseRateLimitInt(value string) *int64 {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed < 0 {
		return nil
	}
	return &parsed
}
