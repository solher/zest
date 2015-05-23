package interfaces

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
)

type PermissionGate struct {
	next   *httptreemux.HandlerFunc
	routes *RouteDirectory
	render AbstractRender
}

func NewPermissionGate(next *httptreemux.HandlerFunc, routes *RouteDirectory, render AbstractRender) *PermissionGate {
	return &PermissionGate{next: next, routes: routes, render: render}
}

func (c *PermissionGate) Handler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// sessionCtx := context.Get(r, "currentSession")
	// var role string
	//
	// if sessionCtx == nil {
	// 	role = "guest"
	// } else {
	// 	session := sessionCtx.(domain.Session)
	//
	// 	if session.Account.IsAdmin {
	// 		role = "admin"
	// 	} else {
	// 		role = "authenticated"
	// 	}
	// }
	//
	// if !c.permissions[role].IsGranted(c.next) {
	// 	c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, internalerrors.InsufficentPermissions)
	// 	return
	// }

	(*c.next)(w, r, params)
}
