package ressources

import (
	"database/sql"
	"strings"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/interfaces"
	"github.com/solher/zest/internalerrors"
	"github.com/solher/zest/usecases"
	"github.com/solher/zest/utils"
)

func init() {
	usecases.DependencyDirectory.Register(NewAccountRepo)
}

type AccountRepo struct {
	store interfaces.AbstractGormStore
}

func NewAccountRepo(store interfaces.AbstractGormStore) *AccountRepo {
	return &AccountRepo{store: store}
}

func (r *AccountRepo) Create(accounts []domain.Account) ([]domain.Account, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, account := range accounts {
		err := db.Create(&account).Error
		if err != nil {
			transaction.Rollback()

			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}

		accounts[i] = account
	}

	transaction.Commit()
	return accounts, nil
}

func (r *AccountRepo) CreateOne(account *domain.Account) (*domain.Account, error) {
	r.Create([]domain.Account{*account})
	return account, nil
}

func (r *AccountRepo) Find(context usecases.QueryContext) ([]domain.Account, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	accounts := []domain.Account{}

	err = query.Find(&accounts).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return accounts, nil
}

func (r *AccountRepo) FindByID(id int, context usecases.QueryContext) (*domain.Account, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	account := domain.Account{}

	err = query.Where(utils.ToDBName("accounts")+".id = ?", id).First(&account).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.InsufficentPermissions
		}

		return nil, internalerrors.DatabaseError
	}

	return &account, nil
}

func (r *AccountRepo) Update(accounts []domain.Account, context usecases.QueryContext) ([]domain.Account, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	for _, account := range accounts {
		queryCopy := *query
		oldAccount := domain.Account{}

		err = queryCopy.Where(utils.ToDBName("accounts")+".id = ?", account.ID).First(&oldAccount).Error
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return nil, internalerrors.InsufficentPermissions
			}

			return nil, internalerrors.DatabaseError
		}

		err = r.store.GetDB().Model(&oldAccount).Updates(&account).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}
	}

	transaction.Commit()
	return accounts, nil
}

func (r *AccountRepo) UpdateByID(id int, account *domain.Account,
	context usecases.QueryContext) (*domain.Account, error) {

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	oldAccount := domain.Account{}

	err = query.Where(utils.ToDBName("accounts")+".id = ?", id).First(&oldAccount).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.InsufficentPermissions
		}

		return nil, internalerrors.DatabaseError
	}

	err = r.store.GetDB().Model(&oldAccount).Updates(&account).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return account, nil
}

func (r *AccountRepo) DeleteAll(context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	accounts := []domain.Account{}
	err = query.Find(&accounts).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	accountIDs := []int{}
	for _, account := range accounts {
		accountIDs = append(accountIDs, account.ID)
	}

	err = r.store.GetDB().Delete(&accounts, utils.ToDBName("accounts")+".id IN (?)", accountIDs).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *AccountRepo) DeleteByID(id int, context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	account := &domain.Account{}

	err = query.Where(utils.ToDBName("accounts")+".id = ?", id).First(&account).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return internalerrors.InsufficentPermissions
		}

		return internalerrors.DatabaseError
	}

	err = r.store.GetDB().Where(utils.ToDBName("accounts")+".id = ?", account.ID).Delete(domain.Account{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *AccountRepo) Raw(query string, values ...interface{}) (*sql.Rows, error) {
	db := r.store.GetDB()

	rows, err := db.Raw(query, values...).Rows()
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return rows, nil
}
