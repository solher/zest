package middlewares

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/solher/zest/apierrors"
	"github.com/solher/zest/internalerrors"
)

type AbstractRender interface {
	JSONError(w http.ResponseWriter, status int, apiError *apierrors.APIError, err error)
	JSON(w http.ResponseWriter, status int, object interface{})
}

type Recovery struct {
	Logger    *log.Logger
	Render    AbstractRender
	StackAll  bool
	StackSize int
}

func NewRecovery(render AbstractRender) *Recovery {
	return &Recovery{
		Logger:    log.New(os.Stdout, "", 0),
		Render:    render,
		StackAll:  false,
		StackSize: 1024 * 8,
	}
}

func (rec *Recovery) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]

			f := "PANIC: %s\n%s"
			rec.Logger.Printf(f, err, stack)

			rec.Render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, internalerrors.Undefined)
		}
	}()

	next(w, r)
}
