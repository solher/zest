package interfaces

import (
	"net/http"

	"github.com/Solher/auth-scaffold/apierrors"
	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/internalerrors"
	"github.com/dimfeld/httptreemux"
	"github.com/gorilla/context"
)

type PermissionGate struct {
	accountInter AbstractAccountInter
	next         *httptreemux.HandlerFunc
	routes       *RouteDirectory
	render       AbstractRender
}

func NewPermissionGate(accountInter AbstractAccountInter, next *httptreemux.HandlerFunc,
	routes *RouteDirectory, render AbstractRender) *PermissionGate {
	return &PermissionGate{accountInter: accountInter, next: next, routes: routes, render: render}
}

func (c *PermissionGate) Handler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	sessionCtx := context.Get(r, "currentSession")

	accountID := 0

	if sessionCtx != nil {
		session := sessionCtx.(domain.Session)
		accountID = session.AccountID
	}

	dirKey, err := c.routes.GetKey(c.next)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	roleNames, err := c.accountInter.GetGrantedRoles(accountID, dirKey.Ressource, dirKey.Method)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	if len(roleNames) == 0 {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, internalerrors.InsufficentPermissions)
		return
	}

	if len(roleNames) == 1 && roleNames[0] == "Owner" {
		relations := domain.ModelDirectory.FindPathToOwner(dirKey.Ressource)
		context.Set(r, "ownerRelations", relations)
	}

	(*c.next)(w, r, params)
}
