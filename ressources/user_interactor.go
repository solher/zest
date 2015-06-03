// Generated by: main
// TypeWriter: interactor
// Directive: +gen on User

package ressources

import (
	"database/sql"
	"time"

	"github.com/Solher/zest/domain"
	"github.com/Solher/zest/usecases"
)

type AbstractUserRepo interface {
	Create(users []domain.User) ([]domain.User, error)
	CreateOne(user *domain.User) (*domain.User, error)
	Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.User, error)
	FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.User, error)
	Update(users []domain.User, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.User, error)
	UpdateByID(id int, user *domain.User, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.User, error)
	DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error
	DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error
	Raw(query string, values ...interface{}) (*sql.Rows, error)
}

type UserInter struct {
	repo AbstractUserRepo
}

func NewUserInter(repo AbstractUserRepo) *UserInter {
	return &UserInter{repo: repo}
}

func (i *UserInter) BeforeSave(user *domain.User) error {
	user.ID = 0
	user.CreatedAt = time.Time{}
	user.UpdatedAt = time.Time{}

	err := user.ScopeModel()
	if err != nil {
		return err
	}

	return nil
}

func (i *UserInter) Create(users []domain.User) ([]domain.User, error) {
	var err error

	for k := range users {
		err := i.BeforeSave(&users[k])
		if err != nil {
			return nil, err
		}
	}

	users, err = i.repo.Create(users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (i *UserInter) CreateOne(user *domain.User) (*domain.User, error) {
	err := i.BeforeSave(user)
	if err != nil {
		return nil, err
	}

	user, err = i.repo.CreateOne(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (i *UserInter) Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.User, error) {
	users, err := i.repo.Find(filter, ownerRelations)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (i *UserInter) FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.User, error) {
	user, err := i.repo.FindByID(id, filter, ownerRelations)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (i *UserInter) Upsert(users []domain.User, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.User, error) {
	usersToUpdate := []domain.User{}
	usersToCreate := []domain.User{}

	for k := range users {
		err := i.BeforeSave(&users[k])
		if err != nil {
			return nil, err
		}

		if users[k].ID != 0 {
			usersToUpdate = append(usersToUpdate, users[k])
		} else {
			usersToCreate = append(usersToCreate, users[k])
		}
	}

	usersToUpdate, err := i.repo.Update(usersToUpdate, filter, ownerRelations)
	if err != nil {
		return nil, err
	}

	usersToCreate, err = i.repo.Create(usersToCreate)
	if err != nil {
		return nil, err
	}

	return append(usersToUpdate, usersToCreate...), nil
}

func (i *UserInter) UpsertOne(user *domain.User, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.User, error) {
	err := i.BeforeSave(user)
	if err != nil {
		return nil, err
	}

	if user.ID != 0 {
		user, err = i.repo.UpdateByID(user.ID, user, filter, ownerRelations)
	} else {
		user, err = i.repo.CreateOne(user)
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (i *UserInter) UpdateByID(id int, user *domain.User,
	filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.User, error) {

	err := i.BeforeSave(user)
	if err != nil {
		return nil, err
	}

	user, err = i.repo.UpdateByID(id, user, filter, ownerRelations)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (i *UserInter) DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error {
	err := i.repo.DeleteAll(filter, ownerRelations)
	if err != nil {
		return err
	}

	return nil
}

func (i *UserInter) DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error {
	err := i.repo.DeleteByID(id, filter, ownerRelations)
	if err != nil {
		return err
	}

	return nil
}
