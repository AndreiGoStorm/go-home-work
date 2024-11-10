package internalhttp

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(next http.Handler, logg *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logg.Error(err.Error())
		}
		logg.Info(buildLoggingString(r, ip, start))
		next.ServeHTTP(w, r)
	})
}

func buildLoggingString(r *http.Request, ip string, start time.Time) string {
	var builder strings.Builder
	builder.WriteString("ip - ")
	builder.WriteString(ip)
	builder.WriteString(", method - ")
	builder.WriteString(r.Method)
	builder.WriteString(", url - ")
	builder.WriteString(r.URL.String())
	builder.WriteString(", ")
	builder.WriteString(r.Proto)
	builder.WriteString(", user-agent - ")
	builder.WriteString(r.Header.Get("user-agent"))
	builder.WriteString(", latency - ")
	builder.WriteString(time.Since(start).String())

	return builder.String()
}
