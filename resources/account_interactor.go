package resources

import (
	"database/sql"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/usecases"
)

func init() {
	usecases.DependencyDirectory.Register(NewAccountInter)
	usecases.DependencyDirectory.Register(PopulateAccountInter)
}

type AbstractAccountRepo interface {
	Create(accounts []domain.Account) ([]domain.Account, error)
	CreateOne(account *domain.Account) (*domain.Account, error)
	Find(context usecases.QueryContext) ([]domain.Account, error)
	FindByID(id int, context usecases.QueryContext) (*domain.Account, error)
	Update(accounts []domain.Account, context usecases.QueryContext) ([]domain.Account, error)
	UpdateByID(id int, account *domain.Account, context usecases.QueryContext) (*domain.Account, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
	Raw(query string, values ...interface{}) (*sql.Rows, error)
}

type AccountInter struct {
	repo             AbstractAccountRepo
	userInter        AbstractUserInter
	sessionInter     AbstractSessionInter
	roleMappingInter AbstractRoleMappingInter
}

func NewAccountInter(repo AbstractAccountRepo, userInter AbstractUserInter, sessionInter AbstractSessionInter, roleMappingInter AbstractRoleMappingInter) *AccountInter {
	return &AccountInter{repo: repo, userInter: userInter, sessionInter: sessionInter, roleMappingInter: roleMappingInter}
}

func PopulateAccountInter(accountInter *AccountInter, repo AbstractAccountRepo, userInter AbstractUserInter, sessionInter AbstractSessionInter, roleMappingInter AbstractRoleMappingInter) {
	if accountInter.repo == nil {
		accountInter.repo = repo
	}

	if accountInter.userInter == nil {
		accountInter.userInter = userInter
	}

	if accountInter.sessionInter == nil {
		accountInter.sessionInter = sessionInter
	}

	if accountInter.roleMappingInter == nil {
		accountInter.roleMappingInter = roleMappingInter
	}
}

func (i *AccountInter) Create(accounts []domain.Account) ([]domain.Account, error) {
	accounts, err := i.BeforeCreate(accounts)
	if err != nil {
		return nil, err
	}

	accounts, err = i.repo.Create(accounts)
	if err != nil {
		return nil, err
	}

	accounts, err = i.AfterCreate(accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (i *AccountInter) CreateOne(account *domain.Account) (*domain.Account, error) {
	accounts, err := i.Create([]domain.Account{*account})
	if err != nil {
		return nil, err
	}

	return &accounts[0], nil
}

func (i *AccountInter) Find(context usecases.QueryContext) ([]domain.Account, error) {
	accounts, err := i.repo.Find(context)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (i *AccountInter) FindByID(id int, context usecases.QueryContext) (*domain.Account, error) {
	account, err := i.repo.FindByID(id, context)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (i *AccountInter) Upsert(accounts []domain.Account, context usecases.QueryContext) ([]domain.Account, error) {
	accountsToUpdate := []domain.Account{}
	accountsToCreate := []domain.Account{}

	for k := range accounts {
		if accounts[k].ID != 0 {
			accountsToUpdate = append(accountsToUpdate, accounts[k])
		} else {
			accountsToCreate = append(accountsToCreate, accounts[k])
		}
	}

	accountsToUpdate, err := i.BeforeUpdate(accountsToUpdate)
	if err != nil {
		return nil, err
	}

	accountsToUpdate, err = i.repo.Update(accountsToUpdate, context)
	if err != nil {
		return nil, err
	}

	accountsToUpdate, err = i.AfterUpdate(accountsToUpdate)
	if err != nil {
		return nil, err
	}

	accountsToCreate, err = i.BeforeCreate(accountsToCreate)
	if err != nil {
		return nil, err
	}

	accountsToCreate, err = i.repo.Create(accountsToCreate)
	if err != nil {
		return nil, err
	}

	accountsToCreate, err = i.AfterCreate(accountsToCreate)
	if err != nil {
		return nil, err
	}

	return append(accountsToUpdate, accountsToCreate...), nil
}

func (i *AccountInter) UpsertOne(account *domain.Account, context usecases.QueryContext) (*domain.Account, error) {
	accounts, err := i.Upsert([]domain.Account{*account}, context)
	if err != nil {
		return nil, err
	}

	return &accounts[0], nil
}

func (i *AccountInter) UpdateByID(id int, account *domain.Account,
	context usecases.QueryContext) (*domain.Account, error) {

	accounts, err := i.BeforeUpdate([]domain.Account{*account})
	if err != nil {
		return nil, err
	}

	account = &accounts[0]

	account, err = i.repo.UpdateByID(id, account, context)
	if err != nil {
		return nil, err
	}

	accounts, err = i.AfterUpdate([]domain.Account{*account})
	if err != nil {
		return nil, err
	}

	return &accounts[0], nil
}

func (i *AccountInter) DeleteAll(context usecases.QueryContext) error {
	accounts, err := i.repo.Find(context)
	if err != nil {
		return err
	}

	accounts, err = i.BeforeDelete(accounts)
	if err != nil {
		return err
	}

	err = i.repo.DeleteAll(context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete(accounts)
	if err != nil {
		return err
	}

	accountIds := []int{}
	for _, account := range accounts {
		accountIds = append(accountIds, account.ID)
	}

	filter := &usecases.Filter{
		Where: map[string]interface{}{"accountId": accountIds},
	}

	err = i.userInter.DeleteAll(usecases.QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	err = i.sessionInter.DeleteAll(usecases.QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	err = i.roleMappingInter.DeleteAll(usecases.QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	return nil
}

func (i *AccountInter) DeleteByID(id int, context usecases.QueryContext) error {
	account, err := i.repo.FindByID(id, context)
	if err != nil {
		return err
	}

	accounts, err := i.BeforeDelete([]domain.Account{*account})
	if err != nil {
		return err
	}

	account = &accounts[0]

	err = i.repo.DeleteByID(id, context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete([]domain.Account{*account})
	if err != nil {
		return err
	}

	filter := &usecases.Filter{
		Where: map[string]interface{}{"accountId": id},
	}

	err = i.userInter.DeleteAll(usecases.QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	err = i.sessionInter.DeleteAll(usecases.QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	err = i.roleMappingInter.DeleteAll(usecases.QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	return nil
}
