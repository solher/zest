package domain

import "time"

func init() {
	relations := []Relation{
		{Related: "roleMappings"},
		{Related: "aclMappings"},
	}

	ModelDirectory.Register(Role{}, "roles", relations)
}

var RoleRelated *Relation

//+gen access controller:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID"
type Role struct {
	GormModel
	Name         string        `json:"name,omitempty"`
	RoleMappings []RoleMapping `json:"roleMappings,omitempty"`
	AclMappings  []AclMapping  `json:"aclMappings,omitempty"`
}

func (m *Role) ScopeModel(_ int) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}
