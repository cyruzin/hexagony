package middleware

import (
	"hexagony/lib/clog"
	"net/http"
	"strings"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var sb strings.Builder

		sb.WriteString(
			"method: " + r.Method + ", " + "url: " + r.URL.String() + ", " +
				"agent: " + r.UserAgent() + ", " + "referer: " + r.Referer() + ", " +
				"proto: " + r.Proto + ", " + "remote_address: " + r.RemoteAddr + ", " +
				"latency: " + time.Since(start).String(),
		)

		clog.Info(sb.String())

		next.ServeHTTP(w, r)
	})
}
