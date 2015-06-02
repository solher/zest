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

//+gen repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw"
type Account struct {
	GormModel
	Users        []User        `json:"users,omitempty"`
	Sessions     []Session     `json:"sessions,omitempty"`
	RoleMappings []RoleMapping `json:"roleMappings,omitempty"`
}

func scopeAccount(m *Account) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}

func (m *Account) ValidateCreate() error {
	return nil
}

func (m *Account) ValidateUpdate() error {
	return nil
}

func (m *Account) ValidateDelete() error {
	return nil
}

func (m *Account) BeforeCreate() error {
	scopeAccount(m)
	return nil
}

func (m *Account) AfterCreate() error {
	return nil
}

func (m *Account) BeforeUpdate() error {
	scopeAccount(m)
	return nil
}

func (m *Account) AfterUpdate() error {
	return nil
}

func (m *Account) BeforeDelete() error {
	return nil
}

func (m *Account) AfterDelete() error {
	return nil
}
