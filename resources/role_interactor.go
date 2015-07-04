package resources

import (
	"database/sql"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/usecases"
)

func init() {
	usecases.DependencyDirectory.Register(NewRoleInter)
	usecases.DependencyDirectory.Register(PopulateRoleInter)
}

type AbstractRoleRepo interface {
	Create(roles []domain.Role) ([]domain.Role, error)
	CreateOne(role *domain.Role) (*domain.Role, error)
	Find(context usecases.QueryContext) ([]domain.Role, error)
	FindByID(id int, context usecases.QueryContext) (*domain.Role, error)
	Update(roles []domain.Role, context usecases.QueryContext) ([]domain.Role, error)
	UpdateByID(id int, role *domain.Role, context usecases.QueryContext) (*domain.Role, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
	Raw(query string, values ...interface{}) (*sql.Rows, error)
}

type RoleInter struct {
	repo             AbstractRoleRepo
	roleMappingInter AbstractRoleMappingInter
	aclMappingInter  AbstractAclMappingInter
}

func NewRoleInter(repo AbstractRoleRepo, roleMappingInter AbstractRoleMappingInter, aclMappingInter AbstractAclMappingInter) *RoleInter {
	return &RoleInter{repo: repo, roleMappingInter: roleMappingInter, aclMappingInter: aclMappingInter}
}

func PopulateRoleInter(roleInter *RoleInter, repo AbstractRoleRepo, roleMappingInter AbstractRoleMappingInter, aclMappingInter AbstractAclMappingInter) {
	if roleInter.repo == nil {
		roleInter.repo = repo
	}

	if roleInter.roleMappingInter == nil {
		roleInter.roleMappingInter = roleMappingInter
	}

	if roleInter.aclMappingInter == nil {
		roleInter.aclMappingInter = aclMappingInter
	}
}

func (i *RoleInter) Create(roles []domain.Role) ([]domain.Role, error) {
	roles, err := i.BeforeCreate(roles)
	if err != nil {
		return nil, err
	}

	roles, err = i.repo.Create(roles)
	if err != nil {
		return nil, err
	}

	roles, err = i.AfterCreate(roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (i *RoleInter) CreateOne(role *domain.Role) (*domain.Role, error) {
	roles, err := i.Create([]domain.Role{*role})
	if err != nil {
		return nil, err
	}

	return &roles[0], nil
}

func (i *RoleInter) Find(context usecases.QueryContext) ([]domain.Role, error) {
	roles, err := i.repo.Find(context)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (i *RoleInter) FindByID(id int, context usecases.QueryContext) (*domain.Role, error) {
	role, err := i.repo.FindByID(id, context)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (i *RoleInter) Upsert(roles []domain.Role, context usecases.QueryContext) ([]domain.Role, error) {
	rolesToUpdate := []domain.Role{}
	rolesToCreate := []domain.Role{}

	for k := range roles {
		if roles[k].ID != 0 {
			rolesToUpdate = append(rolesToUpdate, roles[k])
		} else {
			rolesToCreate = append(rolesToCreate, roles[k])
		}
	}

	rolesToUpdate, err := i.BeforeUpdate(rolesToUpdate)
	if err != nil {
		return nil, err
	}

	rolesToUpdate, err = i.repo.Update(rolesToUpdate, context)
	if err != nil {
		return nil, err
	}

	rolesToUpdate, err = i.AfterUpdate(rolesToUpdate)
	if err != nil {
		return nil, err
	}

	rolesToCreate, err = i.BeforeCreate(rolesToCreate)
	if err != nil {
		return nil, err
	}

	rolesToCreate, err = i.repo.Create(rolesToCreate)
	if err != nil {
		return nil, err
	}

	rolesToCreate, err = i.AfterCreate(rolesToCreate)
	if err != nil {
		return nil, err
	}

	return append(rolesToUpdate, rolesToCreate...), nil
}

func (i *RoleInter) UpsertOne(role *domain.Role, context usecases.QueryContext) (*domain.Role, error) {
	roles, err := i.Upsert([]domain.Role{*role}, context)
	if err != nil {
		return nil, err
	}

	return &roles[0], nil
}

func (i *RoleInter) UpdateByID(id int, role *domain.Role,
	context usecases.QueryContext) (*domain.Role, error) {

	roles, err := i.BeforeUpdate([]domain.Role{*role})
	if err != nil {
		return nil, err
	}

	role = &roles[0]

	role, err = i.repo.UpdateByID(id, role, context)
	if err != nil {
		return nil, err
	}

	roles, err = i.AfterUpdate([]domain.Role{*role})
	if err != nil {
		return nil, err
	}

	return &roles[0], nil
}

func (i *RoleInter) DeleteAll(context usecases.QueryContext) error {
	roles, err := i.repo.Find(context)
	if err != nil {
		return err
	}

	roles, err = i.BeforeDelete(roles)
	if err != nil {
		return err
	}

	err = i.repo.DeleteAll(context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete(roles)
	if err != nil {
		return err
	}

	roleIds := []int{}
	for _, role := range roles {
		roleIds = append(roleIds, role.ID)
	}

	filter := &usecases.Filter{
		Where: map[string]interface{}{"roleId": roleIds},
	}

	err = i.roleMappingInter.DeleteAll(usecases.QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	err = i.aclMappingInter.DeleteAll(usecases.QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	return nil
}

func (i *RoleInter) DeleteByID(id int, context usecases.QueryContext) error {
	role, err := i.repo.FindByID(id, context)
	if err != nil {
		return err
	}

	roles, err := i.BeforeDelete([]domain.Role{*role})
	if err != nil {
		return err
	}

	role = &roles[0]

	err = i.repo.DeleteByID(id, context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete([]domain.Role{*role})
	if err != nil {
		return err
	}

	filter := &usecases.Filter{
		Where: map[string]interface{}{"roleId": id},
	}

	err = i.roleMappingInter.DeleteAll(usecases.QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	err = i.aclMappingInter.DeleteAll(usecases.QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	return nil
}
