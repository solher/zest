package middlewares

import (
	"net/http"

	"github.com/Solher/auth-scaffold/ressources"
	"github.com/gorilla/context"
)

type Sessions struct {
	interactor ressources.AbstractAccountInter
}

func NewSessions(accountRepo ressources.AbstractAccountRepo, userRepo ressources.AbstractUserRepo, sessionRepo ressources.AbstractSessionRepo) *Sessions {
	accountInteractor := ressources.NewAccountInter(accountRepo, userRepo, sessionRepo)

	return &Sessions{interactor: accountInteractor}
}

func (s *Sessions) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authToken := r.Header.Get("Authorization")

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
