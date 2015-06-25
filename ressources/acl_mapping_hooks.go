package ressources

import (
	"time"

	"github.com/solher/zest/domain"
)

func (i *AclMappingInter) scopeModel(aclMapping *domain.AclMapping) {
	aclMapping.ID = 0
	aclMapping.CreatedAt = time.Time{}
	aclMapping.UpdatedAt = time.Time{}
	aclMapping.Acl = domain.Acl{}
	aclMapping.Role = domain.Role{}
}

func (i *AclMappingInter) refreshCache(aclMapping *domain.AclMapping) error {
	err := i.permissionCacheInter.RefreshAcl(aclMapping.AclID)
	if err != nil {
		return err
	}

	return nil
}

func (i *AclMappingInter) BeforeCreate(aclMappings []domain.AclMapping) ([]domain.AclMapping, error) {
	for k := range aclMappings {
		i.scopeModel(&aclMappings[k])
	}
	return aclMappings, nil
}

func (i *AclMappingInter) AfterCreate(aclMappings []domain.AclMapping) ([]domain.AclMapping, error) {
	for k := range aclMappings {
		err := i.refreshCache(&aclMappings[k])
		if err != nil {
			return nil, err
		}
	}
	return aclMappings, nil
}

func (i *AclMappingInter) BeforeUpdate(aclMappings []domain.AclMapping) ([]domain.AclMapping, error) {
	for k := range aclMappings {
		i.scopeModel(&aclMappings[k])
	}
	return aclMappings, nil
}

func (i *AclMappingInter) AfterUpdate(aclMappings []domain.AclMapping) ([]domain.AclMapping, error) {
	for k := range aclMappings {
		err := i.refreshCache(&aclMappings[k])
		if err != nil {
			return nil, err
		}
	}
	return aclMappings, nil
}

func (i *AclMappingInter) BeforeDelete(aclMappings []domain.AclMapping) ([]domain.AclMapping, error) {
	return aclMappings, nil
}

func (i *AclMappingInter) AfterDelete(aclMappings []domain.AclMapping) ([]domain.AclMapping, error) {
	for k := range aclMappings {
		err := i.refreshCache(&aclMappings[k])
		if err != nil {
			return nil, err
		}
	}
	return aclMappings, nil
}
