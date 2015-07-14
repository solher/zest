// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "\n[Zest] ", 0)}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	path := r.URL.Path
	// Start timer
	start := time.Now()

	// Process request
	next(rw, r)

	if strings.Contains(path, "explorer") || strings.Contains(path, "favicon") {
		return
	}

	// Stop timer
	end := time.Now()
	latency := end.Sub(start)

	clientIP := r.RemoteAddr
	method := r.Method
	res := rw.(negroni.ResponseWriter)
	statusCode := res.Status()
	statusColor := colorForStatus(statusCode)
	methodColor := colorForMethod(method)

	l.Printf("%v | %s %3d %s | %v | %s | %s %s %s %s",
		end.Format("2006/01/02 - 15:04:05"),
		statusColor, statusCode, reset,
		latency,
		clientIP,
		methodColor, method, reset,
		path,
	)
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}
