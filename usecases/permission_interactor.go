package usecases

import "github.com/solher/zest/domain"

type AbstractAclMappingInter interface {
	Find(context QueryContext) ([]domain.AclMapping, error)
	Create(aclMappings []domain.AclMapping) ([]domain.AclMapping, error)
}

type AbstractAclInter interface {
	Find(context QueryContext) ([]domain.Acl, error)
	CreateOne(acl *domain.Acl) (*domain.Acl, error)
}

type AbstractRoleMappingInter interface {
	Find(context QueryContext) ([]domain.RoleMapping, error)
	Create(roleMappings []domain.RoleMapping) ([]domain.RoleMapping, error)
}

type AbstractRoleInter interface {
	Find(context QueryContext) ([]domain.Role, error)
	Create(roles []domain.Role) ([]domain.Role, error)
}

type PermissionInter struct {
	aclInter         AbstractAclInter
	aclMappingInter  AbstractAclMappingInter
	roleInter        AbstractRoleInter
	roleMappingInter AbstractRoleMappingInter
}

func NewPermissionInter(aclInter AbstractAclInter, aclMappingInter AbstractAclMappingInter,
	roleInter AbstractRoleInter, roleMappingInter AbstractRoleMappingInter) *PermissionInter {

	return &PermissionInter{aclInter: aclInter, aclMappingInter: aclMappingInter,
		roleInter: roleInter, roleMappingInter: roleMappingInter}
}

func (i *PermissionInter) SetRole(accountID int, roles ...string) error {
	roleMappings := []domain.RoleMapping{}

	filter := &Filter{
		Limit: 1,
	}

	for _, role := range roles {
		filter.Where = map[string]interface{}{"name": role}

		roles, err := i.roleInter.Find(QueryContext{Filter: filter})
		if err != nil {
			return err
		}

		roleMappings = append(roleMappings, domain.RoleMapping{AccountID: accountID, RoleID: roles[0].ID})
	}

	_, err := i.roleMappingInter.Create(roleMappings)
	if err != nil {
		return err
	}

	return nil
}

func (i *PermissionInter) SetAcl(resource, method string, roles ...string) error {
	filter := &Filter{
		Limit: 1,
		Where: map[string]interface{}{"resource": resource, "method": method},
	}

	acls, err := i.aclInter.Find(QueryContext{Filter: filter})
	if err != nil {
		return err
	}

	aclMappings := []domain.AclMapping{}

	for _, role := range roles {
		filter.Where = map[string]interface{}{"name": role}

		roles, err := i.roleInter.Find(QueryContext{Filter: filter})
		if err != nil {
			return err
		}

		aclMappings = append(aclMappings, domain.AclMapping{RoleID: roles[0].ID, AclID: acls[0].ID})
	}

	_, err = i.aclMappingInter.Create(aclMappings)
	if err != nil {
		return err
	}

	return nil
}

func (i *PermissionInter) RefreshFromRoutes(routes map[DirectoryKey]Route) error {
	for dirKey, route := range routes {
		if !route.CheckPermissions {
			continue
		}

		filter := &Filter{
			Where: map[string]interface{}{
				"resource": dirKey.Ressource,
				"method":    dirKey.Method,
			},
		}

		acls, err := i.aclInter.Find(QueryContext{Filter: filter})
		if err != nil {
			return err
		}

		if len(acls) == 0 {
			acl := &domain.Acl{Ressource: dirKey.Ressource, Method: dirKey.Method}

			acl, err := i.aclInter.CreateOne(acl)
			if err != nil {
				return err
			}

			switch acl.Ressource {
			case "accounts", "sessions", "users", "acls", "aclMappings", "roles", "roleMappings":
				switch acl.Method {
				case "Signin", "Signout", "Signup":
					i.SetAcl(dirKey.Ressource, dirKey.Method, "Admin", "Anyone")
				case "FindByID":
					i.SetAcl(dirKey.Ressource, dirKey.Method, "Admin", "Owner")
				default:
					i.SetAcl(dirKey.Ressource, dirKey.Method, "Admin")
				}
			default:
				i.SetAcl(dirKey.Ressource, dirKey.Method, "Admin", "Owner")
			}
		}
	}

	return nil
}
