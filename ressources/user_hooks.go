package ressources

import (
	"time"

	"github.com/solher/zest/domain"
)

func (i *UserInter) scopeModel(user *domain.User) {
	user.ID = 0
	user.CreatedAt = time.Time{}
	user.UpdatedAt = time.Time{}
	user.Account = domain.Account{}
}

func (i *UserInter) BeforeCreate(users []domain.User) ([]domain.User, error) {
	for k := range users {
		i.scopeModel(&users[k])
	}
	return users, nil
}

func (i *UserInter) AfterCreate(users []domain.User) ([]domain.User, error) {
	return users, nil
}

func (i *UserInter) BeforeUpdate(users []domain.User) ([]domain.User, error) {
	for k := range users {
		i.scopeModel(&users[k])
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
