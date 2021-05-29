// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/6

// Package logger
package logger

type option struct {
	path  string
	level string
	skip  int
}

// Path gives path set log's path
func Path(path string) func(*option) {
	return func(o *option) { o.path = path }
}

// Level gives level set log's level
func Level(level string) func(*option) {
	return func(o *option) { o.level = level }
}
