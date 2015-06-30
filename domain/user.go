package domain

func init() {
	relations := []DBRelation{
		{Related: "Accounts", Fk: "AccountId"},
	}

	ModelDirectory.Register(User{}, "users", relations)
}

type User struct {
	GormModel
	FirstName string  `json:"firstName,omitempty"`
	LastName  string  `json:"lastName,omitempty"`
	Password  string  `json:"password,omitempty"`
	Email     string  `json:"email,omitempty" sql:"unique"`
	AccountID int     `json:"AccountId,omitempty" sql:"index"`
	Account   Account `json:"Account,omitempty"`
}

func (m *User) SetRelatedID(idKey string, id int) {
	switch idKey {
	case "AccountID":
		m.AccountID = id
	}
}

func (m *User) BeforeRender() {
	m.Password = ""
	m.Account.BeforeRender()
}
