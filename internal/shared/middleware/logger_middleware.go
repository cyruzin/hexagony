package middleware

import (
	"hexagony/pkg/clog"
	"net/http"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		clog.Custom(map[string]interface{}{
			"host":      r.Host,
			"method":    r.Method,
			"url":       r.URL.String(),
			"agent":     r.UserAgent(),
			"referer":   r.Referer(),
			"proto":     r.Proto,
			"remote_ip": r.RemoteAddr,
		})

		next.ServeHTTP(w, r)
	})
}
