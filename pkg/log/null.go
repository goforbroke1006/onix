package log

// NewNullLogger creates logger instance that do nothing (suitable for tests).
func NewNullLogger() *nullLogger { // nolint:golint,revive
	return &nullLogger{}
}

type nullLogger struct{}

var _ Logger = &nullLogger{}

func (n nullLogger) WithErr(err error) Logger { // nolint:ireturn
	return n
}

func (n nullLogger) Info(args ...interface{}) {}

func (n nullLogger) Infof(format string, args ...interface{}) {}

func (n nullLogger) Warn(args ...interface{}) {}

func (n nullLogger) Fatal(args ...interface{}) {}
