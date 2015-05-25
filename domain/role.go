package domain

import "time"

func init() {
	ModelDirectory.Register(Role{})
	RoleRelated = &Related{ModelName: "roleMapping"}
}

var RoleRelated *Related

//+gen access controller:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID"
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
