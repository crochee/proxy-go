package e

import (
	"fmt"
	"strconv"
)

type Code interface {
	Status() int
	Detail() string
	ErrorCode() string
}

type ErrorCode string

func (e ErrorCode) Status() int {
	status, err := strconv.Atoi(string(e)[6:9])
	if err != nil {
		panic(err)
	}
	return status
}

func (e ErrorCode) ErrorCode() string {
	return string(e)
}

func (e ErrorCode) Detail() string {
	msg, ok := errorList[e]
	if !ok {
		if msg, ok = errorList[Unknown]; !ok {
			panic(fmt.Sprintf("not define %s", e))
		}
	}
	return msg
}
