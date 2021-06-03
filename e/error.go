package e

type serviceError struct {
	Code
	Message string
}

func (s *serviceError) Error() string {
	if s.Message != "" {
		return s.Message
	}
	return s.Code.String()
}

func New(code Code, message string) error {
	return &serviceError{
		Code:    code,
		Message: message,
	}
}

func NewCode(code Code) error {
	return &serviceError{
		Code: code,
	}
}
