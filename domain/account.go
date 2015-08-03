package domain

func init() {
	relations := []DBRelation{
		{Related: "users"},
		{Related: "sessions"},
		{Related: "roleMappings"},
	}

	ModelDirectory.Register(Account{}, "accounts", relations)
}

type Account struct {
	GormModel
	Users        []User        `json:"users,omitempty"`
	Sessions     []Session     `json:"sessions,omitempty"`
	RoleMappings []RoleMapping `json:"roleMappings,omitempty"`
}

func (m *Account) SetRelatedID(idKey string, id int) {
	switch idKey {
	}
}

func (m *Account) BeforeRender() {
	for i := range m.Users {
		(&m.Users[i]).BeforeRender()
	}

	for i := range m.Sessions {
		(&m.Sessions[i]).BeforeRender()
	}

	for i := range m.RoleMappings {
		(&m.RoleMappings[i]).BeforeRender()
	}
}
