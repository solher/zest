package zest

import (
	"net/http"

	"github.com/unrolled/render"
)

// Render is a unrolled/render based JSON/XML/HTML renderer, customized to increase
// the expressiveness of Zest API error rendering.
type Render struct {
	renderer *render.Render
}

// NewRender returns a new instance of Render.
func NewRender() *Render {
	return &Render{renderer: render.New()}
}

// JSONError forges and writes an APIError into the response writer.
func (r *Render) JSONError(w http.ResponseWriter, status int, apiError *APIError, err error) {
	if apiError == nil {
		apiError = &APIError{}
	}

	if err != nil {
		apiError.Raw = err.Error()
	}

	apiError.Status = status

	r.renderer.JSON(w, status, apiError)
}

// JSON writes the argument object into the response writer.
func (r *Render) JSON(w http.ResponseWriter, status int, object interface{}) {
	if object == nil {
		w.WriteHeader(status)
	} else {
		r.renderer.JSON(w, status, object)
	}
}
