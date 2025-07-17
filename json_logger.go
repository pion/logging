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

// JSONLeveledLogger provides JSON structured logging using Go's slog library
type JSONLeveledLogger struct {
	level  LogLevel
	writer *loggerWriter
	logger *slog.Logger
	scope  string
}

// NewJSONLeveledLoggerForScope returns a configured JSON LeveledLogger
func NewJSONLeveledLoggerForScope(scope string, level LogLevel, writer io.Writer) *JSONLeveledLogger {
	if writer == nil {
		writer = os.Stderr
	}

	// Create a JSON handler with custom options
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: slog.Level(-8), // Allow all levels, filter ourselves
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize timestamp format
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   slog.TimeKey,
					Value: slog.StringValue(a.Value.Time().Format(time.RFC3339)),
				}
			}
			return a
		},
	})

	logger := slog.New(handler)

	return &JSONLeveledLogger{
		level:  level,
		writer: &loggerWriter{output: writer},
		logger: logger,
		scope:  scope,
	}
}

// WithOutput is a chainable configuration function which sets the logger's
// logging output to the supplied io.Writer.
func (jl *JSONLeveledLogger) WithOutput(output io.Writer) *JSONLeveledLogger {
	jl.writer.SetOutput(output)
	// Recreate the logger with the new writer
	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: slog.Level(-8),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   slog.TimeKey,
					Value: slog.StringValue(a.Value.Time().Format(time.RFC3339)),
				}
			}
			return a
		},
	})
	jl.logger = slog.New(handler)
	return jl
}

// SetLevel sets the logger's logging level.
func (jl *JSONLeveledLogger) SetLevel(newLevel LogLevel) {
	jl.level.Set(newLevel)
}

// logf is the internal logging function that handles level checking and formatting
func (jl *JSONLeveledLogger) logf(level slog.Level, msg string, args ...any) {
	if jl.level.Get() < jl.logLevelToPionLevel(level) {
		return
	}

	// Create structured log entry
	attrs := []any{
		"scope", jl.scope,
		"level", jl.pionLevelToString(jl.logLevelToPionLevel(level)),
	}

	// Add any additional arguments as key-value pairs
	if len(args) > 0 {
		attrs = append(attrs, args...)
	}

	jl.logger.Log(context.Background(), level, msg, attrs...)
}

// logfWithFormat formats the message and calls logf
func (jl *JSONLeveledLogger) logfWithFormat(level slog.Level, format string, args ...any) {
	if jl.level.Get() < jl.logLevelToPionLevel(level) {
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
		"level", jl.pionLevelToString(jl.logLevelToPionLevel(level)),
	}

	jl.logger.Log(context.Background(), level, msg, attrs...)
}

// Helper function to convert slog levels to Pion log levels
func (jl *JSONLeveledLogger) logLevelToPionLevel(level slog.Level) LogLevel {
	switch level {
	case slog.Level(-8): // slog.LevelTrace is -8
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

// Helper function to convert Pion log levels to string
func (jl *JSONLeveledLogger) pionLevelToString(level LogLevel) string {
	switch level {
	case LogLevelTrace:
		return "TRACE"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelDisabled:
		return "DISABLED"
	default:
		return "UNKNOWN"
	}
}

// Trace emits the preformatted message if the logger is at or below LogLevelTrace.
func (jl *JSONLeveledLogger) Trace(msg string) {
	jl.logf(slog.Level(-8), msg) // slog.LevelTrace is -8
}

// Tracef formats and emits a message if the logger is at or below LogLevelTrace.
func (jl *JSONLeveledLogger) Tracef(format string, args ...any) {
	jl.logfWithFormat(slog.Level(-8), format, args...) // slog.LevelTrace is -8
}

// Debug emits the preformatted message if the logger is at or below LogLevelDebug.
func (jl *JSONLeveledLogger) Debug(msg string) {
	jl.logf(slog.LevelDebug, msg)
}

// Debugf formats and emits a message if the logger is at or below LogLevelDebug.
func (jl *JSONLeveledLogger) Debugf(format string, args ...any) {
	jl.logfWithFormat(slog.LevelDebug, format, args...)
}

// Info emits the preformatted message if the logger is at or below LogLevelInfo.
func (jl *JSONLeveledLogger) Info(msg string) {
	jl.logf(slog.LevelInfo, msg)
}

// Infof formats and emits a message if the logger is at or below LogLevelInfo.
func (jl *JSONLeveledLogger) Infof(format string, args ...any) {
	jl.logfWithFormat(slog.LevelInfo, format, args...)
}

// Warn emits the preformatted message if the logger is at or below LogLevelWarn.
func (jl *JSONLeveledLogger) Warn(msg string) {
	jl.logf(slog.LevelWarn, msg)
}

// Warnf formats and emits a message if the logger is at or below LogLevelWarn.
func (jl *JSONLeveledLogger) Warnf(format string, args ...any) {
	jl.logfWithFormat(slog.LevelWarn, format, args...)
}

// Error emits the preformatted message if the logger is at or below LogLevelError.
func (jl *JSONLeveledLogger) Error(msg string) {
	jl.logf(slog.LevelError, msg)
}

// Errorf formats and emits a message if the logger is at or below LogLevelError.
func (jl *JSONLeveledLogger) Errorf(format string, args ...any) {
	jl.logfWithFormat(slog.LevelError, format, args...)
}

// JSONLoggerFactory defines levels by scopes and creates new JSONLeveledLogger
type JSONLoggerFactory struct {
	Writer          io.Writer
	DefaultLogLevel LogLevel
	ScopeLevels     map[string]LogLevel
}

// NewJSONLoggerFactory creates a new JSONLoggerFactory
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

// NewLogger returns a configured JSON LeveledLogger for the given scope
func (f *JSONLoggerFactory) NewLogger(scope string) LeveledLogger {
	logLevel := f.DefaultLogLevel
	if f.ScopeLevels != nil {
		scopeLevel, found := f.ScopeLevels[scope]

		if found {
			logLevel = scopeLevel
		}
	}

	return NewJSONLeveledLoggerForScope(scope, logLevel, f.Writer)
} 