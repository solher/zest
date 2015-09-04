package zest

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (rec *Recovery) PanicServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec.ServeHTTP(w, r, func(_ http.ResponseWriter, _ *http.Request) { panic("panic") })
}

func (rec *Recovery) WrappedServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec.ServeHTTP(w, r, func(_ http.ResponseWriter, _ *http.Request) {})
}

// TestRecovery runs tests on the recovery middleware.
func TestRecovery(t *testing.T) {
	a := assert.New(t)
	// r := require.New(t)
	recovery := NewRecovery()
	recovery.Logger = nil

	a.HTTPBodyContains(recovery.PanicServeHTTP, "GET", "/", url.Values{}, "500")
	a.HTTPBodyContains(recovery.PanicServeHTTP, "GET", "/", url.Values{}, "An internal error occured. Please retry later.")
	a.HTTPBodyContains(recovery.PanicServeHTTP, "GET", "/", url.Values{}, "undefined error")
	a.HTTPBodyContains(recovery.PanicServeHTTP, "GET", "/", url.Values{}, "INTERNAL_SERVER_ERROR")

	a.HTTPBodyNotContains(recovery.WrappedServeHTTP, "GET", "/", url.Values{}, "undefined error")
}
