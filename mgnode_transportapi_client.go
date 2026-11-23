package transport_api_client

import (
	"errors"
	"fmt"
	"net/http"
)

type Err interface {
	Error() error
}

type HTTPResponseCarrier interface {
	Response() *http.Response
	ResponseBody() []byte
}

func ExtractError(resp Err, err error) error {
	if err != nil {
		return err
	}

	if resp != nil {
		if response, body, ok := responseDetails(resp); ok && response != nil && response.StatusCode == http.StatusTooManyRequests {
			cause := resp.Error()
			if cause == nil {
				cause = fmt.Errorf("unexpected HTTP status: %d", response.StatusCode)
			}
			return NewRateLimitError(response, body, cause)
		}
		if resp.Error() != nil {
			return resp.Error()
		}
		if response, _, ok := responseDetails(resp); ok && response != nil && response.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("unexpected HTTP status: %d", response.StatusCode)
		}
	}

	return nil
}

func defaultResponseError(response *ErrorResponse) error {
	if response == nil || len(response.Errors) == 0 {
		return nil
	}
	return errors.New(response.Errors[0])
}

func (r EnsureCustomerResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r ListChannelsResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r ActivateChannelResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r DeactivateChannelResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r UpdateChannelResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r ActivateTemplateResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r DeactivateTemplateResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r UpdateTemplateResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r UploadFileResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r UploadFileByUrlResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r GetFileUrlResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r DeleteMessageResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r SendMessageResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r EditMessageResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r AckMessageResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r SendHistoryMessageResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r DeleteMessageReactionResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r AddMessageReactionResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r MarkMessageReadResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r MarkMessagesReadUntilResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r RestoreMessageResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r GetTemplatesResp) Error() error {
	return defaultResponseError(r.JSONDefault)
}

func (r ListChannelsResp) Response() *http.Response { return r.HTTPResponse }
func (r ListChannelsResp) ResponseBody() []byte     { return r.Body }

func (r ActivateChannelResp) Response() *http.Response { return r.HTTPResponse }
func (r ActivateChannelResp) ResponseBody() []byte     { return r.Body }

func (r DeactivateChannelResp) Response() *http.Response { return r.HTTPResponse }
func (r DeactivateChannelResp) ResponseBody() []byte     { return r.Body }

func (r UpdateChannelResp) Response() *http.Response { return r.HTTPResponse }
func (r UpdateChannelResp) ResponseBody() []byte     { return r.Body }

func (r ActivateTemplateResp) Response() *http.Response { return r.HTTPResponse }
func (r ActivateTemplateResp) ResponseBody() []byte     { return r.Body }

func (r DeactivateTemplateResp) Response() *http.Response { return r.HTTPResponse }
func (r DeactivateTemplateResp) ResponseBody() []byte     { return r.Body }

func (r UpdateTemplateResp) Response() *http.Response { return r.HTTPResponse }
func (r UpdateTemplateResp) ResponseBody() []byte     { return r.Body }

func (r UploadFileResp) Response() *http.Response { return r.HTTPResponse }
func (r UploadFileResp) ResponseBody() []byte     { return r.Body }

func (r UploadFileByUrlResp) Response() *http.Response { return r.HTTPResponse }
func (r UploadFileByUrlResp) ResponseBody() []byte     { return r.Body }

func (r GetFileUrlResp) Response() *http.Response { return r.HTTPResponse }
func (r GetFileUrlResp) ResponseBody() []byte     { return r.Body }

func (r DeleteMessageResp) Response() *http.Response { return r.HTTPResponse }
func (r DeleteMessageResp) ResponseBody() []byte     { return r.Body }

func (r SendMessageResp) Response() *http.Response { return r.HTTPResponse }
func (r SendMessageResp) ResponseBody() []byte     { return r.Body }

func (r EditMessageResp) Response() *http.Response { return r.HTTPResponse }
func (r EditMessageResp) ResponseBody() []byte     { return r.Body }

func (r AckMessageResp) Response() *http.Response { return r.HTTPResponse }
func (r AckMessageResp) ResponseBody() []byte     { return r.Body }

func (r SendHistoryMessageResp) Response() *http.Response { return r.HTTPResponse }
func (r SendHistoryMessageResp) ResponseBody() []byte     { return r.Body }

func (r DeleteMessageReactionResp) Response() *http.Response { return r.HTTPResponse }
func (r DeleteMessageReactionResp) ResponseBody() []byte     { return r.Body }

func (r AddMessageReactionResp) Response() *http.Response { return r.HTTPResponse }
func (r AddMessageReactionResp) ResponseBody() []byte     { return r.Body }

func (r MarkMessageReadResp) Response() *http.Response { return r.HTTPResponse }
func (r MarkMessageReadResp) ResponseBody() []byte     { return r.Body }

func (r MarkMessagesReadUntilResp) Response() *http.Response { return r.HTTPResponse }
func (r MarkMessagesReadUntilResp) ResponseBody() []byte     { return r.Body }

func (r RestoreMessageResp) Response() *http.Response { return r.HTTPResponse }
func (r RestoreMessageResp) ResponseBody() []byte     { return r.Body }

func (r GetTemplatesResp) Response() *http.Response { return r.HTTPResponse }
func (r GetTemplatesResp) ResponseBody() []byte     { return r.Body }
