package domain

import "time"

func init() {
	ModelDirectory.Register(Account{})
}

//+gen repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID"
type Account struct {
	GormModel
	Users    []User    `json:"users,omitempty"`
	Sessions []Session `json:"sessions,omitempty"`
}

func (m *Account) ScopeModel() {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}
