package domain

import "golang.org/x/crypto/bcrypt"

func init() {
	relations := []Relation{
		{Related: "accounts", Fk: "accountId"},
	}

	ModelDirectory.Register(User{}, "users", relations)
}

//+gen access controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID"
type User struct {
	GormModel
	AccountID int     `json:"accountId,omitempty" sql:"index"`
	Account   Account `json:"account,omitempty"`
	FirstName string  `json:"firstName,omitempty"`
	LastName  string  `json:"lastName,omitempty"`
	Password  string  `json:"password,omitempty"`
	Email     string  `json:"email,omitempty" sql:"unique"`
}

func (m *User) SetRelatedID(idKey string, id int) {
	switch idKey {
	case "accountID":
		m.AccountID = id
	}
}

func (m *User) ScopeModel() error {
	if m.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), 0)
		if err != nil {
			return err
		}

		m.Password = string(hashedPassword)
	}

	return nil
}

func (m *User) BeforeRender() {
	m.Password = ""
	m.Account.BeforeRender()
}
