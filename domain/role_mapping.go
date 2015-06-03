package domain

func init() {
	relations := []Relation{
		{Related: "accounts", Fk: "accountId"},
		{Related: "roles", Fk: "roleId"},
	}

	ModelDirectory.Register(RoleMapping{}, "roleMappings", relations)
}

//+gen access controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw"
type RoleMapping struct {
	GormModel
	AccountID int     `json:"accountId,omitempty" sql:"index"`
	Account   Account `json:"account,omitempty"`
	RoleID    int     `json:"roleId,omitempty" sql:"index"`
	Role      Role    `json:"role,omitempty"`
}

func (m *RoleMapping) SetRelatedID(idKey string, id int) {
	switch idKey {
	case "roleID":
		m.RoleID = id
	case "accountID":
		m.AccountID = id
	}
}

func (m *RoleMapping) ScopeModel() error {
	return nil
}

func (m *RoleMapping) BeforeRender() {
	m.Role.BeforeRender()
	m.Account.BeforeRender()
}
