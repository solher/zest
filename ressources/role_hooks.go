package ressources

import (
	"time"

	"github.com/solher/zest/domain"
)

func (i *RoleInter) scopeModel(role *domain.Role) {
	role.ID = 0
	role.CreatedAt = time.Time{}
	role.UpdatedAt = time.Time{}
	role.RoleMappings = []domain.RoleMapping{}
	role.AclMappings = []domain.AclMapping{}
}

func (i *RoleInter) BeforeCreate(roles []domain.Role) ([]domain.Role, error) {
	for k := range roles {
		i.scopeModel(&roles[k])
	}
	return roles, nil
}

func (i *RoleInter) AfterCreate(roles []domain.Role) ([]domain.Role, error) {
	return roles, nil
}

func (i *RoleInter) BeforeUpdate(roles []domain.Role) ([]domain.Role, error) {
	for k := range roles {
		i.scopeModel(&roles[k])
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
