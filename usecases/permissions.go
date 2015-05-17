package usecases

import "github.com/julienschmidt/httprouter"

type (
	Permissions struct {
		permissions []*httprouter.Handle
	}

	PermissionDirectory map[string]*Permissions
)

func addRoles(permissionDir PermissionDirectory) {
	permissionDir["admin"] = &Permissions{}
	permissionDir["authenticated"] = &Permissions{}
	permissionDir["guest"] = &Permissions{}
}

func NewPermissionDirectory() PermissionDirectory {
	permissionDir := PermissionDirectory{}
	addRoles(permissionDir)

	return permissionDir
}

func (p *Permissions) Add(permission *httprouter.Handle) *Permissions {
	p.permissions = append(p.permissions, permission)

	return p
}

func (p *Permissions) IsGranted(permission *httprouter.Handle) bool {
	if p == nil {
		return false
	}

	for _, perm := range p.permissions {
		if perm == permission {
			return true
		}
	}

	return false
}
