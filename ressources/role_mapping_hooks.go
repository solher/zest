package ressources

import (
	"time"

	"github.com/solher/zest/domain"
)

func (i *RoleMappingInter) scopeModel(roleMapping *domain.RoleMapping) {
	roleMapping.ID = 0
	roleMapping.CreatedAt = time.Time{}
	roleMapping.UpdatedAt = time.Time{}
	roleMapping.Account = domain.Account{}
	roleMapping.Role = domain.Role{}
}

func (i *RoleMappingInter) refreshCache(roleMappings *domain.RoleMapping) error {
	err := i.permissionCacheInter.RefreshRole(roleMappings.AccountID)
	if err != nil {
		return err
	}

	return nil
}

func (i *RoleMappingInter) BeforeCreate(roleMappings []domain.RoleMapping) ([]domain.RoleMapping, error) {
	for k := range roleMappings {
		i.scopeModel(&roleMappings[k])
	}
	return roleMappings, nil
}

func (i *RoleMappingInter) AfterCreate(roleMappings []domain.RoleMapping) ([]domain.RoleMapping, error) {
	for k := range roleMappings {
		err := i.refreshCache(&roleMappings[k])
		if err != nil {
			return nil, err
		}
	}
	return roleMappings, nil
}

func (i *RoleMappingInter) BeforeUpdate(roleMappings []domain.RoleMapping) ([]domain.RoleMapping, error) {
	for k := range roleMappings {
		i.scopeModel(&roleMappings[k])
	}
	return roleMappings, nil
}

func (i *RoleMappingInter) AfterUpdate(roleMappings []domain.RoleMapping) ([]domain.RoleMapping, error) {
	for k := range roleMappings {
		err := i.refreshCache(&roleMappings[k])
		if err != nil {
			return nil, err
		}
	}
	return roleMappings, nil
}

func (i *RoleMappingInter) BeforeDelete(roleMappings []domain.RoleMapping) ([]domain.RoleMapping, error) {
	return roleMappings, nil
}

func (i *RoleMappingInter) AfterDelete(roleMappings []domain.RoleMapping) ([]domain.RoleMapping, error) {
	for k := range roleMappings {
		err := i.refreshCache(&roleMappings[k])
		if err != nil {
			return nil, err
		}
	}
	return roleMappings, nil
}
