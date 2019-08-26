package logging_test

import (
	"strings"
	"testing"

	"github.com/pion/logging"
)

func testNoDebugLevel(t *testing.T, logger *logging.Logger, builder *strings.Builder) {
	builder.Reset()
	logger.DebugLvl().Msg("this shouldn't be logged")
	if builder.Len() > 0 {
		t.Error("Debug was logged when it shouldn't have been")
	}
	logger.DebugLvl().Msgf("this shouldn't be logged")
	if builder.Len() > 0 {
		t.Error("Debug was logged when it shouldn't have been")
	}
}

func testDebugLevel(t *testing.T, logger *logging.Logger, builder *strings.Builder) {
	builder.Reset()
	dbgMsg := "this is a debug message"
	logger.DebugLvl().Msg(dbgMsg)
	if !strings.Contains(builder.String(), dbgMsg) {
		t.Errorf("Expected to find %q in %q, but didn't", dbgMsg, builder.String())
	}
	logger.DebugLvl().Msgf(dbgMsg)
	if !strings.Contains(builder.String(), dbgMsg) {
		t.Errorf("Expected to find %q in %q, but didn't", dbgMsg, builder.String())
	}
}

func testWarnLevel(t *testing.T, logger *logging.Logger, builder *strings.Builder) {
	builder.Reset()
	warnMsg := "this is a warning message"
	logger.WarnLvl().Msg(warnMsg)
	if !strings.Contains(builder.String(), warnMsg) {
		t.Errorf("Expected to find %q in %q, but didn't", warnMsg, builder.String())
	}
	logger.WarnLvl().Msgf(warnMsg)
	if !strings.Contains(builder.String(), warnMsg) {
		t.Errorf("Expected to find %q in %q, but didn't", warnMsg, builder.String())
	}
}

func testErrorLevel(t *testing.T, logger *logging.Logger, builder *strings.Builder) {
	builder.Reset()
	errMsg := "this is an error message"
	logger.ErrorLvl().Msg(errMsg)
	if !strings.Contains(builder.String(), errMsg) {
		t.Errorf("Expected to find %q in %q, but didn't", errMsg, builder.String())
	}
	logger.ErrorLvl().Msgf(errMsg)
	if !strings.Contains(builder.String(), errMsg) {
		t.Errorf("Expected to find %q in %q, but didn't", errMsg, builder.String())
	}
}

func TestDefaultLoggerFactory(t *testing.T) {
	writer := &strings.Builder{}

	f := logging.DefaultLoggerFactory{
		Writer:          writer,
		DefaultLogLevel: logging.LogLevelWarn,
		ScopeLevels: map[string]logging.LogLevel{
			"foo": logging.LogLevelDebug,
		},
	}

	bazLogger := f.NewLogger("baz")
	testNoDebugLevel(t, bazLogger, writer)
	testWarnLevel(t, bazLogger, writer)

	fooLogger := f.NewLogger("foo")
	testDebugLevel(t, fooLogger, writer)
}

func TestDefaultLogger(t *testing.T) {
	writer := &strings.Builder{}
	logger := logging.
		NewDefaultLeveledLoggerForScope("test1", logging.LogLevelWarn, writer)

	testNoDebugLevel(t, logger, writer)
	testWarnLevel(t, logger, writer)
	testErrorLevel(t, logger, writer)
}
