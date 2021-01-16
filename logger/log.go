// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Encoder func(zapcore.EncoderConfig) zapcore.Encoder

func NewZap(level zapcore.Level, encoderFunc Encoder, w io.Writer, fields ...zap.Field) *zap.Logger {
	core := zapcore.NewCore(
		encoderFunc(newEncoderConfig()),
		zap.CombineWriteSyncers(zapcore.AddSync(w)),
		level,
	).With(fields)
	//大于error增加堆栈信息
	return zap.New(core).WithOptions(zap.AddCaller(), zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.DPanicLevel))
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}
