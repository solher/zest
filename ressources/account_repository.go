// Generated by: main
// TypeWriter: repository
// Directive: +gen on Account

package ressources

import (
	"strings"

	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/internalerrors"
	"github.com/jinzhu/gorm"
)

type AccountRepo struct {
	store interfaces.AbstractGormStore
}

func NewAccountRepo(store interfaces.AbstractGormStore) *AccountRepo {
	return &AccountRepo{store: store}
}

func (r *AccountRepo) Create(accounts []Account) ([]Account, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, account := range accounts {
		err := db.Create(&account).Error
		if err != nil {
			transaction.Rollback()

			if err == gorm.InvalidSql {

			} else {
				return nil, internalerrors.DatabaseError
			}
		}

		accounts[i] = account
	}

	transaction.Commit()
	return accounts, nil
}

func (r *AccountRepo) CreateOne(account *Account) (*Account, error) {
	db := r.store.GetDB()

	err := db.Create(account).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.ViolatedConstraint
		} else {
			return nil, internalerrors.DatabaseError
		}
	}

	return account, nil
}

func (r *AccountRepo) Find(filter *interfaces.Filter) ([]Account, error) {
	query, err := r.store.BuildQuery(filter)
	if err != nil {
		return nil, err
	}

	accounts := []Account{}

	err = query.Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountRepo) FindByID(id int, filter *interfaces.Filter) (*Account, error) {
	query, err := r.store.BuildQuery(filter)
	if err != nil {
		return nil, err
	}

	account := Account{}

	err = query.First(&account, id).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepo) Upsert(accounts []Account) ([]Account, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, account := range accounts {
		if account.ID != 0 {
			oldUser := Account{}

			err := db.First(&oldUser, account.ID).Updates(account).Error
			if err != nil {
				transaction.Rollback()
				return nil, err
			}
		} else {
			err := db.Create(&account).Error
			if err != nil {
				transaction.Rollback()
				return nil, err
			}
		}

		accounts[i] = account
	}

	transaction.Commit()
	return accounts, nil
}

func (r *AccountRepo) UpsertOne(account *Account) (*Account, error) {
	db := r.store.GetDB()

	if account.ID != 0 {
		oldUser := Account{}

		err := db.First(&oldUser, account.ID).Updates(account).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := db.Create(&account).Error
		if err != nil {
			return nil, err
		}
	}

	return account, nil
}

func (r *AccountRepo) DeleteAll(filter *interfaces.Filter) error {
	query, err := r.store.BuildQuery(filter)
	if err != nil {
		return err
	}

	err = query.Delete(Account{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *AccountRepo) DeleteByID(id int) error {
	db := r.store.GetDB()

	err := db.Delete(&Account{GormModel: domain.GormModel{ID: id}}).Error
	if err != nil {
		return err
	}

	return nil
}
