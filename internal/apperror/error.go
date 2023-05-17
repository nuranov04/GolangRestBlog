package apperror

import (
	"encoding/json"
	"fmt"
)

var (
	ErrorNotFound         = NewAppError("not found", "US-000003", "")
	UserAlreadyExist      = NewAppError("username or email already exist", "US-000005", "")
	PostTitleAlreadyExist = NewAppError("post title already exist", "US-000005", "")
	IdQueryParamError     = NewAppError("param id must be number", "US-00006", "")
	NotCorrectPassword    = NewAppError("password is not correct", "US-0007", "")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(message, code, developerMessage string) *AppError {

	return &AppError{
		Err:              fmt.Errorf(message),
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}

func systemError(developerMessage string) *AppError {
	return NewAppError("system error", "NS-000001", developerMessage)
}

func BadRequestError(message string) *AppError {
	return NewAppError(message, "NS-000002", "some thing wrong with data")
}

func APIError(code, message, developerMessage string) *AppError {
	return NewAppError(message, code, developerMessage)
}

func UnauthorizedError(message string) *AppError {
	return NewAppError(message, "NS-000003", "")
}
