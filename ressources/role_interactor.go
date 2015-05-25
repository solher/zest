// Generated by: main
// TypeWriter: interactor
// Directive: +gen on Role

package ressources

import (
	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/interfaces"
)

type AbstractRoleRepo interface {
	Create(roles []domain.Role) ([]domain.Role, error)
	CreateOne(role *domain.Role) (*domain.Role, error)
	Find(filter *interfaces.Filter) ([]domain.Role, error)
	FindByID(id int, filter *interfaces.Filter) (*domain.Role, error)
	Upsert(roles []domain.Role) ([]domain.Role, error)
	UpsertOne(role *domain.Role) (*domain.Role, error)
	DeleteAll(filter *interfaces.Filter) error
	DeleteByID(id int) error
}

type RoleInter struct {
	repo AbstractRoleRepo
}

func NewRoleInter(repo AbstractRoleRepo) *RoleInter {
	return &RoleInter{repo: repo}
}

func (i *RoleInter) Create(roles []domain.Role) ([]domain.Role, error) {
	roles, err := i.repo.Create(roles)
	return roles, err
}

func (i *RoleInter) CreateOne(role *domain.Role) (*domain.Role, error) {
	role, err := i.repo.CreateOne(role)
	return role, err
}

func (i *RoleInter) Find(filter *interfaces.Filter) ([]domain.Role, error) {
	roles, err := i.repo.Find(filter)
	return roles, err
}

func (i *RoleInter) FindByID(id int, filter *interfaces.Filter) (*domain.Role, error) {
	role, err := i.repo.FindByID(id, filter)
	return role, err
}

func (i *RoleInter) Upsert(roles []domain.Role) ([]domain.Role, error) {
	roles, err := i.repo.Upsert(roles)
	return roles, err
}

func (i *RoleInter) UpsertOne(role *domain.Role) (*domain.Role, error) {
	role, err := i.repo.UpsertOne(role)
	return role, err
}

func (i *RoleInter) DeleteAll(filter *interfaces.Filter) error {
	err := i.repo.DeleteAll(filter)
	return err
}

func (i *RoleInter) DeleteByID(id int) error {
	err := i.repo.DeleteByID(id)
	return err
}
