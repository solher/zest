// Generated by: main
// TypeWriter: repository
// Directive: +gen on Session

package ressources

import (
	"strings"

	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/internalerrors"
)

type SessionRepo struct {
	store interfaces.AbstractGormStore
}

func NewSessionRepo(store interfaces.AbstractGormStore) *SessionRepo {
	return &SessionRepo{store: store}
}

func (r *SessionRepo) Create(sessions []domain.Session) ([]domain.Session, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, session := range sessions {
		err := db.Create(&session).Error
		if err != nil {
			transaction.Rollback()

			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			} else {
				return nil, internalerrors.DatabaseError
			}
		}

		sessions[i] = session
	}

	transaction.Commit()
	return sessions, nil
}

func (r *SessionRepo) CreateOne(session *domain.Session) (*domain.Session, error) {
	db := r.store.GetDB()

	err := db.Create(session).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		} else {
			return nil, internalerrors.DatabaseError
		}
	}

	return session, nil
}

func (r *SessionRepo) Find(filter *interfaces.Filter) ([]domain.Session, error) {
	query, err := r.store.BuildQuery(filter)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	sessions := []domain.Session{}

	err = query.Find(&sessions).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return sessions, nil
}

func (r *SessionRepo) FindByID(id int, filter *interfaces.Filter) (*domain.Session, error) {
	query, err := r.store.BuildQuery(filter)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	session := domain.Session{}

	err = query.First(&session, id).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return &session, nil
}

func (r *SessionRepo) Upsert(sessions []domain.Session) ([]domain.Session, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, session := range sessions {
		if session.ID != 0 {
			oldUser := domain.Session{}

			err := db.First(&oldUser, session.ID).Updates(session).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				} else {
					return nil, internalerrors.DatabaseError
				}
			}
		} else {
			err := db.Create(&session).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				} else {
					return nil, internalerrors.DatabaseError
				}
			}
		}

		sessions[i] = session
	}

	transaction.Commit()
	return sessions, nil
}

func (r *SessionRepo) UpsertOne(session *domain.Session) (*domain.Session, error) {
	db := r.store.GetDB()

	if session.ID != 0 {
		oldUser := domain.Session{}

		err := db.First(&oldUser, session.ID).Updates(session).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			} else {
				return nil, internalerrors.DatabaseError
			}
		}
	} else {
		err := db.Create(&session).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			} else {
				return nil, internalerrors.DatabaseError
			}
		}
	}

	return session, nil
}

func (r *SessionRepo) DeleteAll(filter *interfaces.Filter) error {
	query, err := r.store.BuildQuery(filter)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(domain.Session{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *SessionRepo) DeleteByID(id int) error {
	db := r.store.GetDB()

	err := db.Delete(&domain.Session{GormModel: domain.GormModel{ID: id}}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}
