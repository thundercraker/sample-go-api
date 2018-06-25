package logging

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogging_InterfaceImplementation_StackDriver(t *testing.T) {
	var _ Logging = StackdriverLogger{}       // Verify that T implements I.
	var _ Logging = (*StackdriverLogger)(nil) // Verify that *T implements I.
}

func TestLogging_InterfaceImplementation_Console(t *testing.T) {
	var _ Logging = WriterLogger{}       // Verify that T implements I.
	var _ Logging = (*WriterLogger)(nil) // Verify that *T implements I.
}

func TestLogging_LogWithWriter(t *testing.T) {

	bs := bytes.NewBufferString("")
	buf := bufio.NewWriter(bs)
	logger := NewWriterLogging(buf)

	logger.Verbose("A verbose message")
	logger.Debug("A debug message")
	logger.Info("An Info message")
	logger.Warning("A warning message")
	logger.Error("An error message")
	logger.Critical("A critical message")

	buf.Flush()
	output := bs.String()
	assert.True(t, len(output) > 0)
	assert.Equal(t,
		"[Verbose] A verbose message\n[Debug] A debug message\n[Info] An Info message\n"+
			"[Warning] A warning message\n[Error] An error message\n[Critical] A critical message\n",
		output)
}

func TestLogging_Formatting(t *testing.T) {

	bs := bytes.NewBufferString("")
	buf := bufio.NewWriter(bs)
	logger := NewWriterLogging(buf)

	logger.Critical("The value is %v", "500")
	buf.Flush()
	output := bs.String()
	assert.Equal(t, output, "[Critical] The value is 500\n")

	bs.Reset()
	logger.Critical("The value is %v and index is %f", "500", 232.2323)
	buf.Flush()
	output = bs.String()
	assert.Equal(t, "[Critical] The value is 500 and index is 232.232300\n", output)

	bs.Reset()
	logger.Critical("The value is %v and index is %f and params were %v",
		"500", 232.2323, []string{"data", "todo"})
	buf.Flush()
	output = bs.String()
	assert.Equal(t, "[Critical] The value is 500 and index is 232.232300 "+
		"and params were [data todo]\n", output)
}
