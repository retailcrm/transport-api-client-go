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
        "https://mg-s1.retailcrm.pro/api/bot/v1/",
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
    errors := *response.JSONDefault.Errors
    log.Printf("Error: %s", errors[0])
}

if response.JSON200 != nil {
    log.Printf("Message id: %d", response.JSON200.MessageId)
}
```
