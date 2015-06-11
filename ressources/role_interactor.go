// Generated by: main
// TypeWriter: interactor
// Directive: +gen on Role

package ressources

import (
	"database/sql"
	"time"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/usecases"
)

func init() {
	usecases.DependencyDirectory.Register(NewRoleInter)
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
	repo AbstractRoleRepo
}

func NewRoleInter(repo AbstractRoleRepo) *RoleInter {
	return &RoleInter{repo: repo}
}

func (i *RoleInter) BeforeSave(role *domain.Role) error {
	role.ID = 0
	role.CreatedAt = time.Time{}
	role.UpdatedAt = time.Time{}

	err := role.ScopeModel()
	if err != nil {
		return err
	}

	return nil
}

func (i *RoleInter) Create(roles []domain.Role) ([]domain.Role, error) {
	var err error

	for k := range roles {
		err := i.BeforeSave(&roles[k])
		if err != nil {
			return nil, err
		}
	}

	roles, err = i.repo.Create(roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (i *RoleInter) CreateOne(role *domain.Role) (*domain.Role, error) {
	err := i.BeforeSave(role)
	if err != nil {
		return nil, err
	}

	role, err = i.repo.CreateOne(role)
	if err != nil {
		return nil, err
	}

	return role, nil
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
		err := i.BeforeSave(&roles[k])
		if err != nil {
			return nil, err
		}

		if roles[k].ID != 0 {
			rolesToUpdate = append(rolesToUpdate, roles[k])
		} else {
			rolesToCreate = append(rolesToCreate, roles[k])
		}
	}

	rolesToUpdate, err := i.repo.Update(rolesToUpdate, context)
	if err != nil {
		return nil, err
	}

	rolesToCreate, err = i.repo.Create(rolesToCreate)
	if err != nil {
		return nil, err
	}

	return append(rolesToUpdate, rolesToCreate...), nil
}

func (i *RoleInter) UpsertOne(role *domain.Role, context usecases.QueryContext) (*domain.Role, error) {
	err := i.BeforeSave(role)
	if err != nil {
		return nil, err
	}

	if role.ID != 0 {
		role, err = i.repo.UpdateByID(role.ID, role, context)
	} else {
		role, err = i.repo.CreateOne(role)
	}

	if err != nil {
		return nil, err
	}

	return role, nil
}

func (i *RoleInter) UpdateByID(id int, role *domain.Role,
	context usecases.QueryContext) (*domain.Role, error) {

	err := i.BeforeSave(role)
	if err != nil {
		return nil, err
	}

	role, err = i.repo.UpdateByID(id, role, context)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (i *RoleInter) DeleteAll(context usecases.QueryContext) error {
	err := i.repo.DeleteAll(context)
	if err != nil {
		return err
	}

	return nil
}

func (i *RoleInter) DeleteByID(id int, context usecases.QueryContext) error {
	err := i.repo.DeleteByID(id, context)
	if err != nil {
		return err
	}

	return nil
}
