package domain

func init() {
	relations := []DBRelation{
		{Related: "accounts"},
	}

	ModelDirectory.Register(User{}, "users", relations)
}

type User struct {
	GormModel
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Password  string    `json:"password"`
	Email     string    `json:"email" sql:"unique"`
	Accounts  []Account `json:"accounts,omitempty"`
}

func (m *User) SetRelatedID(idKey string, id int) {
	switch idKey {
	}
}

func (m *User) BeforeRender() {
	m.Password = ""
	for i := range m.Accounts {
		(&m.Accounts[i]).BeforeRender()
	}
}
