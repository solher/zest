package domain

func init() {
	relations := []DBRelation{
		{Related: "RoleMappings"},
		{Related: "AclMappings"},
	}

	ModelDirectory.Register(Role{}, "roles", relations)
}

type Role struct {
	GormModel
	Name         string        `json:"name,omitempty"`
	RoleMappings []RoleMapping `json:"RoleMappings,omitempty"`
	AclMappings  []AclMapping  `json:"AclMappings,omitempty"`
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
