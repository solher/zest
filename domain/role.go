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

func (m *Role) SetRelatedID(idKey string, id int) {
}

func scopeRole(m *Role) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}

func (m *Role) BeforeRender() error {
	roleMappings := m.RoleMappings
	aclMappings := m.AclMappings

	for i := range roleMappings {
		(&roleMappings[i]).BeforeRender()
	}

	for i := range aclMappings {
		(&aclMappings[i]).BeforeRender()
	}

	return nil
}

func (m *Role) BeforeActionCreate() error {
	scopeRole(m)
	return nil
}

func (m *Role) AfterActionCreate() error {
	return nil
}

func (m *Role) BeforeActionUpdate() error {
	scopeRole(m)
	return nil
}

func (m *Role) AfterActionUpdate() error {
	return nil
}

func (m *Role) BeforeActionDelete() error {
	return nil
}

func (m *Role) AfterActionDelete() error {
	return nil
}
