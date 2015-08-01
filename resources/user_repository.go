package resources

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
	usecases.DependencyDirectory.Register(NewUserRepo)
}

type UserRepo struct {
	store interfaces.AbstractGormStore
}

func NewUserRepo(store interfaces.AbstractGormStore) *UserRepo {
	return &UserRepo{store: store}
}

func (r *UserRepo) Create(users []domain.User) ([]domain.User, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, user := range users {
		err := db.Create(&user).Error
		if err != nil {
			transaction.Rollback()

			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}

		users[i] = user
	}

	transaction.Commit()
	return users, nil
}

func (r *UserRepo) CreateOne(user *domain.User) (*domain.User, error) {
	users, err := r.Create([]domain.User{*user})
	if err != nil {
		return nil, err
	}

	return &users[0], nil
}

func (r *UserRepo) Find(context usecases.QueryContext) ([]domain.User, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	var users []domain.User

	err = query.Find(&users).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	if len(users) == 0 {
		users = []domain.User{}
	}

	return users, nil
}

func (r *UserRepo) FindByID(id int, context usecases.QueryContext) (*domain.User, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	user := domain.User{}

	err = query.Where(utils.ToDBName("users")+".id = ?", id).First(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.NotFound
		}

		return nil, internalerrors.DatabaseError
	}

	return &user, nil
}

func (r *UserRepo) Update(users []domain.User, context usecases.QueryContext) ([]domain.User, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	for _, user := range users {
		queryCopy := *query

		dbName := utils.ToDBName("users")
		oldUser := &domain.User{}

		err = queryCopy.Where(dbName+".id = ?", user.ID).First(oldUser).Error
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return nil, internalerrors.NotFound
			}

			return nil, internalerrors.DatabaseError
		}

		user.ID = oldUser.ID
		user.CreatedAt = oldUser.CreatedAt

		err = r.store.GetDB().Save(&user).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}
	}

	transaction.Commit()
	return users, nil
}

func (r *UserRepo) UpdateByID(id int, user *domain.User,
	context usecases.QueryContext) (*domain.User, error) {

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	dbName := utils.ToDBName("users")
	oldUser := &domain.User{}

	err = query.Where(dbName+".id = ?", id).First(oldUser).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.NotFound
		}

		return nil, internalerrors.DatabaseError
	}

	user.ID = oldUser.ID
	user.CreatedAt = oldUser.CreatedAt

	err = r.store.GetDB().Save(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return user, nil
}

func (r *UserRepo) UpdateAttributesByID(id int, attributes map[string]interface{},
	context usecases.QueryContext) (*domain.User, error) {

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	dbName := utils.ToDBName("users")
	user := &domain.User{}

	err = query.Where(dbName+".id = ?", id).First(user).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.NotFound
		}

		return nil, internalerrors.DatabaseError
	}

	err = r.store.GetDB().Model(&user).Updates(attributes).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return user, nil
}

func (r *UserRepo) DeleteAll(context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	users := []domain.User{}
	err = query.Find(&users).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	userIDs := []int{}
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	err = r.store.GetDB().Delete(&users, utils.ToDBName("users")+".id IN (?)", userIDs).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *UserRepo) DeleteByID(id int, context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	user := &domain.User{}

	err = query.Where(utils.ToDBName("users")+".id = ?", id).First(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return internalerrors.NotFound
		}

		return internalerrors.DatabaseError
	}

	err = r.store.GetDB().Where(utils.ToDBName("users")+".id = ?", user.ID).Delete(domain.User{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *UserRepo) Raw(query string, values ...interface{}) (*sql.Rows, error) {
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
