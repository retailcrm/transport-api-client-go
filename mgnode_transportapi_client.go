package transport_api_client

import (
	"errors"
)

func (r ListChannelsResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r ActivateChannelResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r DeactivateChannelResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r UpdateChannelResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r ActivateTemplateResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r DeactivateTemplateResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r UpdateTemplateResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r UploadFileResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r UploadFileByUrlResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r GetFileUrlResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r DeleteMessageResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r SendMessageResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r EditMessageResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r AckMessageResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r SendHistoryMessageResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r DeleteMessageReactionResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r AddMessageReactionResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r MarkMessageReadResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r MarkMessagesReadUntilResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r RestoreMessageResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}

func (r GetTemplatesResp) Error() error {
	if r.JSONDefault != nil {
		return errors.New(r.JSONDefault.Errors[0])
	}

	return nil
}
