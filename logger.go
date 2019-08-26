package logging

import (
	"fmt"
	"io"
	"net"

	"os"
	"strings"
	"sync/atomic"
)

// LogLevel represents the level at which the logger will emit log messages
type LogLevel int32

// Set updates the LogLevel to the supplied value
func (ll *LogLevel) Set(newLevel LogLevel) {
	atomic.StoreInt32((*int32)(ll), int32(newLevel))
}

// Get retrieves the current LogLevel value
func (ll *LogLevel) Get() LogLevel {
	return LogLevel(atomic.LoadInt32((*int32)(ll)))
}

func (ll LogLevel) String() string {
	switch ll {
	case LogLevelDisabled:
		return "Disabled"
	case LogLevelError:
		return "Error"
	case LogLevelWarn:
		return "Warn"
	case LogLevelInfo:
		return "Info"
	case LogLevelDebug:
		return "Debug"
	case LogLevelTrace:
		return "Trace"
	default:
		return "UNKNOWN"
	}
}

const (
	// LogLevelDisabled completely disables logging of any events
	LogLevelDisabled LogLevel = iota
	// LogLevelError is for fatal errors which should be handled by user code,
	// but are logged to ensure that they are seen
	LogLevelError
	// LogLevelWarn is for logging abnormal, but non-fatal library operation
	LogLevelWarn
	// LogLevelInfo is for logging normal library operation (e.g. state transitions, etc.)
	LogLevelInfo
	// LogLevelDebug is for logging low-level library information (e.g. internal operations)
	LogLevelDebug
	// LogLevelTrace is for logging very low-level library information (e.g. network traces)
	LogLevelTrace
)

// LoggerFactory is the basic pion LoggerFactory interface
type LoggerFactory interface {
	NewLogger(scope string) LeveledLogger
}

// Formatter is the actual log entry encoder
type Formatter interface {
	Msg(message string)
	Msgf(format string, args ...interface{})
	Bool(key string, b bool)
	Err(err error)
	Float32(key string, f float32)
	Float64(key string, f float64)
	IPAddr(key string, ip net.IP)
	Int(key string, i int)
	Int16(key string, i int16)
	Int32(key string, i int32)
	Int64(key string, i int64)
	Int8(key string, i int8)
	Str(key, val string)
	Uint(key string, i uint)
	Uint16(key string, i uint16)
	Uint32(key string, i uint32)
	Uint64(key string, i uint64)
	Uint8(key string, i uint8)
}

// FormatterFactory ...
type FormatterFactory func(LogLevel) Formatter

// Logger represents a pion Logger
type Logger struct {
	FormatterFactory FormatterFactory
	Lvl              LogLevel
}

// LeveledLogger is the legacy pion Logger interface
type LeveledLogger = *Logger

// Event is a helper structure to be able to chain events
type Event struct {
	formatter Formatter
}

// InfoLvl creates a new event with info lvl
func (l *Logger) InfoLvl() *Event { return l.newEvent(LogLevelInfo) }

// DebugLvl creates a new event with debug lvl
func (l *Logger) DebugLvl() *Event { return l.newEvent(LogLevelDebug) }

// WarnLvl creates a new event with warn lvl
func (l *Logger) WarnLvl() *Event { return l.newEvent(LogLevelWarn) }

// ErrorLvl creates a new event with error lvl
func (l *Logger) ErrorLvl() *Event { return l.newEvent(LogLevelError) }

// TraceLvl creates a new event with tracelvl
func (l *Logger) TraceLvl() *Event { return l.newEvent(LogLevelTrace) }

// Trace is a legacy method to report a trace lvl entry
func (l *Logger) Trace(msg string) { l.TraceLvl().Msg(msg) }

// Tracef is a legacy method to report a trace lvl entry
func (l *Logger) Tracef(format string, args ...interface{}) { l.TraceLvl().Msgf(format, args...) }

// Debug is a legacy method to report a trace lvl entry
func (l *Logger) Debug(msg string) { l.DebugLvl().Msg(msg) }

// Debugf is a legacy method to report a trace lvl entry
func (l *Logger) Debugf(format string, args ...interface{}) { l.DebugLvl().Msgf(format, args...) }

// Info is a legacy method to report a trace lvl entry
func (l *Logger) Info(msg string) { l.InfoLvl().Msg(msg) }

// Infof is a legacy method to report a trace lvl entry
func (l *Logger) Infof(format string, args ...interface{}) { l.InfoLvl().Msgf(format, args...) }

// Warn is a legacy method to report a trace lvl entry
func (l *Logger) Warn(msg string) { l.WarnLvl().Msg(msg) }

// Warnf is a legacy method to report a trace lvl entry
func (l *Logger) Warnf(format string, args ...interface{}) { l.WarnLvl().Msgf(format, args...) }

// Error is a legacy method to report a trace lvl entry
func (l *Logger) Error(msg string) { l.ErrorLvl().Msg(msg) }

// Errorf is a legacy method to report a trace lvl entry
func (l *Logger) Errorf(format string, args ...interface{}) { l.ErrorLvl().Msgf(format, args...) }

func (l *Logger) newEvent(lvl LogLevel) *Event {
	if l.Lvl < lvl {
		return &Event{formatter: &NoopFormatter{}}
	}

	return &Event{formatter: l.FormatterFactory(lvl)}
}

// Bool ...
func (e *Event) Bool(key string, b bool) *Event {
	e.formatter.Bool(key, b)
	return e
}

// Err ...
func (e *Event) Err(err error) *Event {
	e.formatter.Err(err)
	return e
}

// Float32 ...
func (e *Event) Float32(key string, f float32) *Event {
	e.formatter.Float32(key, f)
	return e
}

// Float64 ...
func (e *Event) Float64(key string, f float64) *Event {
	e.formatter.Float64(key, f)
	return e
}

// IPAddr ...
func (e *Event) IPAddr(key string, ip net.IP) *Event {
	e.formatter.IPAddr(key, ip)
	return e
}

// Int ...
func (e *Event) Int(key string, i int) *Event {
	e.formatter.Int(key, i)
	return e
}

// Int16 ...
func (e *Event) Int16(key string, i int16) *Event {
	e.formatter.Int16(key, i)
	return e
}

// Int32 ...
func (e *Event) Int32(key string, i int32) *Event {
	e.formatter.Int32(key, i)
	return e
}

// Int64 ...
func (e *Event) Int64(key string, i int64) *Event {
	e.formatter.Int64(key, i)
	return e
}

// Int8 ...
func (e *Event) Int8(key string, i int8) *Event {
	e.formatter.Int8(key, i)
	return e
}

// Str ...
func (e *Event) Str(key, val string) *Event {
	e.formatter.Str(key, val)
	return e
}

// Uint ...
func (e *Event) Uint(key string, i uint) *Event {
	e.formatter.Uint(key, i)
	return e
}

// Uint16 ...
func (e *Event) Uint16(key string, i uint16) *Event {
	e.formatter.Uint16(key, i)
	return e
}

// Uint32 ...
func (e *Event) Uint32(key string, i uint32) *Event {
	e.formatter.Uint32(key, i)
	return e
}

// Uint64 ...
func (e *Event) Uint64(key string, i uint64) *Event {
	e.formatter.Uint64(key, i)
	return e
}

// Uint8 ...
func (e *Event) Uint8(key string, i uint8) *Event {
	e.formatter.Uint8(key, i)
	return e
}

// Msg writes the event to the writer
func (e *Event) Msg(message string) {
	e.formatter.Msg(message)
}

// Msgf writes the event to the writer
func (e *Event) Msgf(format string, args ...interface{}) {
	e.formatter.Msgf(format, args...)
}

// DefaultLoggerFactory define levels by scopes and creates new Loggers
type DefaultLoggerFactory struct {
	Writer          io.Writer
	DefaultLogLevel LogLevel
	ScopeLevels     map[string]LogLevel
}

// NewDefaultLoggerFactory creates a new DefaultLoggerFactory
func NewDefaultLoggerFactory() *DefaultLoggerFactory {
	factory := DefaultLoggerFactory{}
	factory.DefaultLogLevel = LogLevelError
	factory.ScopeLevels = make(map[string]LogLevel)

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
			factory.DefaultLogLevel = level
			continue
		}

		scopes := strings.Split(strings.ToLower(env), ",")
		for _, scope := range scopes {
			factory.ScopeLevels[scope] = level
		}
	}

	return &factory
}

// NewDefaultLeveledLoggerForScope returns a configured Logger
func NewDefaultLeveledLoggerForScope(scope string, level LogLevel, writer io.Writer) *Logger {
	if writer == nil {
		writer = os.Stdout
	}
	return &Logger{
		FormatterFactory: func(lvl LogLevel) Formatter {
			return NewStringFormatter(writer, lvl)
		},
		Lvl: level,
	}
}

// NewLogger returns a configured Logger for the given , argsscope
func (f *DefaultLoggerFactory) NewLogger(scope string) LeveledLogger {
	logLevel := f.DefaultLogLevel
	if f.ScopeLevels != nil {
		scopeLevel, found := f.ScopeLevels[scope]

		if found {
			logLevel = scopeLevel
		}
	}

	return NewDefaultLeveledLoggerForScope(scope, logLevel, f.Writer)
}
