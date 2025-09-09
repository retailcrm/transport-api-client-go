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
