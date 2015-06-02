// Generated by: main
// TypeWriter: repository
// Directive: +gen on Role

package ressources

import (
	"database/sql"
	"strings"

	"github.com/Solher/zest/domain"
	"github.com/Solher/zest/interfaces"
	"github.com/Solher/zest/internalerrors"
	"github.com/Solher/zest/usecases"
)

type RoleRepo struct {
	store interfaces.AbstractGormStore
}

func NewRoleRepo(store interfaces.AbstractGormStore) *RoleRepo {
	return &RoleRepo{store: store}
}

func (r *RoleRepo) Create(roles []domain.Role) ([]domain.Role, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, role := range roles {
		err := db.Create(&role).Error
		if err != nil {
			transaction.Rollback()

			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}

		roles[i] = role
	}

	transaction.Commit()
	return roles, nil
}

func (r *RoleRepo) CreateOne(role *domain.Role) (*domain.Role, error) {
	db := r.store.GetDB()

	err := db.Create(role).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return role, nil
}

func (r *RoleRepo) Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.Role, error) {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	roles := []domain.Role{}

	err = query.Find(&roles).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return roles, nil
}

func (r *RoleRepo) FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.Role, error) {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	role := domain.Role{}

	err = query.Where("roles.id = ?", id).First(&role).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return &role, nil
}

func (r *RoleRepo) Update(roles []domain.Role, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.Role, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	for i, role := range roles {
		queryCopy := *query
		oldUser := domain.Role{}

		err := queryCopy.Where("roles.id = ?", role.ID).First(&oldUser).Updates(roles[i]).Error
		if err != nil {
			transaction.Rollback()

			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}
	}

	transaction.Commit()
	return roles, nil
}

func (r *RoleRepo) UpdateByID(id int, role *domain.Role,
	filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.Role, error) {

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	oldUser := domain.Role{}

	err = query.Where("roles.id = ?", id).First(&oldUser).Updates(role).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return role, nil
}

func (r *RoleRepo) DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(domain.Role{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *RoleRepo) DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(&domain.Role{GormModel: domain.GormModel{ID: id}}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *RoleRepo) Raw(query string, values ...interface{}) (*sql.Rows, error) {
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
