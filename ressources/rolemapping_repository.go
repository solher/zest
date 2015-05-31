// Generated by: main
// TypeWriter: repository
// Directive: +gen on RoleMapping

package ressources

import (
	"database/sql"
	"strings"

	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/internalerrors"
	"github.com/Solher/auth-scaffold/usecases"
)

type RoleMappingRepo struct {
	store interfaces.AbstractGormStore
}

func NewRoleMappingRepo(store interfaces.AbstractGormStore) *RoleMappingRepo {
	return &RoleMappingRepo{store: store}
}

func (r *RoleMappingRepo) Create(rolemappings []domain.RoleMapping) ([]domain.RoleMapping, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, rolemapping := range rolemappings {
		err := db.Create(&rolemapping).Error
		if err != nil {
			transaction.Rollback()

			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}

		rolemappings[i] = rolemapping
	}

	transaction.Commit()
	return rolemappings, nil
}

func (r *RoleMappingRepo) CreateOne(rolemapping *domain.RoleMapping) (*domain.RoleMapping, error) {
	db := r.store.GetDB()

	err := db.Create(rolemapping).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return rolemapping, nil
}

func (r *RoleMappingRepo) Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.RoleMapping, error) {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	rolemappings := []domain.RoleMapping{}

	err = query.Find(&rolemappings).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return rolemappings, nil
}

func (r *RoleMappingRepo) FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.RoleMapping, error) {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	rolemapping := domain.RoleMapping{}

	err = query.Where("rolemappings.id = ?", id).First(&rolemapping).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return &rolemapping, nil
}

func (r *RoleMappingRepo) Upsert(rolemappings []domain.RoleMapping, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.RoleMapping, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	for i, rolemapping := range rolemappings {
		queryCopy := *query

		if rolemapping.ID != 0 {
			oldUser := domain.RoleMapping{}

			err := queryCopy.Where("rolemappings.id = ?", rolemapping.ID).First(&oldUser).Updates(rolemapping).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				}

				return nil, internalerrors.DatabaseError
			}
		} else {
			err := db.Create(&rolemapping).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				}

				return nil, internalerrors.DatabaseError
			}
		}

		rolemappings[i] = rolemapping
	}

	transaction.Commit()
	return rolemappings, nil
}

func (r *RoleMappingRepo) UpsertOne(rolemapping *domain.RoleMapping, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.RoleMapping, error) {
	db := r.store.GetDB()

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	if rolemapping.ID != 0 {
		oldUser := domain.RoleMapping{}

		err := query.Where("rolemappings.id = ?", rolemapping.ID).First(&oldUser).Updates(rolemapping).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}
	} else {
		err := db.Create(&rolemapping).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}
	}

	return rolemapping, nil
}

func (r *RoleMappingRepo) UpdateByID(id int, rolemapping *domain.RoleMapping,
	filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.RoleMapping, error) {

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	oldUser := domain.RoleMapping{}

	err = query.Where("rolemappings.id = ?", id).First(&oldUser).Updates(rolemapping).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return rolemapping, nil
}

func (r *RoleMappingRepo) DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(domain.RoleMapping{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *RoleMappingRepo) DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(&domain.RoleMapping{GormModel: domain.GormModel{ID: id}}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *RoleMappingRepo) Raw(query string, values ...interface{}) (*sql.Rows, error) {
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
