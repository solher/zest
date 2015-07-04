package resources

import (
	"database/sql"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/usecases"
)

func init() {
	usecases.DependencyDirectory.Register(NewRoleMappingInter)
	usecases.DependencyDirectory.Register(PopulateRoleMappingInter)
}

type AbstractRoleMappingRepo interface {
	Create(roleMappings []domain.RoleMapping) ([]domain.RoleMapping, error)
	CreateOne(roleMapping *domain.RoleMapping) (*domain.RoleMapping, error)
	Find(context usecases.QueryContext) ([]domain.RoleMapping, error)
	FindByID(id int, context usecases.QueryContext) (*domain.RoleMapping, error)
	Update(roleMappings []domain.RoleMapping, context usecases.QueryContext) ([]domain.RoleMapping, error)
	UpdateByID(id int, roleMapping *domain.RoleMapping, context usecases.QueryContext) (*domain.RoleMapping, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
	Raw(query string, values ...interface{}) (*sql.Rows, error)
}

type RoleMappingInter struct {
	repo                 AbstractRoleMappingRepo
	permissionCacheInter usecases.AbstractPermissionCacheInter
}

func NewRoleMappingInter(repo AbstractRoleMappingRepo, permissionCacheInter usecases.AbstractPermissionCacheInter) *RoleMappingInter {
	return &RoleMappingInter{repo: repo, permissionCacheInter: permissionCacheInter}
}

func PopulateRoleMappingInter(roleMappingInter *RoleMappingInter, repo AbstractRoleMappingRepo, permissionCacheInter usecases.AbstractPermissionCacheInter) {
	if roleMappingInter.repo == nil {
		roleMappingInter.repo = repo
	}

	if roleMappingInter.permissionCacheInter == nil {
		roleMappingInter.permissionCacheInter = permissionCacheInter
	}
}

func (i *RoleMappingInter) Create(roleMappings []domain.RoleMapping) ([]domain.RoleMapping, error) {
	roleMappings, err := i.BeforeCreate(roleMappings)
	if err != nil {
		return nil, err
	}

	roleMappings, err = i.repo.Create(roleMappings)
	if err != nil {
		return nil, err
	}

	roleMappings, err = i.AfterCreate(roleMappings)
	if err != nil {
		return nil, err
	}

	return roleMappings, nil
}

func (i *RoleMappingInter) CreateOne(roleMapping *domain.RoleMapping) (*domain.RoleMapping, error) {
	roleMappings, err := i.Create([]domain.RoleMapping{*roleMapping})
	if err != nil {
		return nil, err
	}

	return &roleMappings[0], nil
}

func (i *RoleMappingInter) Find(context usecases.QueryContext) ([]domain.RoleMapping, error) {
	roleMappings, err := i.repo.Find(context)
	if err != nil {
		return nil, err
	}

	return roleMappings, nil
}

func (i *RoleMappingInter) FindByID(id int, context usecases.QueryContext) (*domain.RoleMapping, error) {
	roleMapping, err := i.repo.FindByID(id, context)
	if err != nil {
		return nil, err
	}

	return roleMapping, nil
}

func (i *RoleMappingInter) Upsert(roleMappings []domain.RoleMapping, context usecases.QueryContext) ([]domain.RoleMapping, error) {
	roleMappingsToUpdate := []domain.RoleMapping{}
	roleMappingsToCreate := []domain.RoleMapping{}

	for k := range roleMappings {
		if roleMappings[k].ID != 0 {
			roleMappingsToUpdate = append(roleMappingsToUpdate, roleMappings[k])
		} else {
			roleMappingsToCreate = append(roleMappingsToCreate, roleMappings[k])
		}
	}

	roleMappingsToUpdate, err := i.BeforeUpdate(roleMappingsToUpdate)
	if err != nil {
		return nil, err
	}

	roleMappingsToUpdate, err = i.repo.Update(roleMappingsToUpdate, context)
	if err != nil {
		return nil, err
	}

	roleMappingsToUpdate, err = i.AfterUpdate(roleMappingsToUpdate)
	if err != nil {
		return nil, err
	}

	roleMappingsToCreate, err = i.BeforeCreate(roleMappingsToCreate)
	if err != nil {
		return nil, err
	}

	roleMappingsToCreate, err = i.repo.Create(roleMappingsToCreate)
	if err != nil {
		return nil, err
	}

	roleMappingsToCreate, err = i.AfterCreate(roleMappingsToCreate)
	if err != nil {
		return nil, err
	}

	return append(roleMappingsToUpdate, roleMappingsToCreate...), nil
}

func (i *RoleMappingInter) UpsertOne(roleMapping *domain.RoleMapping, context usecases.QueryContext) (*domain.RoleMapping, error) {
	roleMappings, err := i.Upsert([]domain.RoleMapping{*roleMapping}, context)
	if err != nil {
		return nil, err
	}

	return &roleMappings[0], nil
}

func (i *RoleMappingInter) UpdateByID(id int, roleMapping *domain.RoleMapping,
	context usecases.QueryContext) (*domain.RoleMapping, error) {

	roleMappings, err := i.BeforeUpdate([]domain.RoleMapping{*roleMapping})
	if err != nil {
		return nil, err
	}

	roleMapping = &roleMappings[0]

	roleMapping, err = i.repo.UpdateByID(id, roleMapping, context)
	if err != nil {
		return nil, err
	}

	roleMappings, err = i.AfterUpdate([]domain.RoleMapping{*roleMapping})
	if err != nil {
		return nil, err
	}

	return &roleMappings[0], nil
}

func (i *RoleMappingInter) DeleteAll(context usecases.QueryContext) error {
	roleMappings, err := i.repo.Find(context)
	if err != nil {
		return err
	}

	roleMappings, err = i.BeforeDelete(roleMappings)
	if err != nil {
		return err
	}

	err = i.repo.DeleteAll(context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete(roleMappings)
	if err != nil {
		return err
	}

	return nil
}

func (i *RoleMappingInter) DeleteByID(id int, context usecases.QueryContext) error {
	roleMapping, err := i.repo.FindByID(id, context)
	if err != nil {
		return err
	}

	roleMappings, err := i.BeforeDelete([]domain.RoleMapping{*roleMapping})
	if err != nil {
		return err
	}

	roleMapping = &roleMappings[0]

	err = i.repo.DeleteByID(id, context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete([]domain.RoleMapping{*roleMapping})
	if err != nil {
		return err
	}

	return nil
}
