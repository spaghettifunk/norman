package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// LogRecord warps a http.ResponseWriter and records the status
type LogRecord struct {
	http.ResponseWriter
	status int
}

func (r *LogRecord) Write(p []byte) (int, error) {
	return r.ResponseWriter.Write(p)
}

// WriteHeader overrides ResponseWriter.WriteHeader to keep track of the response code
func (r *LogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		record := &LogRecord{
			ResponseWriter: rw,
			status:         200,
		}
		h.ServeHTTP(record, r) // serve the original request

		duration := time.Since(start)

		// log request details
		log.Trace().Fields(map[string]interface{}{
			"ident":       r.Host,
			"method":      r.Method,
			"status_code": record.status,
			"url":         r.URL.Path,
			"duration":    duration,
		}).Send()
	}
	return http.HandlerFunc(logFn)
}
