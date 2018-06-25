package server

import "cloud.google.com/go/logging"

// Create a generic interface that allows logging measurement
type Measurement interface {
	Log(name string, timeMillis int64) // log a measurement
}

// a model that represents a single measurement
type MeasurementModel struct {
	Name string // a name to identify the measurement
	Time int64  // the time taken for the code block being measured to execute
}

// blank measurement for development mode
type MeasurementBlank struct {
	_id uint // a uid for the object
}

// create a new instance of the Blank Measurement implementation
func NewBlankMeasurementTool() MeasurementBlank {
	return MeasurementBlank{
		_id: 1,
	}
}

// log a measurement
func (m MeasurementBlank) Log(name string, timeMillis int64) {
	// do nothing
}

// a implementation to provide Stackdriver measurement logging
type MeasurementStackdriver struct {
	flushSize int64
	logCount  int64
	logger    *logging.Logger
}

// create a new instace of the MeasurementStackdriver
func NewMeasurementStackdriver(logger *logging.Logger) MeasurementStackdriver {
	return MeasurementStackdriver{
		flushSize: 1,
		logCount:  0,
		logger:    logger,
	}
}

// log a measurement
func (m MeasurementStackdriver) Log(name string, timeMillis int64) {

	m.logger.Log(logging.Entry{
		Payload: MeasurementModel{
			Name: name,
			Time: timeMillis,
		},
	})

	m.logCount += 1
	if m.logCount >= m.flushSize {
		m.logger.Flush()
		m.logCount = 0
	}
}
