// Generated by: main
// TypeWriter: interactor
// Directive: +gen on RoleMapping

package ressources

import (
	"database/sql"
	"time"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/usecases"
)

func init() {
	usecases.DependencyDirectory.Register(NewRoleMappingInter)
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

func NewRoleMappingInter(repo AbstractRoleMappingRepo) *RoleMappingInter {
	return &RoleMappingInter{repo: repo}
}

func (i *RoleMappingInter) BeforeSave(roleMapping *domain.RoleMapping) error {
	roleMapping.ID = 0
	roleMapping.CreatedAt = time.Time{}
	roleMapping.UpdatedAt = time.Time{}

	err := roleMapping.ScopeModel()
	if err != nil {
		return err
	}

	return nil
}

func (i *RoleMappingInter) AfterModification(roleMapping *domain.RoleMapping) error {
	err := i.permissionCacheInter.RefreshRole(roleMapping.AccountID)
	if err != nil {
		return err
	}

	return nil
}

func (i *RoleMappingInter) Create(roleMappings []domain.RoleMapping) ([]domain.RoleMapping, error) {
	var err error

	for k := range roleMappings {
		err := i.BeforeSave(&roleMappings[k])
		if err != nil {
			return nil, err
		}
	}

	roleMappings, err = i.repo.Create(roleMappings)
	if err != nil {
		return nil, err
	}

	for k := range roleMappings {
		err := i.AfterModification(&roleMappings[k])
		if err != nil {
			return nil, err
		}
	}

	return roleMappings, nil
}

func (i *RoleMappingInter) CreateOne(roleMapping *domain.RoleMapping) (*domain.RoleMapping, error) {
	err := i.BeforeSave(roleMapping)
	if err != nil {
		return nil, err
	}

	roleMapping, err = i.repo.CreateOne(roleMapping)
	if err != nil {
		return nil, err
	}

	err = i.AfterModification(roleMapping)
	if err != nil {
		return nil, err
	}

	return roleMapping, nil
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
		err := i.BeforeSave(&roleMappings[k])
		if err != nil {
			return nil, err
		}

		if roleMappings[k].ID != 0 {
			roleMappingsToUpdate = append(roleMappingsToUpdate, roleMappings[k])
		} else {
			roleMappingsToCreate = append(roleMappingsToCreate, roleMappings[k])
		}
	}

	roleMappingsToUpdate, err := i.repo.Update(roleMappingsToUpdate, context)
	if err != nil {
		return nil, err
	}

	roleMappingsToCreate, err = i.repo.Create(roleMappingsToCreate)
	if err != nil {
		return nil, err
	}

	roleMappings = append(roleMappingsToUpdate, roleMappingsToCreate...)

	for k := range roleMappings {
		err := i.AfterModification(&roleMappings[k])
		if err != nil {
			return nil, err
		}
	}

	return roleMappings, nil
}

func (i *RoleMappingInter) UpsertOne(roleMapping *domain.RoleMapping, context usecases.QueryContext) (*domain.RoleMapping, error) {
	err := i.BeforeSave(roleMapping)
	if err != nil {
		return nil, err
	}

	if roleMapping.ID != 0 {
		roleMapping, err = i.repo.UpdateByID(roleMapping.ID, roleMapping, context)
	} else {
		roleMapping, err = i.repo.CreateOne(roleMapping)
	}

	if err != nil {
		return nil, err
	}

	err = i.AfterModification(roleMapping)
	if err != nil {
		return nil, err
	}

	return roleMapping, nil
}

func (i *RoleMappingInter) UpdateByID(id int, roleMapping *domain.RoleMapping,
	context usecases.QueryContext) (*domain.RoleMapping, error) {

	err := i.BeforeSave(roleMapping)
	if err != nil {
		return nil, err
	}

	roleMapping, err = i.repo.UpdateByID(id, roleMapping, context)
	if err != nil {
		return nil, err
	}

	err = i.AfterModification(roleMapping)
	if err != nil {
		return nil, err
	}

	return roleMapping, nil
}

func (i *RoleMappingInter) DeleteAll(context usecases.QueryContext) error {
	context.Filter.Fields = nil

	roleMappings, err := i.repo.Find(context)
	if err != nil {
		return err
	}

	err = i.repo.DeleteAll(context)
	if err != nil {
		return err
	}

	for k := range roleMappings {
		err := i.AfterModification(&roleMappings[k])
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *RoleMappingInter) DeleteByID(id int, context usecases.QueryContext) error {
	context.Filter.Fields = nil

	roleMapping, err := i.repo.FindByID(id, context)
	if err != nil {
		return err
	}

	err = i.repo.DeleteByID(id, context)
	if err != nil {
		return err
	}

	err = i.AfterModification(roleMapping)
	if err != nil {
		return err
	}

	return nil
}
