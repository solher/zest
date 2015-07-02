package domain

func init() {
	relations := []DBRelation{
		{Related: "accounts", Fk: "accountId"},
		{Related: "roles", Fk: "roleId"},
	}

	ModelDirectory.Register(RoleMapping{}, "roleMappings", relations)
}

type RoleMapping struct {
	GormModel
	AccountID int     `json:"accountId" sql:"index"`
	Account   Account `json:"account,omitempty"`
	RoleID    int     `json:"roleId" sql:"index"`
	Role      Role    `json:"role,omitempty"`
}

func (m *RoleMapping) SetRelatedID(idKey string, id int) {
	switch idKey {
	case "accountID":
		m.AccountID = id
	case "roleID":
		m.RoleID = id
	}
}

func (m *RoleMapping) BeforeRender() {
	m.Account.BeforeRender()
	m.Role.BeforeRender()
}
