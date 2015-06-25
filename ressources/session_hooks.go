package ressources

import (
	"time"

	"github.com/solher/zest/domain"
)

func (i *SessionInter) scopeModel(session *domain.Session) {
	session.ID = 0
	session.CreatedAt = time.Time{}
	session.UpdatedAt = time.Time{}
	session.Account = domain.Account{}
}

func (i *SessionInter) refreshCache(session *domain.Session) error {
	err := i.sessionCacheInter.RefreshSession(session.ID)
	if err != nil {
		return err
	}
	return nil
}

func (i *SessionInter) removeFromCache(session *domain.Session) error {
	err := i.sessionCacheInter.Remove(session.AuthToken)
	if err != nil {
		return err
	}
	return nil
}

func (i *SessionInter) BeforeCreate(sessions []domain.Session) ([]domain.Session, error) {
	for k := range sessions {
		i.scopeModel(&sessions[k])
	}
	return sessions, nil
}

func (i *SessionInter) AfterCreate(sessions []domain.Session) ([]domain.Session, error) {
	for k := range sessions {
		err := i.refreshCache(&sessions[k])
		if err != nil {
			return nil, err
		}
	}
	return sessions, nil
}

func (i *SessionInter) BeforeUpdate(sessions []domain.Session) ([]domain.Session, error) {
	for k := range sessions {
		i.scopeModel(&sessions[k])
	}
	return sessions, nil
}

func (i *SessionInter) AfterUpdate(sessions []domain.Session) ([]domain.Session, error) {
	for k := range sessions {
		err := i.refreshCache(&sessions[k])
		if err != nil {
			return nil, err
		}
	}
	return sessions, nil
}

func (i *SessionInter) BeforeDelete(sessions []domain.Session) ([]domain.Session, error) {
	return sessions, nil
}

func (i *SessionInter) AfterDelete(sessions []domain.Session) ([]domain.Session, error) {
	for k := range sessions {
		err := i.removeFromCache(&sessions[k])
		if err != nil {
			return nil, err
		}
	}
	return sessions, nil
}
