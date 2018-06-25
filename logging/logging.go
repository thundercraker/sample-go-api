/*
	Package logging provides interfaces and implementations that provide a layer
	of abstraction between the actual logging implementation and the services that need
	to log information in a project

	The main interface provides 6 levels of message importance. From least to most critical they are:
	- VERBOSE
	- DEBUG
	- INFO
	- WARNING
	- INFO
	- WARNING
 	- CRITICAL
*/
package logging

// provides a logging abstraction
type Logging interface {
	Log(level LogLevel, message string, vars []interface{}) // the base logging method
	Verbose(message string, vars ...interface{})            // logs Verbose messages
	Debug(message string, vars ...interface{})              // logs Debug messages
	Info(message string, vars ...interface{})               // logs Info messages
	Warning(message string, vars ...interface{})            //logs Warning messages
	Error(message string, vars ...interface{})              // logs Error messages
	Critical(message string, vars ...interface{})           // logs Critical messages
}

// LogLevel indicates the severity of the log it represents
type LogLevel int32

const (
	// (0) The log entry has no assigned severity level
	LogLevel_VERBOSE LogLevel = 0
	// (100) Debug or trace information
	LogLevel_DEBUG LogLevel = 100
	// (200) Routine information, such as ongoing status or performance
	LogLevel_INFO LogLevel = 200
	// (300) Warning events might cause problems
	LogLevel_WARNING LogLevel = 300
	// (400) Error event, request failure
	LogLevel_ERROR LogLevel = 400
	// (500) Critical event, app outage
	LogLevel_CRITICAL LogLevel = 500
)
