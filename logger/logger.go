// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *Logger

type Logger struct {
	level       zapcore.Level
	path        string
	logger      *zap.Logger
	loggerSugar *zap.SugaredLogger
}

// SetLoggerWriter return Logger writer
func SetLoggerWriter(path string) io.Writer {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    20,    //单个日志文件最大MaxSize*M大小 // megabytes
		MaxAge:     30,    //days
		MaxBackups: 50,    //备份数量
		Compress:   false, //不压缩
		LocalTime:  true,  //备份名采用本地时间
	}
}

// InitLogger init Logger
func InitLogger(opts ...Option) {
	logger = &Logger{}
	for _, opt := range opts {
		opt.Apply(logger)
	}
	if logger.path == "" {
		logger.logger = NewZap(logger.level, zapcore.NewConsoleEncoder, os.Stdout)
	} else {
		logger.logger = NewZap(logger.level, zapcore.NewConsoleEncoder, SetLoggerWriter(logger.path))
	}
	logger.loggerSugar = logger.logger.Sugar()
}

func Infof(format string, v ...interface{}) {
	if logger != nil {
		logger.loggerSugar.Infof(format, v...)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l != nil {
		l.loggerSugar.Infof(format, v...)
	}
}

func Info(message string) {
	if logger != nil {
		logger.logger.Info(message)
	}
}

func (l *Logger) Info(message string) {
	if l != nil {
		l.logger.Info(message)
	}
}

func Debugf(format string, v ...interface{}) {
	if logger != nil {
		logger.loggerSugar.Debugf(format, v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l != nil {
		l.loggerSugar.Debugf(format, v...)
	}
}

func Debug(message string) {
	if logger != nil {
		logger.logger.Debug(message)
	}
}

func (l *Logger) Debug(message string) {
	if l != nil {
		l.logger.Debug(message)
	}
}

func Warnf(format string, v ...interface{}) {
	if logger != nil {
		logger.loggerSugar.Warnf(format, v...)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l != nil {
		l.loggerSugar.Warnf(format, v...)
	}
}

func Warn(message string) {
	if logger != nil {
		logger.logger.Warn(message)
	}
}

func (l *Logger) Warn(message string) {
	if l != nil {
		l.logger.Warn(message)
	}
}

func Errorf(format string, v ...interface{}) {
	if logger != nil {
		logger.loggerSugar.Errorf(format, v...)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l != nil {
		l.loggerSugar.Errorf(format, v...)
	}
}

func Error(message string) {
	if logger != nil {
		logger.logger.Error(message)
	}
}

func (l *Logger) Error(message string) {
	if l != nil {
		l.logger.Error(message)
	}
}

func Fatalf(format string, v ...interface{}) {
	if logger != nil {
		logger.loggerSugar.Fatalf(format, v...)
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l != nil {
		l.loggerSugar.Fatalf(format, v...)
	}
}

func Fatal(message string) {
	if logger != nil {
		logger.logger.Fatal(message)
	}
}

func (l *Logger) Fatal(message string) {
	if l != nil {
		l.logger.Fatal(message)
	}
}
