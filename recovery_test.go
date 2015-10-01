package zest

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	emptyFunc http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {}
	panicFunc http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) { panic("panic") }
)

// TestRecovery runs tests on the recovery middleware.
func TestRecovery(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	rec := NewRecovery()
	rec.Logger = nil

	recorder := httptest.NewRecorder()

	r.NotPanics(func() {
		rec.ServeHTTP(recorder, (*http.Request)(nil), emptyFunc)
	})

	a.NotContains(recorder.Body.String(), "undefined error")

	buff := bytes.NewBufferString("")
	rec.Logger = log.New(buff, "[Zest] ", 0)
	rec.ServeHTTP(recorder, (*http.Request)(nil), panicFunc)

	r.NotEqual(0, buff.Len())
	a.Contains(recorder.Body.String(), "500")
	a.Contains(recorder.Body.String(), "An internal error occured. Please retry later.")
	a.Contains(recorder.Body.String(), "undefined error")
	a.Contains(recorder.Body.String(), "UNDEFINED_ERROR")
}
