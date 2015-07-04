package resources

import (
	"database/sql"
	"strings"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/interfaces"
	"github.com/solher/zest/internalerrors"
	"github.com/solher/zest/usecases"
	"github.com/solher/zest/utils"
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
	sessions, err := r.Create([]domain.Session{*session})
	if err != nil {
		return nil, err
	}

	return &sessions[0], nil
}

func (r *SessionRepo) Find(context usecases.QueryContext) ([]domain.Session, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	var sessions []domain.Session

	err = query.Find(&sessions).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	if len(sessions) == 0 {
		sessions = []domain.Session{}
	}

	return sessions, nil
}

func (r *SessionRepo) FindByID(id int, context usecases.QueryContext) (*domain.Session, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	session := domain.Session{}

	err = query.Where(utils.ToDBName("sessions")+".id = ?", id).First(&session).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.NotFound
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

	for _, session := range sessions {
		queryCopy := *query

		dbName := utils.ToDBName("sessions")

		err = queryCopy.Where(dbName+".id = ?", session.ID).First(&domain.Session{}).Error
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return nil, internalerrors.NotFound
			}

			return nil, internalerrors.DatabaseError
		}

		err = r.store.GetDB().Where(dbName+".id = ?", session.ID).Model(&domain.Session{}).Updates(&session).Error
		if err != nil {
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

	dbName := utils.ToDBName("sessions")

	err = query.Where(dbName+".id = ?", id).First(&domain.Session{}).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.NotFound
		}

		return nil, internalerrors.DatabaseError
	}

	err = r.store.GetDB().Where(dbName+".id = ?", id).Model(&domain.Session{}).Updates(&session).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return session, nil
}

func (r *SessionRepo) DeleteAll(context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	sessions := []domain.Session{}
	err = query.Find(&sessions).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	sessionIDs := []int{}
	for _, session := range sessions {
		sessionIDs = append(sessionIDs, session.ID)
	}

	err = r.store.GetDB().Delete(&sessions, utils.ToDBName("sessions")+".id IN (?)", sessionIDs).Error
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

	err = query.Where(utils.ToDBName("sessions")+".id = ?", id).First(&session).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return internalerrors.NotFound
		}

		return internalerrors.DatabaseError
	}

	err = r.store.GetDB().Where(utils.ToDBName("sessions")+".id = ?", session.ID).Delete(domain.Session{}).Error
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

	return rows, nil
}
