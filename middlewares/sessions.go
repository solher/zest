package middlewares

import (
	"net/http"

	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/ressources"
	"github.com/gorilla/context"
)

type Sessions struct {
	interactor ressources.AbstractSessionInter
}

func NewSessions(store interfaces.AbstractGormStore) *Sessions {
	sessionRepository := ressources.NewSessionRepo(store)
	sessionInteractor := ressources.NewSessionInter(sessionRepository)

	return &Sessions{interactor: sessionInteractor}
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
		session, _ := s.interactor.CurrentFromToken(authToken)
		if session != nil {
			context.Set(r, "currentSession", *session)
		}
	}

	next(w, r)

	context.Clear(r)
}
