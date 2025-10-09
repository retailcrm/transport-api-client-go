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

func TestListChannelsWithResponse(t *testing.T) {
	t.Parallel()

	active := BooleanTrue
	channelTypes := ChannelTypeQuery{ChannelTypeTelegram, ChannelTypeWhatsapp}

	testCases := []struct {
		name          string
		params        *ListChannelsParams
		expectedQuery string
	}{
		{
			name: "all parameters",
			params: &ListChannelsParams{
				ID:     123,
				Active: &active,
				Types:  channelTypes,
			},
			expectedQuery: "active=true&id=123&types=telegram&types=whatsapp",
		},
		{
			name:          "empty parameters",
			params:        &ListChannelsParams{},
			expectedQuery: "id=0&types=",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "GET", req.Method)
				assert.Contains(t, req.URL.Path, "/channels")
				assert.Equal(t, tc.expectedQuery, req.URL.RawQuery)

				body := `{
					"data": [
						{
							"id": 1,
							"type": "telegram",
							"name": "Test Channel",
							"is_active": true,
							"settings": {},
							"activated_at": "2024-12-31T00:00:00.000000Z",
							"created_at": "2024-12-31T00:00:00.000000Z"
						}
					]
				}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			})

			client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
			require.NoError(t, err)

			resp, err := client.ListChannelsWithResponse(context.Background(), tc.params)
			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, resp.StatusCode())
		})
	}
}

func TestActivateChannelWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/channels")

		body := `{"data": {"id": 123}}`
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.ActivateChannelWithResponse(context.Background(), ActivateChannelJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
}

func TestActivateChannelWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/channels")

		body := `{"data": {"id": 123}}`
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.ActivateChannelWithBodyWithResponse(context.Background(), "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
}

func TestDeactivateChannelWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "DELETE", req.Method)
		assert.Contains(t, req.URL.Path, "/channels/123")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.DeactivateChannelWithResponse(context.Background(), 123)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestUpdateChannelWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "PUT", req.Method)
		assert.Contains(t, req.URL.Path, "/channels/123")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.UpdateChannelWithResponse(context.Background(), 123, UpdateChannelJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestUpdateChannelWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "PUT", req.Method)
		assert.Contains(t, req.URL.Path, "/channels/123")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.UpdateChannelWithBodyWithResponse(context.Background(), 123, "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestActivateTemplateWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/channels/123/templates")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.ActivateTemplateWithResponse(context.Background(), 123, ActivateTemplateJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
}

func TestActivateTemplateWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/channels/123/templates")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.ActivateTemplateWithBodyWithResponse(context.Background(), 123, "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
}

func TestDeactivateTemplateWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "DELETE", req.Method)
		assert.Contains(t, req.URL.Path, "/channels/123/templates/test_template")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.DeactivateTemplateWithResponse(context.Background(), 123, "test_template")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestUpdateTemplateWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "PUT", req.Method)
		assert.Contains(t, req.URL.Path, "/channels/123/templates/test_template")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.UpdateTemplateWithResponse(context.Background(), 123, "test_template", UpdateTemplateJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestUpdateTemplateWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "PUT", req.Method)
		assert.Contains(t, req.URL.Path, "/channels/123/templates/test_template")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.UpdateTemplateWithBodyWithResponse(context.Background(), 123, "test_template", "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestUploadFileWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/files/upload")

		body := `{"data": {"id": "file-123"}}`
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.UploadFileWithBodyWithResponse(context.Background(), "multipart/form-data", strings.NewReader("file data"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
}

func TestUploadFileByUrlWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/files/upload_by_url")

		body := `{"data": {"id": "file-123"}}`
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.UploadFileByUrlWithResponse(context.Background(), UploadFileByUrlJSONRequestBody{
		Url: "https://example.com/file.png",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
}

func TestUploadFileByUrlWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/files/upload_by_url")

		body := `{"data": {"id": "file-123"}}`
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.UploadFileByUrlWithBodyWithResponse(context.Background(), "application/json", strings.NewReader(`{"url": "https://example.com/file.png"}`))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
}

func TestGetFileUrlWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "GET", req.Method)
		assert.Contains(t, req.URL.Path, "/files/")

		body := `{"data": {"url": "https://example.com/file.png"}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.GetFileUrlWithResponse(context.Background(), "file-123")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestSendMessageWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages")

		body := `{"data": {"id": 123}}`
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.SendMessageWithResponse(context.Background(), SendMessageJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
}

func TestSendMessageWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages")

		body := `{"data": {"id": 123}}`
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.SendMessageWithBodyWithResponse(context.Background(), "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
}

func TestEditMessageWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "PUT", req.Method)
		assert.Contains(t, req.URL.Path, "/messages")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.EditMessageWithResponse(context.Background(), EditMessageJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestEditMessageWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "PUT", req.Method)
		assert.Contains(t, req.URL.Path, "/messages")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.EditMessageWithBodyWithResponse(context.Background(), "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestDeleteMessageWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "DELETE", req.Method)
		assert.Contains(t, req.URL.Path, "/messages")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.DeleteMessageWithResponse(context.Background(), DeleteMessageJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestDeleteMessageWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "DELETE", req.Method)
		assert.Contains(t, req.URL.Path, "/messages")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.DeleteMessageWithBodyWithResponse(context.Background(), "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestAckMessageWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/ack")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.AckMessageWithResponse(context.Background(), AckMessageJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestAckMessageWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/ack")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.AckMessageWithBodyWithResponse(context.Background(), "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestSendHistoryMessageWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/history")

		body := `{"data": {"id": 123}}`
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.SendHistoryMessageWithResponse(context.Background(), SendHistoryMessageJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
}

func TestSendHistoryMessageWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/history")

		body := `{"data": {"id": 123}}`
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.SendHistoryMessageWithBodyWithResponse(context.Background(), "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
}

func TestAddMessageReactionWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/reaction")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.AddMessageReactionWithResponse(context.Background(), AddMessageReactionJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestAddMessageReactionWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/reaction")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.AddMessageReactionWithBodyWithResponse(context.Background(), "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestDeleteMessageReactionWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "DELETE", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/reaction")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.DeleteMessageReactionWithResponse(context.Background(), DeleteMessageReactionJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestDeleteMessageReactionWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "DELETE", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/reaction")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.DeleteMessageReactionWithBodyWithResponse(context.Background(), "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestMarkMessageReadWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/read")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.MarkMessageReadWithResponse(context.Background(), MarkMessageReadJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestMarkMessageReadWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/read")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.MarkMessageReadWithBodyWithResponse(context.Background(), "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestMarkMessagesReadUntilWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/read_until")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.MarkMessagesReadUntilWithResponse(context.Background(), MarkMessagesReadUntilJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestMarkMessagesReadUntilWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/read_until")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.MarkMessagesReadUntilWithBodyWithResponse(context.Background(), "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestRestoreMessageWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/restore")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.RestoreMessageWithResponse(context.Background(), RestoreMessageJSONRequestBody{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestRestoreMessageWithBodyWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/messages/restore")

		body := `{"data": {"success": true}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.RestoreMessageWithBodyWithResponse(context.Background(), "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestGetTemplatesWithResponse(t *testing.T) {
	t.Parallel()

	mockDoer := DoerFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "GET", req.Method)
		assert.Contains(t, req.URL.Path, "/templates")

		body := `{
			"data": [
				{
					"code": "template1",
					"name": "Template 1"
				}
			]
		}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithResponses("https://example.com", WithHTTPClient(mockDoer))
	require.NoError(t, err)

	resp, err := client.GetTemplatesWithResponse(context.Background())
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}
