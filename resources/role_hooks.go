package resources

import (
	"time"

	"github.com/solher/zest/domain"
)

func (i *RoleInter) scopeModel(role *domain.Role) error {
	role.CreatedAt = time.Time{}
	role.UpdatedAt = time.Time{}
	role.RoleMappings = []domain.RoleMapping{}
	role.AclMappings = []domain.AclMapping{}

	return nil
}

func (i *RoleInter) BeforeCreate(roles []domain.Role) ([]domain.Role, error) {
	for k := range roles {
		roles[k].ID = 0
		err := i.scopeModel(&roles[k])
		if err != nil {
			return nil, err
		}
	}
	return roles, nil
}

func (i *RoleInter) AfterCreate(roles []domain.Role) ([]domain.Role, error) {
	return roles, nil
}

func (i *RoleInter) BeforeUpdate(roles []domain.Role) ([]domain.Role, error) {
	for k := range roles {
		err := i.scopeModel(&roles[k])
		if err != nil {
			return nil, err
		}
	}
	return roles, nil
}

func (i *RoleInter) AfterUpdate(roles []domain.Role) ([]domain.Role, error) {
	return roles, nil
}

func (i *RoleInter) BeforeDelete(roles []domain.Role) ([]domain.Role, error) {
	return roles, nil
}

func (i *RoleInter) AfterDelete(roles []domain.Role) ([]domain.Role, error) {
	return roles, nil
}
