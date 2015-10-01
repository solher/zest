package zest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLogger runs tests on the logger middleware.
func TestLogger(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	logger := NewLogger()
	buff := bytes.NewBufferString("")
	logger.SetOutput(buff)

	n := negroni.New()
	n.Use(logger)

	codes := []int{100, 200, 300, 400, 500}
	for _, code := range codes {
		n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(code)
		}))

		req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
		a.NoError(err)

		n.ServeHTTP(httptest.NewRecorder(), req)
		r.NotEqual(0, buff.Len())
		a.Contains(buff.String(), "GET")
		a.Contains(buff.String(), "/foobar")
	}

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "OTHER"}
	for _, method := range methods {
		req, err := http.NewRequest(method, "http://localhost:3000/foobar", nil)
		a.NoError(err)

		n.ServeHTTP(httptest.NewRecorder(), req)
		r.NotEqual(0, buff.Len())
		a.Contains(buff.String(), method)
		a.Contains(buff.String(), "/foobar")
	}
}
