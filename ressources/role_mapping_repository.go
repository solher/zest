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
	usecases.DependencyDirectory.Register(NewRoleMappingRepo)
}

type RoleMappingRepo struct {
	store interfaces.AbstractGormStore
}

func NewRoleMappingRepo(store interfaces.AbstractGormStore) *RoleMappingRepo {
	return &RoleMappingRepo{store: store}
}

func (r *RoleMappingRepo) Create(roleMappings []domain.RoleMapping) ([]domain.RoleMapping, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	for i, roleMapping := range roleMappings {
		err := db.Create(&roleMapping).Error
		if err != nil {
			transaction.Rollback()

			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}

		roleMappings[i] = roleMapping
	}

	transaction.Commit()
	return roleMappings, nil
}

func (r *RoleMappingRepo) CreateOne(roleMapping *domain.RoleMapping) (*domain.RoleMapping, error) {
	r.Create([]domain.RoleMapping{*roleMapping})
	return roleMapping, nil
}

func (r *RoleMappingRepo) Find(context usecases.QueryContext) ([]domain.RoleMapping, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	roleMappings := []domain.RoleMapping{}

	err = query.Find(&roleMappings).Error
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	return roleMappings, nil
}

func (r *RoleMappingRepo) FindByID(id int, context usecases.QueryContext) (*domain.RoleMapping, error) {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	roleMapping := domain.RoleMapping{}

	err = query.Where(utils.ToDBName("roleMappings")+".id = ?", id).First(&roleMapping).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.InsufficentPermissions
		}

		return nil, internalerrors.DatabaseError
	}

	return &roleMapping, nil
}

func (r *RoleMappingRepo) Update(roleMappings []domain.RoleMapping, context usecases.QueryContext) ([]domain.RoleMapping, error) {
	db := r.store.GetDB()
	transaction := db.Begin()

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	for _, roleMapping := range roleMappings {
		queryCopy := *query
		oldRoleMapping := domain.RoleMapping{}

		err = queryCopy.Where(utils.ToDBName("roleMappings")+".id = ?", roleMapping.ID).First(&oldRoleMapping).Error
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return nil, internalerrors.InsufficentPermissions
			}

			return nil, internalerrors.DatabaseError
		}

		err = r.store.GetDB().Model(&oldRoleMapping).Updates(&roleMapping).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}
	}

	transaction.Commit()
	return roleMappings, nil
}

func (r *RoleMappingRepo) UpdateByID(id int, roleMapping *domain.RoleMapping,
	context usecases.QueryContext) (*domain.RoleMapping, error) {

	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return nil, internalerrors.DatabaseError
	}

	oldRoleMapping := domain.RoleMapping{}

	err = query.Where(utils.ToDBName("roleMappings")+".id = ?", id).First(&oldRoleMapping).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, internalerrors.InsufficentPermissions
		}

		return nil, internalerrors.DatabaseError
	}

	err = r.store.GetDB().Model(&oldRoleMapping).Updates(&roleMapping).Error
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return roleMapping, nil
}

func (r *RoleMappingRepo) DeleteAll(context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	roleMappings := []domain.RoleMapping{}
	err = query.Find(&roleMappings).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	roleMappingIDs := []int{}
	for _, roleMapping := range roleMappings {
		roleMappingIDs = append(roleMappingIDs, roleMapping.ID)
	}

	err = r.store.GetDB().Delete(&roleMappings, utils.ToDBName("roleMappings")+".id IN (?)", roleMappingIDs).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *RoleMappingRepo) DeleteByID(id int, context usecases.QueryContext) error {
	query, err := r.store.BuildQuery(context.Filter, context.OwnerRelations)
	if err != nil {
		return internalerrors.DatabaseError
	}

	roleMapping := &domain.RoleMapping{}

	err = query.Where(utils.ToDBName("roleMappings")+".id = ?", id).First(&roleMapping).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return internalerrors.InsufficentPermissions
		}

		return internalerrors.DatabaseError
	}

	err = r.store.GetDB().Where(utils.ToDBName("roleMappings")+".id = ?", roleMapping.ID).Delete(domain.RoleMapping{}).Error
	if err != nil {
		return internalerrors.DatabaseError
	}

	return nil
}

func (r *RoleMappingRepo) Raw(query string, values ...interface{}) (*sql.Rows, error) {
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
