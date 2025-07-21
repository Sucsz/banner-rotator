package middleware

import (
	"net/http"
	"time"

	"github.com/Sucsz/banner-rotator/internal/log"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// logEntry реализует chimw.LogEntry
type logEntry struct {
	logger zerolog.Logger
}

func (e *logEntry) Write(
	status, bytes int,
	header http.Header,
	elapsed time.Duration,
	_ interface{}, // extra не используется
) {
	e.logger.Info().
		Int("status", status).
		Int("bytes", bytes).
		Dur("latency", elapsed).
		Msg("access")
}

func (e *logEntry) Panic(v interface{}, stack []byte) {
	e.logger.Error().
		Interface("panic", v).
		Bytes("stack", stack).
		Msg("panic in handler")
}

type formatter struct{}

func (f *formatter) NewLogEntry(r *http.Request) chimw.LogEntry {
	requestID := chimw.GetReqID(r.Context())

	logger := log.WithComponent("http").With().
		Str("request_id", requestID).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remote_addr", r.RemoteAddr).
		Str("user_agent", r.UserAgent()).
		Str("referer", r.Referer()).
		Logger()

	return &logEntry{logger: logger}
}

// RequestLogger — middleware для structured access‑логов через zerolog
func RequestLogger(next http.Handler) http.Handler {
	return chimw.RequestLogger(&formatter{})(next)
}
