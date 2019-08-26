package logging

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// StringFormatter writes a string with format key1=value1 key2=value2 .. into the writer
type StringFormatter struct {
	b      strings.Builder
	writer io.Writer
}

// NewStringFormatter creates a new string formatter
func NewStringFormatter(writer io.Writer, lvl LogLevel) Formatter {
	return &StringFormatter{b: strings.Builder{}, writer: writer}
}

// Bool ...
func (f *StringFormatter) Bool(key string, b bool) {
	var v string

	if b {
		v = "true"
	} else {
		v = "false"
	}

	f.Str(key, v)
}

// Err ...
func (f *StringFormatter) Err(err error) {
	f.Str("error", err.Error())
}

// Float32 ...
func (f *StringFormatter) Float32(key string, v float32) {
	f.Str(key, strconv.FormatFloat(float64(v), 'E', -1, 32))
}

// Float64 ...
func (f *StringFormatter) Float64(key string, v float64) {
	f.Str(key, strconv.FormatFloat(v, 'E', -1, 64))
}

// IPAddr ...
func (f *StringFormatter) IPAddr(key string, ip net.IP) {
	f.Str(key, ip.String())
}

// Int ...
func (f *StringFormatter) Int(key string, i int) {
	f.Str(key, strconv.FormatInt(int64(i), 10))
}

// Int16 ...
func (f *StringFormatter) Int16(key string, i int16) {
	f.Str(key, strconv.FormatInt(int64(i), 10))
}

// Int32 ...
func (f *StringFormatter) Int32(key string, i int32) {
	f.Str(key, strconv.FormatInt(int64(i), 10))
}

// Int64 ...
func (f *StringFormatter) Int64(key string, i int64) {
	f.Str(key, strconv.FormatInt(i, 10))
}

// Int8 ...
func (f *StringFormatter) Int8(key string, i int8) {
	f.Str(key, strconv.FormatInt(int64(i), 10))
}

// Str ...
func (f *StringFormatter) Str(key, val string) {
	if f.b.Len() > 0 {
		f.writeString(" ")
	}

	f.writeString(key)
	f.writeString("=")
	f.writeString(val)
}

// Uint ...
func (f *StringFormatter) Uint(key string, i uint) {
	f.Str(key, strconv.FormatUint(uint64(i), 10))
}

// Uint16 ...
func (f *StringFormatter) Uint16(key string, i uint16) {
	f.Str(key, strconv.FormatUint(uint64(i), 10))
}

// Uint32 ...
func (f *StringFormatter) Uint32(key string, i uint32) {
	f.Str(key, strconv.FormatUint(uint64(i), 10))
}

// Uint64 ...
func (f *StringFormatter) Uint64(key string, i uint64) {
	f.Str(key, strconv.FormatUint(i, 10))
}

// Uint8 ...
func (f *StringFormatter) Uint8(key string, i uint8) {
	f.Str(key, strconv.FormatUint(uint64(i), 10))
}

// Msg ...
func (f *StringFormatter) Msg(message string) {
	if f.b.Len() > 0 {
		f.writeString(" ")
	}

	f.writeString(message)
	f.writeString("\n")

	_, err := f.writer.Write([]byte(f.b.String()))
	if err != nil {
		fmt.Printf("error writing log %s\n", err.Error())
	}
}

// Msgf ...
func (f *StringFormatter) Msgf(format string, args ...interface{}) {
	if f.b.Len() > 0 {
		f.writeString(" ")
	}

	f.writeString(fmt.Sprintf(format, args...))
	f.writeString("\n")

	_, err := f.writer.Write([]byte(f.b.String()))
	if err != nil {
		fmt.Printf("error writing log %s\n", err.Error())
	}
}

func (f *StringFormatter) writeString(s string) {
	if _, err := f.b.WriteString(s); err != nil {
		fmt.Printf("error writing string to builder %s\n", err.Error())
	}
}

// NoopFormatter is a no-op formatter
type NoopFormatter struct{}

// Bool ...
func (f *NoopFormatter) Bool(key string, b bool) {}

// Err ...
func (f *NoopFormatter) Err(err error) {}

// Float32 ...
func (f *NoopFormatter) Float32(key string, v float32) {}

// Float64 ...
func (f *NoopFormatter) Float64(key string, v float64) {}

// IPAddr ...
func (f *NoopFormatter) IPAddr(key string, ip net.IP) {}

// Int ...
func (f *NoopFormatter) Int(key string, i int) {}

// Int16 ...
func (f *NoopFormatter) Int16(key string, i int16) {}

// Int32 ...
func (f *NoopFormatter) Int32(key string, i int32) {}

// Int64 ...
func (f *NoopFormatter) Int64(key string, i int64) {}

// Int8 ...
func (f *NoopFormatter) Int8(key string, i int8) {}

// Str ...
func (f *NoopFormatter) Str(key, val string) {}

// Uint ...
func (f *NoopFormatter) Uint(key string, i uint) {}

// Uint16 ...
func (f *NoopFormatter) Uint16(key string, i uint16) {}

// Uint32 ...
func (f *NoopFormatter) Uint32(key string, i uint32) {}

// Uint64 ...
func (f *NoopFormatter) Uint64(key string, i uint64) {}

// Uint8 ...
func (f *NoopFormatter) Uint8(key string, i uint8) {}

// Msg ...
func (f *NoopFormatter) Msg(message string) {}

// Msgf ...
func (f *NoopFormatter) Msgf(format string, args ...interface{}) {}
