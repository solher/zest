package domain

func init() {
	relations := []DBRelation{
		{Related: "accounts", Fk: "accountId"},
	}

	ModelDirectory.Register(User{}, "users", relations)
}

type User struct {
	GormModel
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Password  string  `json:"password"`
	Email     string  `json:"email" sql:"unique"`
	AccountID int     `json:"accountId" sql:"index"`
	Account   Account `json:"account,omitempty"`
}

func (m *User) SetRelatedID(idKey string, id int) {
	switch idKey {
	case "accountID":
		m.AccountID = id
	}
}

func (m *User) BeforeRender() {
	m.Password = ""
	m.Account.BeforeRender()
}
