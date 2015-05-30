package usecases

import "github.com/Solher/auth-scaffold/domain"

type AclCacheKey struct {
	Ressource, Method string
}

type AbstractCacheStore interface {
	Add(key interface{}, value interface{}) error
	Remove(key interface{}) error
	Get(key interface{}) (interface{}, error)
	Purge() error
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
