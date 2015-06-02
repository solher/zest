package usecases

import (
	"github.com/Solher/zest/domain"
	"github.com/Solher/zest/utils"
)

type AclCacheKey struct {
	Ressource, Method string
}

type AbstractAccountRepo interface {
	Find(filter *Filter, ownerRelations []domain.Relation) ([]domain.Account, error)
	FindByID(id int, filter *Filter, ownerRelations []domain.Relation) (*domain.Account, error)
}

type AbstractAclRepo interface {
	Find(filter *Filter, ownerRelations []domain.Relation) ([]domain.Acl, error)
	FindByID(id int, filter *Filter, ownerRelations []domain.Relation) (*domain.Acl, error)
}

type PermissionCacheInter struct {
	accountRepo         AbstractAccountRepo
	aclRepo             AbstractAclRepo
	roleCache, aclCache AbstractCacheStore
}

func NewPermissionCacheInter(accountRepo AbstractAccountRepo, aclRepo AbstractAclRepo,
	roleCache, aclCache AbstractCacheStore) *PermissionCacheInter {

	return &PermissionCacheInter{accountRepo: accountRepo, aclRepo: aclRepo,
		roleCache: roleCache, aclCache: aclCache}
}

func (i *PermissionCacheInter) GetPermissionRoles(accountID int, ressource, method string) ([]string, error) {
	roleNames := []string{}

	cachedAccountRoles, err := i.roleCache.Get(accountID)
	if err != nil {
		return nil, err
	}

	cachedAclRoles, err := i.aclCache.Get(AclCacheKey{Ressource: ressource, Method: method})
	if err != nil {
		return nil, err
	}

	accountRoles := cachedAccountRoles.([]string)
	aclRoles := cachedAclRoles.([]string)

	if accountID == 0 {
		accountRoles = append(accountRoles, "Guest", "Anyone")
	} else {
		accountRoles = append(accountRoles, "Authenticated", "Owner", "Anyone")
	}

	for _, role := range accountRoles {
		if utils.ContainsStr(aclRoles, role) {
			roleNames = append(roleNames, role)
		}
	}

	return roleNames, nil
}

func (i *PermissionCacheInter) Refresh() error {
	filter := &Filter{
		Include: []interface{}{
			map[string]interface{}{
				"relation": "roleMappings",
				"include":  []interface{}{"role"},
			},
		},
	}

	accounts, err := i.accountRepo.Find(filter, nil)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		roleMappings := account.RoleMappings
		roleNames := []string{}

		for _, roleMapping := range roleMappings {
			roleNames = append(roleNames, roleMapping.Role.Name)
		}

		i.roleCache.Add(account.ID, roleNames)
	}

	filter = &Filter{
		Include: []interface{}{
			map[string]interface{}{
				"relation": "aclMappings",
				"include":  []interface{}{"role"},
			},
		},
	}

	acls, err := i.aclRepo.Find(filter, nil)
	if err != nil {
		return err
	}

	for _, acl := range acls {
		aclMappings := acl.AclMappings
		roleNames := []string{}

		for _, aclMapping := range aclMappings {
			roleNames = append(roleNames, aclMapping.Role.Name)
		}

		i.aclCache.Add(AclCacheKey{Ressource: acl.Ressource, Method: acl.Method}, roleNames)
	}

	return nil
}

func (i *PermissionCacheInter) RefreshRole(accountID int) error {
	filter := &Filter{
		Include: []interface{}{
			map[string]interface{}{
				"relation": "roleMappings",
				"include":  []interface{}{"role"},
			},
		},
	}

	account, err := i.accountRepo.FindByID(accountID, filter, nil)
	if err != nil {
		return err
	}

	roleMappings := account.RoleMappings
	roleNames := []string{}

	for _, roleMapping := range roleMappings {
		roleNames = append(roleNames, roleMapping.Role.Name)
	}

	i.roleCache.Add(account.ID, roleNames)

	return nil
}

func (i *PermissionCacheInter) RefreshAcl(aclID int) error {
	filter := &Filter{
		Include: []interface{}{
			map[string]interface{}{
				"relation": "aclMappings",
				"include":  []interface{}{"role"},
			},
		},
	}

	acl, err := i.aclRepo.FindByID(aclID, filter, nil)
	if err != nil {
		return err
	}

	aclMappings := acl.AclMappings
	roleNames := []string{}

	for _, aclMapping := range aclMappings {
		roleNames = append(roleNames, aclMapping.Role.Name)
	}

	i.aclCache.Add(AclCacheKey{Ressource: acl.Ressource, Method: acl.Method}, roleNames)

	return nil
}
