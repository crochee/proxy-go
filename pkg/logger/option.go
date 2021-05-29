// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/6

// Package logger
package logger

import "go.uber.org/zap"

type option struct {
	path   string
	level  string
	skip   int
	fields []zap.Field
}

// Path gives path set log's path
func Path(path string) func(*option) {
	return func(o *option) { o.path = path }
}

// Level gives level set log's level
func Level(level string) func(*option) {
	return func(o *option) { o.level = level }
}

// WithFields gives fields set log's fields
func WithFields(fields ...zap.Field) func(*option) {
	return func(o *option) { o.fields = fields }
}
