package interfaces

import (
	"net/http"

	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/utils"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

type PermissionGate struct {
	Next httprouter.Handle
}

func NewPermissionGate(next httprouter.Handle) *PermissionGate {
	return &PermissionGate{Next: next}
}

func (c *PermissionGate) Handler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	sessionCtx := context.Get(r, "currentSession")

	if sessionCtx == nil {
		utils.Dump("ROLE: GUEST")
	} else {
		session := sessionCtx.(domain.Session)

		if session.Account.IsAdmin {
			utils.Dump("ROLE: ADMIN")
		} else {
			utils.Dump("ROLE: AUTHENTICATED")
		}
	}

	c.Next(w, r, params)
}
