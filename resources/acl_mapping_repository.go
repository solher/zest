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
	usecases.DependencyDirectory.Register(NewAclMappingRepo)
}

type AclMappingRepo struct {
	store interfaces.AbstractGormStore
}

func NewAclMappingRepo(store interfaces.AbstractGormStore) *AclMappingRepo {
	return &AclMappingRepo{store: store}
}

func (r *AclMappingRepo) Create(aclMappings []domain.AclMapping) ([]domain.AclMapping, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, aclMapping := range aclMappings {
		err := db.Create(&aclMapping).Error
		if err != nil {
			transaction.Rollback()

			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}

		aclMappings[i] = aclMapping
	}

	transaction.Commit()
	return aclMappings, nil
}

func (r *AclMappingRepo) CreateOne(aclMapping *domain.AclMapping) (*domain.AclMapping, error) {
	aclMappings, err := r.Create([]domain.AclMapping{*aclMapping})
	if err != nil {
		return nil, err
	}

	return &aclMappings[0], nil
}

func (r *AclMappingRepo) Find(context usecases.QueryContext) ([]domain.AclMapping, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	var aclMappings []domain.AclMapping

	err = query.Find(&aclMappings).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	if len(aclMappings) == 0 {
		aclMappings = []domain.AclMapping{}
	}

	return aclMappings, nil
}

func (r *AclMappingRepo) FindByID(id int, context usecases.QueryContext) (*domain.AclMapping, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	aclMapping := domain.AclMapping{}

	err = query.Where(utils.ToDBName("aclMappings")+".id = ?", id).First(&aclMapping).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.NotFound
		}

		return nil, internalerrors.DatabaseError
	}

	return &aclMapping, nil
}

func (r *AclMappingRepo) Update(aclMappings []domain.AclMapping, context usecases.QueryContext) ([]domain.AclMapping, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	for _, aclMapping := range aclMappings {
		queryCopy := *query

		dbName := utils.ToDBName("aclMappings")
		oldAclMapping := &domain.AclMapping{}

		err = queryCopy.Where(dbName+".id = ?", aclMapping.ID).First(oldAclMapping).Error
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return nil, internalerrors.NotFound
			}

			return nil, internalerrors.DatabaseError
		}

		aclMapping.ID = oldAclMapping.ID
		aclMapping.CreatedAt = oldAclMapping.CreatedAt

		err = r.store.GetDB().Save(&aclMapping).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}
	}

	transaction.Commit()
	return aclMappings, nil
}

func (r *AclMappingRepo) UpdateByID(id int, aclMapping *domain.AclMapping,
	context usecases.QueryContext) (*domain.AclMapping, error) {

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	dbName := utils.ToDBName("aclMappings")
	oldAclMapping := &domain.AclMapping{}

	err = query.Where(dbName+".id = ?", id).First(oldAclMapping).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.NotFound
		}

		return nil, internalerrors.DatabaseError
	}

	aclMapping.ID = oldAclMapping.ID
	aclMapping.CreatedAt = oldAclMapping.CreatedAt

	err = r.store.GetDB().Save(&aclMapping).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return aclMapping, nil
}

func (r *AclMappingRepo) DeleteAll(context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	aclMappings := []domain.AclMapping{}
	err = query.Find(&aclMappings).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	aclMappingIDs := []int{}
	for _, aclMapping := range aclMappings {
		aclMappingIDs = append(aclMappingIDs, aclMapping.ID)
	}

	err = r.store.GetDB().Delete(&aclMappings, utils.ToDBName("aclMappings")+".id IN (?)", aclMappingIDs).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *AclMappingRepo) DeleteByID(id int, context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	aclMapping := &domain.AclMapping{}

	err = query.Where(utils.ToDBName("aclMappings")+".id = ?", id).First(&aclMapping).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return internalerrors.NotFound
		}

		return internalerrors.DatabaseError
	}

	err = r.store.GetDB().Where(utils.ToDBName("aclMappings")+".id = ?", aclMapping.ID).Delete(domain.AclMapping{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *AclMappingRepo) Raw(query string, values ...interface{}) (*sql.Rows, error) {
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
