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

// jsonLeveledLogger provides JSON structured logging using Go's slog library.
type jsonLeveledLogger struct {
	level  LogLevel
	writer *loggerWriter
	logger *slog.Logger
	scope  string
}

var _ LeveledLogger = (*jsonLeveledLogger)(nil)

func (jl *jsonLeveledLogger) Slog() *slog.Logger {
	return jl.logger
}

// newJSONLeveledLoggerForScope returns a configured JSON LeveledLogger.
func newJSONLeveledLoggerForScope(scope string, level LogLevel, writer io.Writer) *jsonLeveledLogger {
	if writer == nil {
		writer = os.Stderr
	}

	// Create a JSON handler with custom options
	lw := &loggerWriter{output: writer}
	logger := slog.New(newJSONHandlerHelper(lw))

	return &jsonLeveledLogger{
		level:  level,
		writer: lw,
		logger: logger,
		scope:  scope,
	}
}

// WithOutput is a chainable configuration function which sets the logger's
// logging output to the supplied io.Writer.
func (jl *jsonLeveledLogger) WithOutput(output io.Writer) LeveledLogger {
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
func (jl *jsonLeveledLogger) SetLevel(newLevel LogLevel) {
	jl.level.Set(newLevel)
}

// logf is the internal logging function that handles level checking and formatting.
func (jl *jsonLeveledLogger) logf(level slog.Level, msg string, args ...any) {
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
func (jl *jsonLeveledLogger) logfWithFormatf(level slog.Level, format string, args ...any) {
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
func (jl *jsonLeveledLogger) Trace(msg string) {
	jl.logf(slog.Level(-8), msg) // slog.LevelTrace is -8
}

// Tracef formats and emits a message if the logger is at or below LogLevelTrace.
func (jl *jsonLeveledLogger) Tracef(format string, args ...any) {
	jl.logfWithFormatf(slog.Level(-8), format, args...) // slog.LevelTrace is -8
}

// Debug emits the preformatted message if the logger is at or below LogLevelDebug.
func (jl *jsonLeveledLogger) Debug(msg string) {
	jl.logf(slog.LevelDebug, msg)
}

// Debugf formats and emits a message if the logger is at or below LogLevelDebug.
func (jl *jsonLeveledLogger) Debugf(format string, args ...any) {
	jl.logfWithFormatf(slog.LevelDebug, format, args...)
}

// Info emits the preformatted message if the logger is at or below LogLevelInfo.
func (jl *jsonLeveledLogger) Info(msg string) {
	jl.logf(slog.LevelInfo, msg)
}

// Infof formats and emits a message if the logger is at or below LogLevelInfo.
func (jl *jsonLeveledLogger) Infof(format string, args ...any) {
	jl.logfWithFormatf(slog.LevelInfo, format, args...)
}

// Warn emits the preformatted message if the logger is at or below LogLevelWarn.
func (jl *jsonLeveledLogger) Warn(msg string) {
	jl.logf(slog.LevelWarn, msg)
}

// Warnf formats and emits a message if the logger is at or below LogLevelWarn.
func (jl *jsonLeveledLogger) Warnf(format string, args ...any) {
	jl.logfWithFormatf(slog.LevelWarn, format, args...)
}

// Error emits the preformatted message if the logger is at or below LogLevelError.
func (jl *jsonLeveledLogger) Error(msg string) {
	jl.logf(slog.LevelError, msg)
}

// Errorf formats and emits a message if the logger is at or below LogLevelError.
func (jl *jsonLeveledLogger) Errorf(format string, args ...any) {
	jl.logfWithFormatf(slog.LevelError, format, args...)
}

// jsonLoggerFactory defines levels by scopes and creates new jsonLeveledLogger.
type jsonLoggerFactory struct {
	writer          io.Writer
	defaultLogLevel LogLevel
	scopeLevels     map[string]LogLevel
}

var _ LoggerFactory = (*jsonLoggerFactory)(nil)

// JSONLoggerFactoryOption configures the JSON LoggerFactory.
type JSONLoggerFactoryOption func(*jsonLoggerFactory)

// WithJSONWriter overrides the writer used by JSON loggers.
func WithJSONWriter(writer io.Writer) JSONLoggerFactoryOption {
	return func(factory *jsonLoggerFactory) {
		if writer == nil {
			factory.writer = os.Stderr

			return
		}

		factory.writer = writer
	}
}

// WithJSONDefaultLevel overrides the default log level used by JSON loggers.
func WithJSONDefaultLevel(level LogLevel) JSONLoggerFactoryOption {
	return func(factory *jsonLoggerFactory) {
		factory.defaultLogLevel = level
	}
}

// WithJSONScopeLevels sets specific log levels for scopes, overriding env values.
func WithJSONScopeLevels(levels map[string]LogLevel) JSONLoggerFactoryOption {
	return func(factory *jsonLoggerFactory) {
		if levels == nil {
			return
		}

		if factory.scopeLevels == nil {
			factory.scopeLevels = make(map[string]LogLevel, len(levels))
		}

		for scope, level := range levels {
			factory.scopeLevels[strings.ToLower(scope)] = level
		}
	}
}

// NewJSONLoggerFactory creates a new LoggerFactory that emits JSON logs.
func NewJSONLoggerFactory(options ...JSONLoggerFactoryOption) LoggerFactory {
	factory := newJSONLoggerFactory()

	for _, option := range options {
		if option == nil {
			continue
		}

		option(factory)
	}

	return factory
}

// newJSONLoggerFactory creates a new JSON LoggerFactory.
func newJSONLoggerFactory() *jsonLoggerFactory {
	factory := jsonLoggerFactory{}
	factory.defaultLogLevel = LogLevelError
	factory.scopeLevels = make(map[string]LogLevel)
	factory.writer = os.Stderr

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
			if factory.defaultLogLevel < level {
				factory.defaultLogLevel = level
			}

			continue
		}

		scopes := strings.Split(strings.ToLower(env), ",")
		for _, scope := range scopes {
			factory.scopeLevels[scope] = level
		}
	}

	return &factory
}

// NewLogger returns a configured JSON LeveledLogger for the given scope.
func (f *jsonLoggerFactory) NewLogger(scope string) LeveledLogger {
	logLevel := f.defaultLogLevel
	if f.scopeLevels != nil {
		if scopeLevel, found := f.scopeLevels[scope]; found {
			logLevel = scopeLevel
		}
	}

	return newJSONLeveledLoggerForScope(scope, logLevel, f.writer)
}
