package e

import "net/http"

type Code interface {
	Code() string
	Status() int
	English() string
	Chinese() string
	String() string
}

type success struct{}

func (s success) Code() string {
	return "CPTS.CPTS-BUILD.000000"
}

func (s success) String() string {
	return s.English()
}

func (s success) Status() int {
	return http.StatusOK
}

func (s success) English() string {
	return "success"
}

func (s success) Chinese() string {
	return "成功"
}

type unknown struct{}

func (u unknown) Code() string {
	return "CPTS.CPTS-BUILD.000001"
}

func (u unknown) String() string {
	return u.English()
}

func (u unknown) Status() int {
	return http.StatusInternalServerError
}

func (u unknown) English() string {
	return "An unknown error occurred"
}

func (u unknown) Chinese() string {
	return "发生未知错误"
}
