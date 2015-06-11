package infrastructure

import (
	"net/http"

	"github.com/solher/zest/apierrors"
	"github.com/unrolled/render"
)

type Render struct {
	renderer *render.Render
}

func NewRender() *Render {
	return &Render{renderer: render.New()}
}

func (r *Render) JSONError(w http.ResponseWriter, status int, apiError *apierrors.APIError, err error) {
	r.renderer.JSON(w, status, apierrors.Make(*apiError, status, err))
}

func (r *Render) JSON(w http.ResponseWriter, status int, object interface{}) {
	if object == nil {
		w.WriteHeader(status)
	} else {
		r.renderer.JSON(w, status, object)
	}
}
