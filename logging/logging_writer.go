package logging

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

// implementation with io writers
type WriterLogger struct {
	writer *bufio.Writer
}

// creates a new implementation that writes to the given writer
func NewConsoleLogging() WriterLogger {
	return WriterLogger{
		writer: bufio.NewWriter(os.Stdout),
	}
}

// creates a new implementation that writes to the given writer
func NewWriterLogging(writer *bufio.Writer) WriterLogger {
	return WriterLogger{
		writer: writer,
	}
}

// creates a new interface that writer to a blank writer, for use during testing
// then the result of the logs may not be necessary
func NewMockLogging() WriterLogger {
	return WriterLogger{
		writer: bufio.NewWriter(bytes.NewBufferString("")),
	}
}

// write log with the given LogLevel, message and object
func (bundle WriterLogger) Log(level LogLevel, message string, vars []interface{}) {

	var severity string
	switch level {
	case LogLevel_VERBOSE:
		severity = "Verbose"
		break
	case LogLevel_DEBUG:
		severity = "Debug"
		break
	case LogLevel_INFO:
		severity = "Info"
		break
	case LogLevel_WARNING:
		severity = "Warning"
		break
	case LogLevel_ERROR:
		severity = "Error"
		break
	case LogLevel_CRITICAL:
		severity = "Critical"
		break
	}
	payload := fmt.Sprintf("[%s] %s", severity, message)
	if len(vars) > 0 {
		payload = fmt.Sprintf(payload, vars...)
	}
	bundle.writer.WriteString(payload + "\n")
}

// write Debug log with the given message and object
func (bundle WriterLogger) Debug(message string, vars ...interface{}) {
	bundle.Log(LogLevel_DEBUG, message, vars)
}

// write Verbose log with the given message and object
func (bundle WriterLogger) Verbose(message string, vars ...interface{}) {
	bundle.Log(LogLevel_VERBOSE, message, vars)
}

// write Info log with the given message and object
func (bundle WriterLogger) Info(message string, vars ...interface{}) {
	bundle.Log(LogLevel_INFO, message, vars)
}

// write Warning log with the given message and object
func (bundle WriterLogger) Warning(message string, vars ...interface{}) {
	bundle.Log(LogLevel_WARNING, message, vars)
}

// write Error log with the given message and object
func (bundle WriterLogger) Error(message string, vars ...interface{}) {
	bundle.Log(LogLevel_ERROR, message, vars)
}

// write Critical log with the given message and object
func (bundle WriterLogger) Critical(message string, vars ...interface{}) {
	bundle.Log(LogLevel_CRITICAL, message, vars)
}
