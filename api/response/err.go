// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/4

package response

import (
	"fmt"
	"net/http"
	"strconv"

	"proxy-go/internal"
)

type ProxyError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func (p *ProxyError) Error() string {
	buffer := internal.GetBuffer()
	buffer.AppendString("code:")
	buffer.AppendString(p.Code.Value())
	buffer.AppendString(" message:")
	buffer.AppendString(p.Message)
	result := buffer.String()
	buffer.Free()
	return result
}

// Error according to the given message structure returns an ProxyError
func Error(code interface{}, message string) *ProxyError {
	return &ProxyError{
		Code:    &ProxyCode{code},
		Message: message,
	}
}

// ErrorWiths according to the given message and error structure returns an ProxyError
func ErrorWiths(code interface{}, err error, message string) *ProxyError {
	buffer := internal.GetBuffer()
	buffer.AppendString("error:")
	buffer.AppendString(err.Error())
	buffer.AppendString(",message:")
	buffer.AppendString(message)
	message = buffer.String()
	buffer.Free()
	return &ProxyError{
		Code:    &ProxyCode{code},
		Message: message,
	}
}

// ErrorWith according to the given error structure returns an ProxyError
func ErrorWith(code interface{}, err error) *ProxyError {
	return &ProxyError{
		Code:    &ProxyCode{code},
		Message: err.Error(),
	}
}

type Code interface {
	Status() int
	Value() string
}

type ProxyCode struct {
	value interface{}
}

func (p *ProxyCode) Status() int {
	status := http.StatusInternalServerError
	switch value := p.value.(type) {
	case string:
		if tempStatus, err := strconv.Atoi(value[:3]); err == nil {
			status = tempStatus
		}
	case int:
		if value >= 100 && value < 600 {
			status = value
		}
	default:
	}
	return status
}

func (p *ProxyCode) Value() string {
	return fmt.Sprint(p.value)
}
