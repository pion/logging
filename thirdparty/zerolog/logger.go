package zerolog

import (
	"net"

	"github.com/pion/logging"
	"github.com/rs/zerolog"
)

type ZerologFormatter struct {
	event *zerolog.Event
}

func (f *ZerologFormatter) Bool(key string, b bool) {
	f.event = f.event.Bool(key, b)
}

func (f *ZerologFormatter) Err(err error) {
	f.event = f.event.Err(err)
}

func (f *ZerologFormatter) Float32(key string, v float32) {
	f.event = f.event.Float32(key, v)
}

func (f *ZerologFormatter) Float64(key string, v float64) {
	f.event = f.event.Float64(key, v)
}

func (f *ZerologFormatter) IPAddr(key string, ip net.IP) {
	f.event = f.event.IPAddr(key, ip)
}

func (f *ZerologFormatter) Int(key string, i int) {
	f.event = f.event.Int(key, i)
}

func (f *ZerologFormatter) Int16(key string, i int16) {
	f.event = f.event.Int16(key, i)
}

func (f *ZerologFormatter) Int32(key string, i int32) {
	f.event = f.event.Int32(key, i)
}

func (f *ZerologFormatter) Int64(key string, i int64) {
	f.event = f.event.Int64(key, i)
}

func (f *ZerologFormatter) Int8(key string, i int8) {
	f.event = f.event.Int8(key, i)
}

func (f *ZerologFormatter) Str(key, val string) {
	f.event = f.event.Str(key, val)
}

func (f *ZerologFormatter) Uint(key string, i uint) {
	f.event = f.event.Uint(key, i)
}

func (f *ZerologFormatter) Uint16(key string, i uint16) {
	f.event = f.event.Uint16(key, i)
}

func (f *ZerologFormatter) Uint32(key string, i uint32) {
	f.event = f.event.Uint32(key, i)
}

func (f *ZerologFormatter) Uint64(key string, i uint64) {
	f.event = f.event.Uint64(key, i)
}

func (f *ZerologFormatter) Uint8(key string, i uint8) {
	f.event = f.event.Uint8(key, i)
}

func (f *ZerologFormatter) Msg(message string) {
	f.event.Msg(message)
	f.event = nil
}

func (f *ZerologFormatter) Msgf(format string, args ...interface{}) {
	f.event.Msgf(format, args...)
	f.event = nil
}

func NewZerologFormatter(logger zerolog.Logger, lvl logging.LogLevel) logging.Formatter {
	var event *zerolog.Event

	switch lvl {
	case logging.LogLevelError:
		event = logger.Error()
	case logging.LogLevelWarn:
		event = logger.Warn()
	case logging.LogLevelInfo:
		event = logger.Info()
	default:
		event = logger.Debug()
	}

	return &ZerologFormatter{event: event}
}
