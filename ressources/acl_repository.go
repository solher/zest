// Generated by: main
// TypeWriter: repository
// Directive: +gen on Acl

package ressources

import (
	"database/sql"
	"strings"

	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/internalerrors"
)

type AclRepo struct {
	store interfaces.AbstractGormStore
}

func NewAclRepo(store interfaces.AbstractGormStore) *AclRepo {
	return &AclRepo{store: store}
}

func (r *AclRepo) Create(acls []domain.Acl) ([]domain.Acl, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, acl := range acls {
		err := db.Create(&acl).Error
		if err != nil {
			transaction.Rollback()

			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			} else {
				return nil, internalerrors.DatabaseError
			}
		}

		acls[i] = acl
	}

	transaction.Commit()
	return acls, nil
}

func (r *AclRepo) CreateOne(acl *domain.Acl) (*domain.Acl, error) {
	db := r.store.GetDB()

	err := db.Create(acl).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		} else {
			return nil, internalerrors.DatabaseError
		}
	}

	return acl, nil
}

func (r *AclRepo) Find(filter *interfaces.Filter, ownerRelations []domain.Relation) ([]domain.Acl, error) {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	acls := []domain.Acl{}

	err = query.Find(&acls).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return acls, nil
}

func (r *AclRepo) FindByID(id int, filter *interfaces.Filter, ownerRelations []domain.Relation) (*domain.Acl, error) {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	acl := domain.Acl{}

	err = query.Where("acls.id = ?", id).First(&acl).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return &acl, nil
}

func (r *AclRepo) Upsert(acls []domain.Acl, filter *interfaces.Filter, ownerRelations []domain.Relation) ([]domain.Acl, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	for i, acl := range acls {
		queryCopy := *query

		if acl.ID != 0 {
			oldUser := domain.Acl{}

			err := queryCopy.Where("acls.id = ?", acl.ID).First(&oldUser).Updates(acl).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				} else {
					return nil, internalerrors.DatabaseError
				}
			}
		} else {
			err := db.Create(&acl).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				} else {
					return nil, internalerrors.DatabaseError
				}
			}
		}

		acls[i] = acl
	}

	transaction.Commit()
	return acls, nil
}

func (r *AclRepo) UpsertOne(acl *domain.Acl, filter *interfaces.Filter, ownerRelations []domain.Relation) (*domain.Acl, error) {
	db := r.store.GetDB()

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	if acl.ID != 0 {
		oldUser := domain.Acl{}

		err := query.Where("acls.id = ?", acl.ID).First(&oldUser).Updates(acl).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			} else {
				return nil, internalerrors.DatabaseError
			}
		}
	} else {
		err := db.Create(&acl).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			} else {
				return nil, internalerrors.DatabaseError
			}
		}
	}

	return acl, nil
}

func (r *AclRepo) UpdateByID(id int, acl *domain.Acl,
	filter *interfaces.Filter, ownerRelations []domain.Relation) (*domain.Acl, error) {

	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	oldUser := domain.Acl{}

	err = query.Where("acls.id = ?", id).First(&oldUser).Updates(acl).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		} else {
			return nil, internalerrors.DatabaseError
		}
	}

	return acl, nil
}

func (r *AclRepo) DeleteAll(filter *interfaces.Filter, ownerRelations []domain.Relation) error {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(domain.Acl{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *AclRepo) DeleteByID(id int, filter *interfaces.Filter, ownerRelations []domain.Relation) error {
	query, err := r.store.BuildQuery(filter, ownerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	err = query.Delete(&domain.Acl{GormModel: domain.GormModel{ID: id}}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *AclRepo) Raw(query string, values ...interface{}) (*sql.Rows, error) {
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
