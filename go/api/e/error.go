package e

import "fmt"

type ResponseError struct {
	Code
	Message string
}

func (s *ResponseError) Error() string {
	if s.Message != "" {
		return s.Message
	}
	return s.Code.English()
}

func NewMsg(code Code, message string) error {
	return &ResponseError{
		Code:    code,
		Message: message,
	}
}

func New(code Code) error {
	return NewMsg(code, "")
}

func NewFormat(code Code, format string, v ...interface{}) error {
	return NewMsg(code, fmt.Sprintf(format, v...))
}
