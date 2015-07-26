package domain

func init() {
	relations := []DBRelation{
		{Related: "users", Fk: "userId"},
		{Related: "sessions"},
		{Related: "roleMappings"},
	}

	ModelDirectory.Register(Account{}, "accounts", relations)
}

type Account struct {
	GormModel
	UserID       int           `json:"userId" sql:"index"`
	User         User          `json:"user,omitempty"`
	Sessions     []Session     `json:"sessions,omitempty"`
	RoleMappings []RoleMapping `json:"roleMappings,omitempty"`
}

func (m *Account) SetRelatedID(idKey string, id int) {
	switch idKey {
	case "userID":
		m.UserID = id
	}
}

func (m *Account) BeforeRender() {
	m.User.BeforeRender()

	for i := range m.Sessions {
		(&m.Sessions[i]).BeforeRender()
	}

	for i := range m.RoleMappings {
		(&m.RoleMappings[i]).BeforeRender()
	}
}
