// Generated by: main
// TypeWriter: repository
// Directive: +gen on Session

package ressources

import (
	"database/sql"
	"strings"

	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/internalerrors"
)

type SessionRepo struct {
	store interfaces.AbstractGormStore
	cache interfaces.AbstractLRUCacheStore
}

func NewSessionRepo(store interfaces.AbstractGormStore, cache interfaces.AbstractLRUCacheStore) *SessionRepo {
	return &SessionRepo{store: store, cache: cache}
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
			}

			return nil, internalerrors.DatabaseError
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
		}

		return nil, internalerrors.DatabaseError
	}

	return session, nil
}

func (r *SessionRepo) Find(filter *interfaces.Filter, ownerRelations []domain.Relation) ([]domain.Session, error) {
	query, err := r.store.BuildQuery(filter, ownerRelations)
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

func (r *SessionRepo) FindByID(id int, filter *interfaces.Filter, ownerRelations []domain.Relation) (*domain.Session, error) {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	session := domain.Session{}

	err = query.Where("sessions.id = ?", id).First(&session).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return &session, nil
}

func (r *SessionRepo) Upsert(sessions []domain.Session, filter *interfaces.Filter, ownerRelations []domain.Relation) ([]domain.Session, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	for i, session := range sessions {
		queryCopy := *query

		if session.ID != 0 {
			oldUser := domain.Session{}

			authToken := session.AuthToken

			err := queryCopy.Where("sessions.id = ?", session.ID).First(&oldUser).Updates(session).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				}

				return nil, internalerrors.DatabaseError
			}

			r.cache.Remove(authToken)
			err = r.cache.Add(session.AuthToken, session)
			if err != nil {
				return nil, internalerrors.DatabaseError
			}
		} else {
			err := db.Create(&session).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				}

				return nil, internalerrors.DatabaseError
			}
		}

		sessions[i] = session
	}

	transaction.Commit()
	return sessions, nil
}

func (r *SessionRepo) UpsertOne(session *domain.Session, filter *interfaces.Filter, ownerRelations []domain.Relation) (*domain.Session, error) {
	db := r.store.GetDB()

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	if session.ID != 0 {
		oldUser := domain.Session{}

		authToken := session.AuthToken

		err := query.Where("sessions.id = ?", session.ID).First(&oldUser).Updates(session).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}

		r.cache.Remove(authToken)
		err = r.cache.Add(session.AuthToken, session)
		if err != nil {
			return nil, internalerrors.DatabaseError
		}
	} else {
		err := db.Create(&session).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}
	}

	return session, nil
}

func (r *SessionRepo) UpdateByID(id int, session *domain.Session,
	filter *interfaces.Filter, ownerRelations []domain.Relation) (*domain.Session, error) {

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	oldUser := domain.Session{}

	authToken := session.AuthToken

	err = query.Where("sessions.id = ?", id).First(&oldUser).Updates(session).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	r.cache.Remove(authToken)
	err = r.cache.Add(session.AuthToken, session)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return session, nil
}

func (r *SessionRepo) DeleteAll(filter *interfaces.Filter, ownerRelations []domain.Relation) error {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(domain.Session{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = r.cache.Purge()
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *SessionRepo) DeleteByID(id int, filter *interfaces.Filter, ownerRelations []domain.Relation) error {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(&domain.Session{GormModel: domain.GormModel{ID: id}}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = r.cache.Purge()
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *SessionRepo) Raw(query string, values ...interface{}) (*sql.Rows, error) {
	db := r.store.GetDB()

	rows, err := db.Raw(query, values...).Rows()
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	err = r.cache.Purge()
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return rows, nil
}
