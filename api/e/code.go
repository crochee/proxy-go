package e

import (
	"fmt"
	"strconv"
)

type Code interface {
	Status() int
	English() string
	Chinese() string
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

func (e ErrorCode) English() string {
	return e.detail().E
}

func (e ErrorCode) Chinese() string {
	return e.detail().C
}

func (e ErrorCode) ErrorCode() string {
	return string(e)
}

func (e ErrorCode) detail() Detail {
	msg, ok := errorList[e]
	if !ok {
		if msg, ok = errorList[Unknown]; !ok {
			panic(fmt.Sprintf("not define %s", e))
		}
	}
	return msg
}
