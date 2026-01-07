// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

// JSONLogger is an optional extension interface so users
// can type-assert a LeveledLogger to JSONLogger to access slog.
type JSONLogger interface {
	LeveledLogger
	Slog() *slog.Logger
}

// JSONLeveledLogger provides JSON structured logging using Go's slog library.
type JSONLeveledLogger struct {
	level  LogLevel
	writer *loggerWriter
	logger *slog.Logger
	scope  string
}

var _ JSONLogger = (*JSONLeveledLogger)(nil)

func (jl *JSONLeveledLogger) Slog() *slog.Logger {
	return jl.logger
}

// NewJSONLeveledLoggerForScope returns a configured JSON LeveledLogger.
func NewJSONLeveledLoggerForScope(scope string, level LogLevel, writer io.Writer) *JSONLeveledLogger {
	if writer == nil {
		writer = os.Stderr
	}

	// Create a JSON handler with custom options
	lw := &loggerWriter{output: writer}
	logger := slog.New(newJSONHandlerHelper(lw))

	return &JSONLeveledLogger{
		level:  level,
		writer: lw,
		logger: logger,
		scope:  scope,
	}
}

// WithOutput is a chainable configuration function which sets the logger's
// logging output to the supplied io.Writer.
func (jl *JSONLeveledLogger) WithOutput(output io.Writer) *JSONLeveledLogger {
	if output == nil {
		output = os.Stderr
	}
	jl.writer.SetOutput(output)

	return jl
}

// newJSONHandlerHelper creates a new JSON slog.Handler with custom formatting.
func newJSONHandlerHelper(w io.Writer) slog.Handler {
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.Level(-8), // Allow all levels, filter ourselves
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			// Customize timestamp format
			switch attr.Key {
			case slog.TimeKey:
				attr.Value = slog.StringValue(attr.Value.Time().Format(time.RFC3339))

				return attr

			case slog.LevelKey:
				if lvl, ok := attr.Value.Any().(slog.Level); ok {
					attr.Value = slogLevelToSlogStringValue(lvl)

					return attr
				}

				// if slog changes representation then leave it alone.
				return attr

			default:
				return attr
			}
		},
	})
}

// SetLevel sets the logger's logging level.
func (jl *JSONLeveledLogger) SetLevel(newLevel LogLevel) {
	jl.level.Set(newLevel)
}

// logf is the internal logging function that handles level checking and formatting.
func (jl *JSONLeveledLogger) logf(level slog.Level, msg string, args ...any) {
	if jl.level.Get() < logLevelToPionLevel(level) {
		return
	}

	// Create structured log entry
	attrs := []any{
		"scope", jl.scope,
	}

	// Add any additional arguments as key-value pairs
	if len(args) > 0 {
		attrs = append(attrs, args...)
	}

	jl.logger.Log(context.Background(), level, msg, attrs...)
}

// logfWithFormatf formats the message and calls logf.
func (jl *JSONLeveledLogger) logfWithFormatf(level slog.Level, format string, args ...any) {
	if jl.level.Get() < logLevelToPionLevel(level) {
		return
	}

	// Format the message
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	// Create structured log entry
	attrs := []any{
		"scope", jl.scope,
	}

	jl.logger.Log(context.Background(), level, msg, attrs...)
}

// Convert slog record levels to the exact strings you want in JSON.
func slogLevelToSlogStringValue(level slog.Level) slog.Value {
	switch level {
	case slog.Level(-8): // trace
		return slog.StringValue("TRACE")
	case slog.LevelDebug:
		return slog.StringValue("DEBUG")
	case slog.LevelInfo:
		return slog.StringValue("INFO")
	case slog.LevelWarn:
		return slog.StringValue("WARN")
	case slog.LevelError:
		return slog.StringValue("ERROR")
	default:
		return slog.StringValue("UNKNOWN")
	}
}

// Helper to convert slog levels to Pion log levels.
func logLevelToPionLevel(level slog.Level) LogLevel {
	switch level {
	case slog.Level(-8): // trace
		return LogLevelTrace
	case slog.LevelDebug:
		return LogLevelDebug
	case slog.LevelInfo:
		return LogLevelInfo
	case slog.LevelWarn:
		return LogLevelWarn
	case slog.LevelError:
		return LogLevelError
	default:
		return LogLevelDisabled
	}
}

// Trace emits the preformatted message if the logger is at or below LogLevelTrace.
func (jl *JSONLeveledLogger) Trace(msg string) {
	jl.logf(slog.Level(-8), msg) // slog.LevelTrace is -8
}

// Tracef formats and emits a message if the logger is at or below LogLevelTrace.
func (jl *JSONLeveledLogger) Tracef(format string, args ...any) {
	jl.logfWithFormatf(slog.Level(-8), format, args...) // slog.LevelTrace is -8
}

// Debug emits the preformatted message if the logger is at or below LogLevelDebug.
func (jl *JSONLeveledLogger) Debug(msg string) {
	jl.logf(slog.LevelDebug, msg)
}

// Debugf formats and emits a message if the logger is at or below LogLevelDebug.
func (jl *JSONLeveledLogger) Debugf(format string, args ...any) {
	jl.logfWithFormatf(slog.LevelDebug, format, args...)
}

// Info emits the preformatted message if the logger is at or below LogLevelInfo.
func (jl *JSONLeveledLogger) Info(msg string) {
	jl.logf(slog.LevelInfo, msg)
}

// Infof formats and emits a message if the logger is at or below LogLevelInfo.
func (jl *JSONLeveledLogger) Infof(format string, args ...any) {
	jl.logfWithFormatf(slog.LevelInfo, format, args...)
}

// Warn emits the preformatted message if the logger is at or below LogLevelWarn.
func (jl *JSONLeveledLogger) Warn(msg string) {
	jl.logf(slog.LevelWarn, msg)
}

// Warnf formats and emits a message if the logger is at or below LogLevelWarn.
func (jl *JSONLeveledLogger) Warnf(format string, args ...any) {
	jl.logfWithFormatf(slog.LevelWarn, format, args...)
}

// Error emits the preformatted message if the logger is at or below LogLevelError.
func (jl *JSONLeveledLogger) Error(msg string) {
	jl.logf(slog.LevelError, msg)
}

// Errorf formats and emits a message if the logger is at or below LogLevelError.
func (jl *JSONLeveledLogger) Errorf(format string, args ...any) {
	jl.logfWithFormatf(slog.LevelError, format, args...)
}

// JSONLoggerFactory defines levels by scopes and creates new JSONLeveledLogger.
type JSONLoggerFactory struct {
	Writer          io.Writer
	DefaultLogLevel LogLevel
	ScopeLevels     map[string]LogLevel
}

// NewJSONLoggerFactory creates a new JSONLoggerFactory.
func NewJSONLoggerFactory() *JSONLoggerFactory {
	factory := JSONLoggerFactory{}
	factory.DefaultLogLevel = LogLevelError
	factory.ScopeLevels = make(map[string]LogLevel)
	factory.Writer = os.Stderr

	logLevels := map[string]LogLevel{
		"DISABLE": LogLevelDisabled,
		"ERROR":   LogLevelError,
		"WARN":    LogLevelWarn,
		"INFO":    LogLevelInfo,
		"DEBUG":   LogLevelDebug,
		"TRACE":   LogLevelTrace,
	}

	for name, level := range logLevels {
		env := os.Getenv(fmt.Sprintf("PION_LOG_%s", name))

		if env == "" {
			env = os.Getenv(fmt.Sprintf("PIONS_LOG_%s", name))
		}

		if env == "" {
			continue
		}

		if strings.ToLower(env) == "all" {
			if factory.DefaultLogLevel < level {
				factory.DefaultLogLevel = level
			}

			continue
		}

		scopes := strings.Split(strings.ToLower(env), ",")
		for _, scope := range scopes {
			factory.ScopeLevels[scope] = level
		}
	}

	return &factory
}

// NewLogger returns a configured JSON LeveledLogger for the given scope.
func (f *JSONLoggerFactory) NewLogger(scope string) LeveledLogger {
	logLevel := f.DefaultLogLevel
	if f.ScopeLevels != nil {
		if scopeLevel, found := f.ScopeLevels[scope]; found {
			logLevel = scopeLevel
		}
	}

	return NewJSONLeveledLoggerForScope(scope, logLevel, f.Writer)
}
