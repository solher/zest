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
	usecases.DependencyDirectory.Register(NewAclRepo)
}

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
			}

			return nil, internalerrors.DatabaseError
		}

		acls[i] = acl
	}

	transaction.Commit()
	return acls, nil
}

func (r *AclRepo) CreateOne(acl *domain.Acl) (*domain.Acl, error) {
	acls, err := r.Create([]domain.Acl{*acl})
	if err != nil {
		return nil, err
	}

	return &acls[0], nil
}

func (r *AclRepo) Find(context usecases.QueryContext) ([]domain.Acl, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	var acls []domain.Acl

	err = query.Find(&acls).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	if len(acls) == 0 {
		acls = []domain.Acl{}
	}

	return acls, nil
}

func (r *AclRepo) FindByID(id int, context usecases.QueryContext) (*domain.Acl, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	acl := domain.Acl{}

	err = query.Where(utils.ToDBName("acls")+".id = ?", id).First(&acl).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.NotFound
		}

		return nil, internalerrors.DatabaseError
	}

	return &acl, nil
}

func (r *AclRepo) Update(acls []domain.Acl, context usecases.QueryContext) ([]domain.Acl, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	for _, acl := range acls {
		queryCopy := *query

		dbName := utils.ToDBName("acls")

		err = queryCopy.Where(dbName+".id = ?", acl.ID).First(&domain.Acl{}).Error
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return nil, internalerrors.NotFound
			}

			return nil, internalerrors.DatabaseError
		}

		err = r.store.GetDB().Where(dbName+".id = ?", acl.ID).Model(&domain.Acl{}).Updates(&acl).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}
	}

	transaction.Commit()
	return acls, nil
}

func (r *AclRepo) UpdateByID(id int, acl *domain.Acl,
	context usecases.QueryContext) (*domain.Acl, error) {

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	dbName := utils.ToDBName("acls")

	err = query.Where(dbName+".id = ?", id).First(&domain.Acl{}).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.NotFound
		}

		return nil, internalerrors.DatabaseError
	}

	err = r.store.GetDB().Where(dbName+".id = ?", id).Model(&domain.Acl{}).Updates(&acl).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return acl, nil
}

func (r *AclRepo) DeleteAll(context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	acls := []domain.Acl{}
	err = query.Find(&acls).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	aclIDs := []int{}
	for _, acl := range acls {
		aclIDs = append(aclIDs, acl.ID)
	}

	err = r.store.GetDB().Delete(&acls, utils.ToDBName("acls")+".id IN (?)", aclIDs).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *AclRepo) DeleteByID(id int, context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	acl := &domain.Acl{}

	err = query.Where(utils.ToDBName("acls")+".id = ?", id).First(&acl).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return internalerrors.NotFound
		}

		return internalerrors.DatabaseError
	}

	err = r.store.GetDB().Where(utils.ToDBName("acls")+".id = ?", acl.ID).Delete(domain.Acl{}).Error
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
		}

		return nil, internalerrors.DatabaseError
	}

	return rows, nil
}
