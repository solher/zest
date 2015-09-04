package zest

import (
	"net/http"

	"github.com/unrolled/render"
)

type Render struct {
	renderer *render.Render
}

func NewRender() *Render {
	return &Render{renderer: render.New()}
}

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

func (r *Render) JSON(w http.ResponseWriter, status int, object interface{}) {
	if object == nil {
		w.WriteHeader(status)
	} else {
		r.renderer.JSON(w, status, object)
	}
}
