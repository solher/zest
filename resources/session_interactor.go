package resources

import (
	"database/sql"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/usecases"
)

func init() {
	usecases.DependencyDirectory.Register(NewSessionInter)
	usecases.DependencyDirectory.Register(PopulateSessionInter)
}

type AbstractSessionRepo interface {
	Create(sessions []domain.Session) ([]domain.Session, error)
	CreateOne(session *domain.Session) (*domain.Session, error)
	Find(context usecases.QueryContext) ([]domain.Session, error)
	FindByID(id int, context usecases.QueryContext) (*domain.Session, error)
	Update(sessions []domain.Session, context usecases.QueryContext) ([]domain.Session, error)
	UpdateByID(id int, session *domain.Session, context usecases.QueryContext) (*domain.Session, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
	Raw(query string, values ...interface{}) (*sql.Rows, error)
}

type SessionInter struct {
	repo              AbstractSessionRepo
	sessionCacheInter usecases.AbstractSessionCacheInter
}

func NewSessionInter(repo AbstractSessionRepo, sessionCacheInter usecases.AbstractSessionCacheInter) *SessionInter {
	return &SessionInter{repo: repo, sessionCacheInter: sessionCacheInter}
}

func PopulateSessionInter(sessionInter *SessionInter, repo AbstractSessionRepo, sessionCacheInter usecases.AbstractSessionCacheInter) {
	if sessionInter.repo == nil {
		sessionInter.repo = repo
	}

	if sessionInter.sessionCacheInter == nil {
		sessionInter.sessionCacheInter = sessionCacheInter
	}
}

func (i *SessionInter) Create(sessions []domain.Session) ([]domain.Session, error) {
	sessions, err := i.BeforeCreate(sessions)
	if err != nil {
		return nil, err
	}

	sessions, err = i.repo.Create(sessions)
	if err != nil {
		return nil, err
	}

	sessions, err = i.AfterCreate(sessions)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (i *SessionInter) CreateOne(session *domain.Session) (*domain.Session, error) {
	sessions, err := i.Create([]domain.Session{*session})
	if err != nil {
		return nil, err
	}

	return &sessions[0], nil
}

func (i *SessionInter) Find(context usecases.QueryContext) ([]domain.Session, error) {
	sessions, err := i.repo.Find(context)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (i *SessionInter) FindByID(id int, context usecases.QueryContext) (*domain.Session, error) {
	session, err := i.repo.FindByID(id, context)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (i *SessionInter) Upsert(sessions []domain.Session, context usecases.QueryContext) ([]domain.Session, error) {
	sessionsToUpdate := []domain.Session{}
	sessionsToCreate := []domain.Session{}

	for k := range sessions {
		if sessions[k].ID != 0 {
			sessionsToUpdate = append(sessionsToUpdate, sessions[k])
		} else {
			sessionsToCreate = append(sessionsToCreate, sessions[k])
		}
	}

	sessionsToUpdate, err := i.BeforeUpdate(sessionsToUpdate)
	if err != nil {
		return nil, err
	}

	sessionsToUpdate, err = i.repo.Update(sessionsToUpdate, context)
	if err != nil {
		return nil, err
	}

	sessionsToUpdate, err = i.AfterUpdate(sessionsToUpdate)
	if err != nil {
		return nil, err
	}

	sessionsToCreate, err = i.BeforeCreate(sessionsToCreate)
	if err != nil {
		return nil, err
	}

	sessionsToCreate, err = i.repo.Create(sessionsToCreate)
	if err != nil {
		return nil, err
	}

	sessionsToCreate, err = i.AfterCreate(sessionsToCreate)
	if err != nil {
		return nil, err
	}

	return append(sessionsToUpdate, sessionsToCreate...), nil
}

func (i *SessionInter) UpsertOne(session *domain.Session, context usecases.QueryContext) (*domain.Session, error) {
	sessions, err := i.Upsert([]domain.Session{*session}, context)
	if err != nil {
		return nil, err
	}

	return &sessions[0], nil
}

func (i *SessionInter) UpdateByID(id int, session *domain.Session,
	context usecases.QueryContext) (*domain.Session, error) {

	sessions, err := i.BeforeUpdate([]domain.Session{*session})
	if err != nil {
		return nil, err
	}

	session = &sessions[0]

	session, err = i.repo.UpdateByID(id, session, context)
	if err != nil {
		return nil, err
	}

	sessions, err = i.AfterUpdate([]domain.Session{*session})
	if err != nil {
		return nil, err
	}

	return &sessions[0], nil
}

func (i *SessionInter) DeleteAll(context usecases.QueryContext) error {
	sessions, err := i.repo.Find(context)
	if err != nil {
		return err
	}

	sessions, err = i.BeforeDelete(sessions)
	if err != nil {
		return err
	}

	err = i.repo.DeleteAll(context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete(sessions)
	if err != nil {
		return err
	}

	return nil
}

func (i *SessionInter) DeleteByID(id int, context usecases.QueryContext) error {
	session, err := i.repo.FindByID(id, context)
	if err != nil {
		return err
	}

	sessions, err := i.BeforeDelete([]domain.Session{*session})
	if err != nil {
		return err
	}

	session = &sessions[0]

	err = i.repo.DeleteByID(id, context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete([]domain.Session{*session})
	if err != nil {
		return err
	}

	return nil
}
