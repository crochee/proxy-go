// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/7/16

// Package e
package e

const (
	Success ErrorCode = "PROXY.2000000"

	// 1000~1999 系统级别
	Unknown ErrorCode = "PROXY.5001000"

	// 2000~2999 服务级别
	UserNotLogin ErrorCode = "PROXY.4012000"
)

var errorList = map[ErrorCode]string{
	Success: "Success",

	Unknown: "Unknown error",

	UserNotLogin: "User isn't login",
}
