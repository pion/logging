// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package logging_test

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/pion/logging"
	"github.com/stretchr/testify/assert"
)

func testJSONLoggerLevels(t *testing.T, logger *logging.JSONLeveledLogger) {
	t.Helper()

	var outBuf bytes.Buffer
	logger.WithOutput(&outBuf)

	// Test Info level
	infoMsg := "this is an info message"
	logger.Info(infoMsg)
	output := outBuf.String()
	assert.True(t, strings.Contains(output, infoMsg), "Expected to find %q in %q", infoMsg, output)
	assert.True(t, strings.Contains(output, `"level":"INFO"`), "Expected JSON to contain INFO level")
	assert.True(t, strings.Contains(output, `"scope":"test"`), "Expected JSON to contain scope")

	// Test Debug level
	outBuf.Reset()
	debugMsg := "this is a debug message"
	logger.Debug(debugMsg)
	output = outBuf.String()
	assert.True(t, strings.Contains(output, debugMsg), "Expected to find %q in %q", debugMsg, output)
	assert.True(t, strings.Contains(output, `"level":"DEBUG"`), "Expected JSON to contain DEBUG level")

	// Test Warn level
	outBuf.Reset()
	warnMsg := "this is a warning message"
	logger.Warn(warnMsg)
	output = outBuf.String()
	assert.True(t, strings.Contains(output, warnMsg), "Expected to find %q in %q", warnMsg, output)
	assert.True(t, strings.Contains(output, `"level":"WARN"`), "Expected JSON to contain WARN level")

	// Test Error level
	outBuf.Reset()
	errMsg := "this is an error message"
	logger.Error(errMsg)
	output = outBuf.String()
	assert.True(t, strings.Contains(output, errMsg), "Expected to find %q in %q", errMsg, output)
	assert.True(t, strings.Contains(output, `"level":"ERROR"`), "Expected JSON to contain ERROR level")

	// Test Trace level
	outBuf.Reset()
	traceMsg := "this is a trace message"
	logger.Trace(traceMsg)
	output = outBuf.String()
	assert.True(t, strings.Contains(output, traceMsg), "Expected to find %q in %q", traceMsg, output)
	assert.True(t, strings.Contains(output, `"level":"TRACE"`), "Expected JSON to contain TRACE level")
}

func testJSONLoggerFormatting(t *testing.T, logger *logging.JSONLeveledLogger) {
	t.Helper()

	var outBuf bytes.Buffer
	logger.WithOutput(&outBuf)

	// Test formatted messages
	formatMsg := "formatted message with %s"
	arg := "argument"
	logger.Infof(formatMsg, arg)
	output := outBuf.String()
	expectedMsg := "formatted message with argument"
	assert.True(t, strings.Contains(output, expectedMsg), "Expected to find %q in %q", expectedMsg, output)
}

func testJSONLoggerLevelFiltering(t *testing.T, logger *logging.JSONLeveledLogger) {
	t.Helper()

	var outBuf bytes.Buffer
	logger.WithOutput(&outBuf)

	// Set level to WARN, so DEBUG and INFO should be filtered
	logger.SetLevel(logging.LogLevelWarn)

	// These should not be logged
	logger.Debug("debug message")
	logger.Info("info message")
	assert.Equal(t, 0, outBuf.Len(), "Debug and Info messages should not be logged at WARN level")

	// These should be logged
	logger.Warn("warn message")
	logger.Error("error message")
	output := outBuf.String()
	assert.True(t, strings.Contains(output, "warn message"), "Warn message should be logged")
	assert.True(t, strings.Contains(output, "error message"), "Error message should be logged")
}

func TestJSONLogger(t *testing.T) {
	logger := logging.NewJSONLeveledLoggerForScope("test", logging.LogLevelTrace, os.Stderr)

	testJSONLoggerLevels(t, logger)
	testJSONLoggerFormatting(t, logger)
	testJSONLoggerLevelFiltering(t, logger)
}

func TestJSONLoggerFactory(t *testing.T) {
	factory := logging.JSONLoggerFactory{
		Writer:          os.Stderr,
		DefaultLogLevel: logging.LogLevelWarn,
		ScopeLevels: map[string]logging.LogLevel{
			"foo": logging.LogLevelDebug,
		},
	}

	logger := factory.NewLogger("baz")
	bazLogger, ok := logger.(*logging.JSONLeveledLogger)
	assert.True(t, ok, "Invalid logger type")

	// Test that baz logger respects WARN level
	var outBuf bytes.Buffer
	bazLogger.WithOutput(&outBuf)
	bazLogger.Debug("debug message")
	assert.Equal(t, 0, outBuf.Len(), "Debug message should not be logged at WARN level")

	logger = factory.NewLogger("foo")
	fooLogger, ok := logger.(*logging.JSONLeveledLogger)
	assert.True(t, ok, "Invalid logger type")

	// Test that foo logger respects DEBUG level
	outBuf.Reset()
	fooLogger.WithOutput(&outBuf)
	fooLogger.Debug("debug message")
	output := outBuf.String()
	assert.True(t, strings.Contains(output, "debug message"), "Debug message should be logged at DEBUG level")
}

func TestNewJSONLoggerFactory(t *testing.T) {
	factory := logging.NewJSONLoggerFactory()

	disabled := factory.NewLogger("DISABLE")
	errorLevel := factory.NewLogger("ERROR")
	warnLevel := factory.NewLogger("WARN")
	infoLevel := factory.NewLogger("INFO")
	debugLevel := factory.NewLogger("DEBUG")
	traceLevel := factory.NewLogger("TRACE")

	disabledLogger, ok := disabled.(*logging.JSONLeveledLogger)
	assert.True(t, ok, "Missing disabled logger")

	errorLogger, ok := errorLevel.(*logging.JSONLeveledLogger)
	assert.True(t, ok, "Missing error logger")

	_, ok = warnLevel.(*logging.JSONLeveledLogger)
	assert.True(t, ok, "Missing warn logger")

	_, ok = infoLevel.(*logging.JSONLeveledLogger)
	assert.True(t, ok, "Missing info logger")

	_, ok = debugLevel.(*logging.JSONLeveledLogger)
	assert.True(t, ok, "Missing debug logger")

	_, ok = traceLevel.(*logging.JSONLeveledLogger)
	assert.True(t, ok, "Missing trace logger")

	// Test that all loggers are properly configured
	var outBuf bytes.Buffer
	disabledLogger.WithOutput(&outBuf)
	disabledLogger.Info("test message")
	assert.Equal(t, 0, outBuf.Len(), "Disabled logger should not log anything")

	outBuf.Reset()
	errorLogger.WithOutput(&outBuf)
	errorLogger.Error("error message")
	output := outBuf.String()
	assert.True(t, strings.Contains(output, "error message"), "Error logger should log error messages")
}

func TestJSONLoggerStructuredOutput(t *testing.T) {
	logger := logging.NewJSONLeveledLoggerForScope("test-scope", logging.LogLevelInfo, os.Stderr)
	var outBuf bytes.Buffer
	logger.WithOutput(&outBuf)

	logger.Info("test message")
	output := outBuf.String()

	// Verify it's valid JSON
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(output), &jsonData)
	assert.NoError(t, err, "Output should be valid JSON")

	// Verify required fields
	assert.Contains(t, jsonData, "time", "JSON should contain time field")
	assert.Contains(t, jsonData, "level", "JSON should contain level field")
	assert.Contains(t, jsonData, "msg", "JSON should contain msg field")
	assert.Contains(t, jsonData, "scope", "JSON should contain scope field")

	// Verify values
	assert.Equal(t, "INFO", jsonData["level"], "Level should be INFO")
	assert.Equal(t, "test message", jsonData["msg"], "Message should match")
	assert.Equal(t, "test-scope", jsonData["scope"], "Scope should match")
} 