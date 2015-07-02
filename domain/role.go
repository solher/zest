package domain

func init() {
	relations := []DBRelation{
		{Related: "roleMappings"},
		{Related: "aclMappings"},
	}

	ModelDirectory.Register(Role{}, "roles", relations)
}

type Role struct {
	GormModel
	Name         string        `json:"name"`
	RoleMappings []RoleMapping `json:"roleMappings,omitempty"`
	AclMappings  []AclMapping  `json:"aclMappings,omitempty"`
}

func (m *Role) SetRelatedID(idKey string, id int) {
	switch idKey {
	}
}

func (m *Role) BeforeRender() {
	for i := range m.RoleMappings {
		(&m.RoleMappings[i]).BeforeRender()
	}

	for i := range m.AclMappings {
		(&m.AclMappings[i]).BeforeRender()
	}
}
