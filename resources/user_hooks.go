package resources

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/solher/zest/domain"
)

func (i *UserInter) scopeModel(user *domain.User) error {
	user.CreatedAt = time.Time{}
	user.UpdatedAt = time.Time{}
	user.Account = domain.Account{}

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
		if err != nil {
			return err
		}

		user.Password = string(hashedPassword)
	}

	return nil
}

func (i *UserInter) BeforeCreate(users []domain.User) ([]domain.User, error) {
	for k := range users {
		users[k].ID = 0
		err := i.scopeModel(&users[k])
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

func (i *UserInter) AfterCreate(users []domain.User) ([]domain.User, error) {
	return users, nil
}

func (i *UserInter) BeforeUpdate(users []domain.User) ([]domain.User, error) {
	for k := range users {
		err := i.scopeModel(&users[k])
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

func (i *UserInter) AfterUpdate(users []domain.User) ([]domain.User, error) {
	return users, nil
}

func (i *UserInter) BeforeDelete(users []domain.User) ([]domain.User, error) {
	return users, nil
}

func (i *UserInter) AfterDelete(users []domain.User) ([]domain.User, error) {
	return users, nil
}
