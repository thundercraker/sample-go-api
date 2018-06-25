package logging

import (
	"cloud.google.com/go/logging"
	"fmt"
)

// implementation of the logging interface to log to Google Stackdriver, this implementation
// is an encapsulation of cloud.google.com/go/logging's logging implementation
type StackdriverLogger struct {
	flushSize uint64
	logCount  uint64
	logger    *logging.Logger
}

// creates a new instance of the stack driver logger
func NewStackdriverLogger(logger *logging.Logger) StackdriverLogger {
	return StackdriverLogger{
		logger:    logger,
		flushSize: 1,
		logCount:  0,
	}
}

// write log with the given LogLevel, message and object
func (bundle StackdriverLogger) Log(level LogLevel, message string, vars []interface{}) {

	var severity logging.Severity
	switch level {
	case LogLevel_VERBOSE:
		severity = logging.Debug
		break
	case LogLevel_DEBUG:
		severity = logging.Debug
		break
	case LogLevel_INFO:
		severity = logging.Info
		break
	case LogLevel_WARNING:
		severity = logging.Warning
		break
	case LogLevel_ERROR:
		severity = logging.Error
		break
	case LogLevel_CRITICAL:
		severity = logging.Critical
		break
	}

	payload := fmt.Sprintf("[%s] %s", severity, message)
	if len(vars) > 0 {
		payload = fmt.Sprintf(payload, vars...)
	}
	bundle.logger.Log(logging.Entry{
		Severity: severity,
		Payload:  payload,
	})
	bundle.logCount += 1
	if bundle.logCount >= bundle.flushSize {
		bundle.logger.Flush()
		bundle.logCount = 0
	}
}

// write Debug log with the given message and object
func (bundle StackdriverLogger) Debug(message string, vars ...interface{}) {
	bundle.Log(LogLevel_DEBUG, message, vars)
}

// write Verbose log with the given message and object
func (bundle StackdriverLogger) Verbose(message string, vars ...interface{}) {
	bundle.Log(LogLevel_VERBOSE, message, vars)
}

// write Info log with the given message and object
func (bundle StackdriverLogger) Info(message string, vars ...interface{}) {
	bundle.Log(LogLevel_INFO, message, vars)
}

// write Warning log with the given message and object
func (bundle StackdriverLogger) Warning(message string, vars ...interface{}) {
	bundle.Log(LogLevel_WARNING, message, vars)
}

// write Error log with the given message and object
func (bundle StackdriverLogger) Error(message string, vars ...interface{}) {
	bundle.Log(LogLevel_ERROR, message, vars)
}

// write Critical log with the given message and object
func (bundle StackdriverLogger) Critical(message string, vars ...interface{}) {
	bundle.Log(LogLevel_CRITICAL, message, vars)
}
