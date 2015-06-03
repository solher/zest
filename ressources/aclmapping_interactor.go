// Generated by: main
// TypeWriter: interactor
// Directive: +gen on AclMapping

package ressources

import (
	"database/sql"
	"time"

	"github.com/Solher/zest/domain"
	"github.com/Solher/zest/usecases"
)

type AbstractAclMappingRepo interface {
	Create(aclMappings []domain.AclMapping) ([]domain.AclMapping, error)
	CreateOne(aclMapping *domain.AclMapping) (*domain.AclMapping, error)
	Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.AclMapping, error)
	FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.AclMapping, error)
	Update(aclMappings []domain.AclMapping, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.AclMapping, error)
	UpdateByID(id int, aclMapping *domain.AclMapping, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.AclMapping, error)
	DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error
	DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error
	Raw(query string, values ...interface{}) (*sql.Rows, error)
}

type AclMappingInter struct {
	repo                 AbstractAclMappingRepo
	permissionCacheInter AbstractPermissionCacheInter
}

func NewAclMappingInter(repo AbstractAclMappingRepo) *AclMappingInter {
	return &AclMappingInter{repo: repo}
}

func (i *AclMappingInter) BeforeSave(aclMapping *domain.AclMapping) error {
	aclMapping.ID = 0
	aclMapping.CreatedAt = time.Time{}
	aclMapping.UpdatedAt = time.Time{}

	err := aclMapping.ScopeModel()
	if err != nil {
		return err
	}

	return nil
}

func (i *AclMappingInter) AfterModification(aclMapping *domain.AclMapping) error {
	err := i.permissionCacheInter.RefreshAcl(aclMapping.AclID)
	if err != nil {
		return err
	}

	return nil
}

func (i *AclMappingInter) Create(aclMappings []domain.AclMapping) ([]domain.AclMapping, error) {
	var err error

	for k := range aclMappings {
		err := i.BeforeSave(&aclMappings[k])
		if err != nil {
			return nil, err
		}
	}

	aclMappings, err = i.repo.Create(aclMappings)
	if err != nil {
		return nil, err
	}

	for k := range aclMappings {
		err := i.AfterModification(&aclMappings[k])
		if err != nil {
			return nil, err
		}
	}

	return aclMappings, nil
}

func (i *AclMappingInter) CreateOne(aclMapping *domain.AclMapping) (*domain.AclMapping, error) {
	err := i.BeforeSave(aclMapping)
	if err != nil {
		return nil, err
	}

	aclMapping, err = i.repo.CreateOne(aclMapping)
	if err != nil {
		return nil, err
	}

	err = i.AfterModification(aclMapping)
	if err != nil {
		return nil, err
	}

	return aclMapping, nil
}

func (i *AclMappingInter) Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.AclMapping, error) {
	aclMappings, err := i.repo.Find(filter, ownerRelations)
	if err != nil {
		return nil, err
	}

	return aclMappings, nil
}

func (i *AclMappingInter) FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.AclMapping, error) {
	aclMapping, err := i.repo.FindByID(id, filter, ownerRelations)
	if err != nil {
		return nil, err
	}

	return aclMapping, nil
}

func (i *AclMappingInter) Upsert(aclMappings []domain.AclMapping, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.AclMapping, error) {
	aclMappingsToUpdate := []domain.AclMapping{}
	aclMappingsToCreate := []domain.AclMapping{}

	for k := range aclMappings {
		err := i.BeforeSave(&aclMappings[k])
		if err != nil {
			return nil, err
		}

		if aclMappings[k].ID != 0 {
			aclMappingsToUpdate = append(aclMappingsToUpdate, aclMappings[k])
		} else {
			aclMappingsToCreate = append(aclMappingsToCreate, aclMappings[k])
		}
	}

	aclMappingsToUpdate, err := i.repo.Update(aclMappingsToUpdate, filter, ownerRelations)
	if err != nil {
		return nil, err
	}

	aclMappingsToCreate, err = i.repo.Create(aclMappingsToCreate)
	if err != nil {
		return nil, err
	}

	aclMappings = append(aclMappingsToUpdate, aclMappingsToCreate...)

	for k := range aclMappings {
		err := i.AfterModification(&aclMappings[k])
		if err != nil {
			return nil, err
		}
	}

	return aclMappings, nil
}

func (i *AclMappingInter) UpsertOne(aclMapping *domain.AclMapping, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.AclMapping, error) {
	err := i.BeforeSave(aclMapping)
	if err != nil {
		return nil, err
	}

	if aclMapping.ID != 0 {
		aclMapping, err = i.repo.UpdateByID(aclMapping.ID, aclMapping, filter, ownerRelations)
	} else {
		aclMapping, err = i.repo.CreateOne(aclMapping)
	}

	if err != nil {
		return nil, err
	}

	err = i.AfterModification(aclMapping)
	if err != nil {
		return nil, err
	}

	return aclMapping, nil
}

func (i *AclMappingInter) UpdateByID(id int, aclMapping *domain.AclMapping,
	filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.AclMapping, error) {

	err := i.BeforeSave(aclMapping)
	if err != nil {
		return nil, err
	}

	aclMapping, err = i.repo.UpdateByID(id, aclMapping, filter, ownerRelations)
	if err != nil {
		return nil, err
	}

	err = i.AfterModification(aclMapping)
	if err != nil {
		return nil, err
	}

	return aclMapping, nil
}

func (i *AclMappingInter) DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error {
	filter.Fields = nil

	aclMappings, err := i.repo.Find(filter, ownerRelations)
	if err != nil {
		return err
	}

	err = i.repo.DeleteAll(filter, ownerRelations)
	if err != nil {
		return err
	}

	for k := range aclMappings {
		err := i.AfterModification(&aclMappings[k])
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *AclMappingInter) DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error {
	filter.Fields = nil

	aclMapping, err := i.repo.FindByID(id, filter, ownerRelations)
	if err != nil {
		return err
	}

	err = i.repo.DeleteByID(id, filter, ownerRelations)
	if err != nil {
		return err
	}

	err = i.AfterModification(aclMapping)
	if err != nil {
		return err
	}

	return nil
}
