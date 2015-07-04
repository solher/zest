package resources

import (
	"time"

	"github.com/solher/zest/domain"
)

func (i *AclInter) scopeModel(acl *domain.Acl) error {
	acl.CreatedAt = time.Time{}
	acl.UpdatedAt = time.Time{}
	acl.AclMappings = []domain.AclMapping{}

	return nil
}

func (i *AclInter) BeforeCreate(acls []domain.Acl) ([]domain.Acl, error) {
	for k := range acls {
		acls[k].ID = 0
		err := i.scopeModel(&acls[k])
		if err != nil {
			return nil, err
		}
	}
	return acls, nil
}

func (i *AclInter) AfterCreate(acls []domain.Acl) ([]domain.Acl, error) {
	return acls, nil
}

func (i *AclInter) BeforeUpdate(acls []domain.Acl) ([]domain.Acl, error) {
	for k := range acls {
		err := i.scopeModel(&acls[k])
		if err != nil {
			return nil, err
		}
	}
	return acls, nil
}

func (i *AclInter) AfterUpdate(acls []domain.Acl) ([]domain.Acl, error) {
	return acls, nil
}

func (i *AclInter) BeforeDelete(acls []domain.Acl) ([]domain.Acl, error) {
	return acls, nil
}

func (i *AclInter) AfterDelete(acls []domain.Acl) ([]domain.Acl, error) {
	return acls, nil
}
