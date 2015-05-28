// Generated by: main
// TypeWriter: repository
// Directive: +gen on User

package ressources

import (
	"database/sql"
	"strings"

	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/internalerrors"
)

type UserRepo struct {
	store interfaces.AbstractGormStore
}

func NewUserRepo(store interfaces.AbstractGormStore) *UserRepo {
	return &UserRepo{store: store}
}

func (r *UserRepo) Create(users []domain.User) ([]domain.User, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, user := range users {
		err := db.Create(&user).Error
		if err != nil {
			transaction.Rollback()

			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			} else {
				return nil, internalerrors.DatabaseError
			}
		}

		users[i] = user
	}

	transaction.Commit()
	return users, nil
}

func (r *UserRepo) CreateOne(user *domain.User) (*domain.User, error) {
	db := r.store.GetDB()

	err := db.Create(user).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		} else {
			return nil, internalerrors.DatabaseError
		}
	}

	return user, nil
}

func (r *UserRepo) Find(filter *interfaces.Filter, ownerRelations []domain.Relation) ([]domain.User, error) {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	users := []domain.User{}

	err = query.Find(&users).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return users, nil
}

func (r *UserRepo) FindByID(id int, filter *interfaces.Filter, ownerRelations []domain.Relation) (*domain.User, error) {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	user := domain.User{}

	err = query.Where("users.id = ?", id).First(&user).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return &user, nil
}

func (r *UserRepo) Upsert(users []domain.User, filter *interfaces.Filter, ownerRelations []domain.Relation) ([]domain.User, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	for i, user := range users {
		queryCopy := *query

		if user.ID != 0 {
			oldUser := domain.User{}

			err := queryCopy.Where("users.id = ?", user.ID).First(&oldUser).Updates(user).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				} else {
					return nil, internalerrors.DatabaseError
				}
			}
		} else {
			err := db.Create(&user).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				} else {
					return nil, internalerrors.DatabaseError
				}
			}
		}

		users[i] = user
	}

	transaction.Commit()
	return users, nil
}

func (r *UserRepo) UpsertOne(user *domain.User, filter *interfaces.Filter, ownerRelations []domain.Relation) (*domain.User, error) {
	db := r.store.GetDB()

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	if user.ID != 0 {
		oldUser := domain.User{}

		err := query.Where("users.id = ?", user.ID).First(&oldUser).Updates(user).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			} else {
				return nil, internalerrors.DatabaseError
			}
		}
	} else {
		err := db.Create(&user).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			} else {
				return nil, internalerrors.DatabaseError
			}
		}
	}

	return user, nil
}

func (r *UserRepo) DeleteAll(filter *interfaces.Filter, ownerRelations []domain.Relation) error {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(domain.User{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *UserRepo) DeleteByID(id int, filter *interfaces.Filter, ownerRelations []domain.Relation) error {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(&domain.User{GormModel: domain.GormModel{ID: id}}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *UserRepo) Raw(query string, values ...interface{}) (*sql.Rows, error) {
	db := r.store.GetDB()

	rows, err := db.Raw(query, values...).Rows()
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		} else {
			return nil, internalerrors.DatabaseError
		}
	}

	return rows, nil
}
