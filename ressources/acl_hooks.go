package ressources

import (
	"time"

	"github.com/solher/zest/domain"
)

func (i *AclInter) scopeModel(acl *domain.Acl) {
	acl.ID = 0
	acl.CreatedAt = time.Time{}
	acl.UpdatedAt = time.Time{}
	acl.AclMappings = []domain.AclMapping{}
}

func (i *AclInter) BeforeCreate(acls []domain.Acl) ([]domain.Acl, error) {
	for k := range acls {
		i.scopeModel(&acls[k])
	}
	return acls, nil
}

func (i *AclInter) AfterCreate(acls []domain.Acl) ([]domain.Acl, error) {
	return acls, nil
}

func (i *AclInter) BeforeUpdate(acls []domain.Acl) ([]domain.Acl, error) {
	for k := range acls {
		i.scopeModel(&acls[k])
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
