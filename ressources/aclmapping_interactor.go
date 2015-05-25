// Generated by: main
// TypeWriter: interactor
// Directive: +gen on AclMapping

package ressources

import (
	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/interfaces"
)

type AbstractAclMappingRepo interface {
	Create(aclmappings []domain.AclMapping) ([]domain.AclMapping, error)
	CreateOne(aclmapping *domain.AclMapping) (*domain.AclMapping, error)
	Find(filter *interfaces.Filter) ([]domain.AclMapping, error)
	FindByID(id int, filter *interfaces.Filter) (*domain.AclMapping, error)
	Upsert(aclmappings []domain.AclMapping) ([]domain.AclMapping, error)
	UpsertOne(aclmapping *domain.AclMapping) (*domain.AclMapping, error)
	DeleteAll(filter *interfaces.Filter) error
	DeleteByID(id int) error
}

type AclMappingInter struct {
	repo AbstractAclMappingRepo
}

func NewAclMappingInter(repo AbstractAclMappingRepo) *AclMappingInter {
	return &AclMappingInter{repo: repo}
}

func (i *AclMappingInter) Create(aclmappings []domain.AclMapping) ([]domain.AclMapping, error) {
	aclmappings, err := i.repo.Create(aclmappings)
	return aclmappings, err
}

func (i *AclMappingInter) CreateOne(aclmapping *domain.AclMapping) (*domain.AclMapping, error) {
	aclmapping, err := i.repo.CreateOne(aclmapping)
	return aclmapping, err
}

func (i *AclMappingInter) Find(filter *interfaces.Filter) ([]domain.AclMapping, error) {
	aclmappings, err := i.repo.Find(filter)
	return aclmappings, err
}

func (i *AclMappingInter) FindByID(id int, filter *interfaces.Filter) (*domain.AclMapping, error) {
	aclmapping, err := i.repo.FindByID(id, filter)
	return aclmapping, err
}

func (i *AclMappingInter) Upsert(aclmappings []domain.AclMapping) ([]domain.AclMapping, error) {
	aclmappings, err := i.repo.Upsert(aclmappings)
	return aclmappings, err
}

func (i *AclMappingInter) UpsertOne(aclmapping *domain.AclMapping) (*domain.AclMapping, error) {
	aclmapping, err := i.repo.UpsertOne(aclmapping)
	return aclmapping, err
}

func (i *AclMappingInter) DeleteAll(filter *interfaces.Filter) error {
	err := i.repo.DeleteAll(filter)
	return err
}

func (i *AclMappingInter) DeleteByID(id int) error {
	err := i.repo.DeleteByID(id)
	return err
}
