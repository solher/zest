package domain

import "time"

func init() {
	ModelDirectory.Register(User{})
}

//+gen access controller:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID"
type User struct {
	GormModel
	AccountID int    `json:"accountId,omitempty" sql:"index"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Password  string `json:"-"`
	Email     string `json:"email,omitempty" sql:"unique"`
}

func (m *User) ScopeModel(accountID int) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt

	if accountID != 0 {
		m.AccountID = accountID
	}
}
