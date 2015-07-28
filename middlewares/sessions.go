package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/solher/zest/domain"
)

type AbstractAccountInter interface {
	CurrentSessionFromToken(authToken string) (*domain.Session, error)
}

type Sessions struct {
	interactor AbstractAccountInter
}

func NewSessions(accountInteractor AbstractAccountInter) *Sessions {
	return &Sessions{interactor: accountInteractor}
}

func (s *Sessions) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authorizations := strings.Split(r.Header.Get("Authorization"), ",")
	var authToken string

	for _, v := range authorizations {
		fmt.Sscanf(strings.Trim(v, " "), "AuthToken authToken=%s", &authToken)
		if authToken != "" {
			break
		}
	}

	if authToken == "" {
		cookie, err := r.Cookie("authToken")
		if err == nil {
			authToken = cookie.Value
		}
	}

	if authToken != "" {
		session, _ := s.interactor.CurrentSessionFromToken(authToken)
		if session != nil {
			context.Set(r, "currentSession", *session)
		}
	}

	next(w, r)

	context.Clear(r)
}
