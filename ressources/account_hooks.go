package ressources

import (
	"time"

	"github.com/solher/zest/domain"
)

func (i *AccountInter) scopeModel(account *domain.Account) {
	account.ID = 0
	account.CreatedAt = time.Time{}
	account.UpdatedAt = time.Time{}
	account.Users = []domain.User{}
	account.Sessions = []domain.Session{}
	account.RoleMappings = []domain.RoleMapping{}
}

func (i *AccountInter) BeforeCreate(accounts []domain.Account) ([]domain.Account, error) {
	for k := range accounts {
		i.scopeModel(&accounts[k])
	}
	return accounts, nil
}

func (i *AccountInter) AfterCreate(accounts []domain.Account) ([]domain.Account, error) {
	return accounts, nil
}

func (i *AccountInter) BeforeUpdate(accounts []domain.Account) ([]domain.Account, error) {
	for k := range accounts {
		i.scopeModel(&accounts[k])
	}
	return accounts, nil
}

func (i *AccountInter) AfterUpdate(accounts []domain.Account) ([]domain.Account, error) {
	return accounts, nil
}

func (i *AccountInter) BeforeDelete(accounts []domain.Account) ([]domain.Account, error) {
	return accounts, nil
}

func (i *AccountInter) AfterDelete(accounts []domain.Account) ([]domain.Account, error) {
	return accounts, nil
}
