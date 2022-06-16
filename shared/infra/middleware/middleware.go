package middleware

import (
	"fmt"
	"hexagony/lib/clog"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		clog.Info(
			fmt.Sprintf(
				"method: %s, url: %s, agent: %s, referer: %s, proto: %s, remote_address: %s, latency: %s",
				r.Method, r.URL.String(), r.UserAgent(), r.Referer(), r.Proto, r.RemoteAddr, time.Since(start),
			))

		next.ServeHTTP(w, r)
	})
}
