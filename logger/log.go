// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package logger

import (
	"io"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newZap(level string, encoderFunc func(zapcore.EncoderConfig) zapcore.Encoder,
	skip int, w io.Writer, fields ...zap.Field) *zap.Logger {
	core := zapcore.NewCore(
		encoderFunc(newEncoderConfig()),
		zap.CombineWriteSyncers(zapcore.AddSync(w)),
		newLevel(level),
	).With(fields) // 自带node 信息
	// 大于error增加堆栈信息
	return zap.New(core).WithOptions(zap.AddCaller(), zap.AddCallerSkip(skip),
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

func newLevel(level string) zapcore.Level {
	l := zap.InfoLevel
	if temp, ok := map[string]zapcore.Level{
		"DEBUG":  zap.DebugLevel,
		"INFO":   zap.InfoLevel,
		"WARN":   zap.WarnLevel,
		"ERROR":  zap.ErrorLevel,
		"DPANIC": zap.DPanicLevel,
		"PANIC":  zap.PanicLevel,
		"FATAL":  zap.FatalLevel,
	}[strings.ToUpper(level)]; ok {
		l = temp
	}
	return l
}
