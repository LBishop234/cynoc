package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type LogLevel string

var (
	TRACE LogLevel = "trace"
	DEBUG LogLevel = "debug"
	INFO  LogLevel = "info"
	WARN  LogLevel = "warn"
	ERROR LogLevel = "error"
	FATAL LogLevel = "fatal"
)

var Log zerolog.Logger = zerolog.Nop()

func InitLogger(logLevel LogLevel) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	Log = zerolog.New(os.Stdout).
		With().Timestamp().Logger().
		Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC822,
		})

	var zerologLevel zerolog.Level
	switch logLevel {
	case TRACE:
		zerologLevel = zerolog.TraceLevel
		Log = Log.With().Caller().Logger()

	case DEBUG:
		zerologLevel = zerolog.DebugLevel
		Log = Log.With().Caller().Logger()

	case INFO:
		zerologLevel = zerolog.InfoLevel

	case WARN:
		zerologLevel = zerolog.WarnLevel

	case ERROR:
		zerologLevel = zerolog.ErrorLevel

	case FATAL:
		zerologLevel = zerolog.FatalLevel

	default:
		zerologLevel = zerolog.InfoLevel
	}

	Log = Log.Level(zerologLevel)

	Log.Info().Str("log_level", Log.GetLevel().String()).Msg("initialised logger pkg")
}
