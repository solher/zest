package zest

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubController struct {
	r      *Render
	status int
	obj    interface{}
}

type stubErrController struct {
	r      *Render
	status int
	apiErr *APIError
	err    error
}

func (c *stubController) Handler(w http.ResponseWriter, r *http.Request) {
	c.r.JSON(w, c.status, c.obj)
}

func (c *stubErrController) Handler(w http.ResponseWriter, r *http.Request) {
	c.r.JSONError(w, c.status, c.apiErr, c.err)
}

// TestJSON runs tests on the render JSON method.
func TestJSON(t *testing.T) {
	a := assert.New(t)
	// r := require.New(t)
	render := NewRender()

	str := "foobar"

	c := &stubController{r: render, status: http.StatusOK, obj: str}
	a.HTTPBodyContains(c.Handler, "GET", "/", url.Values{}, "foobar")

	c.obj = &str
	a.HTTPBodyContains(c.Handler, "GET", "/", url.Values{}, "foobar")

	c.obj = nil
	a.HTTPBodyNotContains(c.Handler, "GET", "/", url.Values{}, "foobar")
}

// TestJSONError runs tests on the render JSONError method.
func TestJSONError(t *testing.T) {
	a := assert.New(t)
	// r := require.New(t)
	render := NewRender()

	apiErr := &APIError{Description: "An error occured.", ErrorCode: "ERROR"}
	c := &stubErrController{r: render, status: http.StatusInternalServerError, apiErr: apiErr, err: errors.New("undefined error")}
	a.HTTPBodyContains(c.Handler, "GET", "/", url.Values{}, "500")
	a.HTTPBodyContains(c.Handler, "GET", "/", url.Values{}, "An error occured.")
	a.HTTPBodyContains(c.Handler, "GET", "/", url.Values{}, "undefined error")
	a.HTTPBodyContains(c.Handler, "GET", "/", url.Values{}, "ERROR")

	c.apiErr = nil
	a.HTTPBodyContains(c.Handler, "GET", "/", url.Values{}, "500")
	a.HTTPBodyNotContains(c.Handler, "GET", "/", url.Values{}, "An error occured.")
	a.HTTPBodyContains(c.Handler, "GET", "/", url.Values{}, "undefined error")
	a.HTTPBodyNotContains(c.Handler, "GET", "/", url.Values{}, "ERROR")

	c.err = nil
	a.HTTPBodyContains(c.Handler, "GET", "/", url.Values{}, "500")
	a.HTTPBodyNotContains(c.Handler, "GET", "/", url.Values{}, "An error occured.")
	a.HTTPBodyNotContains(c.Handler, "GET", "/", url.Values{}, "undefined error")
	a.HTTPBodyNotContains(c.Handler, "GET", "/", url.Values{}, "ERROR")
}
