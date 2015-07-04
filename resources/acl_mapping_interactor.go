package resources

import (
	"database/sql"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/usecases"
)

func init() {
	usecases.DependencyDirectory.Register(NewAclMappingInter)
	usecases.DependencyDirectory.Register(PopulateAclMappingInter)
}

type AbstractAclMappingRepo interface {
	Create(aclMappings []domain.AclMapping) ([]domain.AclMapping, error)
	CreateOne(aclMapping *domain.AclMapping) (*domain.AclMapping, error)
	Find(context usecases.QueryContext) ([]domain.AclMapping, error)
	FindByID(id int, context usecases.QueryContext) (*domain.AclMapping, error)
	Update(aclMappings []domain.AclMapping, context usecases.QueryContext) ([]domain.AclMapping, error)
	UpdateByID(id int, aclMapping *domain.AclMapping, context usecases.QueryContext) (*domain.AclMapping, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
	Raw(query string, values ...interface{}) (*sql.Rows, error)
}

type AclMappingInter struct {
	repo                 AbstractAclMappingRepo
	permissionCacheInter usecases.AbstractPermissionCacheInter
}

func NewAclMappingInter(repo AbstractAclMappingRepo, permissionCacheInter usecases.AbstractPermissionCacheInter) *AclMappingInter {
	return &AclMappingInter{repo: repo, permissionCacheInter: permissionCacheInter}
}

func PopulateAclMappingInter(aclMappingInter *AclMappingInter, repo AbstractAclMappingRepo, permissionCacheInter usecases.AbstractPermissionCacheInter) {
	if aclMappingInter.repo == nil {
		aclMappingInter.repo = repo
	}

	if aclMappingInter.permissionCacheInter == nil {
		aclMappingInter.permissionCacheInter = permissionCacheInter
	}
}

func (i *AclMappingInter) Create(aclMappings []domain.AclMapping) ([]domain.AclMapping, error) {
	aclMappings, err := i.BeforeCreate(aclMappings)
	if err != nil {
		return nil, err
	}

	aclMappings, err = i.repo.Create(aclMappings)
	if err != nil {
		return nil, err
	}

	aclMappings, err = i.AfterCreate(aclMappings)
	if err != nil {
		return nil, err
	}

	return aclMappings, nil
}

func (i *AclMappingInter) CreateOne(aclMapping *domain.AclMapping) (*domain.AclMapping, error) {
	aclMappings, err := i.Create([]domain.AclMapping{*aclMapping})
	if err != nil {
		return nil, err
	}

	return &aclMappings[0], nil
}

func (i *AclMappingInter) Find(context usecases.QueryContext) ([]domain.AclMapping, error) {
	aclMappings, err := i.repo.Find(context)
	if err != nil {
		return nil, err
	}

	return aclMappings, nil
}

func (i *AclMappingInter) FindByID(id int, context usecases.QueryContext) (*domain.AclMapping, error) {
	aclMapping, err := i.repo.FindByID(id, context)
	if err != nil {
		return nil, err
	}

	return aclMapping, nil
}

func (i *AclMappingInter) Upsert(aclMappings []domain.AclMapping, context usecases.QueryContext) ([]domain.AclMapping, error) {
	aclMappingsToUpdate := []domain.AclMapping{}
	aclMappingsToCreate := []domain.AclMapping{}

	for k := range aclMappings {
		if aclMappings[k].ID != 0 {
			aclMappingsToUpdate = append(aclMappingsToUpdate, aclMappings[k])
		} else {
			aclMappingsToCreate = append(aclMappingsToCreate, aclMappings[k])
		}
	}

	aclMappingsToUpdate, err := i.BeforeUpdate(aclMappingsToUpdate)
	if err != nil {
		return nil, err
	}

	aclMappingsToUpdate, err = i.repo.Update(aclMappingsToUpdate, context)
	if err != nil {
		return nil, err
	}

	aclMappingsToUpdate, err = i.AfterUpdate(aclMappingsToUpdate)
	if err != nil {
		return nil, err
	}

	aclMappingsToCreate, err = i.BeforeCreate(aclMappingsToCreate)
	if err != nil {
		return nil, err
	}

	aclMappingsToCreate, err = i.repo.Create(aclMappingsToCreate)
	if err != nil {
		return nil, err
	}

	aclMappingsToCreate, err = i.AfterCreate(aclMappingsToCreate)
	if err != nil {
		return nil, err
	}

	return append(aclMappingsToUpdate, aclMappingsToCreate...), nil
}

func (i *AclMappingInter) UpsertOne(aclMapping *domain.AclMapping, context usecases.QueryContext) (*domain.AclMapping, error) {
	aclMappings, err := i.Upsert([]domain.AclMapping{*aclMapping}, context)
	if err != nil {
		return nil, err
	}

	return &aclMappings[0], nil
}

func (i *AclMappingInter) UpdateByID(id int, aclMapping *domain.AclMapping,
	context usecases.QueryContext) (*domain.AclMapping, error) {

	aclMappings, err := i.BeforeUpdate([]domain.AclMapping{*aclMapping})
	if err != nil {
		return nil, err
	}

	aclMapping = &aclMappings[0]

	aclMapping, err = i.repo.UpdateByID(id, aclMapping, context)
	if err != nil {
		return nil, err
	}

	aclMappings, err = i.AfterUpdate([]domain.AclMapping{*aclMapping})
	if err != nil {
		return nil, err
	}

	return &aclMappings[0], nil
}

func (i *AclMappingInter) DeleteAll(context usecases.QueryContext) error {
	aclMappings, err := i.repo.Find(context)
	if err != nil {
		return err
	}

	aclMappings, err = i.BeforeDelete(aclMappings)
	if err != nil {
		return err
	}

	err = i.repo.DeleteAll(context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete(aclMappings)
	if err != nil {
		return err
	}

	return nil
}

func (i *AclMappingInter) DeleteByID(id int, context usecases.QueryContext) error {
	aclMapping, err := i.repo.FindByID(id, context)
	if err != nil {
		return err
	}

	aclMappings, err := i.BeforeDelete([]domain.AclMapping{*aclMapping})
	if err != nil {
		return err
	}

	aclMapping = &aclMappings[0]

	err = i.repo.DeleteByID(id, context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete([]domain.AclMapping{*aclMapping})
	if err != nil {
		return err
	}

	return nil
}
