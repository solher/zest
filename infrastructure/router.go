package infrastructure

import (
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/julienschmidt/httprouter"
)

type Router struct {
	router *httprouter.Router
}

func NewRouter() *Router {
	return &Router{router: httprouter.New()}
}

func (r *Router) Handle(method string, path string, handle interfaces.Handle) {
	r.router.Handle(method, path, httprouter.Handle(handle))
}
