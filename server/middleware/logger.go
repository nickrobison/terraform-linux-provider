package middleware

import (
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/pkgerrors"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func SetupLogging(writer io.Writer, level zerolog.Level) {
	zerolog.ErrorMarshalFunc = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var output io.Writer = zerolog.ConsoleWriter{
		Out:        writer,
		TimeFormat: time.UnixDate,
	}

	log = zerolog.New(output).
		Level(level).
		With().Timestamp().Logger()
}

// Logger returns a logger to use
// Make sure SetupLogging is called first
func Logger() zerolog.Logger {
	return log
}

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
