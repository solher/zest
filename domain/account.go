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

func (m *Account) BeforeActionCreate() error {
	scopeAccount(m)
	return nil
}

func (m *Account) AfterActionCreate() error {
	return nil
}

func (m *Account) BeforeActionUpdate() error {
	scopeAccount(m)
	return nil
}

func (m *Account) AfterActionUpdate() error {
	return nil
}

func (m *Account) BeforeActionDelete() error {
	return nil
}

func (m *Account) AfterActionDelete() error {
	return nil
}
