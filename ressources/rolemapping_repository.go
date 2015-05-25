// Generated by: main
// TypeWriter: repository
// Directive: +gen on RoleMapping

package ressources

import (
	"strings"

	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/internalerrors"
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
			} else {
				return nil, internalerrors.DatabaseError
			}
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
		} else {
			return nil, internalerrors.DatabaseError
		}
	}

	return rolemapping, nil
}

func (r *RoleMappingRepo) Find(filter *interfaces.Filter) ([]domain.RoleMapping, error) {
	query, err := r.store.BuildQuery(filter)
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

func (r *RoleMappingRepo) FindByID(id int, filter *interfaces.Filter) (*domain.RoleMapping, error) {
	query, err := r.store.BuildQuery(filter)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	rolemapping := domain.RoleMapping{}

	err = query.First(&rolemapping, id).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return &rolemapping, nil
}

func (r *RoleMappingRepo) Upsert(rolemappings []domain.RoleMapping) ([]domain.RoleMapping, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, rolemapping := range rolemappings {
		if rolemapping.ID != 0 {
			oldUser := domain.RoleMapping{}

			err := db.First(&oldUser, rolemapping.ID).Updates(rolemapping).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				} else {
					return nil, internalerrors.DatabaseError
				}
			}
		} else {
			err := db.Create(&rolemapping).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				} else {
					return nil, internalerrors.DatabaseError
				}
			}
		}

		rolemappings[i] = rolemapping
	}

	transaction.Commit()
	return rolemappings, nil
}

func (r *RoleMappingRepo) UpsertOne(rolemapping *domain.RoleMapping) (*domain.RoleMapping, error) {
	db := r.store.GetDB()

	if rolemapping.ID != 0 {
		oldUser := domain.RoleMapping{}

		err := db.First(&oldUser, rolemapping.ID).Updates(rolemapping).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			} else {
				return nil, internalerrors.DatabaseError
			}
		}
	} else {
		err := db.Create(&rolemapping).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			} else {
				return nil, internalerrors.DatabaseError
			}
		}
	}

	return rolemapping, nil
}

func (r *RoleMappingRepo) DeleteAll(filter *interfaces.Filter) error {
	query, err := r.store.BuildQuery(filter)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(domain.RoleMapping{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *RoleMappingRepo) DeleteByID(id int) error {
	db := r.store.GetDB()

	err := db.Delete(&domain.RoleMapping{GormModel: domain.GormModel{ID: id}}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}
