package middlewares

import (
	"net/http"

	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/ressources/sessions"
	"github.com/gorilla/context"
)

type Sessions struct {
	interactor *sessions.Interactor
}

func NewSessions(store interfaces.GormStore) *Sessions {
	sessionsRepository := sessions.NewRepository(store)
	sessionsInteractor := sessions.NewInteractor(sessionsRepository)

	return &Sessions{interactor: sessionsInteractor}
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
		session, user, _ := s.interactor.CurrentSession("toto")
		context.Set(r, "currentSession", session)
		context.Set(r, "currentUser", user)
	}

	next(w, r)

	context.Clear(r)
}
