package domain

import "time"

func init() {
	relations := []Relation{
		{Related: "aclMappings"},
	}

	ModelDirectory.Register(Acl{}, "acls", relations)
}

//+gen access controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID"
type Acl struct {
	GormModel
	Ressource   string       `json:"ressource,omitempty"`
	Method      string       `json:"method,omitempty"`
	AclMappings []AclMapping `json:"aclMappings,omitempty"`
}

func scopeAcl(m *Acl) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}

func (m *Acl) ValidateCreate() error {
	return nil
}

func (m *Acl) ValidateUpdate() error {
	return nil
}

func (m *Acl) ValidateDelete() error {
	return nil
}

func (m *Acl) BeforeCreate() error {
	scopeAcl(m)
	return nil
}

func (m *Acl) AfterCreate() error {
	return nil
}

func (m *Acl) BeforeUpdate() error {
	scopeAcl(m)
	return nil
}

func (m *Acl) AfterUpdate() error {
	return nil
}

func (m *Acl) BeforeDelete() error {
	return nil
}

func (m *Acl) AfterDelete() error {
	return nil
}
