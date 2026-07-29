package transport_api_client

import (
	"context"
	"net/http"
	"runtime/debug"
	"sync"
)

const modulePath = "github.com/retailcrm/transport-api-client-go"

var defaultUserAgent = sync.OnceValue(func() string {
	version := "unknown"
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		if buildInfo.Main.Path == modulePath && buildInfo.Main.Version != "" {
			version = buildInfo.Main.Version
		}
		for _, dependency := range buildInfo.Deps {
			if dependency.Path == modulePath && dependency.Version != "" {
				version = dependency.Version
				break
			}
		}
	}
	return "retailcrm-transport-api-client-go/" + version
})

func DefaultUserAgent() string { return defaultUserAgent() }

func withDefaultUserAgent(_ context.Context, req *http.Request) error {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", DefaultUserAgent())
	}

	return nil
}

func WithUserAgent(userAgent string) ClientOption {
	return WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		if userAgent != "" {
			req.Header.Set("User-Agent", userAgent)
		}
		return nil
	})
}
