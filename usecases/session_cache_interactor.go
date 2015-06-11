package usecases

import "github.com/solher/zest/domain"

type AbstractSessionRepo interface {
	Find(context QueryContext) ([]domain.Session, error)
	FindByID(id int, context QueryContext) (*domain.Session, error)
}

type SessionCacheInter struct {
	sessionRepo  AbstractSessionRepo
	sessionCache AbstractCacheStore
}

func NewSessionCacheInter(sessionRepo AbstractSessionRepo, sessionCache AbstractCacheStore) *SessionCacheInter {
	return &SessionCacheInter{sessionRepo: sessionRepo, sessionCache: sessionCache}
}

func (i *SessionCacheInter) Add(authToken string, session domain.Session) error {
	i.sessionCache.Add(authToken, session)

	return nil
}

func (i *SessionCacheInter) Remove(authToken string) error {
	i.sessionCache.Remove(authToken)

	return nil
}

func (i *SessionCacheInter) Get(authToken string) (domain.Session, error) {
	value, err := i.sessionCache.Get(authToken)

	if err != nil {
		return domain.Session{}, err
	}

	return value.(domain.Session), nil
}

func (i *SessionCacheInter) Refresh() error {
	filter := &Filter{
		Limit: i.sessionCache.MaxSize(),
		Order: "updatedAt DESC",
	}

	sessions, err := i.sessionRepo.Find(QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	for _, session := range sessions {
		i.sessionCache.Add(session.AuthToken, session)
	}

	return nil
}

func (i *SessionCacheInter) RefreshSession(sessionID int) error {
	session, err := i.sessionRepo.FindByID(sessionID, QueryContext{})
	if err != nil {
		return err
	}

	i.sessionCache.Add(session.AuthToken, session)

	return nil
}
