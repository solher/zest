package interfaces

import (
	"net/http"

	"github.com/Solher/auth-scaffold/apierrors"
	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/internalerrors"
	"github.com/Solher/auth-scaffold/usecases"
	"github.com/dimfeld/httptreemux"
	"github.com/gorilla/context"
)

type PermissionGate struct {
	next        *httptreemux.HandlerFunc
	routes      *RouteDirectory
	permissions usecases.PermissionDirectory
	render      AbstractRender
}

func NewPermissionGate(next *httptreemux.HandlerFunc, routes *RouteDirectory, permissions usecases.PermissionDirectory, render AbstractRender) *PermissionGate {
	return &PermissionGate{next: next, routes: routes, permissions: permissions, render: render}
}

func (c *PermissionGate) Handler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	sessionCtx := context.Get(r, "currentSession")
	var role string

	if sessionCtx == nil {
		role = "guest"
	} else {
		session := sessionCtx.(domain.Session)

		if session.Account.IsAdmin {
			role = "admin"
		} else {
			role = "authenticated"
		}
	}

	if !c.permissions[role].IsGranted(c.next) {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, internalerrors.InsufficentPermissions)
		return
	}

	(*c.next)(w, r, params)
}
