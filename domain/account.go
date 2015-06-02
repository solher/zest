package domain

func init() {
	relations := []Relation{
		{Related: "users"},
		{Related: "sessions"},
		{Related: "roleMappings"},
	}

	ModelDirectory.Register(Account{}, "accounts", relations)
}

//+gen hooks repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw"
type Account struct {
	GormModel
	Users        []User        `json:"users,omitempty"`
	Sessions     []Session     `json:"sessions,omitempty"`
	RoleMappings []RoleMapping `json:"roleMappings,omitempty"`
}

func (m *Account) SetRelatedID(idKey string, id int) {
}

func (m *Account) ScopeModel() error {
	return nil
}

func (m *Account) BeforeRender() {
	users := m.Users
	sessions := m.Sessions
	roleMappings := m.RoleMappings

	for i := range users {
		(&users[i]).BeforeRender()
	}

	for i := range sessions {
		(&sessions[i]).BeforeRender()
	}

	for i := range roleMappings {
		(&roleMappings[i]).BeforeRender()
	}
}
