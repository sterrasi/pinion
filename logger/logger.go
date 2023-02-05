package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"strings"
	"time"
)

var rootLogger zerolog.Logger

// Config for Logging
type Config struct {
	UseUnstructuredLog bool           `ini:",omitempty"`
	LogLevel           *zerolog.Level `ini:",omitempty"`
}

// Trace log event
func Trace() *zerolog.Event {
	return rootLogger.Trace()
}

// Debug log event
func Debug() *zerolog.Event {
	return rootLogger.Debug()
}

// Info log event
func Info() *zerolog.Event {
	return rootLogger.Info()
}

// Warn log event
func Warn() *zerolog.Event {
	return rootLogger.Warn()
}

// Error log event
func Error() *zerolog.Event {
	return rootLogger.Error()
}

// ConfigureLogging will configure the root logger with the given log level and unstructured flag.
// TODO: add the ability to create and return other loggers for in order to turn on/off logging of specific app features
func ConfigureLogging(level zerolog.Level, unstructured bool) {

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(level)

	var writer io.Writer

	if unstructured {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.UnixDate}
		consoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		consoleWriter.FormatMessage = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("msg:'%s'", i))
		}
		consoleWriter.FormatFieldName = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s:", i))
		}
		consoleWriter.FormatFieldValue = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		}
		writer = consoleWriter
	} else {
		writer = os.Stdout
	}

	rootLogger = zerolog.New(writer).With().Timestamp().Logger()
}
