// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONLoggerLevels(t *testing.T) {
	logger := NewJSONLeveledLoggerForScope("test", LogLevelTrace, os.Stderr)

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
}

func TestJSONLoggerFormatting(t *testing.T) {
	logger := NewJSONLeveledLoggerForScope("test", LogLevelTrace, os.Stderr)

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

func TestJSONLoggerLevelFiltering(t *testing.T) {
	logger := NewJSONLeveledLoggerForScope("test", LogLevelTrace, os.Stderr)

	var outBuf bytes.Buffer
	logger.WithOutput(&outBuf)

	// Set level to WARN, so DEBUG and INFO should be filtered
	logger.SetLevel(LogLevelWarn)

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

func TestJSONLoggerFactory(t *testing.T) {
	factory := JSONLoggerFactory{
		Writer:          os.Stderr,
		DefaultLogLevel: LogLevelWarn,
		ScopeLevels: map[string]LogLevel{
			"foo": LogLevelDebug,
		},
	}

	logger := factory.NewLogger("baz")
	bazLogger, ok := logger.(*JSONLeveledLogger)
	assert.True(t, ok, "Invalid logger type")

	// Test that baz logger respects WARN level
	var outBuf bytes.Buffer
	bazLogger.WithOutput(&outBuf)
	bazLogger.Debug("debug message")
	assert.Equal(t, 0, outBuf.Len(), "Debug message should not be logged at WARN level")

	logger = factory.NewLogger("foo")
	fooLogger, ok := logger.(*JSONLeveledLogger)
	assert.True(t, ok, "Invalid logger type")

	// Test that foo logger respects DEBUG level
	outBuf.Reset()
	fooLogger.WithOutput(&outBuf)
	fooLogger.Debug("debug message")
	output := outBuf.String()
	assert.True(t, strings.Contains(output, "debug message"), "Debug message should be logged at DEBUG level")
}

func TestNewJSONLoggerFactory(t *testing.T) {
	factory := NewJSONLoggerFactory()

	disabled := factory.NewLogger("DISABLE")
	errorLevel := factory.NewLogger("ERROR")
	warnLevel := factory.NewLogger("WARN")
	infoLevel := factory.NewLogger("INFO")
	debugLevel := factory.NewLogger("DEBUG")
	traceLevel := factory.NewLogger("TRACE")

	disabledLogger, ok := disabled.(*JSONLeveledLogger)
	assert.True(t, ok, "Missing disabled logger")

	errorLogger, ok := errorLevel.(*JSONLeveledLogger)
	assert.True(t, ok, "Missing error logger")

	_, ok = warnLevel.(*JSONLeveledLogger)
	assert.True(t, ok, "Missing warn logger")

	_, ok = infoLevel.(*JSONLeveledLogger)
	assert.True(t, ok, "Missing info logger")

	_, ok = debugLevel.(*JSONLeveledLogger)
	assert.True(t, ok, "Missing debug logger")

	_, ok = traceLevel.(*JSONLeveledLogger)
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

func TestJSONLoggerTraceOutput(t *testing.T) {
	logger := NewJSONLeveledLoggerForScope("trace-scope", LogLevelTrace, os.Stderr)
	var outBuf bytes.Buffer
	logger.WithOutput(&outBuf)

	logger.Trace("test message")
	output := outBuf.String()

	// Verify it's valid JSON
	var jsonData map[string]any
	err := json.Unmarshal([]byte(output), &jsonData)
	assert.NoError(t, err, "Output should be valid JSON")

	// Verify required fields
	assert.Contains(t, jsonData, "time", "JSON should contain time field")
	assert.Contains(t, jsonData, "level", "JSON should contain level field")
	assert.Contains(t, jsonData, "msg", "JSON should contain msg field")
	assert.Contains(t, jsonData, "scope", "JSON should contain scope field")

	// Verify values
	assert.Equal(t, "TRACE", jsonData["level"], "Level should be TRACE")
	assert.Equal(t, "test message", jsonData["msg"], "Message should match")
	assert.Equal(t, "trace-scope", jsonData["scope"], "Scope should match")
}

func TestJSONLoggerStructuredOutput(t *testing.T) {
	logger := NewJSONLeveledLoggerForScope("test-scope", LogLevelInfo, os.Stderr)
	var outBuf bytes.Buffer
	logger.WithOutput(&outBuf)

	logger.Info("test message")
	output := outBuf.String()

	// Verify it's valid JSON
	var jsonData map[string]any
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

func TestJSONLeveledLogger_logf_IncludesAdditionalArgs(t *testing.T) {
	factory := NewJSONLoggerFactory()
	factory.Writer = os.Stderr
	factory.DefaultLogLevel = LogLevelTrace

	l := factory.NewLogger("test-scope")
	jl, ok := l.(*JSONLeveledLogger)
	assert.True(t, ok, "Invalid logger type")

	var outBuf bytes.Buffer
	jl.WithOutput(&outBuf)

	args := []any{
		"method", "GET",
		"path", "/users",
		"duration_ms", 15,
		"ok", true,
	}

	jl.logf(slog.LevelInfo, "Processing request", args...)

	raw := strings.TrimSpace(outBuf.String())

	var jsonData map[string]any
	err := json.Unmarshal([]byte(raw), &jsonData)
	assert.NoError(t, err, "Output should be valid JSON")

	// base fields
	assert.Equal(t, "Processing request", jsonData["msg"])
	assert.Equal(t, "INFO", jsonData["level"])
	assert.Equal(t, "test-scope", jsonData["scope"])

	// additional args should appear as structured fields
	assert.Equal(t, "GET", jsonData["method"])
	assert.Equal(t, "/users", jsonData["path"])
	assert.EqualValues(t, 15, jsonData["duration_ms"]) // json.Unmarshal numbers -> float64
	assert.Equal(t, true, jsonData["ok"])
}

func clearLogEnv(t *testing.T) {
	t.Helper()

	for _, name := range []string{"DISABLE", "ERROR", "WARN", "INFO", "DEBUG", "TRACE"} {
		t.Setenv("PION_LOG_"+name, "")
		t.Setenv("PIONS_LOG_"+name, "")
	}
}

func TestNewJSONLoggerFactory_AllSetsDefaultToMaxLevel(t *testing.T) {
	clearLogEnv(t)

	t.Setenv("PION_LOG_INFO", "All")
	t.Setenv("PION_LOG_DEBUG", "ALL")
	t.Setenv("PION_LOG_TRACE", "all")

	factory := NewJSONLoggerFactory()

	assert.Equal(t, LogLevelTrace, factory.DefaultLogLevel)
	assert.Equal(t, 0, len(factory.ScopeLevels))
}

func TestNewJSONLoggerFactory_AllDoesNotLowerDefaultLevel(t *testing.T) {
	clearLogEnv(t)

	t.Setenv("PION_LOG_DISABLE", "all")

	factory := NewJSONLoggerFactory()
	assert.Equal(t, LogLevelError, factory.DefaultLogLevel)
}

func TestNewJSONLoggerFactory_ScopesAreSplitAndLowercased(t *testing.T) {
	clearLogEnv(t)

	t.Setenv("PION_LOG_DEBUG", "Foo,BAR")

	factory := NewJSONLoggerFactory()

	assert.Equal(t, LogLevelError, factory.DefaultLogLevel)

	assert.Equal(t, LogLevelDebug, factory.ScopeLevels["foo"])
	assert.Equal(t, LogLevelDebug, factory.ScopeLevels["bar"])
}

func TestNewJSONLoggerFactory_AllAndScopedInteract(t *testing.T) {
	clearLogEnv(t)

	t.Setenv("PION_LOG_WARN", "all")

	t.Setenv("PION_LOG_DEBUG", "foo")

	factory := NewJSONLoggerFactory()

	assert.Equal(t, LogLevelWarn, factory.DefaultLogLevel)
	assert.Equal(t, LogLevelDebug, factory.ScopeLevels["foo"])

	foo := factory.NewLogger("foo").(*JSONLeveledLogger) //nolint:forcetypeassert
	bar := factory.NewLogger("bar").(*JSONLeveledLogger) //nolint:forcetypeassert

	assert.Equal(t, LogLevelDebug, foo.level.Get(), "scope override should win")
	assert.Equal(t, LogLevelWarn, bar.level.Get(), "default should apply when no scope override")
}

func TestNewJSONLoggerFactory_Fallback(t *testing.T) {
	clearLogEnv(t)

	t.Setenv("PION_LOG_INFO", "")
	t.Setenv("PIONS_LOG_INFO", "all")

	factory := NewJSONLoggerFactory()
	assert.Equal(t, LogLevelInfo, factory.DefaultLogLevel)
}
