# Transport API Client

A Go client library for interacting with the Message Gateway Transport API. This client provides methods for managing channels, templates, and message operations in a messaging platform.

## Installation

```bash
go get github.com/retailcrm/transport-api-client-go
```

## Usage

### Initializing the Client

```go
package main

import (
    "github.com/retailcrm/transport-api-client-go"
    "log"
)

func main() {
    client, err := transport_api_client.NewClientWithResponses(
        "https://mg-s1.retailcrm.pro/api/transport/v1/",
        transport_api_client.WithTransportToken("TRANSPORT_TOKEN"),
    )

    if err != nil {
        log.Fatalf("Error creating client: %v", err)

        return
    }
}
```

### REST API Examples

#### Sending a Message

```go
response, err := client.SendMessageWithResponse(
    context.Background(),
    transport_api_client.SendMessageJSONRequestBody{},
)

if err != nil {
    log.Fatalf("Error sending message: %v", err)
}

if response.JSONDefault != nil {
    log.Printf("Error: %s", response.JSONDefault.Errors[0])
}

if response.JSON200 != nil {
    log.Printf("Message id: %d", response.JSON200.MessageId)
}
```

#### Handling Webhooks

```go
func webhookHandler(c *gin.Context) {
	var webhookRequest WebhookRequest

	if err := c.ShouldBindJSON(&webhookRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	webhookData, err := webhookRequest.ValueByDiscriminator()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch v := webhookData.(type) {
	case WebhookMessageDeleted:
		deleteMessageHandler(c, v)
	case WebhookMessageSent:
		sentMessageHandler(c, v)
	case WebhookMessageRead:
		readMessageHandler(c, v)
	case WebhookMessageUpdated:
		updatedMessageHandler(c, v)
	case WebhookMessageReactionAdd:
		addReactionHandler(c, v)
	case WebhookMessageReactionDelete:
		deleteReactionHandler(c, v)
	case WebhookTemplateCreate:
		createTemplateHandler(c, v)
	case WebhookTemplateDelete:
		deleteTemplateHandler(c, v)
	case WebhookTemplateUpdate:
		updateTemplateHandler(c, v)
	}
}

func deleteMessageHandler(c *gin.Context, data WebhookMessageDeleted) {
    // Handle message deletion logic here
	c.JSON(http.StatusOK, WebhookEmptyResponse{})
}

func sentMessageHandler(c *gin.Context, data WebhookMessageSent) {
	externalChatID := "external-chat-id-12345"
	externalCustomerID := "external-customer-id-67890"
	externalMessageID := "external-message-id-abcde"

	c.JSON(http.StatusOK, WebhookSentMessageResponseData{
		Async:              false,
		ExternalChatID:     &externalChatID,
		ExternalCustomerID: &externalCustomerID,
		ExternalMessageID:  &externalMessageID,
	})
}

// Other handler functions would be defined similarly...
```

### Client with Logging and Rate Limiting

The library supports **middleware** to wrap HTTP requests.
Typical use cases are logging and rate limiting, but you can also implement your own (e.g. retries, tracing, headers injection).

#### Example

```go
package main

import (
    "context"
    "google.golang.org/grpc/internal/channelz"
    "log"
    "os"
    "time"

    "github.com/retailcrm/transport-api-client-go"
)

func main() {
    // standard Go logger
    stdLogger := log.New(os.Stdout, "[transport-api] ", log.LstdFlags)
    logger := transport_api_client.NewDefaultLogger(stdLogger)

    // rate limiter: 2 requests per second, burst up to 5
    limiter := transport_api_client.NewDefaultLimiter(2, 5)

    // create client with middlewares
    client, err := transport_api_client.NewClientWithResponses(
        "https://api.example.com",
        transport_api_client.WithMiddlewares(
            transport_api_client.Logging(logger),
            transport_api_client.Limiter(limiter),
        ),
    )
    if err != nil {
        log.Fatalf("failed to create client: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    resp, err := client.SendMessageWithResponse(
        ctx,
        transport_api_client.SendMessageRequest{},
    )
    if err != nil {
        log.Fatalf("request failed: %v", err)
    }

    log.Printf("status: %s", resp.Status)
}
```

#### Middleware execution order

Middlewares are applied **in the order they are passed** to `WithMiddlewares`.
That means the following code:

```go
WithMiddlewares(
    Logging(logger),
    Limiter(limiter),
)
```

will wrap the underlying HTTP client like this:

```
Request
   │
   ▼
Logging (start := time.Now())
   │           └── measures total time:
   │                 - waiting in Limiter
   │                 - network request
   │                 - response handling
   ▼
Limiter (may wait before sending)
   │
   ▼
Transport (http.Client → real HTTP request)
   │
   ▼
Logging (dur := time.Since(start))
```

So:

1. The request goes through `Logging` first,
2. then through `Limiter`,
3. and finally reaches the underlying HTTP transport.

### Writing Your Own Middleware

A middleware has the signature:

```go
type Middleware func(HttpRequestDoer) HttpRequestDoer
```

It receives the next `HttpRequestDoer` in the chain and must return a new one.
This allows you to implement **cross-cutting concerns** like logging, tracing, caching, retries, etc.

Middleware typically has three phases:

1. **Before** — runs before calling `next.Do(req)` (e.g. inject headers, modify context).
2. **Do** — forwards the request to the next middleware or transport.
3. **After** — runs after the response is received or an error occurred.

---

#### Example: Request ID Middleware

```go
func RequestIDMiddleware() transport_api_client.Middleware {
	return func(next transport_api_client.HttpRequestDoer) transport_api_client.HttpRequestDoer {
		return transport_api_client.DoerFunc(func(req *http.Request) (*http.Response, error) {
			// BEFORE: add a request ID into context and header
			reqID := uuid.New().String()
			ctx := context.WithValue(req.Context(), "requestID", reqID)
			req = req.WithContext(ctx)
			req.Header.Set("X-Request-ID", reqID)

			// DO: pass to the next middleware / transport
			resp, err := next.Do(req)

			// AFTER: log result with the request ID
			if err != nil {
				log.Printf("[req:%s] failed: %v", reqID, err)
				return nil, err
			}
			log.Printf("[req:%s] completed with status %d", reqID, resp.StatusCode)
			return resp, nil
		})
	}
}
```

#### Usage

```go
client, err := transport_api_client.NewClientWithResponses(
	"https://api.example.com",
	transport_api_client.WithMiddlewares(
		RequestIDMiddleware(),
		transport_api_client.Logging(logger),
	),
)
```

This middleware:

* **Before**: generates a unique request ID, stores it in the request context, and sets the `X-Request-ID` header.
* **Do**: forwards the request to the next handler.
* **After**: logs the outcome together with the request ID.

