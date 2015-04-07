package server

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
)

// LoggerMiddleware is a middleware handler that logs the request as it goes in and the response as it goes out.
type LoggerMiddleware struct {
	// Logger is the log.Logger instance used to log messages with the Logger middleware
	Logger *logrus.Logger
}

// NewMiddleware returns a new *LoggerMiddleware, yay!
func NewLoggerMiddleware() *LoggerMiddleware {
	return NewCustomMiddleware(logrus.InfoLevel, &logrus.TextFormatter{})
}

// NewCustomMiddleware builds a *LoggerMiddleware with the given level and formatter
func NewCustomMiddleware(level logrus.Level, formatter logrus.Formatter) *LoggerMiddleware {
	log := logrus.New()
	log.Level = level
	log.Formatter = formatter

	return &LoggerMiddleware{Logger: log}
}

func (l *LoggerMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	l.Logger.WithFields(logrus.Fields{
		"method":    r.Method,
		"request":   r.RequestURI,
		"remote":    r.RemoteAddr,
		"requestID": r.Header.Get("X-Request-Id"),
	}).Info("started handling request")

	next(rw, r)

	latency := time.Since(start)
	res := rw.(negroni.ResponseWriter)
	l.Logger.WithFields(logrus.Fields{
		"status":      res.Status(),
		"method":      r.Method,
		"request":     r.RequestURI,
		"remote":      r.RemoteAddr,
		"text_status": http.StatusText(res.Status()),
		"took":        latency,
		"requestID":   r.Header.Get("X-Request-Id"),
		"latency":     latency.Nanoseconds(),
	}).Info("completed handling request")
}
