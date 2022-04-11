package log

// Logger describe project-level logger.
type Logger interface {
	WithErr(err error) Logger
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
}
