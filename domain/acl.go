package domain

import "time"

func init() {
	relations := []Relation{
		{Related: "aclMappings"},
	}

	ModelDirectory.Register(Acl{}, "acls", relations)
}

//+gen access controller:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID"
type Acl struct {
	GormModel
	Ressource   string       `json:"ressource,omitempty"`
	Method      string       `json:"method,omitempty"`
	AclMappings []AclMapping `json:"aclMappings,omitempty"`
}

func (m *Acl) ScopeModel(_ int) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}
