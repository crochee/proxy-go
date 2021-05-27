// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package logger

import (
	"io"
	"os"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultLogSizeM int = 20
	DefaultMaxZip   int = 50
	MaxLogDays      int = 30
)

func setLoggerWriter(path string) io.Writer {
	if path == "" {
		return os.Stdout
	}
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    DefaultLogSizeM, //单个日志文件最大MaxSize*M大小 // megabytes
		MaxAge:     MaxLogDays,      //days
		MaxBackups: DefaultMaxZip,   //备份数量
		Compress:   false,           //不压缩
		LocalTime:  true,            //备份名采用本地时间
	}
}

// NewLogger 初始化日志对象
//
// @param: path 日志路径
// @param: level 日志等级
func NewLogger(opts ...func(*option)) *Logger {
	l := &Logger{
		option: option{
			path:  "",
			level: "INFO",
			skip:  1,
		},
	}
	for _, opt := range opts {
		opt(&l.option)
	}
	var encode func(zapcore.EncoderConfig) zapcore.Encoder
	if l.option.path == "" {
		encode = zapcore.NewConsoleEncoder
	} else {
		encode = zapcore.NewJSONEncoder
	}
	l.Logger = newZap(l.option.level, encode, l.option.skip, setLoggerWriter(l.option.path))
	l.LoggerSugar = l.Logger.Sugar()

	return l
}

type Logger struct {
	Logger      *zap.Logger
	LoggerSugar *zap.SugaredLogger
	option
}

// Debugf 打印Debug信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.LoggerSugar.Debugf(format, v...)
}

// Debug 打印Debug信息
//
// @param: message 格式信息
func (l *Logger) Debug(message string) {
	l.Logger.Debug(message)
}

// Infof 打印Info信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Infof(format string, v ...interface{}) {
	l.LoggerSugar.Infof(format, v...)
}

// Info 打印Info信息
//
// @param: message 格式信息
func (l *Logger) Info(message string) {
	l.Logger.Info(message)
}

// Warnf 打印Warn信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.LoggerSugar.Warnf(format, v...)
}

// Warn 打印Warn信息
//
// @param: message 信息
func (l *Logger) Warn(message string) {
	l.Logger.Warn(message)
}

// Errorf 打印Error信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.LoggerSugar.Errorf(format, v...)
}

// Error 打印Error信息
//
// @param: message 信息
func (l *Logger) Error(message string) {
	l.Logger.Error(message)
}

// Fatalf 打印Fatalf信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.LoggerSugar.Errorf(format, v...)
}

// Fatal 打印Fatal信息
//
// @param: message 信息
func (l *Logger) Fatal(message string) {
	l.Logger.Error(message)
}

func (l *Logger) Sync() error {
	var resultErr error
	if err := l.Logger.Sync(); err != nil {
		resultErr = multierr.Append(resultErr, err)
	}
	if err := l.LoggerSugar.Sync(); err != nil {
		resultErr = multierr.Append(resultErr, err)
	}
	return resultErr
}
