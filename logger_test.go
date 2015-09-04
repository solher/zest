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
	n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(httptest.NewRecorder(), req)
	r.NotEqual(buff.Len(), 0)
	a.Contains(buff.String(), "GET")
	a.Contains(buff.String(), "/foobar")
}
