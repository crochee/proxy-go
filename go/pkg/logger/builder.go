package logger

type Builder interface {
	Debugf(format string, v ...interface{})
	Debug(message string)
	Infof(format string, v ...interface{})
	Info(message string)
	Warnf(format string, v ...interface{})
	Warn(message string)
	Errorf(format string, v ...interface{})
	Error(message string)
	Fatalf(format string, v ...interface{})
	Fatal(message string)
	Sync() error
}

type NoLogger struct {
}

func (n NoLogger) Debugf(string, ...interface{}) {
}

func (n NoLogger) Debug(string) {
}

func (n NoLogger) Infof(string, ...interface{}) {
}

func (n NoLogger) Info(string) {
}

func (n NoLogger) Warnf(string, ...interface{}) {
}

func (n NoLogger) Warn(string) {
}

func (n NoLogger) Errorf(string, ...interface{}) {
}

func (n NoLogger) Error(string) {
}

func (n NoLogger) Fatalf(string, ...interface{}) {
}

func (n NoLogger) Fatal(string) {
}

func (n NoLogger) Sync() error {
	return nil
}
