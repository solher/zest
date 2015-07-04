package resources

import (
	"time"

	"github.com/solher/zest/domain"
)

func (i *AccountInter) scopeModel(account *domain.Account) error {
	account.CreatedAt = time.Time{}
	account.UpdatedAt = time.Time{}
	account.Users = []domain.User{}
	account.Sessions = []domain.Session{}
	account.RoleMappings = []domain.RoleMapping{}

	return nil
}

func (i *AccountInter) BeforeCreate(accounts []domain.Account) ([]domain.Account, error) {
	for k := range accounts {
		accounts[k].ID = 0
		err := i.scopeModel(&accounts[k])
		if err != nil {
			return nil, err
		}
	}
	return accounts, nil
}

func (i *AccountInter) AfterCreate(accounts []domain.Account) ([]domain.Account, error) {
	return accounts, nil
}

func (i *AccountInter) BeforeUpdate(accounts []domain.Account) ([]domain.Account, error) {
	for k := range accounts {
		err := i.scopeModel(&accounts[k])
		if err != nil {
			return nil, err
		}
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
