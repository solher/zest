package domain

import "time"

func init() {
	relations := []Relation{
		{Related: "users"},
		{Related: "sessions"},
		{Related: "roleMappings"},
	}

	ModelDirectory.Register(Account{}, "accounts", relations)
}

//+gen repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID,Raw"
type Account struct {
	GormModel
	Users        []User        `json:"users,omitempty"`
	Sessions     []Session     `json:"sessions,omitempty"`
	RoleMappings []RoleMapping `json:"roleMappings,omitempty"`
}

func (m *Account) ScopeModel() {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}
