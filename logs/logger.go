package logs

import (
	"io/fs"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger = zerolog.Nop()

var osOpenFile func(name string, flag int, perm fs.FileMode) (*os.File, error) = os.OpenFile

func PrepareLogger(level string) {
	// initialize zerolog config
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	// wrap format and build logger
	logger = zerolog.New(os.Stderr).Level(parseLogLevel(level)).With().Timestamp().Logger()
}

func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
