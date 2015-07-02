package domain

func init() {
	relations := []DBRelation{
		{Related: "acls", Fk: "aclId"},
		{Related: "roles", Fk: "roleId"},
	}

	ModelDirectory.Register(AclMapping{}, "aclMappings", relations)
}

type AclMapping struct {
	GormModel
	AclID  int  `json:"aclId" sql:"index"`
	Acl    Acl  `json:"acl,omitempty"`
	RoleID int  `json:"roleId" sql:"index"`
	Role   Role `json:"role,omitempty"`
}

func (m *AclMapping) SetRelatedID(idKey string, id int) {
	switch idKey {
	case "aclID":
		m.AclID = id
	case "roleID":
		m.RoleID = id
	}
}

func (m *AclMapping) BeforeRender() {
	m.Acl.BeforeRender()
	m.Role.BeforeRender()
}
