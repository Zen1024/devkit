package log

func Debug(msg string) {
	logger.Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(msg string) {
	logger.Info(msg)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(msg string) {
	logger.Warn(msg)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(msg string) {
	logger.Error(msg)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(msg string) {
	logger.Fatal(msg)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func (l *Logger) Debug(msg string) {
	l.logf(msg, leveldebug)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(format, leveldebug, args...)
}

func (l *Logger) Info(msg string) {
	l.logf(msg, levelinfo)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(format, levelinfo, args...)
}

func (l *Logger) Warn(msg string) {
	l.logf(msg, levelwarn)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logf(format, levelwarn, args...)
}

func (l *Logger) Error(msg string) {
	l.logf(msg, levelerror)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(format, levelerror, args...)
}

func (l *Logger) Fatal(msg string) {
	l.logf(msg, levelfatal)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logf(format, levelfatal, args...)
}
