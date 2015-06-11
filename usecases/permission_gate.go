package usecases

import (
	"net/http"

	"github.com/solher/zest/apierrors"
	"github.com/solher/zest/domain"
	"github.com/solher/zest/internalerrors"
	"github.com/gorilla/context"
)

type PermissionGate struct {
	ressource, method string
	accountInter      AbstractAccountInter
	render            AbstractRender
	next              HandlerFunc
}

func NewPermissionGate(ressource, method string, accountInter AbstractAccountInter, render AbstractRender,
	next HandlerFunc) *PermissionGate {
	return &PermissionGate{ressource: ressource, method: method, accountInter: accountInter, render: render, next: next}
}

func (p *PermissionGate) Handler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	sessionCtx := context.Get(r, "currentSession")

	accountID := 0

	if sessionCtx != nil {
		session := sessionCtx.(domain.Session)
		accountID = session.AccountID
	}

	roleNames, err := p.accountInter.GetGrantedRoles(accountID, p.ressource, p.method)
	if err != nil {
		p.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	if len(roleNames) == 0 {
		p.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, internalerrors.InsufficentPermissions)
		return
	}

	if len(roleNames) == 1 && roleNames[0] == "Owner" {
		relations := domain.ModelDirectory.FindPathToOwner(p.ressource)
		context.Set(r, "ownerRelations", relations)
	}

	p.next(w, r, params)
}
