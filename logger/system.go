// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/3

package logger

var systemLogger Builder = NoLogger{}

// InitSystemLogger 初始化系统级日志对象
//
// @param: path 日志路径
// @param: level 日志等级
func InitSystemLogger(opts ...func(*option)) {
	opts = append(opts, func(option *option) {
		option.skip = 2
	})
	systemLogger = NewLogger(opts...)
}

// Debugf 打印Debug信息
//
// @param: format 格式信息
// @param: v 参数信息
func Debugf(format string, v ...interface{}) {
	systemLogger.Debugf(format, v...)
}

// Debug 打印Debug信息
//
// @param: message 信息
func Debug(message string) {
	systemLogger.Debug(message)
}

// Infof 打印Info信息
//
// @param: format 格式信息
// @param: v 参数信息
func Infof(format string, v ...interface{}) {
	systemLogger.Infof(format, v...)

}

// Info 打印Info信息
//
// @param: message 信息
func Info(message string) {
	systemLogger.Info(message)
}

// Warnf 打印Warn信息
//
// @param: format 格式信息
// @param: v 参数信息
func Warnf(format string, v ...interface{}) {
	systemLogger.Warnf(format, v...)
}

// Warn 打印Warn信息
//
// @param: message 信息
func Warn(message string) {
	systemLogger.Warn(message)
}

// Errorf 打印Error信息
//
// @param: format 格式信息
// @param: v 参数信息
func Errorf(format string, v ...interface{}) {
	systemLogger.Errorf(format, v...)
}

// Error 打印Error信息
//
// @param: message 信息
func Error(message string) {
	systemLogger.Error(message)
}

// Fatalf 打印Fatal信息
//
// @param: format 格式信息
// @param: v 参数信息
func Fatalf(format string, v ...interface{}) {
	systemLogger.Fatalf(format, v...)
}

// Fatal 打印Fatal信息
//
// @param: message 信息
func Fatal(message string) {
	systemLogger.Fatal(message)
}

// Exit 打印系统退出信息
//
// @param: message 信息
func Exit(message string) {
	systemLogger.Info(message)
	_ = systemLogger.Sync()
}
