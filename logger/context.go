// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option interface {
	Apply(*Logger)
}

type Level string

func (l Level) Apply(log *Logger) {
	if log == nil {
		return
	}
	if coreLevel, ok := map[Level]zapcore.Level{
		DebugLevel:  zap.DebugLevel,
		InfoLevel:   zap.InfoLevel,
		WarnLevel:   zap.WarnLevel,
		ErrorLevel:  zap.ErrorLevel,
		DpanicLevel: zap.DPanicLevel,
		PanicLevel:  zap.PanicLevel,
		FatalLevel:  zap.FatalLevel,
	}[l]; ok {
		log.level = coreLevel
		return
	}
	log.level = zap.InfoLevel
}

const (
	DebugLevel  Level = "DEBUG"
	InfoLevel   Level = "INFO"
	WarnLevel   Level = "WARN"
	ErrorLevel  Level = "ERROR"
	DpanicLevel Level = "DPANIC"
	PanicLevel  Level = "PANIC"
	FatalLevel  Level = "FATAL"
)

type LogPath string

func (l LogPath) Apply(log *Logger) {
	if log == nil {
		return
	}
	log.path = string(l)
}

type Enable bool

func (e Enable) Apply(log *Logger) {
	if !e {
		log = (*Logger)(nil)
	}
}

type loggerKey struct{}

// With Adds fields.
func With(ctx context.Context, opts ...Option) context.Context {
	InitLogger(opts...)
	return context.WithValue(ctx, loggerKey{}, logger)
}

// FromContext Gets the logger from context.
func FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return nil
	}

	l, ok := ctx.Value(loggerKey{}).(*Logger)
	if !ok {
		return nil
	}
	return l
}
