package domain

import "time"

func init() {
	relations := []Relation{
		{Related: "roleMappings"},
		{Related: "aclMappings"},
	}

	ModelDirectory.Register(Role{}, "roles", relations)
}

//+gen access controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID"
type Role struct {
	GormModel
	Name         string        `json:"name,omitempty"`
	RoleMappings []RoleMapping `json:"roleMappings,omitempty"`
	AclMappings  []AclMapping  `json:"aclMappings,omitempty"`
}

func scopeRole(m *Role) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}

func (m *Role) ValidateCreate() error {
	return nil
}

func (m *Role) ValidateUpdate() error {
	return nil
}

func (m *Role) ValidateDelete() error {
	return nil
}

func (m *Role) BeforeCreate() error {
	scopeRole(m)
	return nil
}

func (m *Role) AfterCreate() error {
	return nil
}

func (m *Role) BeforeUpdate() error {
	scopeRole(m)
	return nil
}

func (m *Role) AfterUpdate() error {
	return nil
}

func (m *Role) BeforeDelete() error {
	return nil
}

func (m *Role) AfterDelete() error {
	return nil
}
