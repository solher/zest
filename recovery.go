package zest

import (
	"errors"
	"log"
	"net/http"
	"os"
	"runtime"
)

type Recovery struct {
	Logger    *log.Logger
	Render    *Render
	StackAll  bool
	StackSize int
}

func NewRecovery() *Recovery {
	return &Recovery{
		Logger:    log.New(os.Stdout, "", 0),
		Render:    NewRender(),
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

			err := &APIError{Description: "An internal error occured. Please retry later.", ErrorCode: "INTERNAL_SERVER_ERROR"}

			rec.Render.JSONError(w, http.StatusInternalServerError, err, errors.New("Undefined error."))
		}
	}()

	next(w, r)
}
