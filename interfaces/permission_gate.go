package interfaces

import (
	"net/http"

	"github.com/Solher/auth-scaffold/apierrors"
	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/internalerrors"
	"github.com/Solher/auth-scaffold/utils"
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

	roleNames, err := c.accountInter.GetGrantedRoles(accountID, dirKey.Ressources, dirKey.Method)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	if len(roleNames) == 0 {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, internalerrors.InsufficentPermissions)
		return
	}

	utils.Dump(roleNames)

	//
	// if !c.permissions[role].IsGranted(c.next) {
	// 	c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, internalerrors.InsufficentPermissions)
	// 	return
	// }

	(*c.next)(w, r, params)
}
