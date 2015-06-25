package ressources

import (
	"database/sql"
	"strings"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/interfaces"
	"github.com/solher/zest/internalerrors"
	"github.com/solher/zest/usecases"
)

func init() {
	usecases.DependencyDirectory.Register(NewSessionRepo)
}

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
			}

			return nil, internalerrors.DatabaseError
		}

		sessions[i] = session
	}

	transaction.Commit()
	return sessions, nil
}

func (r *SessionRepo) CreateOne(session *domain.Session) (*domain.Session, error) {
	r.Create([]domain.Session{*session})
	return session, nil
}

func (r *SessionRepo) Find(context usecases.QueryContext) ([]domain.Session, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
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

func (r *SessionRepo) FindByID(id int, context usecases.QueryContext) (*domain.Session, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	session := domain.Session{}

	err = query.Where("sessions.id = ?", id).First(&session).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.InsufficentPermissions
		}

		return nil, internalerrors.DatabaseError
	}

	return &session, nil
}

func (r *SessionRepo) Update(sessions []domain.Session, context usecases.QueryContext) ([]domain.Session, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	for i, session := range sessions {
		queryCopy := *query
		oldSession := domain.Session{}

		err := queryCopy.Where("sessions.id = ?", session.ID).First(&oldSession).Updates(sessions[i]).Error
		if err != nil {
			transaction.Rollback()

			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}
	}

	transaction.Commit()
	return sessions, nil
}

func (r *SessionRepo) UpdateByID(id int, session *domain.Session,
	context usecases.QueryContext) (*domain.Session, error) {

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	oldSession := domain.Session{}

	err = query.Where("sessions.id = ?", id).First(&oldSession).Updates(session).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.InsufficentPermissions
		}

		return nil, internalerrors.DatabaseError
	}

	if session.ID == 0 {
		return nil, internalerrors.InsufficentPermissions
	}

	return session, nil
}

func (r *SessionRepo) DeleteAll(context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(domain.Session{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *SessionRepo) DeleteByID(id int, context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	session := &domain.Session{}

	err = query.Where("sessions.id = ?", id).First(&session).Delete(domain.Session{}).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return internalerrors.InsufficentPermissions
		}

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

	return rows, nil
}
