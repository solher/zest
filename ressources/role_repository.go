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
	usecases.DependencyDirectory.Register(NewRoleRepo)
}

type RoleRepo struct {
	store interfaces.AbstractGormStore
}

func NewRoleRepo(store interfaces.AbstractGormStore) *RoleRepo {
	return &RoleRepo{store: store}
}

func (r *RoleRepo) Create(roles []domain.Role) ([]domain.Role, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, role := range roles {
		err := db.Create(&role).Error
		if err != nil {
			transaction.Rollback()

			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}

		roles[i] = role
	}

	transaction.Commit()
	return roles, nil
}

func (r *RoleRepo) CreateOne(role *domain.Role) (*domain.Role, error) {
	r.Create([]domain.Role{*role})
	return role, nil
}

func (r *RoleRepo) Find(context usecases.QueryContext) ([]domain.Role, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	roles := []domain.Role{}

	err = query.Find(&roles).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return roles, nil
}

func (r *RoleRepo) FindByID(id int, context usecases.QueryContext) (*domain.Role, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	role := domain.Role{}

	err = query.Where(utils.ToDBName("roles")+".id = ?", id).First(&role).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.InsufficentPermissions
		}

		return nil, internalerrors.DatabaseError
	}

	return &role, nil
}

func (r *RoleRepo) Update(roles []domain.Role, context usecases.QueryContext) ([]domain.Role, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	for _, role := range roles {
		queryCopy := *query
		oldRole := domain.Role{}

		err = queryCopy.Where(utils.ToDBName("roles")+".id = ?", role.ID).First(&oldRole).Error
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return nil, internalerrors.InsufficentPermissions
			}

			return nil, internalerrors.DatabaseError
		}

		err = r.store.GetDB().Model(&oldRole).Updates(&role).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}
	}

	transaction.Commit()
	return roles, nil
}

func (r *RoleRepo) UpdateByID(id int, role *domain.Role,
	context usecases.QueryContext) (*domain.Role, error) {

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	oldRole := domain.Role{}

	err = query.Where(utils.ToDBName("roles")+".id = ?", id).First(&oldRole).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.InsufficentPermissions
		}

		return nil, internalerrors.DatabaseError
	}

	err = r.store.GetDB().Model(&oldRole).Updates(&role).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return role, nil
}

func (r *RoleRepo) DeleteAll(context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	roles := []domain.Role{}
	err = query.Find(&roles).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	roleIDs := []int{}
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	err = r.store.GetDB().Delete(&roles, utils.ToDBName("roles")+".id IN (?)", roleIDs).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *RoleRepo) DeleteByID(id int, context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	role := &domain.Role{}

	err = query.Where(utils.ToDBName("roles")+".id = ?", id).First(&role).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return internalerrors.InsufficentPermissions
		}

		return internalerrors.DatabaseError
	}

	err = r.store.GetDB().Where(utils.ToDBName("roles")+".id = ?", role.ID).Delete(domain.Role{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *RoleRepo) Raw(query string, values ...interface{}) (*sql.Rows, error) {
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
