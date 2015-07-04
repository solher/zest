package resources

import (
	"time"

	"github.com/solher/zest/domain"
)

func (i *AclMappingInter) scopeModel(aclMapping *domain.AclMapping) error {
	aclMapping.CreatedAt = time.Time{}
	aclMapping.UpdatedAt = time.Time{}
	aclMapping.Acl = domain.Acl{}
	aclMapping.Role = domain.Role{}

	return nil
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
		aclMappings[k].ID = 0
		err := i.scopeModel(&aclMappings[k])
		if err != nil {
			return nil, err
		}
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
		err := i.scopeModel(&aclMappings[k])
		if err != nil {
			return nil, err
		}
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
