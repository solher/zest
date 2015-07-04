package resources

import (
	"database/sql"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/usecases"
)

func init() {
	usecases.DependencyDirectory.Register(NewUserInter)
	usecases.DependencyDirectory.Register(PopulateUserInter)
}

type AbstractUserRepo interface {
	Create(users []domain.User) ([]domain.User, error)
	CreateOne(user *domain.User) (*domain.User, error)
	Find(context usecases.QueryContext) ([]domain.User, error)
	FindByID(id int, context usecases.QueryContext) (*domain.User, error)
	Update(users []domain.User, context usecases.QueryContext) ([]domain.User, error)
	UpdateByID(id int, user *domain.User, context usecases.QueryContext) (*domain.User, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
	Raw(query string, values ...interface{}) (*sql.Rows, error)
}

type UserInter struct {
	repo AbstractUserRepo
}

func NewUserInter(repo AbstractUserRepo) *UserInter {
	return &UserInter{repo: repo}
}

func PopulateUserInter(userInter *UserInter, repo AbstractUserRepo) {
	if userInter.repo == nil {
		userInter.repo = repo
	}
}

func (i *UserInter) Create(users []domain.User) ([]domain.User, error) {
	users, err := i.BeforeCreate(users)
	if err != nil {
		return nil, err
	}

	users, err = i.repo.Create(users)
	if err != nil {
		return nil, err
	}

	users, err = i.AfterCreate(users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (i *UserInter) CreateOne(user *domain.User) (*domain.User, error) {
	users, err := i.Create([]domain.User{*user})
	if err != nil {
		return nil, err
	}

	return &users[0], nil
}

func (i *UserInter) Find(context usecases.QueryContext) ([]domain.User, error) {
	users, err := i.repo.Find(context)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (i *UserInter) FindByID(id int, context usecases.QueryContext) (*domain.User, error) {
	user, err := i.repo.FindByID(id, context)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (i *UserInter) Upsert(users []domain.User, context usecases.QueryContext) ([]domain.User, error) {
	usersToUpdate := []domain.User{}
	usersToCreate := []domain.User{}

	for k := range users {
		if users[k].ID != 0 {
			usersToUpdate = append(usersToUpdate, users[k])
		} else {
			usersToCreate = append(usersToCreate, users[k])
		}
	}

	usersToUpdate, err := i.BeforeUpdate(usersToUpdate)
	if err != nil {
		return nil, err
	}

	usersToUpdate, err = i.repo.Update(usersToUpdate, context)
	if err != nil {
		return nil, err
	}

	usersToUpdate, err = i.AfterUpdate(usersToUpdate)
	if err != nil {
		return nil, err
	}

	usersToCreate, err = i.BeforeCreate(usersToCreate)
	if err != nil {
		return nil, err
	}

	usersToCreate, err = i.repo.Create(usersToCreate)
	if err != nil {
		return nil, err
	}

	usersToCreate, err = i.AfterCreate(usersToCreate)
	if err != nil {
		return nil, err
	}

	return append(usersToUpdate, usersToCreate...), nil
}

func (i *UserInter) UpsertOne(user *domain.User, context usecases.QueryContext) (*domain.User, error) {
	users, err := i.Upsert([]domain.User{*user}, context)
	if err != nil {
		return nil, err
	}

	return &users[0], nil
}

func (i *UserInter) UpdateByID(id int, user *domain.User,
	context usecases.QueryContext) (*domain.User, error) {

	users, err := i.BeforeUpdate([]domain.User{*user})
	if err != nil {
		return nil, err
	}

	user = &users[0]

	user, err = i.repo.UpdateByID(id, user, context)
	if err != nil {
		return nil, err
	}

	users, err = i.AfterUpdate([]domain.User{*user})
	if err != nil {
		return nil, err
	}

	return &users[0], nil
}

func (i *UserInter) DeleteAll(context usecases.QueryContext) error {
	users, err := i.repo.Find(context)
	if err != nil {
		return err
	}

	users, err = i.BeforeDelete(users)
	if err != nil {
		return err
	}

	err = i.repo.DeleteAll(context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete(users)
	if err != nil {
		return err
	}

	return nil
}

func (i *UserInter) DeleteByID(id int, context usecases.QueryContext) error {
	user, err := i.repo.FindByID(id, context)
	if err != nil {
		return err
	}

	users, err := i.BeforeDelete([]domain.User{*user})
	if err != nil {
		return err
	}

	user = &users[0]

	err = i.repo.DeleteByID(id, context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete([]domain.User{*user})
	if err != nil {
		return err
	}

	return nil
}
