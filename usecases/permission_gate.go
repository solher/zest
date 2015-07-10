package usecases

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/solher/zest/apierrors"
	"github.com/solher/zest/domain"
	"github.com/solher/zest/internalerrors"
)

type PermissionGate struct {
	resource, method string
	accountInter     AbstractAccountInter
	render           AbstractRender
	next             HandlerFunc
}

func NewPermissionGate(resource, method string, accountInter AbstractAccountInter, render AbstractRender,
	next HandlerFunc) *PermissionGate {
	return &PermissionGate{resource: resource, method: method, accountInter: accountInter, render: render, next: next}
}

func (p *PermissionGate) Handler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	sessionCtx := context.Get(r, "currentSession")

	accountID := 0

	if sessionCtx != nil {
		session := sessionCtx.(domain.Session)
		accountID = session.AccountID
	}

	roleNames, err := p.accountInter.GetGrantedRoles(accountID, p.resource, p.method)
	if err != nil {
		p.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	if len(roleNames) == 0 {
		p.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, internalerrors.NotFound)
		return
	}

	if len(roleNames) == 1 && roleNames[0] == "Owner" {
		if context.Get(r, "lastResource") == nil && (p.method == "Create" || p.method == "Upsert") {
			p.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, internalerrors.NotFound)
			return
		}

		relations := domain.ModelDirectory.FindPathToOwner(p.resource)
		context.Set(r, "ownerRelations", relations)
	}

	context.Set(r, "method", p.method)
	context.Set(r, "resource", p.resource)

	p.next(w, r, params)
}
