package resources

import (
	"database/sql"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/usecases"
)

func init() {
	usecases.DependencyDirectory.Register(NewAclInter)
	usecases.DependencyDirectory.Register(PopulateAclInter)
}

type AbstractAclRepo interface {
	Create(acls []domain.Acl) ([]domain.Acl, error)
	CreateOne(acl *domain.Acl) (*domain.Acl, error)
	Find(context usecases.QueryContext) ([]domain.Acl, error)
	FindByID(id int, context usecases.QueryContext) (*domain.Acl, error)
	Update(acls []domain.Acl, context usecases.QueryContext) ([]domain.Acl, error)
	UpdateByID(id int, acl *domain.Acl, context usecases.QueryContext) (*domain.Acl, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
	Raw(query string, values ...interface{}) (*sql.Rows, error)
}

type AclInter struct {
	repo            AbstractAclRepo
	aclMappingInter AbstractAclMappingInter
}

func NewAclInter(repo AbstractAclRepo, aclMappingInter AbstractAclMappingInter) *AclInter {
	return &AclInter{repo: repo, aclMappingInter: aclMappingInter}
}

func PopulateAclInter(aclInter *AclInter, repo AbstractAclRepo, aclMappingInter AbstractAclMappingInter) {
	if aclInter.repo == nil {
		aclInter.repo = repo
	}

	if aclInter.aclMappingInter == nil {
		aclInter.aclMappingInter = aclMappingInter
	}
}

func (i *AclInter) Create(acls []domain.Acl) ([]domain.Acl, error) {
	acls, err := i.BeforeCreate(acls)
	if err != nil {
		return nil, err
	}

	acls, err = i.repo.Create(acls)
	if err != nil {
		return nil, err
	}

	acls, err = i.AfterCreate(acls)
	if err != nil {
		return nil, err
	}

	return acls, nil
}

func (i *AclInter) CreateOne(acl *domain.Acl) (*domain.Acl, error) {
	acls, err := i.Create([]domain.Acl{*acl})
	if err != nil {
		return nil, err
	}

	return &acls[0], nil
}

func (i *AclInter) Find(context usecases.QueryContext) ([]domain.Acl, error) {
	acls, err := i.repo.Find(context)
	if err != nil {
		return nil, err
	}

	return acls, nil
}

func (i *AclInter) FindByID(id int, context usecases.QueryContext) (*domain.Acl, error) {
	acl, err := i.repo.FindByID(id, context)
	if err != nil {
		return nil, err
	}

	return acl, nil
}

func (i *AclInter) Upsert(acls []domain.Acl, context usecases.QueryContext) ([]domain.Acl, error) {
	aclsToUpdate := []domain.Acl{}
	aclsToCreate := []domain.Acl{}

	for k := range acls {
		if acls[k].ID != 0 {
			aclsToUpdate = append(aclsToUpdate, acls[k])
		} else {
			aclsToCreate = append(aclsToCreate, acls[k])
		}
	}

	aclsToUpdate, err := i.BeforeUpdate(aclsToUpdate)
	if err != nil {
		return nil, err
	}

	aclsToUpdate, err = i.repo.Update(aclsToUpdate, context)
	if err != nil {
		return nil, err
	}

	aclsToUpdate, err = i.AfterUpdate(aclsToUpdate)
	if err != nil {
		return nil, err
	}

	aclsToCreate, err = i.BeforeCreate(aclsToCreate)
	if err != nil {
		return nil, err
	}

	aclsToCreate, err = i.repo.Create(aclsToCreate)
	if err != nil {
		return nil, err
	}

	aclsToCreate, err = i.AfterCreate(aclsToCreate)
	if err != nil {
		return nil, err
	}

	return append(aclsToUpdate, aclsToCreate...), nil
}

func (i *AclInter) UpsertOne(acl *domain.Acl, context usecases.QueryContext) (*domain.Acl, error) {
	acls, err := i.Upsert([]domain.Acl{*acl}, context)
	if err != nil {
		return nil, err
	}

	return &acls[0], nil
}

func (i *AclInter) UpdateByID(id int, acl *domain.Acl,
	context usecases.QueryContext) (*domain.Acl, error) {

	acls, err := i.BeforeUpdate([]domain.Acl{*acl})
	if err != nil {
		return nil, err
	}

	acl = &acls[0]

	acl, err = i.repo.UpdateByID(id, acl, context)
	if err != nil {
		return nil, err
	}

	acls, err = i.AfterUpdate([]domain.Acl{*acl})
	if err != nil {
		return nil, err
	}

	return &acls[0], nil
}

func (i *AclInter) DeleteAll(context usecases.QueryContext) error {
	acls, err := i.repo.Find(context)
	if err != nil {
		return err
	}

	acls, err = i.BeforeDelete(acls)
	if err != nil {
		return err
	}

	err = i.repo.DeleteAll(context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete(acls)
	if err != nil {
		return err
	}

	aclIds := []int{}
	for _, acl := range acls {
		aclIds = append(aclIds, acl.ID)
	}

	filter := &usecases.Filter{
		Where: map[string]interface{}{"aclId": aclIds},
	}

	err = i.aclMappingInter.DeleteAll(usecases.QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	return nil
}

func (i *AclInter) DeleteByID(id int, context usecases.QueryContext) error {
	acl, err := i.repo.FindByID(id, context)
	if err != nil {
		return err
	}

	acls, err := i.BeforeDelete([]domain.Acl{*acl})
	if err != nil {
		return err
	}

	acl = &acls[0]

	err = i.repo.DeleteByID(id, context)
	if err != nil {
		return err
	}

	_, err = i.AfterDelete([]domain.Acl{*acl})
	if err != nil {
		return err
	}

	filter := &usecases.Filter{
		Where: map[string]interface{}{"aclId": id},
	}

	err = i.aclMappingInter.DeleteAll(usecases.QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	return nil
}
